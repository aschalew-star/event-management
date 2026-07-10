// handlers/auth.go
package handlers

import (
	"context"
	"fmt"
	"net/http"

	"event-management/backend/graphql"
	"event-management/backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ============================================
// REQUEST TYPES (Hasura Envelope Matching)
// ============================================

type HasuraRegisterPayload struct {
	Input struct {
		Input RegisterRequest `json:"input"`
	} `json:"input"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type HasuraLoginPayload struct {
	Input struct {
		Input LoginRequest `json:"input"`
	} `json:"input"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func Register(c *gin.Context) {
	var payload HasuraRegisterPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed: " + err.Error(),
		})
		return
	}

	input := payload.Input.Input

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error during encryption",
		})
		return
	}

	// Check if user exists using QueryRaw
	checkQuery := `query CheckUser($email: String!) {
users(where: { email: { _eq: $email } }, limit: 1) {
id
email
}
}`

	var checkResp struct {
		Users []struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"users"`
	}

	checkErr := graphql.QueryRaw(c.Request.Context(), checkQuery, map[string]interface{}{
		"email": input.Email,
	}, &checkResp)

	if checkErr != nil {
		fmt.Println("LOGGING ERROR - REGISTER LOOKUP FAILED:", checkErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to verify existing profile data matching rules: " + checkErr.Error(),
		})
		return
	}

	if len(checkResp.Users) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "An account with this email already exists",
		})
		return
	}

	// Create user using MutateRaw
	createMutation := `mutation CreateUser($email: String!, $password: String!, $name: String!, $role: String!) {
insert_users_one(object: { email: $email, password: $password, name: $name, role: $role }) {
id
email
name
role
created_at
}
}`

	var createResp struct {
		InsertUsersOne struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			Role      string `json:"role"`
			CreatedAt string `json:"created_at"`
		} `json:"insert_users_one"`
	}

	mutateErr := graphql.MutateRaw(c.Request.Context(), createMutation, map[string]interface{}{
		"email":    input.Email,
		"password": string(hashedPassword),
		"name":     input.Name,
		"role":     "user",
	}, &createResp)

	if mutateErr != nil {
		fmt.Println("LOGGING ERROR - REGISTER MUTATION WRITE FAILED:", mutateErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create database user record entry: " + mutateErr.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(
		createResp.InsertUsersOne.ID,
		createResp.InsertUsersOne.Email,
		createResp.InsertUsersOne.Role,
	)

	if err != nil {
		fmt.Println("LOGGING ERROR - REGISTER TOKEN GENERATION FAILED:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error during session creation tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": createResp.InsertUsersOne.ID,
		"token":   token,
		"message": "User registered successfully",
	})
}

// ============================================

func Login(c *gin.Context) {
	var payload HasuraLoginPayload

	// 1. Validate input structure wrapped by Hasura
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("LOGGING ERROR - JSON BINDING FAILED:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed: " + err.Error(),
		})
		return
	}

	input := payload.Input.Input

	// 2. Query user from database using QueryRaw
	getUserQuery := `query GetUser($email: String!) {
users(where: { email: { _eq: $email } }, limit: 1) {
id
email
name
role
password
}
}`

	var getUserResp struct {
		Users []struct {
			ID       string `json:"id"`
			Email    string `json:"email"`
			Name     string `json:"name"`
			Role     string `json:"role"`
			Password string `json:"password"`
		} `json:"users"`
	}

	queryErr := graphql.QueryRaw(c.Request.Context(), getUserQuery, map[string]interface{}{
		"email": input.Email,
	}, &getUserResp)

	if queryErr != nil {
		fmt.Println("LOGGING ERROR - HASURA GRAPHQL QUERY FAILED:", queryErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database query connection failed: " + queryErr.Error(),
		})
		return
	}

	// 3. Check if user exists
	if len(getUserResp.Users) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email or password is incorrect",
		})
		return
	}

	user := getUserResp.Users[0]

	// 4. Compare password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email or password is incorrect",
		})
		return
	}

	// 5. Generate JWT token
	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)

	if err != nil {
		fmt.Println("LOGGING ERROR - TOKEN GENERATION FAILED:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error generating token",
		})
		return
	}

	// 6. Return success response
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user_id": user.ID,
		"message": "Login successful",
	})
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func GetUserByID(ctx context.Context, userID string) (*UserResponse, error) {
	getUserQuery := `query GetUserByID($id: uuid!) {
users_by_pk(id: $id) {
id
name
email
role
}
}`

	var getUserResp struct {
		UsersByPk struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"users_by_pk"`
	}

	err := graphql.QueryRaw(ctx, getUserQuery, map[string]interface{}{
		"id": userID,
	}, &getUserResp)

	if err != nil {
		return nil, err
	}

	if getUserResp.UsersByPk.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return &UserResponse{
		ID:    getUserResp.UsersByPk.ID,
		Name:  getUserResp.UsersByPk.Name,
		Email: getUserResp.UsersByPk.Email,
		Role:  getUserResp.UsersByPk.Role,
	}, nil
}
