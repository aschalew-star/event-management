package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-management/backend/internal/handler"
	"event-management/backend/internal/repository"
	"event-management/backend/internal/router"
	"event-management/backend/internal/service"
	"event-management/backend/pkg/cloudinary"
	"event-management/backend/pkg/graphql"
	"github.com/joho/godotenv"
)

func main() {

	//load env file
	godotenv.Load()
	// Initialize GraphQL client
	endpoint := os.Getenv("HASURA_GRAPHQL_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8080/v1/graphql"
	}

	secret := os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")
	if secret == "" {
		secret = "myadminsecretkey"
	}

	gqlClient := graphql.NewClient(endpoint, secret)
	graphql.SetClient(gqlClient)

	// Initialize Cloudinary client
	cloud, err := cloudinary.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	// Initialize repositories
	eventRepo := repository.NewEventRepository()
	imageRepo := repository.NewImageRepository()
	tagRepo := repository.NewTagRepository()
	userRepo := repository.NewUserRepository()
	paymentRepo := repository.NewPaymentRepository()
	notificationRepo := repository.NewNotificationRepository()

	// Initialize services
	eventService := service.NewEventService(eventRepo, imageRepo, tagRepo, cloud)
	authService := service.NewAuthService(userRepo)
	imageService := service.NewImageService(imageRepo, cloud)
	paymentService := service.NewPaymentService(paymentRepo)
	emailService := service.NewEmailService()
	notificationService := service.NewNotificationService(notificationRepo, emailService)

	// Initialize handlers
	eventHandler := handler.NewEventHandler(eventService)
	authHandler := handler.NewAuthHandler(authService)
	imageHandler := handler.NewImageHandler(imageService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// Setup router
	r := router.SetupRouter(
		eventHandler,
		authHandler,
		imageHandler,
		paymentHandler,
		notificationHandler,
	)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
