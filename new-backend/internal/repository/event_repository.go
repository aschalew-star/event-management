package repository

import (
	"context"
	"fmt"

	"event-management/backend/internal/domain"
	"event-management/backend/pkg/graphql"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) (string, error)
	GetByID(ctx context.Context, id string) (*domain.Event, error)
	ValidateCategory(ctx context.Context, categoryID string) (bool, error)
}

type eventRepository struct {
	// No need to store gql client since we use package-level functions
}

func NewEventRepository() EventRepository {
	return &eventRepository{}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) (string, error) {
	mutation := `mutation CreateEvent($object: events_insert_input!) {
        insert_events_one(object: $object) {
            id
            title
            event_date
        }
    }`

	vars := map[string]interface{}{
		"object": map[string]interface{}{
			"title":       event.Title,
			"description": event.Description,
			"category_id": event.CategoryID,
			"price":       event.Price,
			"is_free":     event.IsFree,
			"venue":       event.Venue,
			"address":     event.Address,
			"latitude":    event.Latitude,
			"longitude":   event.Longitude,
			"event_date":  event.EventDate,
			"start_time":  event.StartTime,
			"end_time":    event.EndTime,
			"status":      event.Status,
			"user_id":     event.UserID,
		},
	}

	var resp struct {
		InsertEventsOne struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			EventDate string `json:"event_date"`
		} `json:"insert_events_one"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return "", fmt.Errorf("failed to create event: %w", err)
	}

	if resp.InsertEventsOne.ID == "" {
		return "", fmt.Errorf("no ID returned from event creation")
	}

	return resp.InsertEventsOne.ID, nil
}

func (r *eventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	query := `query GetEvent($id: uuid!) {
        events(where: { id: { _eq: $id } }, limit: 1) {
            id
            title
            description
            category_id
            price
            is_free
            venue
            address
            latitude
            longitude
            event_date
            start_time
            end_time
            status
            user_id
            created_at
            updated_at
        }
    }`

	vars := map[string]interface{}{
		"id": id,
	}

	var resp struct {
		Events []domain.Event `json:"events"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Events) == 0 {
		return nil, nil
	}

	return &resp.Events[0], nil
}

func (r *eventRepository) ValidateCategory(ctx context.Context, categoryID string) (bool, error) {
	query := `query CheckCategory($id: uuid!) {
        categories(where: { id: { _eq: $id } }, limit: 1) {
            id
            name
        }
    }`

	vars := map[string]interface{}{
		"id": categoryID,
	}

	var resp struct {
		Categories []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"categories"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return false, err
	}

	return len(resp.Categories) > 0, nil
}
