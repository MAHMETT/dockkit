package errors

import "fmt"

type DockkitError struct {
	Code    string
	Message string
	Detail  string
}

func (e *DockkitError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code, message string) *DockkitError {
	return &DockkitError{Code: code, Message: message}
}

func Wrap(code, message, detail string) *DockkitError {
	return &DockkitError{Code: code, Message: message, Detail: detail}
}

var (
	ErrDockerNotRunning = New("DOCKER_NOT_RUNNING", "Docker is not running")
	ErrPortConflict     = New("PORT_CONFLICT", "Port already in use")
	ErrContainerExists  = New("CONTAINER_EXISTS", "Container already exists")
	ErrConfigCorrupted  = New("CONFIG_CORRUPTED", "Config file is corrupted")
	ErrHubRateLimited   = New("HUB_RATE_LIMITED", "Docker Hub rate limited")
	ErrHubOffline       = New("HUB_OFFLINE", "Docker Hub unreachable")
	ErrPermissionDenied = New("PERMISSION_DENIED", "Permission denied")
)
