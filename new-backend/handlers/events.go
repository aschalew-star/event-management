package handlers

import (
	"fmt"
	"net/http"
	"os"

	"event-management/backend/graphql"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type uuid string
type numeric string
type timestamp string

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
	Tags          []string `json:"tags"`
}

type CreateEventInputWrapper struct {
	Input CreateEventRequest `json:"input"`
}

type HasuraCreateEventActionPayload struct {
	Input            CreateEventInputWrapper `json:"input"`
	SessionVariables map[string]string       `json:"session_variables"`
}

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
	var imageUrl string = ""

	// 🔍 CRITICAL RUNTIME TERMINAL PRINTS
	fmt.Println("============ HASURA PAYLOAD DEBUG ============")
	if input.CategoryID != nil {
		fmt.Printf("👉 category_id from payload: %q\n", *input.CategoryID)
	} else {
		fmt.Println("👉 category_id from payload: nil")
	}
	fmt.Printf("👉 x-hasura-user-id from session: %q\n", userID)
	fmt.Println("===============================================")

	// 🛑 DEFENSIVE SECURITY GUARDRAIL: Reject if user_id is missing entirely
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authorization failed: missing user session id from Hasura metadata context headers",
			"success": false,
		})
		return
	}

	// 1. Cloudinary Asset Pipeline
	if input.FeaturedImage != "" {
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

		uploadResult, err := cld.Upload.Upload(c.Request.Context(), input.FeaturedImage, uploader.UploadParams{
			Folder: "events",
		})
		if err != nil {
			fmt.Println("LOGGING ERROR - CLOUDINARY UPLOAD FAILED:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to upload file asset to Cloudinary hosting environment",
				"success": false,
			})
			return
		}

		imageUrl = uploadResult.SecureURL
	}

	// 2. Strict empty-string stripping logic
	var categoryIDInterface interface{} = nil
	if input.CategoryID != nil && *input.CategoryID != "" && *input.CategoryID != `""` {
		categoryIDInterface = uuid(*input.CategoryID)
	}

	var startTimeInterface interface{} = nil
	if input.StartTime != nil && *input.StartTime != "" {
		startTimeInterface = *input.StartTime
	}

	var endTimeInterface interface{} = nil
	if input.EndTime != nil && *input.EndTime != "" {
		endTimeInterface = *input.EndTime
	}

	// 3. Format event_date as ISO 8601 timestamp string
	var eventDateFormatted string
	if input.EventDate != "" {
		eventDateFormatted = input.EventDate
		if len(eventDateFormatted) == 10 { // "YYYY-MM-DD"
			eventDateFormatted = eventDateFormatted + "T00:00:00Z"
		}
	}

	// 4. Mutation schema map
	var mutation struct {
		InsertEventsOne struct {
			ID string `graphql:"id"`
		} `graphql:"insert_events_one(object:{title:$title,description:$description,category_id:$category_id,price:$price,is_free:$is_free,venue:$venue,address:$address,latitude:$latitude,longitude:$longitude,event_date:$event_date,start_time:$start_time,end_time:$end_time,status:$status,user_id:$user_id,featured_image:$featured_image})"`
	}

	// 5. Execute Mutation
	mutateErr := graphql.Mutate(c.Request.Context(), &mutation, map[string]interface{}{
		"title":          input.Title,
		"description":    input.Description,
		"category_id":    categoryIDInterface,
		"price":          numeric(fmt.Sprintf("%.2f", input.Price)),
		"is_free":        input.IsFree,
		"venue":          input.Venue,
		"address":        input.Address,
		"latitude":       numeric(fmt.Sprintf("%.14f", input.Latitude)),
		"longitude":      numeric(fmt.Sprintf("%.14f", input.Longitude)),
		"event_date":     timestamp(eventDateFormatted),
		"start_time":     startTimeInterface,
		"end_time":       endTimeInterface,
		"status":         input.Status,
		"user_id":        uuid(userID), // Guaranteed safe by guardrail above
		"featured_image": imageUrl,
	})

	if mutateErr != nil {
		fmt.Println("LOGGING ERROR - EVENT MUTATION WRITE FAILED:", mutateErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create event entry: " + mutateErr.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             mutation.InsertEventsOne.ID,
		"message":        "Event created successfully",
		"success":        true,
		"featured_image": imageUrl,
	})
}
