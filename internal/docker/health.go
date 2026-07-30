package docker

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus represents the health status of a container.
type HealthStatus struct {
	Status        string // healthy, unhealthy, starting, none
	FailingStreak int
	LastCheck     time.Time
	Log           []HealthCheckLog
}

// HealthCheckLog represents a single health check log entry.
type HealthCheckLog struct {
	Start    time.Time
	End      time.Time
	ExitCode int
	Output   string
}

// GetHealthStatus returns the health status of a container.
func (c *Client) GetHealthStatus(ctx context.Context, nameOrID string) (*HealthStatus, error) {
	info, err := c.GetContainer(ctx, nameOrID)
	if err != nil {
		return nil, err
	}

	if info.Health == "" {
		return &HealthStatus{Status: "none"}, nil
	}

	return &HealthStatus{
		Status: info.Health,
	}, nil
}

// IsHealthy checks if a container is healthy.
func (c *Client) IsHealthy(ctx context.Context, nameOrID string) (bool, error) {
	status, err := c.GetHealthStatus(ctx, nameOrID)
	if err != nil {
		return false, err
	}
	return status.Status == "healthy", nil
}

// WaitForHealthy waits for a container to become healthy with timeout.
func WaitForHealthy(ctx context.Context, client *Client, nameOrID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for %s to be healthy", nameOrID)
			}

			info, err := client.GetContainer(ctx, nameOrID)
			if err != nil {
				continue
			}

			switch info.Health {
			case "healthy":
				return nil
			case "unhealthy":
				return fmt.Errorf("container %s is unhealthy", nameOrID)
			}

			if info.State == "exited" || info.State == "dead" {
				return fmt.Errorf("container %s exited with code %d", nameOrID, info.ExitCode)
			}
		}
	}
}
