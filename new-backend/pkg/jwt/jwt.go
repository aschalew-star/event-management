package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// HasuraClaims represents the custom claims structure Hasura expects
type HasuraClaims struct {
	XHasuraAllowedRoles []string `json:"x-hasura-allowed-roles"`
	XHasuraDefaultRole  string   `json:"x-hasura-default-role"`
	XHasuraUserID       string   `json:"x-hasura-user-id"`
	XHasuraEmail        string   `json:"x-hasura-email,omitempty"`
}

// Claims represents the full JWT claims with Hasura structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
	// Hasura claims - this is the key Hasura looks for
	Hasura HasuraClaims `json:"https://hasura.io/jwt/claims"`
}

func GenerateToken(userID, email, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	// Set expiration (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "event-management",
		},
		// This is the key Hasura looks for
		Hasura: HasuraClaims{
			XHasuraAllowedRoles: []string{"user", "admin", role},
			XHasuraDefaultRole:  role,
			XHasuraUserID:       userID,
			XHasuraEmail:        email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Validate that Hasura claims exist
		if claims.Hasura.XHasuraUserID == "" {
			return nil, fmt.Errorf("missing hasura claims")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractHasuraClaims gets the Hasura claims from a token string
func ExtractHasuraClaims(tokenString string) (*HasuraClaims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &claims.Hasura, nil
}
