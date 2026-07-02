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

// Note: 'uuid' and 'numeric' custom type definitions are omitted here
// as they are globally managed from handlers/events.go.

// ============================================
// HASURA ACTION PAYLOAD TYPES
// ============================================

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
// GRAPHQL QUERY/MUTATION TYPES
// ============================================

type GetEventByIdQuery struct {
	Events []struct {
		ID     string  `graphql:"id"`
		Title  string  `graphql:"title"`
		Price  float64 `graphql:"price"`
		Status string  `graphql:"status"`
	} `graphql:"events(where:{id:{_eq:$id}}, limit:1)"`
}

type CheckExistingTicketQuery struct {
	Tickets []struct {
		ID string `graphql:"id"`
	} `graphql:"tickets(where:{event_id:{_eq:$event_id}, user_id:{_eq:$user_id}})"`
}

type InsertTicketMutation struct {
	InsertTicketsOne struct {
		ID string `graphql:"id"`
	} `graphql:"insert_tickets_one(object:{event_id:$event_id, user_id:$user_id, quantity:$quantity, total_price:$total_price, status:$status, transaction_ref:$transaction_ref, payment_id:$payment_id})"`
}

type UpdateTicketMutation struct {
	UpdateTickets struct {
		AffectedRows int `graphql:"affected_rows"`
	} `graphql:"update_tickets(where:{transaction_ref:{_eq:$tx_ref}}, _set:{status:$status, payment_id:$payment_id})"`
}

type GetTicketByRefQuery struct {
	Tickets []struct {
		ID         string  `graphql:"id"`
		EventID    string  `graphql:"event_id"`
		Quantity   int     `graphql:"quantity"`
		TotalPrice float64 `graphql:"total_price"`
		Status     string  `graphql:"status"`
	} `graphql:"tickets(where:{transaction_ref:{_eq:$tx_ref}})"`
}

// ============================================
// HELPER: Validate UUID
// ============================================

func isValidUUID(uuidStr string) bool {
	_, err := guuid.Parse(uuidStr)
	return err == nil
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
		})
		return
	}

	args := payload.Input.Input
	userID := payload.SessionVariables["x-hasura-user-id"]
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Missing user metadata execution context",
		})
		return
	}

	if !isValidUUID(args.EventID) {
		fmt.Printf("[ERROR] Invalid UUID format: '%s'\n", args.EventID)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid event ID format",
		})
		return
	}

	fmt.Printf("[DEBUG] ProcessPayment - EventID: '%s', UserID: '%s', Quantity: %d\n",
		args.EventID, userID, args.Quantity)

	// 1. Fetch Event parameters
	var eventQuery GetEventByIdQuery
	err := graphql.Query(c.Request.Context(), &eventQuery, map[string]interface{}{
		"id": uuid(args.EventID),
	})

	if err != nil {
		fmt.Printf("[ERROR] GraphQL query failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Database query failed: %v", err),
			"debug":   err.Error(),
		})
		return
	}

	if len(eventQuery.Events) == 0 {
		fmt.Printf("[ERROR] No event found with ID: %s\n", args.EventID)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Target event record listing parameters not found",
		})
		return
	}

	targetEvent := eventQuery.Events[0]

	if targetEvent.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Booking channels closed for this item",
		})
		return
	}

	calculatedTotalPrice := targetEvent.Price * float64(args.Quantity)

	// 2. Check for existing duplicate ticket
	var duplicateCheck CheckExistingTicketQuery
	err = graphql.Query(c.Request.Context(), &duplicateCheck, map[string]interface{}{
		"event_id": uuid(args.EventID),
		"user_id":  uuid(userID),
	})

	if err != nil {
		fmt.Printf("[WARN] Duplicate check query failed: %v (continuing)\n", err)
	} else if len(duplicateCheck.Tickets) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "A reservation context matching these details already exists",
		})
		return
	}

	// 3. Generate transaction reference key
	txRef := fmt.Sprintf("EVT-%s-%d-%s",
		args.EventID[:8],
		time.Now().Unix(),
		guuid.New().String()[:8])

	// 4. Initialize Chapa session gateway
	chapaBaseURL := os.Getenv("CHAPA_BASE_URL")
	if chapaBaseURL == "" {
		chapaBaseURL = "https://api.chapa.co/v1"
	}

	chapaSecret := os.Getenv("CHAPA_SECRET_KEY")
	if chapaSecret == "" {
		chapaSecret = "CHASECK_TEST-OWOC1Ks6sPaaqOnOeFSv0UwxwxwyR7pO"
	}

	returnURL := os.Getenv("APP_FRONTEND_URL")
	if returnURL == "" {
		returnURL = "http://localhost:3000"
	}

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
		Meta: map[string]interface{}{
			"event_id": args.EventID,
			"user_id":  userID,
			"quantity": args.Quantity,
		},
	}

	reqBody, err := json.Marshal(chapaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to prepare request"})
		return
	}

	httpReq, err := http.NewRequestWithContext(c.Request.Context(), "POST", chapaBaseURL+"/transaction/initialize", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "HTTP initialization failed"})
		return
	}

	httpReq.Header.Set("Authorization", "Bearer "+chapaSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "Failed to reach payment gateway"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var chapaResp ChapaResponse
	if err := json.Unmarshal(respBody, &chapaResp); err != nil || chapaResp.Status != "success" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Gateway reject validation status"})
		return
	}

	// 5. Create pending local ticket reservation record
	var insertMutation InsertTicketMutation
	err = graphql.Mutate(c.Request.Context(), &insertMutation, map[string]interface{}{
		"event_id":        uuid(args.EventID),
		"user_id":         uuid(userID),
		"quantity":        args.Quantity,
		"total_price":     numeric(fmt.Sprintf("%.2f", calculatedTotalPrice)),
		"status":          "pending",
		"transaction_ref": txRef,
		"payment_id":      chapaResp.Data.CheckoutURL,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database sync layer failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment system initialized successfully",
		"data": gin.H{
			"checkout_url":    chapaResp.Data.CheckoutURL,
			"transaction_ref": txRef,
		},
	})
}

// ============================================
// CHAPA CALLBACK HANDLER (WEBHOOK)
// ============================================

func ChapaCallback(c *gin.Context) {
	var txRef, transactionID string // FIXED: Removed unused status variable definition

	// Parse fallback query arguments first (for direct client browser HTTP GET redirects)
	txRef = c.Query("trx_ref")
	transactionID = c.Query("transaction_id")

	// Fall back to decoding incoming json bodies (for secure webhooks via POST background engines)
	if txRef == "" {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil && len(bodyBytes) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			var rawMap map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &rawMap); err == nil {
				if val, ok := rawMap["tx_ref"].(string); ok {
					txRef = val
				}
				if val, ok := rawMap["transaction_id"].(string); ok {
					transactionID = val
				}

				if txRef == "" {
					if dataBlock, ok := rawMap["data"].(map[string]interface{}); ok {
						if val, ok := dataBlock["tx_ref"].(string); ok {
							txRef = val
						}
						if val, ok := dataBlock["transaction_id"].(string); ok {
							transactionID = val
						}
					}
				}
			}
		}
	}

	if txRef == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing reference tracer parameter mapping identifier"})
		return
	}

	chapaBaseURL := os.Getenv("CHAPA_BASE_URL")
	if chapaBaseURL == "" {
		chapaBaseURL = "https://api.chapa.co/v1"
	}
	chapaSecret := os.Getenv("CHAPA_SECRET_KEY")
	if chapaSecret == "" {
		chapaSecret = "CHASECK_TEST-OWOC1Ks6sPaaqOnOeFSv0UwxwxwyR7pO"
	}

	httpReq, err := http.NewRequestWithContext(c.Request.Context(), "GET", chapaBaseURL+"/transaction/verify/"+txRef, nil)
	if err == nil {
		httpReq.Header.Set("Authorization", "Bearer "+chapaSecret)
		resp, err := http.DefaultClient.Do(httpReq)
		if err == nil {
			defer resp.Body.Close()
			respBody, _ := io.ReadAll(resp.Body)
			var verifyResp ChapaVerifyResponse
			if json.Unmarshal(respBody, &verifyResp) == nil {
				finalStatus := "failed"
				if verifyResp.Status == "success" && verifyResp.Data.Status == "success" {
					finalStatus = "confirmed"
				}

				var updateMutation UpdateTicketMutation
				_ = graphql.Mutate(c.Request.Context(), &updateMutation, map[string]interface{}{
					"tx_ref":     txRef,
					"status":     finalStatus,
					"payment_id": transactionID,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Webhook payload completed resolution step"})
}

// ============================================
// VERIFY PAYMENT HANDLER (HASURA ACTION COMPLIANT)
// ============================================

func VerifyPayment(c *gin.Context) {
	txRef := c.Query("tx_ref")

	// Parse deep nested structure matching Hasura Action variables payload format
	if txRef == "" {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil && len(bodyBytes) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var payload map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &payload); err == nil {
				// Safely traverses down standard GraphQL action format layers: input -> input -> transactionRef
				if firstInput, ok := payload["input"].(map[string]interface{}); ok {
					if secondInput, ok := firstInput["input"].(map[string]interface{}); ok {
						if ref, ok := secondInput["transactionRef"].(string); ok {
							txRef = ref
						}
					}
				}
			}
		}
	}

	if txRef == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing reference value payload marker input vector"})
		return
	}

	fmt.Printf("[DEBUG] Verifying payment Action validation: tx_ref=%s\n", txRef)

	var ticketQuery GetTicketByRefQuery
	err := graphql.Query(c.Request.Context(), &ticketQuery, map[string]interface{}{
		"tx_ref": txRef,
	})

	if err != nil || len(ticketQuery.Tickets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Ticket context trace reference not found"})
		return
	}

	ticket := ticketQuery.Tickets[0]

	if ticket.Status == "pending" {
		chapaBaseURL := os.Getenv("CHAPA_BASE_URL")
		if chapaBaseURL == "" {
			chapaBaseURL = "https://api.chapa.co/v1"
		}
		chapaSecret := os.Getenv("CHAPA_SECRET_KEY")
		if chapaSecret == "" {
			chapaSecret = "CHASECK_TEST-OWOC1Ks6sPaaqOnOeFSv0UwxwxwyR7pO"
		}

		httpReq, err := http.NewRequestWithContext(c.Request.Context(), "GET", chapaBaseURL+"/transaction/verify/"+txRef, nil)
		if err == nil {
			httpReq.Header.Set("Authorization", "Bearer "+chapaSecret)
			resp, err := http.DefaultClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
				respBody, _ := io.ReadAll(resp.Body)
				var verifyResp ChapaVerifyResponse
				if err := json.Unmarshal(respBody, &verifyResp); err == nil {
					if verifyResp.Status == "success" && verifyResp.Data.Status == "success" {
						ticket.Status = "confirmed"
						var updateMutation UpdateTicketMutation
						_ = graphql.Mutate(c.Request.Context(), &updateMutation, map[string]interface{}{
							"tx_ref":     txRef,
							"status":     "confirmed",
							"payment_id": verifyResp.Data.Reference,
						})
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification sequence fully resolved",
		"status":  ticket.Status,
		"ticket": gin.H{
			"id":          ticket.ID,
			"event_id":    ticket.EventID,
			"quantity":    ticket.Quantity,
			"total_price": ticket.TotalPrice,
			"status":      ticket.Status,
		},
	})
}
