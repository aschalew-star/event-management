package repository

import (
	"context"
	"fmt"

	"event-management/backend/internal/domain"
	"event-management/backend/pkg/graphql"
)

type userRepository struct{}

func NewUserRepository() domain.UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `query GetUser($email: String!) {
        users(where: { email: { _eq: $email } }, limit: 1) {
            id
            email
            name
            role
            password
            created_at
            updated_at
        }
    }`

	vars := map[string]interface{}{
		"email": email,
	}

	var resp struct {
		Users []domain.User `json:"users"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if len(resp.Users) == 0 {
		return nil, nil
	}

	return &resp.Users[0], nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `query GetUserByID($id: uuid!) {
        users_by_pk(id: $id) {
            id
            name
            email
            role
            created_at
            updated_at
        }
    }`

	vars := map[string]interface{}{
		"id": id,
	}

	var resp struct {
		UsersByPk domain.User `json:"users_by_pk"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}

	if resp.UsersByPk.ID == "" {
		return nil, nil
	}

	return &resp.UsersByPk, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	mutation := `mutation CreateUser($email: String!, $password: String!, $name: String!, $role: String!) {
        insert_users_one(object: { email: $email, password: $password, name: $name, role: $role }) {
            id
            email
            name
            role
            created_at
        }
    }`

	vars := map[string]interface{}{
		"email":    user.Email,
		"password": user.Password,
		"name":     user.Name,
		"role":     user.Role,
	}

	var resp struct {
		InsertUsersOne struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			Role      string `json:"role"`
			CreatedAt string `json:"created_at"`
		} `json:"insert_users_one"`
	}

	err := graphql.MutateRaw(ctx, mutation, vars, &resp)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = resp.InsertUsersOne.ID
	return nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `query CheckUser($email: String!) {
        users(where: { email: { _eq: $email } }, limit: 1) {
            id
            email
        }
    }`

	vars := map[string]interface{}{
		"email": email,
	}

	var resp struct {
		Users []struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"users"`
	}

	err := graphql.QueryRaw(ctx, query, vars, &resp)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return len(resp.Users) > 0, nil
}
