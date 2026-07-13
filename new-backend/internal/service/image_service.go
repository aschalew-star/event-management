package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"event-management/backend/internal/domain"
	"event-management/backend/internal/repository"
	"event-management/backend/pkg/cloudinary"
)

type ImageService interface {
	UploadEventImages(ctx context.Context, eventID string, images []string) (*domain.UploadImageResponse, error)
	DeleteEventImage(ctx context.Context, imageID, userID string) error
	SetFeaturedImage(ctx context.Context, eventID, imageID, userID string) error
}

type imageService struct {
	imageRepo repository.ImageRepository
	cloud     *cloudinary.Client
}

func NewImageService(imageRepo repository.ImageRepository, cloud *cloudinary.Client) ImageService {
	return &imageService{
		imageRepo: imageRepo,
		cloud:     cloud,
	}
}

func (s *imageService) UploadEventImages(ctx context.Context, eventID string, images []string) (*domain.UploadImageResponse, error) {
	// 1. Validate input
	if eventID == "" {
		return nil, fmt.Errorf("event ID is required")
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no images provided to upload")
	}

	fmt.Printf("📸 Processing %d images for event %s\n", len(images), eventID)

	// 2. Check if event has existing images
	hasExistingImages, err := s.imageRepo.HasImages(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing images: %w", err)
	}

	// 3. Upload to Cloudinary with parallel processing
	uploadedImages, err := s.uploadToCloudinaryParallel(ctx, images, "events")
	if err != nil {
		// If some images failed but others succeeded, continue
		if len(uploadedImages) == 0 {
			return nil, fmt.Errorf("all images failed to upload: %w", err)
		}
		fmt.Printf("⚠️ Some images failed to upload, but continuing with %d successful ones\n", len(uploadedImages))
	}

	if len(uploadedImages) == 0 {
		return nil, fmt.Errorf("no images were successfully uploaded")
	}

	fmt.Printf("✅ Successfully uploaded %d images\n", len(uploadedImages))

	// 4. Prepare image objects for database
	var imageObjects []domain.EventImage
	for i, img := range uploadedImages {
		isFeatured := !hasExistingImages && i == 0
		imageObjects = append(imageObjects, domain.EventImage{
			EventID:    eventID,
			ImageURL:   img.URL,
			IsFeatured: isFeatured,
		})
	}

	// 5. Save to database
	urls, err := s.imageRepo.CreateBatch(ctx, imageObjects)
	if err != nil {
		return nil, fmt.Errorf("failed to save images: %w", err)
	}

	return &domain.UploadImageResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully processed and saved %d records", len(uploadedImages)),
		URLs:    urls,
	}, nil
}

func (s *imageService) uploadToCloudinaryParallel(ctx context.Context, images []string, folder string) ([]domain.ImageUploadResult, error) {
	if len(images) == 0 {
		return nil, nil
	}

	// Clean image data (remove base64 prefix)
	cleanedImages := make([]string, len(images))
	for i, img := range images {
		if strings.Contains(img, ";base64,") {
			parts := strings.Split(img, ";base64,")
			if len(parts) == 2 {
				cleanedImages[i] = parts[1]
				fmt.Printf("✅ Cleaned image %d: removed base64 prefix\n", i+1)
			} else {
				cleanedImages[i] = img
			}
		} else {
			cleanedImages[i] = img
		}
	}

	// Upload with shorter timeout
	uploadCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	results, err := s.cloud.UploadMultipleImages(uploadCtx, cleanedImages, folder)
	if err != nil {
		return nil, err
	}

	// Filter successful uploads
	var uploadedImages []domain.ImageUploadResult
	for _, r := range results {
		if r.Error == nil && r.URL != "" {
			uploadedImages = append(uploadedImages, domain.ImageUploadResult{
				URL:      r.URL,
				PublicID: r.PublicID,
			})
			fmt.Printf("✅ Image %d uploaded: %s\n", r.Index+1, r.URL)
		} else if r.Error != nil {
			fmt.Printf("❌ Image %d upload failed: %v\n", r.Index+1, r.Error)
		}
	}

	return uploadedImages, nil
}

func (s *imageService) DeleteEventImage(ctx context.Context, imageID, userID string) error {
	if imageID == "" {
		return fmt.Errorf("image ID is required")
	}

	if userID == "" {
		return fmt.Errorf("unauthorized: missing user session")
	}

	// 1. Get image details
	image, err := s.imageRepo.GetByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	if image == nil {
		return fmt.Errorf("image not found")
	}

	// 2. Verify event ownership
	hasAccess, err := s.imageRepo.VerifyEventOwnership(ctx, image.EventID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}

	if !hasAccess {
		return fmt.Errorf("permission denied: you don't own this event")
	}

	// 3. Delete from database
	if err := s.imageRepo.Delete(ctx, imageID); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	// 4. If there are remaining images and the deleted one was featured, set a new featured image
	if image.IsFeatured {
		remainingImages, err := s.imageRepo.GetByEventID(ctx, image.EventID)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: Failed to get remaining images for event %s: %v\n", image.EventID, err)
		} else if len(remainingImages) > 0 {
			// Set the first remaining image as featured
			if err := s.imageRepo.SetFeatured(ctx, image.EventID, remainingImages[0].ID); err != nil {
				fmt.Printf("Warning: Failed to set new featured image: %v\n", err)
			}
		}
	}

	return nil
}

func (s *imageService) SetFeaturedImage(ctx context.Context, eventID, imageID, userID string) error {
	if eventID == "" {
		return fmt.Errorf("event ID is required")
	}

	if imageID == "" {
		return fmt.Errorf("image ID is required")
	}

	if userID == "" {
		return fmt.Errorf("unauthorized: missing user session")
	}

	// 1. Verify event ownership
	hasAccess, err := s.imageRepo.VerifyEventOwnership(ctx, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}

	if !hasAccess {
		return fmt.Errorf("permission denied: you don't own this event")
	}

	// 2. Verify image belongs to event
	image, err := s.imageRepo.GetByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	if image == nil {
		return fmt.Errorf("image not found")
	}

	if image.EventID != eventID {
		return fmt.Errorf("image does not belong to this event")
	}

	// 3. Set featured image
	if err := s.imageRepo.SetFeatured(ctx, eventID, imageID); err != nil {
		return fmt.Errorf("failed to set featured image: %w", err)
	}

	return nil
}
