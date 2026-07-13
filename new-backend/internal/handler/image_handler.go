package handler

import (
	// "fmt"
	"net/http"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imageService service.ImageService
}

func NewImageHandler(imageService service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// UploadEventImages orchestrates the multi-image parallel uploading pipelines
func (h *ImageHandler) UploadEventImages(c *gin.Context) {
	var payload domain.UploadEventImagesPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload layout: " + err.Error(),
		})
		return
	}

	input := payload.Input
	userID := payload.SessionVariables["x-hasura-user-id"]

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: missing user session",
		})
		return
	}

	eventID := input.EventID
	if eventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Event ID is required",
		})
		return
	}

	// Upload images
	response, err := h.imageService.UploadEventImages(c.Request.Context(), eventID, input.Images)
	if err != nil {
		// Check for specific error types
		errMsg := err.Error()
		if errMsg == "event ID is required" || errMsg == "no images provided to upload" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": errMsg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errMsg,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
		"urls":    response.URLs,
	})
}

// DeleteEventImage handles deletion of an event image
func (h *ImageHandler) DeleteEventImage(c *gin.Context) {
	var payload domain.DeleteEventImagePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload layout: " + err.Error(),
		})
		return
	}

	input := payload.Input
	userID := payload.SessionVariables["x-hasura-user-id"]

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: missing user session",
		})
		return
	}

	if input.ImageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Image ID is required",
		})
		return
	}

	// Delete image
	err := h.imageService.DeleteEventImage(c.Request.Context(), input.ImageID, userID)
	if err != nil {
		errMsg := err.Error()

		switch errMsg {
		case "image ID is required", "image not found":
			c.JSON(http.StatusNotFound, gin.H{
				"message": errMsg,
			})
			return
		case "unauthorized: missing user session":
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": errMsg,
			})
			return
		case "permission denied: you don't own this event":
			c.JSON(http.StatusForbidden, gin.H{
				"message": errMsg,
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": errMsg,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image deleted successfully",
	})
}

// SetFeaturedImage handles setting a featured image for an event
func (h *ImageHandler) SetFeaturedImage(c *gin.Context) {
	var payload domain.SetFeaturedImagePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request payload layout: " + err.Error(),
		})
		return
	}

	input := payload.Input
	userID := payload.SessionVariables["x-hasura-user-id"]

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: missing user session",
		})
		return
	}

	if input.EventID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Event ID is required",
		})
		return
	}

	if input.ImageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Image ID is required",
		})
		return
	}

	// Set featured image
	err := h.imageService.SetFeaturedImage(c.Request.Context(), input.EventID, input.ImageID, userID)
	if err != nil {
		errMsg := err.Error()

		switch errMsg {
		case "event ID is required", "image ID is required", "image not found", "image does not belong to this event":
			c.JSON(http.StatusNotFound, gin.H{
				"message": errMsg,
			})
			return
		case "unauthorized: missing user session":
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": errMsg,
			})
			return
		case "permission denied: you don't own this event":
			c.JSON(http.StatusForbidden, gin.H{
				"message": errMsg,
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": errMsg,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Featured image set successfully",
	})
}
