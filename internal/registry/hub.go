package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	hubAuthURL  = "https://auth.docker.io/token"
	hubTagsURL  = "https://registry.hub.docker.com/v2/library/%s/tags/list"
	defaultTTL  = 24 * time.Hour
	tokenTTL    = 240 * time.Second // Docker Hub tokens valid for 300s, use 240s for safety
)

// HubClient fetches tags from Docker Hub.
type HubClient struct {
	client    *http.Client
	cache     *tagCache
	tokenMu   sync.Mutex
	token     string
	tokenExp  time.Time
}

// NewHubClient creates a new Docker Hub client with caching.
func NewHubClient() *HubClient {
	return &HubClient{
		client: &http.Client{Timeout: 10 * time.Second},
		cache:  newTagCache(defaultTTL),
	}
}

// FetchTags fetches available tags for an image from Docker Hub.
// Returns cached results if available.
func (h *HubClient) FetchTags(image string) ([]TagInfo, error) {
	// Check cache first
	if tags, ok := h.cache.Get(image); ok {
		return tags, nil
	}

	// Get auth token (cached)
	token, err := h.getToken(image)
	if err != nil {
		return nil, fmt.Errorf("getting auth token for %s: %w", image, err)
	}

	// Fetch tags
	tags, err := h.fetchTagsFromHub(image, token)
	if err != nil {
		return nil, err
	}

	// Cache results
	h.cache.Set(image, tags)

	return tags, nil
}

// getToken gets a JWT token from Docker Hub for pulling tags.
// Token is cached for 240 seconds.
func (h *HubClient) getToken(image string) (string, error) {
	h.tokenMu.Lock()
	defer h.tokenMu.Unlock()

	// Return cached token if still valid
	if h.token != "" && time.Now().Before(h.tokenExp) {
		return h.token, nil
	}

	url := fmt.Sprintf("%s?service=registry.docker.io&scope=repository:library/%s:pull", hubAuthURL, image)

	resp, err := h.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	h.token = result.Token
	h.tokenExp = time.Now().Add(tokenTTL)

	return result.Token, nil
}

// fetchTagsFromHub fetches tags from Docker Hub using the auth token.
func (h *HubClient) fetchTagsFromHub(image, token string) ([]TagInfo, error) {
	url := fmt.Sprintf(hubTagsURL, image)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Limit error body read to prevent memory issues
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("hub request failed: %d %s", resp.StatusCode, string(body))
	}

	var result struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Convert to TagInfo (filter out SHA digests)
	var tags []TagInfo
	for _, tag := range result.Tags {
		if strings.HasPrefix(tag, "sha256:") {
			continue
		}
		tags = append(tags, TagInfo{
			Name: tag,
		})
	}

	return tags, nil
}

// ClearCache clears the tag cache and token.
func (h *HubClient) ClearCache() {
	h.cache.Clear()
	h.tokenMu.Lock()
	h.token = ""
	h.tokenExp = time.Time{}
	h.tokenMu.Unlock()
}
