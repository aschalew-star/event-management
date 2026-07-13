package repository

import (
	"context"
	"fmt"

	"event-management/backend/internal/domain"
	"event-management/backend/pkg/graphql"
)

type TagRepository interface {
	GetOrCreate(ctx context.Context, name string) (string, error)
	AssociateWithEvent(ctx context.Context, eventID, tagID string) error
}

type tagRepository struct{}

func NewTagRepository() TagRepository {
	return &tagRepository{}
}

func (r *tagRepository) GetOrCreate(ctx context.Context, name string) (string, error) {
	// First check if tag exists
	query := `query GetTag($name: String!) {
        tags(where: { name: { _eq: $name } }) {
            id
            name
        }
    }`

	vars := map[string]interface{}{
		"name": name,
	}

	var resp struct {
		Tags []domain.Tag `json:"tags"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return "", err
	}

	if len(resp.Tags) > 0 {
		return resp.Tags[0].ID, nil
	}

	// Create new tag
	mutation := `mutation CreateTag($name: String!) {
        insert_tags_one(object: { name: $name }) {
            id
            name
        }
    }`

	var createResp struct {
		InsertTagsOne struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"insert_tags_one"`
	}

	err = graphql.MutateRaw(ctx, mutation, vars, &createResp)
	if err != nil {
		return "", fmt.Errorf("failed to create tag: %w", err)
	}

	return createResp.InsertTagsOne.ID, nil
}

func (r *tagRepository) AssociateWithEvent(ctx context.Context, eventID, tagID string) error {
	mutation := `mutation AssociateTag($event_id: uuid!, $tag_id: uuid!) {
        insert_event_tags(objects: { event_id: $event_id, tag_id: $tag_id }) {
            affected_rows
        }
    }`

	vars := map[string]interface{}{
		"event_id": eventID,
		"tag_id":   tagID,
	}

	var resp struct {
		InsertEventTags struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"insert_event_tags"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return fmt.Errorf("failed to associate tag: %w", err)
	}

	return nil
}
