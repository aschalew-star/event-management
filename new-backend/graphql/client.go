package graphql

import (
	"context"
	"fmt"
	"github.com/hasura/go-graphql-client"
	"net/http"
)

type Client struct {
	client *graphql.Client
	secret string
}

var gqlClient *Client

func NewClient(endpoint, secret string) *Client {
	// create an HTTP client that injects the Hasura admin secret into every request
	httpClient := &http.Client{
		Transport: &authTransport{secret: secret, rt: http.DefaultTransport},
	}
	return &Client{
		client: graphql.NewClient(endpoint, httpClient),
		secret: secret,
	}
}

// set the client for the package-level functions
func SetClient(client *Client) {
	gqlClient = client
}

// authTransport adds the X-Hasura-Admin-Secret header to every request.
type authTransport struct {
	secret string
	rt     http.RoundTripper
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// clone the request to avoid mutating the original
	r2 := req.Clone(req.Context())
	if a.secret != "" {
		r2.Header.Set("X-Hasura-Admin-Secret", a.secret)
	}
	return a.rt.RoundTrip(r2)
}

func (c *Client) Mutate(ctx context.Context, mutation interface{}, vars map[string]interface{}) error {
	return c.client.Mutate(ctx, mutation, vars)
}

func (c *Client) Query(ctx context.Context, query interface{}, vars map[string]interface{}) error {
	return c.client.Query(ctx, query, vars)
}

func Mutate(ctx context.Context, mutation interface{}, vars map[string]interface{}) error {
	if gqlClient == nil {
		return fmt.Errorf("graphql client not initialized")
	}
	return gqlClient.Mutate(ctx, mutation, vars)
}

func Query(ctx context.Context, query interface{}, vars map[string]interface{}) error {
	if gqlClient == nil {
		return fmt.Errorf("graphql client not initialized")
	}
	return gqlClient.Query(ctx, query, vars)
}

// Add these functions to your graphql package

func QueryRaw(ctx context.Context, query string, vars map[string]interface{}, result interface{}) error {
	if gqlClient == nil {
		return fmt.Errorf("graphql client not initialized")
	}
	return gqlClient.QueryRaw(ctx, query, vars, result)
}

func MutateRaw(ctx context.Context, mutation string, vars map[string]interface{}, result interface{}) error {
	if gqlClient == nil {
		return fmt.Errorf("graphql client not initialized")
	}
	return gqlClient.MutateRaw(ctx, mutation, vars, result)
}

// Add these methods to your Client struct
func (c *Client) QueryRaw(ctx context.Context, query string, vars map[string]interface{}, result interface{}) error {
	return c.client.Query(ctx, result, vars)
}

func (c *Client) MutateRaw(ctx context.Context, mutation string, vars map[string]interface{}, result interface{}) error {
	return c.client.Mutate(ctx, result, vars)
}
