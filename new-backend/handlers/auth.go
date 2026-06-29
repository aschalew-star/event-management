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

// ============================================
// REGISTER HANDLER
// ============================================
// ============================================
// REGISTER HANDLER - WITH HARDENED LOGGING
// ============================================

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

	var checkQuery struct {
		Users []struct {
			ID    string `graphql:"id"`
			Email string `graphql:"email"`
		} `graphql:"users(where:{email:{_eq:$email}},limit:1)"`
	}

	// Explicitly named error check to prevent tracking cross-contamination
	checkErr := graphql.Query(c.Request.Context(), &checkQuery, map[string]interface{}{
		"email": input.Email,
	})

	if checkErr != nil {
		fmt.Println("LOGGING ERROR - REGISTER LOOKUP FAILED:", checkErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to verify existing profile data matching rules: " + checkErr.Error(),
		})
		return
	}

	if len(checkQuery.Users) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "An account with this email already exists",
		})
		return
	}

	var mutation struct {
		InsertUsersOne struct {
			ID        string `graphql:"id"`
			Email     string `graphql:"email"`
			Name      string `graphql:"name"`
			Role      string `graphql:"role"`
			CreatedAt string `graphql:"created_at"`
		} `graphql:"insert_users_one(object:{email:$email,password:$password,name:$name,role:$role})"`
	}

	mutateErr := graphql.Mutate(c.Request.Context(), &mutation, map[string]interface{}{
		"email":    input.Email,
		"password": string(hashedPassword),
		"name":     input.Name,
		"role":     "user",
	})

	if mutateErr != nil {
		fmt.Println("LOGGING ERROR - REGISTER MUTATION WRITE FAILED:", mutateErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create database user record entry: " + mutateErr.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(
		mutation.InsertUsersOne.ID,
		mutation.InsertUsersOne.Email,
		mutation.InsertUsersOne.Role,
	)

	if err != nil {
		fmt.Println("LOGGING ERROR - REGISTER TOKEN GENERATION FAILED:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error during session creation tokens",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": mutation.InsertUsersOne.ID,
		"token":   token,
		"message": "User registered successfully",
	})
}

// ============================================
// LOGIN HANDLER
// ============================================
// ============================================
// LOGIN HANDLER - ROBUST & EXPLICIT LOGGING
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

	// 2. Query user from database
	var query struct {
		Users []struct {
			ID       string `graphql:"id"`
			Email    string `graphql:"email"`
			Name     string `graphql:"name"`
			Role     string `graphql:"role"`
			Password string `graphql:"password"`
		} `graphql:"users(where:{email:{_eq:$email}},limit:1)"`
	}

	// Explicitly using a separate err variable name to avoid scoping overlaps
	queryErr := graphql.Query(c.Request.Context(), &query, map[string]interface{}{
		"email": input.Email,
	})

	if queryErr != nil {
		// Look at your Go terminal output when you execute the login!
		fmt.Println("LOGGING ERROR - HASURA GRAPHQL QUERY FAILED:", queryErr.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database query connection failed: " + queryErr.Error(),
		})
		return
	}

	// 3. Check if user exists
	if len(query.Users) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email or password is incorrect",
		})
		return
	}

	user := query.Users[0]

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
	var query struct {
		Users []struct {
			ID    string `graphql:"id"`
			Name  string `graphql:"name"`
			Email string `graphql:"email"`
			Role  string `graphql:"role"`
		} `graphql:"users(where:{id:{_eq:$id}},limit:1)"`
	}

	err := graphql.Query(ctx, &query, map[string]interface{}{
		"id": userID,
	})

	if err != nil {
		return nil, err
	}

	if len(query.Users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := query.Users[0]
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
