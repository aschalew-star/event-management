package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/repository"

	guuid "github.com/google/uuid"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, req domain.ProcessPaymentArgs, userID string) (*domain.ProcessPaymentResponse, error)
	VerifyPayment(ctx context.Context, transactionRef string) (*domain.PaymentVerificationResult, error)
	HandleCallback(ctx context.Context, transactionRef string) (*domain.PaymentVerificationResult, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	httpClient  *http.Client
}

func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// validateEmail checks if email is valid
func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

// cleanPhoneNumber formats phone number for Chapa
func cleanPhoneNumber(phone string) string {
	// Remove spaces and special characters
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// If it doesn't start with + or 0, add country code
	if !strings.HasPrefix(phone, "+") && !strings.HasPrefix(phone, "0") {
		phone = "+251" + phone
	}

	// If it starts with 0, replace with +251
	if strings.HasPrefix(phone, "0") {
		phone = "+251" + phone[1:]
	}

	return phone
}

func (s *paymentService) ProcessPayment(ctx context.Context, req domain.ProcessPaymentArgs, userID string) (*domain.ProcessPaymentResponse, error) {
	// 1. Validate user
	if userID == "" {
		return nil, fmt.Errorf("missing user metadata execution context")
	}

	// 2. Validate event ID
	if req.EventID == "" {
		return nil, fmt.Errorf("event ID is required")
	}

	// 3. Validate required fields
	if req.Email == "" {
		return nil, fmt.Errorf("email is required for payment")
	}

	if !validateEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format: %s", req.Email)
	}

	if req.FirstName == "" {
		return nil, fmt.Errorf("first name is required")
	}

	if req.LastName == "" {
		return nil, fmt.Errorf("last name is required")
	}

	if req.Phone == "" {
		return nil, fmt.Errorf("phone number is required")
	}

	// 4. Get event details
	event, err := s.paymentRepo.GetEventByID(ctx, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	if event.Status != "published" {
		return nil, fmt.Errorf("event is not available for booking")
	}

	// 5. Calculate total price
	totalPrice := event.Price * float64(req.Quantity)

	// 6. Check for duplicate ticket
	hasDuplicate, err := s.paymentRepo.CheckDuplicateTicket(ctx, req.EventID, userID)
	if err != nil {
		fmt.Printf("[WARN] Duplicate check failed: %v\n", err)
	} else if hasDuplicate {
		return nil, fmt.Errorf("you already have a ticket for this event")
	}

	// 7. Generate transaction reference
	txRef := fmt.Sprintf("EVT-%s-%d-%s",
		req.EventID[:8],
		time.Now().Unix(),
		guuid.New().String()[:8])

	// 8. Initialize Chapa payment
	chapaResp, err := s.initializeChapaPayment(ctx, req, event, totalPrice, txRef)
	if err != nil {
		return nil, err
	}

	// 9. Create payment record
	payment := &domain.Payment{
		UserID:         userID,
		Amount:         totalPrice,
		Currency:       "ETB",
		Status:         "pending",
		TransactionRef: txRef,
		PaymentMethod:  "chapa",
	}

	paymentID, err := s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// 10. Create ticket record
	ticket := &domain.Ticket{
		EventID:    req.EventID,
		UserID:     userID,
		PaymentID:  paymentID,
		Quantity:   req.Quantity,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	_, err = s.paymentRepo.CreateTicket(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket record: %w", err)
	}

	return &domain.ProcessPaymentResponse{
		Success:        true,
		Message:        "Payment initialized successfully",
		CheckoutURL:    chapaResp.Data.CheckoutURL,
		TransactionRef: txRef,
	}, nil
}

func (s *paymentService) initializeChapaPayment(ctx context.Context, req domain.ProcessPaymentArgs, event *domain.Event, totalPrice float64, txRef string) (*domain.ChapaResponse, error) {
	chapaBaseURL, chapaSecret := s.getChapaConfig()
	returnURL := s.getReturnURL()

	// Clean and format the data
	email := strings.TrimSpace(strings.ToLower(req.Email))
	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)
	phone := cleanPhoneNumber(req.Phone)

	// Log the request data for debugging
	fmt.Printf("[CHAPA DEBUG] Email: %s, Phone: %s, FirstName: %s, LastName: %s\n",
		email, phone, firstName, lastName)

	chapaReq := domain.ChapaRequest{
		Amount:      totalPrice,
		Currency:    "ETB",
		Email:       email,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phone,
		TxRef:       txRef,
		CallbackURL: s.getCallbackURL(),
		ReturnURL:   fmt.Sprintf("%s/payment/verify?tx_ref=%s", returnURL, txRef),
		Customization: map[string]interface{}{
			"title":       "Event Ticket",
			"description": fmt.Sprintf("Ticket for %s", event.Title),
		},
		Meta: map[string]interface{}{
			"event_id": req.EventID,
			"user_id":  req.Email,
			"quantity": req.Quantity,
		},
	}

	// Log the full request
	reqJSON, _ := json.Marshal(chapaReq)
	fmt.Printf("[CHAPA REQUEST] %s\n", string(reqJSON))

	reqBody, err := json.Marshal(chapaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare payment request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", chapaBaseURL+"/transaction/initialize", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize payment: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+chapaSecret)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Content-Length", fmt.Sprintf("%d", len(reqBody)))

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("payment gateway unavailable: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payment gateway response: %w", err)
	}

	fmt.Printf("[CHAPA RAW RESPONSE] %s\n", string(respBody))

	// Try to parse as ChapaResponse
	var chapaResp domain.ChapaResponse
	if err := json.Unmarshal(respBody, &chapaResp); err != nil {
		return nil, fmt.Errorf("failed to parse payment gateway response: %w", err)
	}

	if chapaResp.Status != "success" {
		// Parse the error message
		var errorResp map[string]interface{}
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			if msg, ok := errorResp["message"]; ok {
				return nil, fmt.Errorf("payment gateway error: %v", msg)
			}
		}
		return nil, fmt.Errorf("payment gateway rejected the transaction: status=%s", chapaResp.Status)
	}

	return &chapaResp, nil
}

func (s *paymentService) getCallbackURL() string {
	callbackURL := os.Getenv("APP_BACKEND_URL")
	if callbackURL == "" {
		callbackURL = "http://localhost:4000"
	}
	return callbackURL + "/api/payment/callback"
}

func (s *paymentService) VerifyPayment(ctx context.Context, transactionRef string) (*domain.PaymentVerificationResult, error) {
	return s.processPaymentVerification(ctx, transactionRef)
}

func (s *paymentService) HandleCallback(ctx context.Context, transactionRef string) (*domain.PaymentVerificationResult, error) {
	return s.processPaymentVerification(ctx, transactionRef)
}

func (s *paymentService) processPaymentVerification(ctx context.Context, transactionRef string) (*domain.PaymentVerificationResult, error) {
	result := &domain.PaymentVerificationResult{
		Success: false,
		Message: "Verification failed",
		Status:  "failed",
	}

	if transactionRef == "" {
		result.Message = "Missing transaction reference"
		return result, nil
	}

	fmt.Printf("[DEBUG] Verifying payment: tx_ref=%s\n", transactionRef)

	// 1. Get payment by transaction reference
	payment, err := s.paymentRepo.GetByTransactionRef(ctx, transactionRef)
	if err != nil {
		result.Message = "Payment not found"
		return result, nil
	}

	if payment == nil {
		result.Message = "Payment not found"
		return result, nil
	}

	result.PaymentID = payment.ID

	// 2. Get ticket by payment_id
	ticket, err := s.paymentRepo.GetTicketByPaymentID(ctx, payment.ID)
	if err != nil || ticket == nil {
		result.Message = "Ticket not found"
		return result, nil
	}

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
		return result, nil
	}

	// 4. Verify with Chapa
	chapaSuccess, chapaMessage := s.verifyTransactionWithChapa(transactionRef)
	if !chapaSuccess {
		result.Message = chapaMessage
		result.Status = "failed"
		return result, nil
	}

	// 5. Update payment status
	if err := s.paymentRepo.UpdateStatus(ctx, transactionRef, "completed"); err != nil {
		result.Message = "Failed to update payment status"
		return result, nil
	}

	// 6. Update ticket status
	if err := s.paymentRepo.UpdateTicketStatus(ctx, payment.ID, "confirmed"); err != nil {
		result.Message = "Failed to update ticket status"
		return result, nil
	}

	// 7. Return success
	result.Success = true
	result.Message = "Payment verified successfully"
	result.Status = "completed"
	result.TicketStatus = "confirmed"

	return result, nil
}

func (s *paymentService) verifyTransactionWithChapa(txRef string) (bool, string) {
	chapaBaseURL, chapaSecret := s.getChapaConfig()

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
	var verifyResp domain.ChapaVerifyResponse
	if err := json.Unmarshal(respBody, &verifyResp); err != nil {
		return false, "Failed to parse verification response"
	}

	if verifyResp.Status == "success" && verifyResp.Data.Status == "success" {
		return true, "Payment verified successfully"
	}

	return false, "Payment verification failed"
}

func (s *paymentService) getChapaConfig() (baseURL, secretKey string) {
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

func (s *paymentService) getReturnURL() string {
	returnURL := os.Getenv("APP_FRONTEND_URL")
	if returnURL == "" {
		returnURL = "http://localhost:3000"
	}
	return returnURL
}
