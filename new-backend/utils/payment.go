package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type PaymentRequest struct {
	Amount    float64 `json:"amount"`
	Email     string  `json:"email"`
	Reference string  `json:"reference"`
	Title     string  `json:"title"`
}

type PaymentResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	Reference   string `json:"reference"`
	CheckoutURL string `json:"checkout_url"`
}

func InitiatePayment(req PaymentRequest) (*PaymentResponse, error) {
	secretKey := os.Getenv("CHAPA_SECRET_KEY")
	if secretKey == "" {
		return &PaymentResponse{
			Status:      "success",
			Reference:   req.Reference,
			CheckoutURL: "http://localhost:3000/payment/success",
		}, nil
	}

	payload := map[string]interface{}{
		"amount":       req.Amount,
		"email":        req.Email,
		"tx_ref":       req.Reference,
		"currency":     "ETB",
		"title":        req.Title,
		"callback_url": "http://localhost:3000/payment/callback",
		"return_url":   "http://localhost:3000/payment/success",
		"cancel_url":   "http://localhost:3000/payment/cancel",
		"customization": map[string]string{
			"title":       req.Title,
			"description": "Event Ticket Purchase",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	request, err := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+secretKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Reference   string `json:"tx_ref"`
			CheckoutURL string `json:"checkout_url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("payment initiation failed: %s", result.Message)
	}

	return &PaymentResponse{
		Status:      result.Status,
		Message:     result.Message,
		Reference:   result.Data.Reference,
		CheckoutURL: result.Data.CheckoutURL,
	}, nil
}
