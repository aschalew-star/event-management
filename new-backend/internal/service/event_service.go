package service

import (
	"context"
	"fmt"
	// "strings"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/repository"
	"event-management/backend/pkg/cloudinary"
	// "event-management/backend/pkg/graphql"
)

type EventService interface {
	CreateEvent(ctx context.Context, req domain.CreateEventRequest, userID string) (*domain.EventResponse, error)
}

type eventService struct {
	eventRepo repository.EventRepository
	imageRepo repository.ImageRepository
	tagRepo   repository.TagRepository
	cloud     *cloudinary.Client
}

func NewEventService(
	eventRepo repository.EventRepository,
	imageRepo repository.ImageRepository,
	tagRepo repository.TagRepository,
	cloud *cloudinary.Client,
) EventService {
	return &eventService{
		eventRepo: eventRepo,
		imageRepo: imageRepo,
		tagRepo:   tagRepo,
		cloud:     cloud,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, req domain.CreateEventRequest, userID string) (*domain.EventResponse, error) {
	// 1. Validate user
	if userID == "" {
		return nil, fmt.Errorf("authorization failed: missing user session id")
	}

	// 2. Handle images
	imageURLs, featuredImageURL, err := s.processImages(ctx, req)
	if err != nil {
		return nil, err
	}

	// 3. Prepare event data
	eventData := s.prepareEventData(req, userID)

	// 4. Validate category if provided
	if eventData.CategoryID != nil && *eventData.CategoryID != "" {
		exists, err := s.eventRepo.ValidateCategory(ctx, *eventData.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate category: %w", err)
		}
		if !exists {
			// Set to null if category doesn't exist
			eventData.CategoryID = nil
			fmt.Printf("⚠️ Category not found, setting to null\n")
		} else {
			fmt.Printf("✅ Category validated: %s\n", *eventData.CategoryID)
		}
	}

	// 5. Create event
	eventID, err := s.eventRepo.Create(ctx, eventData)
	if err != nil {
		return nil, err
	}

	// 6. Insert images
	if len(imageURLs) > 0 {
		var images []domain.EventImage
		for _, img := range imageURLs {
			images = append(images, domain.EventImage{
				EventID:    eventID,
				ImageURL:   img.URL,
				IsFeatured: img.URL == featuredImageURL,
			})
		}

		affected, err := s.imageRepo.CreateBatch(ctx, images)
		if err != nil {
			fmt.Printf("Failed to insert images for event %s: %v\n", eventID, err)
		} else {
			fmt.Printf("✅ Inserted %d images for event %s\n", affected, eventID)
		}
	}

	// 7. Handle tags
	if len(req.Tags) > 0 {
		s.processTags(ctx, eventID, req.Tags)
	}

	// 8. Return response
	return &domain.EventResponse{
		ID:            eventID,
		Message:       "Event created successfully",
		Success:       true,
		FeaturedImage: featuredImageURL,
		TotalImages:   len(imageURLs),
	}, nil
}

func (s *eventService) processImages(ctx context.Context, req domain.CreateEventRequest) ([]domain.ImageUploadResult, string, error) {
	var imageURLs []domain.ImageUploadResult
	var featuredImageURL string

	imagesToUpload := req.Images
	if len(imagesToUpload) == 0 && req.FeaturedImage != "" {
		imagesToUpload = []string{req.FeaturedImage}
	}

	if len(imagesToUpload) == 0 {
		return imageURLs, featuredImageURL, nil
	}

	results, err := s.cloud.UploadMultipleImages(ctx, imagesToUpload, "events")
	if err != nil {
		return nil, "", err
	}

	for _, result := range results {
		if result.Error != nil {
			if result.Error == context.DeadlineExceeded {
				return nil, "", fmt.Errorf("upload for image %d timed out. Please try with smaller images", result.Index+1)
			}
			fmt.Printf("LOGGING ERROR - CLOUDINARY UPLOAD FAILED for image %d: %s\n", result.Index+1, result.Error.Error())
			continue
		}

		isFeatured := result.Index == 0 && featuredImageURL == ""
		if isFeatured {
			featuredImageURL = result.URL
		}

		imageURLs = append(imageURLs, domain.ImageUploadResult{
			URL:      result.URL,
			PublicID: result.PublicID,
		})

		fmt.Printf("✅ Uploaded image %d: %s (featured: %v)\n", result.Index+1, result.URL, isFeatured)
	}

	if len(imageURLs) == 0 {
		return nil, "", fmt.Errorf("failed to upload any images")
	}

	return imageURLs, featuredImageURL, nil
}

func (s *eventService) prepareEventData(req domain.CreateEventRequest, userID string) *domain.Event {
	// Format event date
	eventDateFormatted := req.EventDate
	if eventDateFormatted != "" && len(eventDateFormatted) == 10 {
		eventDateFormatted = eventDateFormatted + "T00:00:00Z"
	}

	// Prepare category ID
	var categoryID *string
	if req.CategoryID != nil && *req.CategoryID != "" && *req.CategoryID != `""` {
		categoryID = req.CategoryID
	}

	// Prepare start/end times
	var startTime *string
	if req.StartTime != nil && *req.StartTime != "" {
		startTime = req.StartTime
	}

	var endTime *string
	if req.EndTime != nil && *req.EndTime != "" {
		endTime = req.EndTime
	}

	return &domain.Event{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  categoryID,
		Price:       req.Price,
		IsFree:      req.IsFree,
		Venue:       req.Venue,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		EventDate:   eventDateFormatted,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      req.Status,
		UserID:      userID,
	}
}

func (s *eventService) processTags(ctx context.Context, eventID string, tags []string) {
	fmt.Printf("📌 Tags to process: %v\n", tags)

	for _, tagName := range tags {
		tagID, err := s.tagRepo.GetOrCreate(ctx, tagName)
		if err != nil {
			fmt.Printf("Error processing tag %s: %v\n", tagName, err)
			continue
		}

		if tagID != "" {
			err = s.tagRepo.AssociateWithEvent(ctx, eventID, tagID)
			if err != nil {
				fmt.Printf("Error associating tag %s with event: %v\n", tagName, err)
			} else {
				fmt.Printf("Associated tag %s with event %s\n", tagName, eventID)
			}
		}
	}
}
