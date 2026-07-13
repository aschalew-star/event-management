package repository

import (
	"context"
	"fmt"

	"event-management/backend/internal/domain"
	"event-management/backend/pkg/graphql"
)

type ImageRepository interface {
	CreateBatch(ctx context.Context, images []domain.EventImage) ([]string, error)
	GetByEventID(ctx context.Context, eventID string) ([]domain.EventImage, error)
	GetByID(ctx context.Context, imageID string) (*domain.EventImage, error)
	Delete(ctx context.Context, imageID string) error
	SetFeatured(ctx context.Context, eventID, imageID string) error
	HasImages(ctx context.Context, eventID string) (bool, error)
	VerifyEventOwnership(ctx context.Context, eventID, userID string) (bool, error)
}

type imageRepository struct{}

func NewImageRepository() ImageRepository {
	return &imageRepository{}
}

func (r *imageRepository) CreateBatch(ctx context.Context, images []domain.EventImage) ([]string, error) {
	if len(images) == 0 {
		return nil, nil
	}

	var imageObjects []map[string]interface{}
	for _, img := range images {
		imageObjects = append(imageObjects, map[string]interface{}{
			"event_id":    img.EventID,
			"image_url":   img.ImageURL,
			"is_featured": img.IsFeatured,
		})
	}

	mutation := `mutation InsertEventImages($objects: [event_images_insert_input!]!) {
        insert_event_images(objects: $objects) {
            returning {
                image_url
            }
        }
    }`

	vars := map[string]interface{}{
		"objects": imageObjects,
	}

	var resp struct {
		InsertEventImages struct {
			Returning []struct {
				ImageURL string `json:"image_url"`
			} `json:"returning"`
		} `json:"insert_event_images"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return nil, fmt.Errorf("database write operation failed: %w", err)
	}

	var urls []string
	for _, r := range resp.InsertEventImages.Returning {
		urls = append(urls, r.ImageURL)
	}

	return urls, nil
}

func (r *imageRepository) GetByEventID(ctx context.Context, eventID string) ([]domain.EventImage, error) {
	query := `query GetEventImages($eventId: uuid!) {
        event_images(where: { event_id: { _eq: $eventId } }) {
            id
            image_url
            is_featured
            public_id
        }
    }`

	vars := map[string]interface{}{
		"eventId": eventID,
	}

	var resp struct {
		EventImages []domain.EventImage `json:"event_images"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, err
	}

	return resp.EventImages, nil
}

func (r *imageRepository) GetByID(ctx context.Context, imageID string) (*domain.EventImage, error) {
	query := `query GetImage($id: uuid!) {
        event_images_by_pk(id: $id) {
            id
            event_id
            image_url
            is_featured
            public_id
        }
    }`

	vars := map[string]interface{}{
		"id": imageID,
	}

	var resp struct {
		EventImagesByPk domain.EventImage `json:"event_images_by_pk"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, err
	}

	if resp.EventImagesByPk.ID == "" {
		return nil, nil
	}

	return &resp.EventImagesByPk, nil
}

func (r *imageRepository) Delete(ctx context.Context, imageID string) error {
	mutation := `mutation DeleteEventImage($id: uuid!) {
        delete_event_images_by_pk(id: $id) {
            id
        }
    }`

	vars := map[string]interface{}{
		"id": imageID,
	}

	var resp struct {
		DeleteEventImagesByPk struct {
			ID string `json:"id"`
		} `json:"delete_event_images_by_pk"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	if resp.DeleteEventImagesByPk.ID == "" {
		return fmt.Errorf("image not found")
	}

	return nil
}

func (r *imageRepository) SetFeatured(ctx context.Context, eventID, imageID string) error {
	// First, unset all featured images for this event
	unsetMutation := `mutation UnsetFeaturedImages($eventId: uuid!) {
        update_event_images(where: { event_id: { _eq: $eventId }, is_featured: { _eq: true } }, _set: { is_featured: false }) {
            affected_rows
        }
    }`

	unsetVars := map[string]interface{}{
		"eventId": eventID,
	}

	var unsetResp struct {
		UpdateEventImages struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"update_event_images"`
	}

	err := graphql.MutateRaw(ctx, unsetMutation, unsetVars, &unsetResp)
	if err != nil {
		return fmt.Errorf("failed to unset featured images: %w", err)
	}

	// Then set the new featured image
	setMutation := `mutation SetFeaturedImage($id: uuid!) {
        update_event_images_by_pk(pk_columns: { id: $id }, _set: { is_featured: true }) {
            id
        }
    }`

	setVars := map[string]interface{}{
		"id": imageID,
	}

	var setResp struct {
		UpdateEventImagesByPk struct {
			ID string `json:"id"`
		} `json:"update_event_images_by_pk"`
	}

	err = graphql.MutateRaw(ctx, setMutation, setVars, &setResp)
	if err != nil {
		return fmt.Errorf("failed to set featured image: %w", err)
	}

	if setResp.UpdateEventImagesByPk.ID == "" {
		return fmt.Errorf("image not found")
	}

	return nil
}

func (r *imageRepository) HasImages(ctx context.Context, eventID string) (bool, error) {
	query := `query CheckEventImages($eventId: uuid!) {
        event_images(where: { event_id: { _eq: $eventId } }, limit: 1) {
            id
        }
    }`

	vars := map[string]interface{}{
		"eventId": eventID,
	}

	var resp struct {
		EventImages []struct {
			ID string `json:"id"`
		} `json:"event_images"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return false, err
	}

	return len(resp.EventImages) > 0, nil
}

func (r *imageRepository) VerifyEventOwnership(ctx context.Context, eventID, userID string) (bool, error) {
	query := `query CheckEvent($id: uuid!, $user_id: uuid!) {
        events(where: { id: { _eq: $id }, user_id: { _eq: $user_id } }) {
            id
        }
    }`

	vars := map[string]interface{}{
		"id":      eventID,
		"user_id": userID,
	}

	var resp struct {
		Events []struct {
			ID string `json:"id"`
		} `json:"events"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return false, err
	}

	return len(resp.Events) > 0, nil
}
