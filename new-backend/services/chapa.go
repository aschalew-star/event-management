package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-api/internal/config"
	"payment-api/internal/models"
	"time"
)

type ChapaService struct {
	client *http.Client
	config *config.Config
}

func NewChapaService() *ChapaService {
	return &ChapaService{
		client: &http.Client{Timeout: 30 * time.Second},
		config: config.AppConfig,
	}
}

func (s *ChapaService) InitializePayment(req models.ChapaRequest) (*models.ChapaResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", s.config.ChapaBaseURL)
	reqBody, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.ChapaSecretKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chapaResp models.ChapaResponse
	if err := json.Unmarshal(body, &chapaResp); err != nil {
		return nil, fmt.Errorf("failed to parse gateway response payload: %w", err)
	}
	return &chapaResp, nil
}

func (s *ChapaService) VerifyPayment(txRef string) (*models.ChapaVerifyResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", s.config.ChapaBaseURL, txRef)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.ChapaSecretKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chapaResp models.ChapaVerifyResponse
	if err := json.Unmarshal(body, &chapaResp); err != nil {
		return nil, fmt.Errorf("failed to decode verification stream: %w", err)
	}
	return &chapaResp, nil
}
