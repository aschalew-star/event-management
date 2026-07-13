package domain

import (
	"context"
	// "time"
)

// Request types
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type HasuraRegisterPayload struct {
	Input struct {
		Input RegisterRequest `json:"input"`
	} `json:"input"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type HasuraLoginPayload struct {
	Input struct {
		Input LoginRequest `json:"input"`
	} `json:"input"`
}

// If using string for timestamps
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"` // "-" hides from JSON
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Or use time.Time with custom unmarshaling (recommended)
// but for quick fix, use string

// Response types
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AuthResponse struct {
	UserID  string `json:"user_id"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

// Repository interface
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
	EmailExists(ctx context.Context, email string) (bool, error)
}
