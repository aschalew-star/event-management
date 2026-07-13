package service

import (
	"context"
	"fmt"

	"event-management/backend/internal/domain"
	"event-management/backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	GetUserByID(ctx context.Context, userID string) (*domain.UserResponse, error)
}

type authService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// ... existing code ...

	// 1. Check if user exists
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to verify existing profile data matching rules: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("an account with this email already exists")
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("internal server error during encryption: %w", err)
	}

	// 3. Create user
	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "user",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create database user record entry: %w", err)
	}

	// Generate JWT token with Hasura claims
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("internal server error during session creation tokens: %w", err)
	}

	return &domain.AuthResponse{
		UserID:  user.ID,
		Token:   token,
		Message: "User registered successfully",
	}, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// 1. Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		fmt.Printf("🔍 ERROR getting user: %v\n", err)
		return nil, fmt.Errorf("database query connection failed: %w", err)
	}

	if user == nil {
		fmt.Printf("🔍 User not found: %s\n", req.Email)
		return nil, fmt.Errorf("email or password is incorrect")
	}

	// Debug: print user info (remove in production)
	fmt.Printf("🔍 User found: ID=%s, Email=%s, Name=%s, Role=%s\n", user.ID, user.Email, user.Name, user.Role)
	fmt.Printf("🔍 Stored password hash: %s\n", user.Password)
	fmt.Printf("🔍 Provided password: %s\n", req.Password)

	// 2. Compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		fmt.Printf("🔍 Password mismatch: %v\n", err)
		return nil, fmt.Errorf("email or password is incorrect")
	}

	fmt.Printf("✅ Password matched successfully\n")

	// 3. Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		fmt.Printf("🔍 Token generation failed: %v\n", err)
		return nil, fmt.Errorf("internal server error generating token: %w", err)
	}

	return &domain.AuthResponse{
		UserID:  user.ID,
		Token:   token,
		Message: "Login successful",
	}, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &domain.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
