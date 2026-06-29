package handlers

import (
	"fmt"
	"net/http"

	"event-management/backend/graphql"

	"github.com/gin-gonic/gin"
)

// ============================================
// REQUEST TYPES
// ============================================

type CreateEventRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
	Price       float64 `json:"price"`
	IsFree      bool    `json:"is_free"`
	Venue       string  `json:"venue"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	EventDate   string  `json:"event_date"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Status      string  `json:"status"`
}

type UpdateEventRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
	Price       float64 `json:"price"`
	IsFree      bool    `json:"is_free"`
	Venue       string  `json:"venue"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	EventDate   string  `json:"event_date"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Status      string  `json:"status"`
}

type EventResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
	Price       float64 `json:"price"`
	IsFree      bool    `json:"is_free"`
	Venue       string  `json:"venue"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	EventDate   string  `json:"event_date"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ============================================
// CREATE EVENT HANDLER
// ============================================

func CreateEvent(c *gin.Context) {
	var input CreateEventRequest

	// 1. Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// 2. Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 3. Create event in database
	var mutation struct {
		InsertEventsOne struct {
			ID          string  `graphql:"id"`
			Title       string  `graphql:"title"`
			Description string  `graphql:"description"`
			Price       float64 `graphql:"price"`
			IsFree      bool    `graphql:"is_free"`
			Venue       string  `graphql:"venue"`
			Address     string  `graphql:"address"`
			EventDate   string  `graphql:"event_date"`
			Status      string  `graphql:"status"`
			CreatedAt   string  `graphql:"created_at"`
		} `graphql:"insert_events_one(object:{title:$title,description:$description,category_id:$category_id,price:$price,is_free:$is_free,venue:$venue,address:$address,latitude:$latitude,longitude:$longitude,event_date:$event_date,start_time:$start_time,end_time:$end_time,status:$status,user_id:$user_id})"`
	}

	err := graphql.Mutate(c.Request.Context(), &mutation, map[string]interface{}{
		"title":       input.Title,
		"description": input.Description,
		"category_id": input.CategoryID,
		"price":       input.Price,
		"is_free":     input.IsFree,
		"venue":       input.Venue,
		"address":     input.Address,
		"latitude":    input.Latitude,
		"longitude":   input.Longitude,
		"event_date":  input.EventDate,
		"start_time":  input.StartTime,
		"end_time":    input.EndTime,
		"status":      input.Status,
		"user_id":     userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create event",
			"details": err.Error(),
		})
		return
	}

	// 4. Return success response
	c.JSON(http.StatusOK, gin.H{
		"id":      mutation.InsertEventsOne.ID,
		"message": "Event created successfully",
	})
}

// ============================================
// UPLOAD EVENT IMAGES HANDLER
// ============================================

type UploadImagesRequest struct {
	Images []string `json:"images" binding:"required,min=1"`
}

func UploadEventImages(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Event ID is required",
		})
		return
	}

	var input UploadImagesRequest

	// 1. Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// 2. Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 3. Check if event exists and user owns it
	var checkQuery struct {
		Events []struct {
			ID     string `graphql:"id"`
			UserID string `graphql:"user_id"`
		} `graphql:"events(where:{id:{_eq:$id}},limit:1)"`
	}

	err := graphql.Query(c.Request.Context(), &checkQuery, map[string]interface{}{
		"id": eventID,
	})

	if err != nil || len(checkQuery.Events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	// 4. Check if user owns the event or is admin
	userRole, _ := c.Get("user_role")
	if checkQuery.Events[0].UserID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to upload images for this event",
		})
		return
	}

	// 5. Here you would process and upload images to cloud storage
	// For now, we'll return dummy URLs
	urls := make([]string, len(input.Images))
	for i, _ := range input.Images {
		// In production, you would upload to S3, Cloudinary, etc.
		// For now, we'll generate dummy URLs
		urls[i] = fmt.Sprintf("https://storage.example.com/events/%s/image_%d.jpg", eventID, i)
	}

	// 6. Update event with image URLs
	var updateMutation struct {
		UpdateEventsByPk struct {
			ID     string   `graphql:"id"`
			Images []string `graphql:"images"`
		} `graphql:"update_events_by_pk(pk_columns:{id:$id},_set:{images:$images})"`
	}

	err = graphql.Mutate(c.Request.Context(), &updateMutation, map[string]interface{}{
		"id":     eventID,
		"images": urls,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update event images",
			"details": err.Error(),
		})
		return
	}

	// 7. Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Images uploaded successfully",
		"urls":    urls,
	})
}
