// graphql/client.go
package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client   *http.Client
	endpoint string
	secret   string
}

var gqlClient *Client

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{}    `json:"data,omitempty"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []map[string]int       `json:"locations"`
	Path       []string               `json:"path"`
	Extensions map[string]interface{} `json:"extensions"`
}

func NewClient(endpoint, secret string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		endpoint: endpoint,
		secret:   secret,
	}
}

func SetClient(client *Client) {
	gqlClient = client
}

// MutateRaw - for raw GraphQL mutation strings
func (c *Client) MutateRaw(ctx context.Context, mutation string, vars map[string]interface{}, result interface{}) error {
	return c.executeRawGraphQL(ctx, mutation, vars, result)
}

// QueryRaw - for raw GraphQL query strings
func (c *Client) QueryRaw(ctx context.Context, query string, vars map[string]interface{}, result interface{}) error {
	return c.executeRawGraphQL(ctx, query, vars, result)
}

// executeRawGraphQL executes raw GraphQL queries/mutations
func (c *Client) executeRawGraphQL(ctx context.Context, query string, vars map[string]interface{}, result interface{}) error {
	request := GraphQLRequest{
		Query:     query,
		Variables: vars,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Log the request for debugging
	fmt.Printf("📤 Sending GraphQL request:\n%s\n", string(jsonData))

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.secret != "" {
		req.Header.Set("X-Hasura-Admin-Secret", c.secret)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response
	fmt.Printf("📥 GraphQL Response Status: %d\n", resp.StatusCode)
	fmt.Printf("📥 GraphQL Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var gqlResp GraphQLResponse
	gqlResp.Data = result
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		errorMsg := gqlResp.Errors[0].Message
		if len(gqlResp.Errors) > 1 {
			for _, err := range gqlResp.Errors[1:] {
				errorMsg += "; " + err.Message
			}
		}
		return fmt.Errorf("graphql error: %s", errorMsg)
	}

	return nil
}

// Package-level functions
func MutateRaw(ctx context.Context, mutation string, vars map[string]interface{}, result interface{}) error {
	if gqlClient == nil {
		return fmt.Errorf("graphql client not initialized")
	}
	return gqlClient.MutateRaw(ctx, mutation, vars, result)
}

func QueryRaw(ctx context.Context, query string, vars map[string]interface{}, result interface{}) error {
	if gqlClient == nil {
		return fmt.Errorf("graphql client not initialized")
	}
	return gqlClient.QueryRaw(ctx, query, vars, result)
}
