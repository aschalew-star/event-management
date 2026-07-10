package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"event-management/backend/graphql"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type ImageUploadResults struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}

// UploadEventImagesPayload matches the flat 'input' container coming from Hasura
type UploadEventImagesPayload struct {
	Input            UploadEventImagesRequest `json:"input"`
	SessionVariables map[string]string        `json:"session_variables"`
}

type UploadEventImagesRequest struct {
	EventID string   `json:"event_id"`
	Images  []string `json:"images"` // Base64 image payload strings
}

// DeleteEventImagePayload matches standard Hasura action formatting
type DeleteEventImagePayload struct {
	Input            DeleteEventImageRequest `json:"input"`
	SessionVariables map[string]string       `json:"session_variables"`
}

type DeleteEventImageRequest struct {
	ImageID string `json:"image_id"`
}

// SetFeaturedImagePayload matches standard Hasura action formatting
type SetFeaturedImagePayload struct {
	Input            SetFeaturedImageRequest `json:"input"`
	SessionVariables map[string]string       `json:"session_variables"`
}

type SetFeaturedImageRequest struct {
	EventID string `json:"event_id"`
	ImageID string `json:"image_id"`
}

// UploadEventImages orchestrates the multi-image parallel uploading pipelines
func UploadEventImages(c *gin.Context) {
	var payload UploadEventImagesPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload layout: " + err.Error()})
		return
	}

	input := payload.Input
	userID := payload.SessionVariables["x-hasura-user-id"]

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: missing user session"})
		return
	}

	eventID := input.EventID
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Event ID is required"})
		return
	}

	// 1. Verify event ownership against Postgres schema criteria
	checkEventQuery := `query CheckEvent($id: uuid!, $user_id: uuid!) {
        events(where: { id: { _eq: $id }, user_id: { _eq: $user_id } }) {
            id
        }
    }`

	checkEventVars := map[string]interface{}{
		"id":      eventID,
		"user_id": userID,
	}

	var checkResp struct {
		Events []struct {
			ID string `json:"id"`
		} `json:"events"`
	}

	err := graphql.QueryRaw(c.Request.Context(), checkEventQuery, checkEventVars, &checkResp)
	if err != nil || len(checkResp.Events) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Event not found or permission denied"})
		return
	}

	if len(input.Images) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No images provided to upload"})
		return
	}

	// 2. Parallel Cloudinary execution workers
	uploadedImages := []ImageUploadResult{}
	var uploadMutex sync.Mutex
	var uploadErrors []string
	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	type uploadJob struct {
		index int
		data  string
	}

	jobs := make(chan uploadJob, len(input.Images))
	results := make(chan ImageUploadResult, len(input.Images))
	errChan := make(chan error, len(input.Images))

	numWorkers := 5
	if len(input.Images) < 5 {
		numWorkers = len(input.Images)
	}
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				url, _, err := uploadToCloudinaryWithContext(ctx, job.data, "events")
				if err != nil {
					errChan <- fmt.Errorf("failed image %d: %v", job.index+1, err)
					continue
				}
				results <- ImageUploadResult{URL: url}
			}
		}()
	}

	go func() {
		for i, imgData := range input.Images {
			if strings.TrimSpace(imgData) != "" {
				jobs <- uploadJob{index: i, data: imgData}
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
		close(errChan)
	}()

	for err := range errChan {
		uploadErrors = append(uploadErrors, err.Error())
	}

	for result := range results {
		uploadMutex.Lock()
		uploadedImages = append(uploadedImages, result)
		uploadMutex.Unlock()
	}

	if len(uploadErrors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to upload %d image(s)", len(uploadErrors)),
		})
		return
	}

	// 3. Inspect table to see if a featured image is needed
	checkImagesQuery := `query CheckEventImages($eventId: uuid!) {
        event_images(where: { event_id: { _eq: $eventId } }, limit: 1) { id }
    }`

	var checkImagesResp struct {
		EventImages []struct{ ID string } `json:"event_images"`
	}

	_ = graphql.QueryRaw(c.Request.Context(), checkImagesQuery, map[string]interface{}{"eventId": eventID}, &checkImagesResp)
	hasExistingImages := len(checkImagesResp.EventImages) > 0

	var imageObjects []map[string]interface{}
	for i, img := range uploadedImages {
		isFeatured := !hasExistingImages && i == 0
		imageObjects = append(imageObjects, map[string]interface{}{
			"event_id":    eventID,
			"image_url":   img.URL,
			"is_featured": isFeatured,
		})
	}

	// 4. Save metadata back to 'event_images' database table
	insertImagesMutation := `mutation InsertEventImages($objects: [event_images_insert_input!]!) {
        insert_event_images(objects: $objects) {
            returning {
                image_url
            }
        }
    }`

	var insertResp struct {
		InsertEventImages struct {
			Returning []struct {
				ImageURL string `json:"image_url"`
			} `json:"returning"`
		} `json:"insert_event_images"`
	}

	err = graphql.MutateRaw(c.Request.Context(), insertImagesMutation, map[string]interface{}{"objects": imageObjects}, &insertResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database write operation failed: " + err.Error()})
		return
	}

	// Extract generated URLs for the GraphQL action response
	var finalUrls []string
	for _, r := range insertResp.InsertEventImages.Returning {
		finalUrls = append(finalUrls, r.ImageURL)
	}

	// Returns exact match schema payload configuration definitions
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Successfully processed and saved %d records", len(uploadedImages)),
		"urls":    finalUrls,
	})
}

func uploadToCloudinaryWithContext(ctx context.Context, imageData string, folder string) (string, string, error) {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		cloudinaryURL = "cloudinary://436377263181225:y53nJ5-e01HR93f2CmcFHKO3lb4@dll1pjjbm"
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return "", "", err
	}

	uploadResult, err := cld.Upload.Upload(ctx, imageData, uploader.UploadParams{
		Folder: folder,
		Eager:  "f_auto,q_auto",
	})
	if err != nil {
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}
