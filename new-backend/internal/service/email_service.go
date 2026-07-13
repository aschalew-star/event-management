package service

import (
	"fmt"
	"os"
	"sync"

	"event-management/backend/internal/domain"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	SendEmail(req *domain.EmailRequest) error
	SendBulkEmail(requests []domain.EmailRequest) []error
}

type emailService struct {
	fromEmail string
	fromName  string
	apiKey    string
}

func NewEmailService() EmailService {
	fromEmail := os.Getenv("FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "noreply@yourapp.com"
	}

	fromName := os.Getenv("FROM_NAME")
	if fromName == "" {
		fromName = "Event Platform"
	}

	apiKey := os.Getenv("SENDGRID_API_KEY")

	return &emailService{
		fromEmail: fromEmail,
		fromName:  fromName,
		apiKey:    apiKey,
	}
}

func (s *emailService) SendEmail(req *domain.EmailRequest) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(req.ToName, req.ToEmail)
	message := mail.NewSingleEmail(from, req.Subject, to, req.PlainText, req.HTMLContent)

	// If no API key, just log (development mode)
	if s.apiKey == "" {
		fmt.Printf("📧 [DEV] Email to %s: %s\n", req.ToEmail, req.Subject)
		fmt.Printf("    Body: %s\n", req.PlainText)
		return nil
	}

	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid returned status %d: %s", response.StatusCode, response.Body)
	}

	return nil
}

func (s *emailService) SendBulkEmail(requests []domain.EmailRequest) []error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(requests))

	// Worker pool: 5 concurrent email senders
	numWorkers := 5
	if len(requests) < 5 {
		numWorkers = len(requests)
	}

	// Create job channel
	jobs := make(chan domain.EmailRequest, len(requests))

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				err := s.SendEmail(&job)
				if err != nil {
					errChan <- fmt.Errorf("worker %d failed: %w", workerID, err)
				}
			}
		}(i)
	}

	// Queue all emails
	for _, req := range requests {
		jobs <- req
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
