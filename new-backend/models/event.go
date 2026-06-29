package models

import (
	"time"

	"github.com/google/uuid"
)

// ============================================
// EVENT MODEL
// ============================================

type Event struct {
	ID          string    `json:"id" db:"id" graphql:"id"`
	Title       string    `json:"title" db:"title" graphql:"title"`
	Description string    `json:"description" db:"description" graphql:"description"`
	CategoryID  *string   `json:"category_id" db:"category_id" graphql:"category_id"`
	Price       *float64  `json:"price" db:"price" graphql:"price"`
	IsFree      bool      `json:"is_free" db:"is_free" graphql:"is_free"`
	Venue       *string   `json:"venue" db:"venue" graphql:"venue"`
	Address     *string   `json:"address" db:"address" graphql:"address"`
	Latitude    *float64  `json:"latitude" db:"latitude" graphql:"latitude"`
	Longitude   *float64  `json:"longitude" db:"longitude" graphql:"longitude"`
	EventDate   *string   `json:"event_date" db:"event_date" graphql:"event_date"`
	StartTime   *string   `json:"start_time" db:"start_time" graphql:"start_time"`
	EndTime     *string   `json:"end_time" db:"end_time" graphql:"end_time"`
	Status      string    `json:"status" db:"status" graphql:"status"`
	Images      []string  `json:"images" db:"images" graphql:"images"`
	UserID      string    `json:"user_id" db:"user_id" graphql:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" graphql:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" graphql:"updated_at"`

	// Relationships (not stored in database)
	Category *Category `json:"category,omitempty" graphql:"category"`
	User     *User     `json:"user,omitempty" graphql:"user"`
}

// ============================================
// CATEGORY MODEL
// ============================================

type Category struct {
	ID          string    `json:"id" db:"id" graphql:"id"`
	Name        string    `json:"name" db:"name" graphql:"name"`
	Description *string   `json:"description" db:"description" graphql:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" graphql:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" graphql:"updated_at"`

	// Relationships
	Events []Event `json:"events,omitempty" graphql:"events"`
}

// ============================================
// USER MODEL
// ============================================

type User struct {
	ID           string    `json:"id" db:"id" graphql:"id"`
	Name         string    `json:"name" db:"name" graphql:"name"`
	Email        string    `json:"email" db:"email" graphql:"email"`
	PasswordHash string    `json:"-" db:"password_hash" graphql:"password_hash"`
	Role         string    `json:"role" db:"role" graphql:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" graphql:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" graphql:"updated_at"`

	// Relationships
	Events []Event `json:"events,omitempty" graphql:"events"`
}

// ============================================
// REQUEST/RESPONSE DTOs
// ============================================

// Create Event Request
type CreateEventRequest struct {
	Title       string   `json:"title" binding:"required" validate:"required,min=3,max=255"`
	Description string   `json:"description" validate:"max=1000"`
	CategoryID  *string  `json:"category_id"`
	Price       *float64 `json:"price" validate:"min=0"`
	IsFree      bool     `json:"is_free"`
	Venue       *string  `json:"venue"`
	Address     *string  `json:"address"`
	Latitude    *float64 `json:"latitude" validate:"min=-90,max=90"`
	Longitude   *float64 `json:"longitude" validate:"min=-180,max=180"`
	EventDate   *string  `json:"event_date"`
	StartTime   *string  `json:"start_time"`
	EndTime     *string  `json:"end_time"`
	Status      *string  `json:"status"`
}

// Update Event Request
type UpdateEventRequest struct {
	Title       *string  `json:"title" validate:"min=3,max=255"`
	Description *string  `json:"description" validate:"max=1000"`
	CategoryID  *string  `json:"category_id"`
	Price       *float64 `json:"price" validate:"min=0"`
	IsFree      *bool    `json:"is_free"`
	Venue       *string  `json:"venue"`
	Address     *string  `json:"address"`
	Latitude    *float64 `json:"latitude" validate:"min=-90,max=90"`
	Longitude   *float64 `json:"longitude" validate:"min=-180,max=180"`
	EventDate   *string  `json:"event_date"`
	StartTime   *string  `json:"start_time"`
	EndTime     *string  `json:"end_time"`
	Status      *string  `json:"status"`
}

// Event Response
type EventResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description *string           `json:"description"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Price       *float64          `json:"price"`
	IsFree      bool              `json:"is_free"`
	Venue       *string           `json:"venue"`
	Address     *string           `json:"address"`
	Latitude    *float64          `json:"latitude"`
	Longitude   *float64          `json:"longitude"`
	EventDate   *string           `json:"event_date"`
	StartTime   *string           `json:"start_time"`
	EndTime     *string           `json:"end_time"`
	Status      string            `json:"status"`
	Images      []string          `json:"images"`
	User        *UserResponse     `json:"user,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Category Response
type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	EventCount  int       `json:"event_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User Response
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================================
// EVENT METHODS
// ============================================

// NewEvent creates a new event instance with default values
func NewEvent(title, userID string) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Title:     title,
		UserID:    userID,
		Status:    "draft",
		IsFree:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsUpcoming checks if the event is upcoming
func (e *Event) IsUpcoming() bool {
	if e.EventDate == nil {
		return false
	}
	eventDate, err := time.Parse("2006-01-02", *e.EventDate)
	if err != nil {
		return false
	}
	return eventDate.After(time.Now())
}

// IsPast checks if the event is past
func (e *Event) IsPast() bool {
	if e.EventDate == nil {
		return false
	}
	eventDate, err := time.Parse("2006-01-02", *e.EventDate)
	if err != nil {
		return false
	}
	return eventDate.Before(time.Now())
}

// IsOngoing checks if the event is happening today
func (e *Event) IsOngoing() bool {
	if e.EventDate == nil {
		return false
	}
	eventDate, err := time.Parse("2006-01-02", *e.EventDate)
	if err != nil {
		return false
	}
	today := time.Now()
	return eventDate.Year() == today.Year() &&
		eventDate.Month() == today.Month() &&
		eventDate.Day() == today.Day()
}

// GetStatus returns the current status of the event
func (e *Event) GetStatus() string {
	if e.Status != "" {
		return e.Status
	}

	if e.IsPast() {
		return "completed"
	}
	if e.IsOngoing() {
		return "ongoing"
	}
	if e.IsUpcoming() {
		return "upcoming"
	}
	return "draft"
}

// HasImages checks if the event has images
func (e *Event) HasImages() bool {
	return len(e.Images) > 0
}

// GetMainImage returns the first image or empty string
func (e *Event) GetMainImage() string {
	if e.HasImages() {
		return e.Images[0]
	}
	return ""
}

// AddImage adds an image URL to the event
func (e *Event) AddImage(imageURL string) {
	e.Images = append(e.Images, imageURL)
	e.UpdatedAt = time.Now()
}

// RemoveImage removes an image URL from the event
func (e *Event) RemoveImage(imageURL string) bool {
	for i, img := range e.Images {
		if img == imageURL {
			e.Images = append(e.Images[:i], e.Images[i+1:]...)
			e.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// ============================================
// CATEGORY METHODS
// ============================================

// NewCategory creates a new category instance
func NewCategory(name string) *Category {
	return &Category{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ============================================
// VALIDATION FUNCTIONS
// ============================================

// Validate validates the event data
func (e *Event) Validate() error {
	if e.Title == "" {
		return ErrTitleRequired
	}
	if len(e.Title) < 3 {
		return ErrTitleTooShort
	}
	if len(e.Title) > 255 {
		return ErrTitleTooLong
	}
	if e.UserID == "" {
		return ErrUserIDRequired
	}
	return nil
}

// ValidateCreateEvent validates create event request
func (r *CreateEventRequest) Validate() error {
	if r.Title == "" {
		return ErrTitleRequired
	}
	if len(r.Title) < 3 {
		return ErrTitleTooShort
	}
	if len(r.Title) > 255 {
		return ErrTitleTooLong
	}
	if r.Price != nil && *r.Price < 0 {
		return ErrInvalidPrice
	}
	if r.Latitude != nil && (*r.Latitude < -90 || *r.Latitude > 90) {
		return ErrInvalidLatitude
	}
	if r.Longitude != nil && (*r.Longitude < -180 || *r.Longitude > 180) {
		return ErrInvalidLongitude
	}
	return nil
}

// ============================================
// CUSTOM ERRORS
// ============================================

var (
	ErrTitleRequired    = &EventError{Code: "TITLE_REQUIRED", Message: "Event title is required"}
	ErrTitleTooShort    = &EventError{Code: "TITLE_TOO_SHORT", Message: "Event title must be at least 3 characters"}
	ErrTitleTooLong     = &EventError{Code: "TITLE_TOO_LONG", Message: "Event title must not exceed 255 characters"}
	ErrUserIDRequired   = &EventError{Code: "USER_ID_REQUIRED", Message: "User ID is required"}
	ErrInvalidPrice     = &EventError{Code: "INVALID_PRICE", Message: "Price must be greater than or equal to 0"}
	ErrInvalidLatitude  = &EventError{Code: "INVALID_LATITUDE", Message: "Latitude must be between -90 and 90"}
	ErrInvalidLongitude = &EventError{Code: "INVALID_LONGITUDE", Message: "Longitude must be between -180 and 180"}
	ErrEventNotFound    = &EventError{Code: "EVENT_NOT_FOUND", Message: "Event not found"}
	ErrUnauthorized     = &EventError{Code: "UNAUTHORIZED", Message: "You are not authorized to perform this action"}
)

type EventError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *EventError) Error() string {
	return e.Message
}

// ============================================
// CONVERTERS
// ============================================

// ToResponse converts Event to EventResponse
func (e *Event) ToResponse() *EventResponse {
	return &EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: &e.Description,
		Price:       e.Price,
		IsFree:      e.IsFree,
		Venue:       e.Venue,
		Address:     e.Address,
		Latitude:    e.Latitude,
		Longitude:   e.Longitude,
		EventDate:   e.EventDate,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Status:      e.GetStatus(),
		Images:      e.Images,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// ToResponseWithRelations converts Event to EventResponse with relationships
func (e *Event) ToResponseWithRelations() *EventResponse {
	resp := e.ToResponse()

	if e.Category != nil {
		resp.Category = e.Category.ToResponse()
	}

	if e.User != nil {
		resp.User = e.User.ToResponse()
	}

	return resp
}

// ToResponse converts Category to CategoryResponse
func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

// ============================================
// SORT AND FILTER HELPERS
// ============================================

// SortOptions for event sorting
type SortOptions struct {
	Field string
	Order string // "asc" or "desc"
}

// FilterOptions for event filtering
type FilterOptions struct {
	CategoryID *string
	Status     *string
	IsFree     *bool
	StartDate  *string
	EndDate    *string
	Search     *string
	UserID     *string
}

// BuildWhereClause builds a GraphQL where clause from filter options
func (f *FilterOptions) BuildWhereClause() map[string]interface{} {
	where := make(map[string]interface{})

	if f.CategoryID != nil {
		where["category_id"] = map[string]interface{}{"_eq": *f.CategoryID}
	}

	if f.Status != nil {
		where["status"] = map[string]interface{}{"_eq": *f.Status}
	}

	if f.IsFree != nil {
		where["is_free"] = map[string]interface{}{"_eq": *f.IsFree}
	}

	if f.StartDate != nil {
		where["event_date"] = map[string]interface{}{"_gte": *f.StartDate}
	}

	if f.EndDate != nil {
		if val, ok := where["event_date"]; ok {
			val.(map[string]interface{})["_lte"] = *f.EndDate
		} else {
			where["event_date"] = map[string]interface{}{"_lte": *f.EndDate}
		}
	}

	if f.Search != nil {
		where["_or"] = []map[string]interface{}{
			{"title": map[string]interface{}{"_ilike": "%" + *f.Search + "%"}},
			{"description": map[string]interface{}{"_ilike": "%" + *f.Search + "%"}},
			{"venue": map[string]interface{}{"_ilike": "%" + *f.Search + "%"}},
			{"address": map[string]interface{}{"_ilike": "%" + *f.Search + "%"}},
		}
	}

	if f.UserID != nil {
		where["user_id"] = map[string]interface{}{"_eq": *f.UserID}
	}

	return where
}
