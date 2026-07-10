package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"event-management/backend/graphql"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
)

type ProcessPaymentArgs struct {
	EventID   string `json:"eventId"`
	Quantity  int    `json:"quantity"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type HasuraProcessPaymentPayload struct {
	Input struct {
		Input ProcessPaymentArgs `json:"input"`
	} `json:"input"`
	SessionVariables map[string]string `json:"session_variables"`
}

type VerifyPaymentArgs struct {
	TransactionRef string `json:"transactionRef"`
}

type HasuraVerifyPaymentPayload struct {
	Input struct {
		Input VerifyPaymentArgs `json:"input"`
	} `json:"input"`
	SessionVariables map[string]string `json:"session_variables"`
}

// ============================================
// CHAPA API MODELS
// ============================================

type ChapaRequest struct {
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Email         string                 `json:"email"`
	FirstName     string                 `json:"first_name"`
	LastName      string                 `json:"last_name"`
	PhoneNumber   string                 `json:"phone_number"`
	TxRef         string                 `json:"tx_ref"`
	CallbackURL   string                 `json:"callback_url"`
	ReturnURL     string                 `json:"return_url"`
	Customization map[string]interface{} `json:"customization"`
	Meta          map[string]interface{} `json:"meta"`
}

type ChapaResponse struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

type ChapaVerifyResponse struct {
	Status string `json:"status"`
	Data   struct {
		Status    string  `json:"status"`
		Reference string  `json:"reference"`
		Amount    float64 `json:"amount"`
	} `json:"data"`
}

// ============================================
// GRAPHQL RESPONSE TYPES
// ============================================

type EventQueryResponse struct {
	Events []struct {
		ID     string  `json:"id"`
		Title  string  `json:"title"`
		Price  float64 `json:"price"`
		Status string  `json:"status"`
	} `json:"events"`
}

type TicketCheckResponse struct {
	Tickets []struct {
		ID string `json:"id"`
	} `json:"tickets"`
}

type TicketInsertResponse struct {
	InsertTicketsOne struct {
		ID string `json:"id"`
	} `json:"insert_tickets_one"`
}

type TicketUpdateResponse struct {
	UpdateTickets struct {
		AffectedRows int `json:"affected_rows"`
	} `json:"update_tickets"`
}

type TicketQueryResponse struct {
	Tickets []struct {
		ID         string  `json:"id"`
		EventID    string  `json:"event_id"`
		Quantity   int     `json:"quantity"`
		TotalPrice float64 `json:"total_price"`
		Status     string  `json:"status"`
	} `json:"tickets"`
}

type PaymentInsertResponse struct {
	InsertPaymentsOne struct {
		ID string `json:"id"`
	} `json:"insert_payments_one"`
}

type PaymentUpdateResponse struct {
	UpdatePayments struct {
		AffectedRows int `json:"affected_rows"`
	} `json:"update_payments"`
}

type PaymentQueryResponse struct {
	Payments []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"payments"`
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func isValidUUID(uuidStr string) bool {
	_, err := guuid.Parse(uuidStr)
	return err == nil
}

func getChapaConfig() (baseURL, secretKey string) {
	baseURL = os.Getenv("CHAPA_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.chapa.co/v1"
	}

	secretKey = os.Getenv("CHAPA_SECRET_KEY")
	if secretKey == "" {
		secretKey = "CHASECK_TEST-OWOC1Ks6sPaaqOnOeFSv0UwxwxwyR7pO"
	}
	return
}

func getReturnURL() string {
	returnURL := os.Getenv("APP_FRONTEND_URL")
	if returnURL == "" {
		returnURL = "http://localhost:3000"
	}
	return returnURL
}

// ============================================
// CORE PAYMENT VERIFICATION LOGIC (Shared)
// ============================================

type PaymentVerificationResult struct {
	Success      bool
	Message      string
	Status       string
	PaymentID    string
	TicketID     string
	EventID      string
	Quantity     int
	TotalPrice   float64
	TicketStatus string
}

func verifyTransactionWithChapa(txRef string) (bool, string) {
	chapaBaseURL, chapaSecret := getChapaConfig()

	httpReq, err := http.NewRequest("GET", chapaBaseURL+"/transaction/verify/"+txRef, nil)
	if err != nil {
		return false, "Failed to create verification request"
	}

	httpReq.Header.Set("Authorization", "Bearer "+chapaSecret)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return false, "Failed to verify with payment gateway"
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var verifyResp ChapaVerifyResponse
	if err := json.Unmarshal(respBody, &verifyResp); err != nil {
		return false, "Failed to parse verification response"
	}

	if verifyResp.Status == "success" && verifyResp.Data.Status == "success" {
		return true, "Payment verified successfully"
	}

	return false, "Payment verification failed"
}

func processPaymentVerification(c *gin.Context, txRef string) PaymentVerificationResult {
	result := PaymentVerificationResult{
		Success: false,
		Message: "Verification failed",
		Status:  "failed",
	}

	if txRef == "" {
		result.Message = "Missing transaction reference"
		return result
	}

	fmt.Printf("[DEBUG] Verifying payment: tx_ref=%s\n", txRef)

	// 1. Get payment by transaction reference
	paymentQuery := `query GetPaymentByRef($tx_ref: String!) {
        payments(where: { transaction_ref: { _eq: $tx_ref } }) {
            id
            status
        }
    }`

	var paymentResp PaymentQueryResponse
	err := graphql.QueryRaw(c.Request.Context(), paymentQuery, map[string]interface{}{
		"tx_ref": txRef,
	}, &paymentResp)

	if err != nil || len(paymentResp.Payments) == 0 {
		result.Message = "Payment not found"
		return result
	}

	payment := paymentResp.Payments[0]
	result.PaymentID = payment.ID

	// 2. Get ticket by payment_id
	ticketQuery := `query GetTicketByPayment($payment_id: uuid!) {
        tickets(where: { payment_id: { _eq: $payment_id } }) {
            id
            event_id
            quantity
            total_price
            status
        }
    }`

	var ticketResp TicketQueryResponse
	err = graphql.QueryRaw(c.Request.Context(), ticketQuery, map[string]interface{}{
		"payment_id": payment.ID,
	}, &ticketResp)

	if err != nil || len(ticketResp.Tickets) == 0 {
		result.Message = "Ticket not found"
		return result
	}

	ticket := ticketResp.Tickets[0]
	result.TicketID = ticket.ID
	result.EventID = ticket.EventID
	result.Quantity = ticket.Quantity
	result.TotalPrice = ticket.TotalPrice
	result.TicketStatus = ticket.Status

	// 3. If payment is already completed, return success
	if payment.Status == "completed" || payment.Status == "confirmed" {
		result.Success = true
		result.Message = "Payment already verified"
		result.Status = payment.Status
		result.TicketStatus = ticket.Status
		return result
	}

	// 4. Verify with Chapa
	chapaSuccess, chapaMessage := verifyTransactionWithChapa(txRef)
	if !chapaSuccess {
		result.Message = chapaMessage
		result.Status = "failed"
		return result
	}

	// 5. Update payment status
	updatePaymentMutation := `mutation UpdatePayment($tx_ref: String!, $status: String!) {
        update_payments(where: { transaction_ref: { _eq: $tx_ref } }, _set: { status: $status }) {
            affected_rows
        }
    }`

	var updatePaymentResp PaymentUpdateResponse
	err = graphql.MutateRaw(c.Request.Context(), updatePaymentMutation, map[string]interface{}{
		"tx_ref": txRef,
		"status": "completed",
	}, &updatePaymentResp)

	if err != nil {
		result.Message = "Failed to update payment status"
		return result
	}

	// 6. Update ticket status
	updateTicketMutation := `mutation UpdateTicket($payment_id: uuid!, $status: String!) {
        update_tickets(where: { payment_id: { _eq: $payment_id } }, _set: { status: $status }) {
            affected_rows
        }
    }`

	var updateTicketResp TicketUpdateResponse
	err = graphql.MutateRaw(c.Request.Context(), updateTicketMutation, map[string]interface{}{
		"payment_id": payment.ID,
		"status":     "confirmed",
	}, &updateTicketResp)

	if err != nil {
		result.Message = "Failed to update ticket status"
		return result
	}

	// 7. Return success
	result.Success = true
	result.Message = "Payment verified successfully"
	result.Status = "completed"
	result.TicketStatus = "confirmed"

	return result
}

// ============================================
// PROCESS PAYMENT HANDLER
// ============================================

func ProcessPayment(c *gin.Context) {
	var payload HasuraProcessPaymentPayload
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
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Missing user metadata execution context",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	if !isValidUUID(args.EventID) {
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

	fmt.Printf("[DEBUG] ProcessPayment - EventID: '%s', UserID: '%s', Quantity: %d\n",
		args.EventID, userID, args.Quantity)

	// 1. Fetch Event
	eventQuery := `query GetEvent($id: uuid!) {
        events(where: { id: { _eq: $id } }, limit: 1) {
            id
            title
            price
            status
        }
    }`

	var eventResp EventQueryResponse
	err := graphql.QueryRaw(c.Request.Context(), eventQuery, map[string]interface{}{
		"id": args.EventID,
	}, &eventResp)

	if err != nil {
		fmt.Printf("[ERROR] GraphQL query failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database query failed",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	if len(eventResp.Events) == 0 {
		fmt.Printf("[ERROR] No event found with ID: %s\n", args.EventID)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Event not found",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	targetEvent := eventResp.Events[0]

	if targetEvent.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Event is not available for booking",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	calculatedTotalPrice := targetEvent.Price * float64(args.Quantity)

	// 2. Check for existing ticket
	duplicateQuery := `query CheckDuplicate($event_id: uuid!, $user_id: uuid!) {
        tickets(where: { 
            event_id: { _eq: $event_id }, 
            user_id: { _eq: $user_id },
            status: { _in: ["pending", "confirmed"] }
        }) {
            id
        }
    }`

	var duplicateResp TicketCheckResponse
	err = graphql.QueryRaw(c.Request.Context(), duplicateQuery, map[string]interface{}{
		"event_id": args.EventID,
		"user_id":  userID,
	}, &duplicateResp)

	if err != nil {
		fmt.Printf("[WARN] Duplicate check failed: %v\n", err)
	} else if len(duplicateResp.Tickets) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "You already have a ticket for this event",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	// 3. Generate transaction reference
	txRef := fmt.Sprintf("EVT-%s-%d-%s",
		args.EventID[:8],
		time.Now().Unix(),
		guuid.New().String()[:8])

	// 4. Initialize Chapa payment
	chapaBaseURL, chapaSecret := getChapaConfig()
	returnURL := getReturnURL()

	chapaReq := ChapaRequest{
		Amount:      calculatedTotalPrice,
		Currency:    "ETB",
		Email:       args.Email,
		FirstName:   args.FirstName,
		LastName:    args.LastName,
		PhoneNumber: args.Phone,
		TxRef:       txRef,
		CallbackURL: "http://localhost:4000/api/payment/callback",
		ReturnURL:   fmt.Sprintf("%s/payment/verify?tx_ref=%s", returnURL, txRef),
		Customization: map[string]interface{}{
			"title":       "Event Ticket",
			"description": fmt.Sprintf("Ticket for %s", targetEvent.Title),
		},
		Meta: map[string]interface{}{
			"event_id": args.EventID,
			"user_id":  userID,
			"quantity": args.Quantity,
		},
	}

	reqBody, err := json.Marshal(chapaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to prepare payment request",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	httpReq, err := http.NewRequestWithContext(c.Request.Context(), "POST", chapaBaseURL+"/transaction/initialize", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to initialize payment",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	httpReq.Header.Set("Authorization", "Bearer "+chapaSecret)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Content-Length", fmt.Sprintf("%d", len(reqBody)))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"message": "Payment gateway unavailable",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// Ground-truth response visibility log
	fmt.Printf("[CHAPA RAW RESPONSE] %s\n", string(respBody))

	var chapaResp ChapaResponse
	if err := json.Unmarshal(respBody, &chapaResp); err != nil || chapaResp.Status != "success" {
		fmt.Printf("[ERROR] Chapa validation rejected request. Body status field: %v\n", chapaResp.Status)

		// Forward Chapa's actual error message to the client to streamline debugging
		var chapaErrorMap map[string]interface{}
		_ = json.Unmarshal(respBody, &chapaErrorMap)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Payment gateway rejected the transaction",
			"details": chapaErrorMap, // Look at this field in your browser console!
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	// 5. Create payment record
	paymentMutation := `mutation InsertPayment($user_id: uuid!, $amount: numeric!, $currency: String!, $status: String!, $transaction_ref: String!, $payment_method: String!) {
        insert_payments_one(object: {
            user_id: $user_id
            amount: $amount
            currency: $currency
            status: $status
            transaction_ref: $transaction_ref
            payment_method: $payment_method
        }) {
            id
        }
	}`

	var paymentResp PaymentInsertResponse
	err = graphql.MutateRaw(c.Request.Context(), paymentMutation, map[string]interface{}{
		"user_id":         userID,
		"amount":          calculatedTotalPrice,
		"currency":        "ETB",
		"status":          "pending",
		"transaction_ref": txRef,
		"payment_method":  "chapa",
	}, &paymentResp)

	if err != nil {
		fmt.Printf("[ERROR] Failed to create payment record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create payment record",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	// 6. Create ticket record
	ticketMutation := `mutation InsertTicket($event_id: uuid!, $user_id: uuid!, $payment_id: uuid!, $quantity: Int!, $total_price: numeric!, $status: String!) {
        insert_tickets_one(object: {
            event_id: $event_id
            user_id: $user_id
            payment_id: $payment_id
            quantity: $quantity
            total_price: $total_price
            status: $status
        }) {
            id
        }
	}`

	var ticketResp TicketInsertResponse
	err = graphql.MutateRaw(c.Request.Context(), ticketMutation, map[string]interface{}{
		"event_id":    args.EventID,
		"user_id":     userID,
		"payment_id":  paymentResp.InsertPaymentsOne.ID,
		"quantity":    args.Quantity,
		"total_price": calculatedTotalPrice,
		"status":      "pending",
	}, &ticketResp)

	if err != nil {
		fmt.Printf("[ERROR] Failed to create ticket: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create ticket record",
			"data": gin.H{
				"checkout_url":    "",
				"transaction_ref": "",
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment initialized successfully",
		"data": gin.H{
			"checkout_url":    chapaResp.Data.CheckoutURL,
			"transaction_ref": txRef,
		},
	})
}

// ============================================
// CHAPA CALLBACK HANDLER
// ============================================

func ChapaCallback(c *gin.Context) {
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

	result := processPaymentVerification(c, txRef)

	c.JSON(http.StatusOK, gin.H{
		"success": result.Success,
		"message": result.Message,
	})
}

// ============================================
// VERIFY PAYMENT HANDLER (Hasura Action)
// ============================================

func VerifyPayment(c *gin.Context) {
	var payload HasuraVerifyPaymentPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		txRef := c.Query("tx_ref")
		if txRef != "" {
			result := processPaymentVerification(c, txRef)
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
	result := processPaymentVerification(c, args.TransactionRef)

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
