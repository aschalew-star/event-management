package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// HasuraClaims defines the specific nested claims Hasura requires
type HasuraClaims struct {
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	DefaultRole  string   `json:"x-hasura-default-role"`
	UserID       string   `json:"x-hasura-user-id"`
}

type Claims struct {
	UserID       string       `json:"user_id"`
	Email        string       `json:"email"`
	Role         string       `json:"role"`
	HasuraClaims HasuraClaims `json:"https://hasura.io/jwt/claims"` // Crucial for Hasura
	jwt.RegisteredClaims
}

// getSecret retrieves the secret, ensuring it perfectly matches Hasura's config string
func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Matching your exact docker-compose string fallback
		secret = "superlongsecretkey_1234567890_abcdefghijklmnopqrstuvwxyz"
	}
	return []byte(secret)
}

// GenerateToken creates a new JWT token recognized by Hasura
func GenerateToken(userID, email, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		HasuraClaims: HasuraClaims{
			AllowedRoles: []string{"user", "anonymous"},
			DefaultRole:  role,
			UserID:       userID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "event-management-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenFromHeader extracts token from Authorization header
func GetTokenFromHeader(header string) (string, error) {
	if len(header) < 7 || header[:7] != "Bearer " {
		return "", errors.New("invalid authorization header format")
	}
	return header[7:], nil
}
