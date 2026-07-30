package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/moby/moby/client"
)

// Client wraps the Docker SDK client.
type Client struct {
	api *client.Client
}

// NewClient creates a new Docker client with environment-based config.
func NewClient() (*Client, error) {
	api, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating docker client: %w", err)
	}
	return &Client{api: api}, nil
}

// Ping checks if Docker daemon is reachable.
func (c *Client) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.api.Ping(ctx, client.PingOptions{})
	if err != nil {
		return fmt.Errorf("docker ping: %w", err)
	}
	return nil
}

// Close closes the Docker client.
func (c *Client) Close() error {
	return c.api.Close()
}

// API returns the underlying Docker API client.
func (c *Client) API() *client.Client {
	return c.api
}
