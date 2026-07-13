package domain

import "time"

// Common types
type UUID string
type Numeric float64
type TimestampTZ string

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

// Domain models
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryID  *string   `json:"category_id"`
	Price       float64   `json:"price"`
	IsFree      bool      `json:"is_free"`
	Venue       string    `json:"venue"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	EventDate   string    `json:"event_date"`
	StartTime   *string   `json:"start_time"`
	EndTime     *string   `json:"end_time"`
	Status      string    `json:"status"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Image struct {
	ID         string `json:"id"`
	EventID    string `json:"event_id"`
	ImageURL   string `json:"image_url"`
	IsFeatured bool   `json:"is_featured"`
	PublicID   string `json:"public_id"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ImageUploadResults struct {
	URL        string
	PublicID   string
	IsFeatured bool
}

// Response types
type EventResponse struct {
	ID            string `json:"id"`
	Message       string `json:"message"`
	Success       bool   `json:"success"`
	FeaturedImage string `json:"featured_image"`
	TotalImages   int    `json:"total_images"`
}
