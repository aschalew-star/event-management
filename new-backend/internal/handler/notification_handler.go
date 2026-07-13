package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) EventNotificationWebhook(c *gin.Context) {
	// 1. Verify webhook secret
	secret := c.GetHeader("X-Event-Trigger-Secret")
	if secret != os.Getenv("EVENT_TRIGGER_SECRET") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// 2. Parse Hasura event payload
	var payload domain.EventTriggerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload: " + err.Error(),
		})
		return
	}
	fmt.Print(payload)
	// 3. Extract event data
	var event domain.EventNotification
	var operation string

	if payload.Event.Op == "INSERT" {
		// Check if user_id is empty
		if payload.Event.New.UserID == "" {
			fmt.Printf("⚠️ Warning: user_id is empty in INSERT event\n")
			c.JSON(http.StatusOK, gin.H{
				"status":  "skipped",
				"message": "No user_id provided, skipping notification",
			})
			return
		}

		event = domain.EventNotification{
			EventID:     payload.Event.New.ID,
			EventTitle:  payload.Event.New.Title,
			EventUserID: payload.Event.New.UserID,
			EventVenue:  payload.Event.New.Venue,
			EventDate:   payload.Event.New.EventDate,
			EventStatus: payload.Event.New.Status,
			Operation:   "created",
		}
		operation = "created"
	} else if payload.Event.Op == "UPDATE" {
		// Check if user_id is empty
		if payload.Event.New.UserID == "" {
			fmt.Printf("⚠️ Warning: user_id is empty in UPDATE event\n")
			c.JSON(http.StatusOK, gin.H{
				"status":  "skipped",
				"message": "No user_id provided, skipping notification",
			})
			return
		}

		event = domain.EventNotification{
			EventID:     payload.Event.New.ID,
			EventTitle:  payload.Event.New.Title,
			EventUserID: payload.Event.New.UserID,
			EventVenue:  payload.Event.New.Venue,
			EventDate:   payload.Event.New.EventDate,
			EventStatus: payload.Event.New.Status,
			Operation:   "updated",
		}
		operation = "updated"
	} else {
		// Ignore DELETE operations
		c.JSON(http.StatusOK, gin.H{
			"status":  "ignored",
			"message": "Only INSERT and UPDATE are processed",
		})
		return
	}

	fmt.Printf("📨 Event %s: %s by user %s\n", operation, event.EventTitle, event.EventUserID)

	// Debug: Print the full event payload
	fmt.Printf("🔍 Event payload: %+v\n", event)

	// 4. Process asynchronously
	go func() {
		ctx := context.Background()
		if err := h.notificationService.ProcessEventNotification(ctx, event); err != nil {
			fmt.Printf("❌ Failed to process notification: %v\n", err)
		}
	}()

	// 5. Return immediately to Hasura
	c.JSON(http.StatusOK, gin.H{
		"status":   "accepted",
		"event_id": event.EventID,
		"message":  fmt.Sprintf("Event %s notification processing started", operation),
	})
}
