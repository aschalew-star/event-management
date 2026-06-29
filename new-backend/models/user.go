package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"fullName"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewUser(email, password, fullName string) (*User, error) {
	if email == "" || password == "" || fullName == "" {
		return nil, errors.New("email, password and full name are required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashed),
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stats     UserStats `json:"stats"`
}

type UserStats struct {
	TotalEvents    int `json:"total_events"`
	TotalBookmarks int `json:"total_bookmarks"`
	TotalTickets   int `json:"total_tickets"`
	TotalFollows   int `json:"total_follows"`
}

type UserUpdateRequest struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}
