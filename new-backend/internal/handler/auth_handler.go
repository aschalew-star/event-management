package handler

import (
	"fmt"
	"net/http"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var payload domain.HasuraRegisterPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed: " + err.Error(),
		})
		return
	}

	input := payload.Input.Input

	response, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		// Check for specific error types
		errMsg := err.Error()

		// Handle duplicate email
		if errMsg == "an account with this email already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": errMsg,
			})
			return
		}

		// Handle internal errors
		fmt.Println("LOGGING ERROR - REGISTER:", errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": response.UserID,
		"token":   response.Token,
		"message": response.Message,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload domain.HasuraLoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("LOGGING ERROR - JSON BINDING FAILED:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed: " + err.Error(),
		})
		return
	}

	input := payload.Input.Input

	response, err := h.authService.Login(c.Request.Context(), input)
	if err != nil {
		errMsg := err.Error()

		// Handle authentication errors
		if errMsg == "email or password is incorrect" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": errMsg,
			})
			return
		}

		// Handle other errors
		fmt.Println("LOGGING ERROR - LOGIN:", errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   response.Token,
		"user_id": response.UserID,
		"message": response.Message,
	})
}
