package registry

// TagInfo holds information about a Docker image tag.
type TagInfo struct {
	Name       string `json:"name"`
	LastUpdated string `json:"last_updated"`
	Size       int64  `json:"size"`
}

// TagResponse is the Docker Hub API response for listing tags.
type TagResponse struct {
	Name string    `json:"name"`
	Tags []TagInfo `json:"tags"`
}

// HubImage represents an image on Docker Hub.
type HubImage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"star_count"`
	Pulls       int    `json:"pull_count"`
	IsOfficial  bool   `json:"is_official"`
	IsAutomated bool   `json:"is_automated"`
}
