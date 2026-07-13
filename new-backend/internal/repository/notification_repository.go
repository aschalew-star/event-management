package repository

import (
	"context"
	"fmt"

	"event-management/backend/internal/domain"
	"event-management/backend/pkg/graphql"
)

type NotificationRepository interface {
	GetUserByID(ctx context.Context, userID string) (*domain.Users, error)
	GetFollowersByUserID(ctx context.Context, userID string) ([]domain.Users, error)
}

type notificationRepository struct{}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{}
}

func (r *notificationRepository) GetUserByID(ctx context.Context, userID string) (*domain.Users, error) {
	query := `query GetUser($id: uuid!) {
        users_by_pk(id: $id) {
            id
            name
            email
            avatar_url
        }
    }`

	vars := map[string]interface{}{
		"id": userID,
	}

	var resp struct {
		UsersByPk struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		} `json:"users_by_pk"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if resp.UsersByPk.ID == "" {
		return nil, fmt.Errorf("user not found")
	}

	return &domain.Users{
		ID:        resp.UsersByPk.ID,
		Name:      resp.UsersByPk.Name,
		Email:     resp.UsersByPk.Email,
		AvatarURL: resp.UsersByPk.AvatarURL,
	}, nil
}

func (r *notificationRepository) GetFollowersByUserID(ctx context.Context, userID string) ([]domain.Users, error) {
	query := `query GetFollowers($userId: uuid!) {
        follows(where: { followed_user_id: { _eq: $userId } }) {
            follower {
                id
                name
                email
                avatar_url
            }
        }
	}
	}`

	vars := map[string]interface{}{
		"userId": userID,
	}

	var resp struct {
		Follows []struct {
			Follower struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				AvatarURL string `json:"avatar_url"`
			} `json:"follower"`
		} `json:"follows"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	users := make([]domain.Users, 0)
	for _, f := range resp.Follows {
		users = append(users, domain.Users{
			ID:        f.Follower.ID,
			Name:      f.Follower.Name,
			Email:     f.Follower.Email,
			AvatarURL: f.Follower.AvatarURL,
		})
	}

	return users, nil
}
