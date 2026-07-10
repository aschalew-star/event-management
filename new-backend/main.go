// main.go
package main

import (
	"event-management/backend/graphql"
	"event-management/backend/handlers"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	fmt.Println("Starting Auth Service Webhook Server...")

	// Get Hasura endpoint from environment variable
	hasuraEndpoint := os.Getenv("HASURA_GRAPHQL_URL")
	if hasuraEndpoint == "" {
		// For local development (outside Docker), use localhost
		hasuraEndpoint = "http://localhost:8080/v1/graphql"
	}

	hasuraSecret := os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")
	if hasuraSecret == "" {
		hasuraSecret = "myadminsecretkey"
	}

	fmt.Printf("Connecting to Hasura at: %s\n", hasuraEndpoint)

	// Set GraphQL client available globally
	gqlClient := graphql.NewClient(hasuraEndpoint, hasuraSecret)
	graphql.SetClient(gqlClient)

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware (if needed)
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Hasura-Admin-Secret, X-Hasura-Role")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Test Endpoint
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from test!",
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
			"hasura": hasuraEndpoint,
		})
	})

	// Setup routes matching Hasura actions
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	event := router.Group("/api")
	{
		event.POST("/events", handlers.CreateEvent)
		event.POST("/payment/initialize", handlers.ProcessPayment)
		event.POST("/payment/verify", handlers.VerifyPayment)
		event.Any("/payment/callback", handlers.ChapaCallback)
		event.POST("events/images", handlers.UploadEventImages)
	}

	//   r.POST("/api/create-event", handlers.CreateEvent)
	// r.POST("/api/update-event", handlers.UpdateEvent) // Add this if not exists

	// // Image management routes (Hasura Actions)
	// r.POST("/api/upload-event-images", handlers.UploadEventImages)
	// r.POST("/api/delete-event-image", handlers.DeleteEventImage)
	// r.POST("/api/set-featured-image", handlers.SetFeaturedImage)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Create a custom server with timeout
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start the server
	fmt.Printf("Server listening on :%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
	}
}
