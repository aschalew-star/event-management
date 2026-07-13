package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/service"

	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var payload domain.HasuraProcessPaymentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Printf("[ERROR] Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Malformed body payload structural syntax",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	args := payload.Input.Input
	userID := payload.SessionVariables["x-hasura-user-id"]

	// Validate UUID
	if args.EventID != "" {
		if _, err := guuid.Parse(args.EventID); err != nil {
			fmt.Printf("[ERROR] Invalid UUID format: '%s'\n", args.EventID)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid event ID format",
				"data": gin.H{
					"checkout_url":    "",
					"transaction_ref": "",
				},
			})
			return
		}
	}

	fmt.Printf("[DEBUG] ProcessPayment - EventID: '%s', UserID: '%s', Quantity: %d\n",
		args.EventID, userID, args.Quantity)

	// Process payment
	response, err := h.paymentService.ProcessPayment(c.Request.Context(), args, userID)
	if err != nil {
		errMsg := err.Error()

		switch {
		case errMsg == "missing user metadata execution context":
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": errMsg,
				"data": gin.H{
					"checkout_url":    "",
					"transaction_ref": "",
				},
			})
			return
		case errMsg == "event not found":
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": errMsg,
				"data": gin.H{
					"checkout_url":    "",
					"transaction_ref": "",
				},
			})
			return
		case errMsg == "event is not available for booking" ||
			errMsg == "you already have a ticket for this event":
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errMsg,
				"data": gin.H{
					"checkout_url":    "",
					"transaction_ref": "",
				},
			})
			return
		default:
			// Check if it's a Chapa error
			if len(errMsg) > 30 && errMsg[:30] == "payment gateway rejected the" {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": errMsg,
					"data": gin.H{
						"checkout_url":    "",
						"transaction_ref": "",
					},
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": errMsg,
				"data": gin.H{
					"checkout_url":    "",
					"transaction_ref": "",
				},
			})
			return
		}
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
		"data": gin.H{
			"checkout_url":    response.CheckoutURL,
			"transaction_ref": response.TransactionRef,
		},
	})
}

func (h *PaymentHandler) ChapaCallback(c *gin.Context) {
	var txRef string

	txRef = c.Query("tx_ref")

	if txRef == "" {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil && len(bodyBytes) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			var rawMap map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &rawMap); err == nil {
				if val, ok := rawMap["tx_ref"].(string); ok {
					txRef = val
				}
				if txRef == "" {
					if dataBlock, ok := rawMap["data"].(map[string]interface{}); ok {
						if val, ok := dataBlock["tx_ref"].(string); ok {
							txRef = val
						}
					}
				}
			}
		}
	}

	if txRef == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Missing transaction reference",
		})
		return
	}

	result, err := h.paymentService.HandleCallback(c.Request.Context(), txRef)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": result.Success,
		"message": result.Message,
	})
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var payload domain.HasuraVerifyPaymentPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		// Try to get tx_ref from query string
		txRef := c.Query("tx_ref")
		if txRef != "" {
			result, err := h.paymentService.VerifyPayment(c.Request.Context(), txRef)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"success": result.Success,
				"message": result.Message,
				"status":  result.Status,
				"ticket": gin.H{
					"id":          result.TicketID,
					"event_id":    result.EventID,
					"quantity":    result.Quantity,
					"total_price": result.TotalPrice,
					"status":      result.TicketStatus,
				},
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	args := payload.Input.Input
	result, err := h.paymentService.VerifyPayment(c.Request.Context(), args.TransactionRef)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": result.Success,
		"message": result.Message,
		"status":  result.Status,
		"ticket": gin.H{
			"id":          result.TicketID,
			"event_id":    result.EventID,
			"quantity":    result.Quantity,
			"total_price": result.TotalPrice,
			"status":      result.TicketStatus,
		},
	})
}
