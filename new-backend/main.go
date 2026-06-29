package main

import (
	"event-management/backend/graphql"
	"event-management/backend/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	gqlClient = graphql.NewClient("http://localhost:8080/v1/graphql", "myadminsecretkey")
)

func main() {
	fmt.Println("Starting Auth Service Webhook Server...")

	// Set GraphQL client available globally
	graphql.SetClient(gqlClient)

	// Initialize Gin router
	router := gin.Default()

	// Test Endpoint
	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from test!",
		})
	})

	// Setup routes matching Hasura actions
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Start the server on port 4000
	router.Run(":4000")
}
