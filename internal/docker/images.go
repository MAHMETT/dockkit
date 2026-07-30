package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/moby/moby/client"
)

// ImageInfo holds parsed image information.
type ImageInfo struct {
	ID       string
	RepoTags []string
	Size     int64
	Created  int64
}

// ListImages returns all local images.
func (c *Client) ListImages(ctx context.Context) ([]ImageInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := c.api.ImageList(ctx, client.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing images: %w", err)
	}

	images := make([]ImageInfo, 0, len(result.Items))
	for _, img := range result.Items {
		images = append(images, ImageInfo{
			ID:       img.ID,
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  img.Created,
		})
	}
	return images, nil
}

// PullImage pulls an image from Docker Hub.
// The returned reader MUST be consumed to completion for the pull to finish.
// Use io.Copy(io.Discard, reader) if you don't need the output.
func (c *Client) PullImage(ctx context.Context, ref string) (io.ReadCloser, error) {
	reader, err := c.api.ImagePull(ctx, ref, client.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("pulling image %s: %w", ref, err)
	}
	return reader, nil
}

// RemoveImage removes a local image.
func (c *Client) RemoveImage(ctx context.Context, ref string) error {
	_, err := c.api.ImageRemove(ctx, ref, client.ImageRemoveOptions{})
	if err != nil {
		return fmt.Errorf("removing image %s: %w", ref, err)
	}
	return nil
}

// ImageExists checks if an image exists locally.
func (c *Client) ImageExists(ctx context.Context, ref string) (bool, error) {
	_, err := c.api.ImageInspect(ctx, ref)
	if err != nil {
		// Docker SDK wraps not-found errors with "No such image"
		if strings.Contains(err.Error(), "No such image") {
			return false, nil
		}
		return false, fmt.Errorf("inspecting image %s: %w", ref, err)
	}
	return true, nil
}
