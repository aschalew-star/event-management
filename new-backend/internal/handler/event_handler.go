package handler

import (
	"fmt"
	"net/http"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var payload domain.HasuraCreateEventActionPayload

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

	// Process event creation
	response, err := h.eventService.CreateEvent(c.Request.Context(), input, userID)
	if err != nil {
		// Check for specific error types
		if err.Error() == "authorization failed: missing user session id" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization failed: missing user session id from Hasura metadata context headers",
				"success": false,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"id":             response.ID,
		"message":        response.Message,
		"success":        response.Success,
		"featured_image": response.FeaturedImage,
		"total_images":   response.TotalImages,
	})
}
