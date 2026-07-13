package router

import (
	"event-management/backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	eventHandler *handler.EventHandler,
	authHandler *handler.AuthHandler,
	imageHandler *handler.ImageHandler,
	paymentHandler *handler.PaymentHandler,
	notificationHandler *handler.NotificationHandler,
) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Auth routes
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// Event routes
	r.POST("/api/events", eventHandler.CreateEvent)

	// Image routes
	r.POST("/api/events/images", imageHandler.UploadEventImages)
	r.DELETE("/api/images/delete", imageHandler.DeleteEventImage)
	r.PUT("/api/images/featured", imageHandler.SetFeaturedImage)

	// Payment routes
	r.POST("/api/payment/initialize", paymentHandler.ProcessPayment)
	r.POST("/api/payment/verify", paymentHandler.VerifyPayment)
	r.POST("/api/payment/callback", paymentHandler.ChapaCallback)

	// Notification routes (webhook from Hasura)
	r.POST("/api/notifications/event", notificationHandler.EventNotificationWebhook)

	return r
}
