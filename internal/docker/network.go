package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/moby/moby/client"
)

// EnsureNetwork creates a Docker network if it doesn't exist.
func (c *Client) EnsureNetwork(ctx context.Context, name string) error {
	exists, err := c.NetworkExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return c.CreateNetwork(ctx, name)
}

// CreateNetwork creates a Docker network.
func (c *Client) CreateNetwork(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.api.NetworkCreate(ctx, name, client.NetworkCreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return fmt.Errorf("creating network %s: %w", name, err)
	}
	return nil
}

// NetworkExists checks if a Docker network exists.
func (c *Client) NetworkExists(ctx context.Context, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.api.NetworkInspect(ctx, name, client.NetworkInspectOptions{})
	if err != nil {
		// Docker SDK returns "No such network" for missing networks
		if strings.Contains(err.Error(), "No such network") ||
			strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, fmt.Errorf("inspecting network %s: %w", name, err)
	}
	return true, nil
}

// RemoveNetwork removes a Docker network.
func (c *Client) RemoveNetwork(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.api.NetworkRemove(ctx, name, client.NetworkRemoveOptions{})
	if err != nil {
		return fmt.Errorf("removing network %s: %w", name, err)
	}
	return nil
}
