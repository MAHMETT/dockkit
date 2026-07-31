package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// ContainerInfo holds parsed container information.
type ContainerInfo struct {
	ID            string
	Name          string
	Image         string
	State         string // running, exited, paused, created, restarting
	Status        string // human-readable (e.g., "Up 5 hours")
	Ports         []PortMapping
	Labels        map[string]string
	Health        string // healthy, unhealthy, starting, ""
	CreatedAt     time.Time
	StartedAt     time.Time
	FinishedAt    time.Time
	ExitCode      int
	RestartCount  int
	Networks      []string
}

// PortMapping represents a host:container port mapping.
type PortMapping struct {
	HostIP        string
	HostPort      string
	ContainerPort string
	Protocol      string // tcp, udp
}

// ListContainers returns all containers (running and stopped).
func (c *Client) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := c.api.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("listing containers: %w", err)
	}

	containers := make([]ContainerInfo, 0, len(result.Items))
	for _, r := range result.Items {
		containers = append(containers, parseContainerSummary(r))
	}
	return containers, nil
}

// GetContainer returns detailed info for a specific container.
func (c *Client) GetContainer(ctx context.Context, nameOrID string) (*ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := c.api.ContainerInspect(ctx, nameOrID, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("inspecting container %s: %w", nameOrID, err)
	}

	cr := result.Container
	info := &ContainerInfo{
		ID:           cr.ID,
		Name:         strings.TrimPrefix(cr.Name, "/"),
		Image:        cr.Config.Image,
		State:        string(cr.State.Status),
		Labels:       cr.Config.Labels,
		RestartCount: cr.RestartCount,
	}

	if cr.State.Health != nil {
		info.Health = string(cr.State.Health.Status)
	}

	if cr.State.StartedAt != "" {
		t, err := time.Parse(time.RFC3339Nano, cr.State.StartedAt)
		if err != nil {
			log.Printf("warning: parsing StartedAt %q: %v", cr.State.StartedAt, err)
		}
		info.StartedAt = t
	}
	if cr.State.FinishedAt != "" {
		t, err := time.Parse(time.RFC3339Nano, cr.State.FinishedAt)
		if err != nil {
			log.Printf("warning: parsing FinishedAt %q: %v", cr.State.FinishedAt, err)
		}
		info.FinishedAt = t
	}

	if cr.Created != "" {
		t, err := time.Parse(time.RFC3339Nano, cr.Created)
		if err != nil {
			log.Printf("warning: parsing Created %q: %v", cr.Created, err)
		}
		info.CreatedAt = t
	}
	info.ExitCode = cr.State.ExitCode

	for netName := range cr.NetworkSettings.Networks {
		info.Networks = append(info.Networks, netName)
	}

	for port, bindings := range cr.NetworkSettings.Ports {
		for _, b := range bindings {
			hostIP := b.HostIP.String()
			if hostIP == "<nil>" || hostIP == "::" {
				hostIP = "0.0.0.0"
			}
			info.Ports = append(info.Ports, PortMapping{
				HostIP:        hostIP,
				HostPort:      b.HostPort,
				ContainerPort: port.Port(),
				Protocol:      string(port.Proto()),
			})
		}
	}

	return info, nil
}

// StartContainer starts a container by name or ID.
func (c *Client) StartContainer(ctx context.Context, nameOrID string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err := c.api.ContainerStart(ctx, nameOrID, client.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("starting container %s: %w", nameOrID, err)
	}
	return nil
}

// StopContainer stops a container with a timeout.
func (c *Client) StopContainer(ctx context.Context, nameOrID string, timeoutSec int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec+5)*time.Second)
	defer cancel()

	timeout := timeoutSec
	_, err := c.api.ContainerStop(ctx, nameOrID, client.ContainerStopOptions{
		Timeout: &timeout,
	})
	if err != nil {
		return fmt.Errorf("stopping container %s: %w", nameOrID, err)
	}
	return nil
}

// RestartContainer restarts a container with a timeout.
func (c *Client) RestartContainer(ctx context.Context, nameOrID string, timeoutSec int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec+5)*time.Second)
	defer cancel()

	timeout := timeoutSec
	_, err := c.api.ContainerRestart(ctx, nameOrID, client.ContainerRestartOptions{
		Timeout: &timeout,
	})
	if err != nil {
		return fmt.Errorf("restarting container %s: %w", nameOrID, err)
	}
	return nil
}

// RemoveContainer removes a container.
func (c *Client) RemoveContainer(ctx context.Context, nameOrID string, force bool) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.api.ContainerRemove(ctx, nameOrID, client.ContainerRemoveOptions{
		Force: force,
	})
	if err != nil {
		return fmt.Errorf("removing container %s: %w", nameOrID, err)
	}
	return nil
}

// StreamLogs streams logs from a container.
// Caller MUST read the returned io.ReadCloser to completion for streaming to work.
func (c *Client) StreamLogs(ctx context.Context, nameOrID string, follow bool) (io.ReadCloser, error) {
	opts := client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Timestamps: true,
	}

	reader, err := c.api.ContainerLogs(ctx, nameOrID, opts)
	if err != nil {
		return nil, fmt.Errorf("streaming logs for %s: %w", nameOrID, err)
	}
	return reader, nil
}

// FindContainerByName finds a container by name or ID.
func (c *Client) FindContainerByName(ctx context.Context, name string) (*ContainerInfo, error) {
	containers, err := c.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	for _, ctr := range containers {
		if ctr.Name == name {
			return &ctr, nil
		}
		if strings.HasPrefix(ctr.ID, name) {
			return &ctr, nil
		}
	}

	info, err := c.GetContainer(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("container %q not found", name)
	}
	return info, nil
}

// parseContainerSummary converts a Docker API Summary to our ContainerInfo.
func parseContainerSummary(s container.Summary) ContainerInfo {
	info := ContainerInfo{
		ID:     s.ID,
		Image:  s.Image,
		State:  string(s.State),
		Status: s.Status,
		Labels: s.Labels,
	}

	if len(s.Names) > 0 {
		info.Name = strings.TrimPrefix(s.Names[0], "/")
	}

	for _, p := range s.Ports {
		info.Ports = append(info.Ports, PortMapping{
			HostIP:        p.IP.String(),
			HostPort:      fmt.Sprintf("%d", p.PublicPort),
			ContainerPort: fmt.Sprintf("%d", p.PrivatePort),
			Protocol:      string(p.Type),
		})
	}

	return info
}

// FilterContainers filters containers by label.
func FilterContainers(containers []ContainerInfo, labelKey, labelValue string) []ContainerInfo {
	var result []ContainerInfo
	for _, c := range containers {
		if v, ok := c.Labels[labelKey]; ok && (labelValue == "" || v == labelValue) {
			result = append(result, c)
		}
	}
	return result
}

// FilterByNamePrefix filters containers whose name starts with prefix.
func FilterByNamePrefix(containers []ContainerInfo, prefix string) []ContainerInfo {
	var result []ContainerInfo
	for _, c := range containers {
		if strings.HasPrefix(c.Name, prefix) {
			result = append(result, c)
		}
	}
	return result
}

// ShortID returns the first 12 characters of a container ID.
func ShortID(id string) string {
	if len(id) >= 12 {
		return id[:12]
	}
	return id
}
