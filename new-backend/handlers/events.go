// handlers/events.go
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

type uuid string
type numeric float64
type timestamptz string

// Request structs
type CreateEventRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	CategoryID    *string  `json:"category_id"`
	Price         float64  `json:"price"`
	IsFree        bool     `json:"is_free"`
	Venue         string   `json:"venue"`
	Address       string   `json:"address"`
	Latitude      float64  `json:"latitude"`
	Longitude     float64  `json:"longitude"`
	EventDate     string   `json:"event_date"`
	StartTime     *string  `json:"start_time"`
	EndTime       *string  `json:"end_time"`
	Status        string   `json:"status"`
	FeaturedImage string   `json:"featured_image"`
	Images        []string `json:"images"`
	Tags          []string `json:"tags"`
}

type CreateEventInputWrapper struct {
	Input CreateEventRequest `json:"input"`
}

type HasuraCreateEventActionPayload struct {
	Input            CreateEventInputWrapper `json:"input"`
	SessionVariables map[string]string       `json:"session_variables"`
}

type ImageUploadResult struct {
	URL        string
	PublicID   string
	IsFeatured bool
}

// GraphQL response structs
type EventMutationResponse struct {
	InsertEventsOne struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		EventDate string `json:"event_date"`
	} `json:"insert_events_one"`
}

type ImageInsertResponse struct {
	InsertEventImages struct {
		AffectedRows int `json:"affected_rows"`
		Returning    []struct {
			ID         string `json:"id"`
			ImageURL   string `json:"image_url"`
			IsFeatured bool   `json:"is_featured"`
		} `json:"returning"`
	} `json:"insert_event_images"`
}

type TagQueryResponse struct {
	Tags []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

type TagCreateResponse struct {
	InsertTagsOne struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"insert_tags_one"`
}

type TagAssociateResponse struct {
	InsertEventTags struct {
		AffectedRows int `json:"affected_rows"`
	} `json:"insert_event_tags"`
}

// Category validation response
type CategoryCheckResponse struct {
	Categories []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
}

// Helper string escaping utility
func escapeGraphQLStr(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}

// Helper to format string values that might be nil/empty safely inline
func formatGraphQLOptionalStr(v interface{}) string {
	if v == nil || v == "" {
		return "null"
	}
	return fmt.Sprintf(`"%s"`, escapeGraphQLStr(fmt.Sprintf("%v", v)))
}

// Helper to format UUID values that might be nil/empty safely inline
func formatGraphQLOptionalUUID(v interface{}) string {
	if v == nil || v == "" {
		return "null"
	}
	return fmt.Sprintf(`"%v"`, v)
}

// validateCategory checks if a category exists in the database
func validateCategory(ctx context.Context, categoryID string) (bool, error) {
	query := `query CheckCategory($id: uuid!) {
categories(where: { id: { _eq: $id } }, limit: 1) {
id
name
}
}`

	var resp CategoryCheckResponse
	err := graphql.QueryRaw(ctx, query, map[string]interface{}{
		"id": categoryID,
	}, &resp)

	if err != nil {
		return false, err
	}

	return len(resp.Categories) > 0, nil
}

// CreateEvent handles the creation of a new event
func CreateEvent(c *gin.Context) {
	var payload HasuraCreateEventActionPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("LOGGING ERROR - JSON BINDING FAILED:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed: " + err.Error(),
			"success": false,
		})
		return
	}

	input := payload.Input.Input
	userID := payload.SessionVariables["x-hasura-user-id"]

	// Debug logging
	fmt.Println("============ HASURA PAYLOAD DEBUG ============")
	if input.CategoryID != nil {
		fmt.Printf("👉 category_id from payload: %q\n", *input.CategoryID)
	} else {
		fmt.Println("👉 category_id from payload: nil")
	}
	fmt.Printf("👉 x-hasura-user-id from session: %q\n", userID)
	fmt.Printf("👉 Number of images: %d\n", len(input.Images))
	fmt.Println("===============================================")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authorization failed: missing user session id from Hasura metadata context headers",
			"success": false,
		})
		return
	}

	// 1. Upload all images to Cloudinary with parallel processing
	var imageURLs []ImageUploadResult
	var featuredImageURL string
	var uploadMutex sync.Mutex

	imagesToUpload := input.Images
	if len(imagesToUpload) == 0 && input.FeaturedImage != "" {
		imagesToUpload = []string{input.FeaturedImage}
	}

	if len(imagesToUpload) > 0 {
		cloudinaryURL := os.Getenv("CLOUDINARY_URL")
		if cloudinaryURL == "" {
			cloudinaryURL = "cloudinary://436377263181225:y53nJ5-e01HR93f2CmcFHKO3lb4@dll1pjjbm"
		}

		cld, err := cloudinary.NewFromURL(cloudinaryURL)
		if err != nil {
			fmt.Println("LOGGING ERROR - CLOUDINARY CONFIG FAILED:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Cloudinary configuration failed",
				"success": false,
			})
			return
		}

		// Define a struct for upload results
		type uploadResultStruct struct {
			index int
			url   string
			id    string
			err   error
		}

		resultChan := make(chan uploadResultStruct, len(imagesToUpload))
		ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
		defer cancel()

		// Upload images in parallel
		var wg sync.WaitGroup
		for i, imageData := range imagesToUpload {
			if strings.TrimSpace(imageData) == "" {
				continue
			}

			wg.Add(1)
			go func(index int, data string) {
				defer wg.Done()

				fmt.Printf("📤 Uploading image %d to Cloudinary...\n", index+1)

				uploadCtx, uploadCancel := context.WithTimeout(ctx, 30*time.Second)
				defer uploadCancel()

				// Upload to Cloudinary
				result, err := cld.Upload.Upload(uploadCtx, data, uploader.UploadParams{
					Folder: "events",
					Eager:  "f_auto,q_auto",
				})
				if err != nil {
					resultChan <- uploadResultStruct{index: index, err: err}
					return
				}
				resultChan <- uploadResultStruct{
					index: index,
					url:   result.SecureURL,
					id:    result.PublicID,
					err:   nil,
				}
			}(i, imageData)
		}

		// Wait for all uploads to complete in a goroutine
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// Collect results
		uploadErrors := 0
		for result := range resultChan {
			if result.err != nil {
				uploadErrors++
				fmt.Printf("LOGGING ERROR - CLOUDINARY UPLOAD FAILED for image %d: %s\n", result.index+1, result.err.Error())

				if result.err == context.DeadlineExceeded {
					c.JSON(http.StatusRequestTimeout, gin.H{
						"message": fmt.Sprintf("Upload for image %d timed out. Please try with smaller images.", result.index+1),
						"success": false,
					})
					return
				}
			} else {
				uploadMutex.Lock()
				isFeatured := result.index == 0 && featuredImageURL == ""
				if isFeatured {
					featuredImageURL = result.url
				}

				imageURLs = append(imageURLs, ImageUploadResult{
					URL:        result.url,
					PublicID:   result.id,
					IsFeatured: isFeatured,
				})
				uploadMutex.Unlock()

				fmt.Printf("✅ Uploaded image %d: %s (featured: %v)\n", result.index+1, result.url, isFeatured)
			}
		}

		// Check if any upload failed
		if uploadErrors > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed to upload %d image(s) to Cloudinary", uploadErrors),
				"success": false,
			})
			return
		}
	}

	// 2. Format event date
	var eventDateFormatted string
	if input.EventDate != "" {
		eventDateFormatted = input.EventDate
		if len(eventDateFormatted) == 10 {
			eventDateFormatted = eventDateFormatted + "T00:00:00Z"
		}
	}

	// 3. Extract and validate category ID
	var categoryIDVal interface{} = nil
	if input.CategoryID != nil && *input.CategoryID != "" && *input.CategoryID != `""` {
		// Validate that the category exists
		exists, err := validateCategory(c.Request.Context(), *input.CategoryID)
		if err != nil {
			fmt.Printf("Error validating category: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to validate category",
				"success": false,
			})
			return
		}

		if exists {
			categoryIDVal = *input.CategoryID
			fmt.Printf("✅ Category validated: %s\n", *input.CategoryID)
		} else {
			fmt.Printf("⚠️ Category %s not found, setting to null\n", *input.CategoryID)
			// Set to null - event will be created without a category
			categoryIDVal = nil
			// OR return error if you want to force a valid category:
			// c.JSON(http.StatusBadRequest, gin.H{
			// 	"message": fmt.Sprintf("Category with ID %s does not exist", *input.CategoryID),
			// 	"success": false,
			// })
			// return
		}
	}

	var startTimeVal interface{} = nil
	if input.StartTime != nil && *input.StartTime != "" {
		startTimeVal = *input.StartTime
	}

	var endTimeVal interface{} = nil
	if input.EndTime != nil && *input.EndTime != "" {
		endTimeVal = *input.EndTime
	}

	// 4. Create event mutation
	eventMutation := `mutation CreateEvent($object: events_insert_input!) {
insert_events_one(object: $object) {
id
title
event_date
}
}`

	eventVars := map[string]interface{}{
		"object": map[string]interface{}{
			"title":       input.Title,
			"description": input.Description,
			"category_id": categoryIDVal,
			"price":       input.Price,
			"is_free":     input.IsFree,
			"venue":       input.Venue,
			"address":     input.Address,
			"latitude":    input.Latitude,
			"longitude":   input.Longitude,
			"event_date":  eventDateFormatted,
			"start_time":  startTimeVal,
			"end_time":    endTimeVal,
			"status":      input.Status,
			"user_id":     userID,
		},
	}

	fmt.Println("🚀 Executing event mutation with variables...")

	var eventResp EventMutationResponse
	err := graphql.MutateRaw(c.Request.Context(), eventMutation, eventVars, &eventResp)
	if err != nil {
		fmt.Println("LOGGING ERROR - EVENT MUTATION WRITE FAILED:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create event entry: " + err.Error(),
			"success": false,
		})
		return
	}

	eventID := eventResp.InsertEventsOne.ID
	if eventID == "" {
		fmt.Println("LOGGING ERROR - Event ID not returned from mutation")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create event: no ID returned",
			"success": false,
		})
		return
	}

	fmt.Printf("✅ Event created with ID: %s\n", eventID)

	// 5. Insert images using variables
	if len(imageURLs) > 0 {
		var imageObjects []map[string]interface{}
		for _, img := range imageURLs {
			imageObjects = append(imageObjects, map[string]interface{}{
				"event_id":    eventID,
				"image_url":   img.URL,
				"is_featured": img.IsFeatured,
			})
		}

		imageMutation := `mutation InsertEventImages($objects: [event_images_insert_input!]!) {
insert_event_images(objects: $objects) {
affected_rows
returning {
id
image_url
is_featured
}
}
}`

		imageVars := map[string]interface{}{
			"objects": imageObjects,
		}

		fmt.Printf("📸 Inserting %d images...\n", len(imageObjects))

		var imageResp ImageInsertResponse
		err := graphql.MutateRaw(c.Request.Context(), imageMutation, imageVars, &imageResp)
		if err != nil {
			fmt.Println("LOGGING ERROR - IMAGE INSERT FAILED:", err.Error())
			fmt.Printf("Failed to insert images for event %s: %v\n", eventID, err)
		} else {
			fmt.Printf("✅ Inserted %d images for event %s\n", imageResp.InsertEventImages.AffectedRows, eventID)
			for _, img := range imageResp.InsertEventImages.Returning {
				fmt.Printf("   - Image ID: %s, URL: %s, Featured: %v\n", img.ID, img.ImageURL, img.IsFeatured)
			}
		}
	}

	// 6. Handle Tags
	if len(input.Tags) > 0 {
		fmt.Printf("📌 Tags to process: %v\n", input.Tags)

		for _, tagName := range input.Tags {
			// First check if tag exists
			getTagQuery := `query GetTag($name: String!) {
tags(where: { name: { _eq: $name } }) {
id
name
}
}`

			tagVars := map[string]interface{}{
				"name": tagName,
			}

			var tagResp TagQueryResponse
			err := graphql.QueryRaw(c.Request.Context(), getTagQuery, tagVars, &tagResp)
			if err != nil {
				fmt.Printf("Error checking tag %s: %v\n", tagName, err)
				continue
			}

			var tagID string
			if len(tagResp.Tags) > 0 {
				tagID = tagResp.Tags[0].ID
				fmt.Printf("Found existing tag: %s (ID: %s)\n", tagName, tagID)
			} else {
				// Create new tag
				createTagMutation := `mutation CreateTag($name: String!) {
insert_tags_one(object: { name: $name }) {
id
name
}
}`

				createTagVars := map[string]interface{}{
					"name": tagName,
				}

				var createTagResp TagCreateResponse
				err := graphql.MutateRaw(c.Request.Context(), createTagMutation, createTagVars, &createTagResp)
				if err != nil {
					fmt.Printf("Error creating tag %s: %v\n", tagName, err)
					continue
				}

				tagID = createTagResp.InsertTagsOne.ID
				fmt.Printf("Created new tag: %s (ID: %s)\n", tagName, tagID)
			}

			if tagID != "" {
				// Associate tag with event
				associateMutation := `mutation AssociateTag($event_id: uuid!, $tag_id: uuid!) {
insert_event_tags(objects: { event_id: $event_id, tag_id: $tag_id }) {
affected_rows
}
}`

				associateVars := map[string]interface{}{
					"event_id": eventID,
					"tag_id":   tagID,
				}

				var associateResp TagAssociateResponse
				err := graphql.MutateRaw(c.Request.Context(), associateMutation, associateVars, &associateResp)
				if err != nil {
					fmt.Printf("Error associating tag %s with event: %v\n", tagName, err)
				} else {
					fmt.Printf("Associated tag %s with event %s\n", tagName, eventID)
				}
			}
		}
	}

	// 7. Return success response
	c.JSON(http.StatusOK, gin.H{
		"id":             eventID,
		"message":        "Event created successfully",
		"success":        true,
		"featured_image": featuredImageURL,
		"total_images":   len(imageURLs),
	})
}

func uploadToCloudinary(ctx context.Context, imageData string, folder string) (string, string, error) {
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
	})
	if err != nil {
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}
