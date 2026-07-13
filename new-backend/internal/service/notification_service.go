package service

import (
	"context"
	"fmt"
	"os"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/repository"
)

type NotificationService interface {
	ProcessEventNotification(ctx context.Context, event domain.EventNotification) error
	SendEventNotification(ctx context.Context, event domain.EventNotification) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	emailService     EmailService
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	emailService EmailService,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		emailService:     emailService,
	}
}

func (s *notificationService) ProcessEventNotification(ctx context.Context, event domain.EventNotification) error {
	// 1. Get event creator
	_, err := s.notificationRepo.GetUserByID(ctx, event.EventUserID)
	if err != nil {
		return fmt.Errorf("failed to get event creator: %w", err)
	}

	// 2. Get followers of the event creator
	followers, err := s.notificationRepo.GetFollowersByUserID(ctx, event.EventUserID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	if len(followers) == 0 {
		fmt.Printf("ℹ️ No followers found for user %s\n", event.EventUserID)
		return nil
	}

	fmt.Printf("📨 Found %d followers to notify\n", len(followers))

	// 3. Send emails to followers
	return s.SendEventNotification(ctx, event)
}

func (s *notificationService) SendEventNotification(ctx context.Context, event domain.EventNotification) error {
	// Get event creator
	creator, err := s.notificationRepo.GetUserByID(ctx, event.EventUserID)
	if err != nil {
		return fmt.Errorf("failed to get event creator: %w", err)
	}

	// Get followers
	followers, err := s.notificationRepo.GetFollowersByUserID(ctx, event.EventUserID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	if len(followers) == 0 {
		return nil
	}

	// Build email requests
	emailRequests := make([]domain.EmailRequest, 0, len(followers))
	for _, follower := range followers {
		if follower.Email == "" {
			continue
		}

		subject := fmt.Sprintf("%s %s an event: %s", creator.Name, event.Operation, event.EventTitle)

		htmlContent := s.buildEmailHTML(creator, &follower, &event)
		plainTextContent := s.buildEmailPlainText(creator, &follower, &event)

		emailRequests = append(emailRequests, domain.EmailRequest{
			ToEmail:     follower.Email,
			ToName:      follower.Name,
			Subject:     subject,
			HTMLContent: htmlContent,
			PlainText:   plainTextContent,
		})
	}

	if len(emailRequests) == 0 {
		fmt.Printf("ℹ️ No valid email addresses found\n")
		return nil
	}

	// Send all emails
	errors := s.emailService.SendBulkEmail(emailRequests)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("❌ Email sending error: %v\n", err)
		}
		return fmt.Errorf("failed to send %d emails", len(errors))
	}

	fmt.Printf("📧 All emails sent for event: %s\n", event.EventTitle)
	return nil
}

func (s *notificationService) buildEmailHTML(creator, follower *domain.Users, event *domain.EventNotification) string {
	operationText := event.Operation
	if operationText == "created" {
		operationText = "created"
	} else if operationText == "updated" {
		operationText = "updated"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .event-details { background: white; padding: 15px; margin: 15px 0; border-radius: 5px; border-left: 4px solid #4CAF50; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .button { display: inline-block; background: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
        .label { font-weight: bold; color: #555; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>📢 Event %s</h2>
        </div>
        <div class="content">
            <p>Hello <strong>%s</strong>,</p>
            <p><strong>%s</strong> has %s a new event:</p>
            
            <div class="event-details">
                <h3>%s</h3>
                <p><span class="label">📍 Venue:</span> %s</p>
                <p><span class="label">📅 Date:</span> %s</p>
                <p><span class="label">📊 Status:</span> %s</p>
                <p><span class="label">👤 Created by:</span> %s</p>
            </div>
            
            <p style="text-align: center; margin-top: 25px;">
                <a href="%s/events/%s" class="button">View Event</a>
            </p>
        </div>
        <div class="footer">
            <p>You're receiving this because you follow %s.</p>
            <p>&copy; 2024 Event Platform. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
		operationText,
		follower.Name,
		creator.Name,
		operationText,
		event.EventTitle,
		event.EventVenue,
		event.EventDate,
		event.EventStatus,
		creator.Name,
		frontendURL,
		event.EventID,
		creator.Name,
	)
}

func (s *notificationService) buildEmailPlainText(creator, follower *domain.Users, event *domain.EventNotification) string {
	operationText := event.Operation
	if operationText == "created" {
		operationText = "created"
	} else if operationText == "updated" {
		operationText = "updated"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	return fmt.Sprintf(`
Hello %s,

%s has %s an event: %s

Event Details:
- Title: %s
- Venue: %s
- Date: %s
- Status: %s
- Created by: %s

View Event: %s/events/%s

You're receiving this because you follow %s.

---
Event Platform
`,
		follower.Name,
		creator.Name,
		operationText,
		event.EventTitle,
		event.EventTitle,
		event.EventVenue,
		event.EventDate,
		event.EventStatus,
		creator.Name,
		frontendURL,
		event.EventID,
		creator.Name,
	)
}
