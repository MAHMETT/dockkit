package errors

import "fmt"

// ErrorCode represents a category of error.
type ErrorCode int

const (
	// Config errors (1xxx)
	ErrCodeConfigNotFound  ErrorCode = 1001
	ErrCodeConfigCorrupted ErrorCode = 1002
	ErrCodeConfigInvalid   ErrorCode = 1003
	ErrCodeConfigWrite     ErrorCode = 1004

	// Docker errors (2xxx)
	ErrCodeDockerNotRunning ErrorCode = 2001
	ErrCodeDockerPermission ErrorCode = 2002
	ErrCodeDockerTimeout    ErrorCode = 2003

	// Container errors (3xxx)
	ErrCodeContainerExists   ErrorCode = 3001
	ErrCodeContainerNotFound ErrorCode = 3002
	ErrCodeContainerFailed   ErrorCode = 3003

	// Port errors (4xxx)
	ErrCodePortConflict   ErrorCode = 4001
	ErrCodePortOccupied   ErrorCode = 4002
	ErrCodePortInvalid    ErrorCode = 4003
	ErrCodePortUnavailable ErrorCode = 4004

	// Template errors (5xxx)
	ErrCodeTemplateNotFound  ErrorCode = 5001
	ErrCodeTemplateInvalid   ErrorCode = 5002
	ErrCodeTemplateCorrupted ErrorCode = 5003

	// Registry errors (6xxx)
	ErrCodeHubRateLimited ErrorCode = 6001
	ErrCodeHubOffline     ErrorCode = 6002
	ErrCodeHubNotFound    ErrorCode = 6003

	// Network errors (7xxx)
	ErrCodeNetworkExists ErrorCode = 7001

	// System errors (9xxx)
	ErrCodePermissionDenied ErrorCode = 9001
	ErrCodeDiskFull         ErrorCode = 9002
	ErrCodeInternal         ErrorCode = 9999
)

// DockkitError is the project-wide error type.
type DockkitError struct {
	Code    ErrorCode
	Message string
	Detail  string
	Inner   error
}

func (e *DockkitError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *DockkitError) Unwrap() error {
	return e.Inner
}

// Is returns true if target matches the error code.
func (e *DockkitError) Is(target error) bool {
	if t, ok := target.(*DockkitError); ok {
		return e.Code == t.Code
	}
	return false
}

// New creates a new DockkitError.
func New(code ErrorCode, message string) *DockkitError {
	return &DockkitError{Code: code, Message: message}
}

// Wrap wraps an existing error with a DockkitError.
func Wrap(code ErrorCode, message string, inner error) *DockkitError {
	return &DockkitError{
		Code:    code,
		Message: message,
		Inner:   inner,
	}
}

// WithDetail adds detail to a DockkitError.
func (e *DockkitError) WithDetail(detail string) *DockkitError {
	return &DockkitError{
		Code:    e.Code,
		Message: e.Message,
		Detail:  detail,
		Inner:   e.Inner,
	}
}

// Predefined errors — use these as sentinel values.
var (
	ErrConfigNotFound  = New(ErrCodeConfigNotFound, "config file not found")
	ErrConfigCorrupted = New(ErrCodeConfigCorrupted, "config file is corrupted")
	ErrConfigInvalid   = New(ErrCodeConfigInvalid, "config validation failed")
	ErrConfigWrite     = New(ErrCodeConfigWrite, "failed to write config")

	ErrDockerNotRunning = New(ErrCodeDockerNotRunning, "Docker is not running")
	ErrDockerPermission = New(ErrCodeDockerPermission, "permission denied for Docker")
	ErrDockerTimeout    = New(ErrCodeDockerTimeout, "Docker operation timed out")

	ErrContainerExists   = New(ErrCodeContainerExists, "container already exists")
	ErrContainerNotFound = New(ErrCodeContainerNotFound, "container not found")
	ErrContainerFailed   = New(ErrCodeContainerFailed, "container failed to start")

	ErrPortConflict    = New(ErrCodePortConflict, "port conflict detected")
	ErrPortOccupied    = New(ErrCodePortOccupied, "port is occupied by another process")
	ErrPortInvalid     = New(ErrCodePortInvalid, "invalid port number")
	ErrPortUnavailable = New(ErrCodePortUnavailable, "no available port in range")

	ErrTemplateNotFound  = New(ErrCodeTemplateNotFound, "template not found")
	ErrTemplateInvalid   = New(ErrCodeTemplateInvalid, "template validation failed")
	ErrTemplateCorrupted = New(ErrCodeTemplateCorrupted, "template file is corrupted")

	ErrHubRateLimited = New(ErrCodeHubRateLimited, "Docker Hub rate limited")
	ErrHubOffline     = New(ErrCodeHubOffline, "Docker Hub unreachable")
	ErrHubNotFound    = New(ErrCodeHubNotFound, "image not found on Docker Hub")

	ErrNetworkExists = New(ErrCodeNetworkExists, "network already exists")

	ErrPermissionDenied = New(ErrCodePermissionDenied, "permission denied")
	ErrDiskFull         = New(ErrCodeDiskFull, "insufficient disk space")
	ErrInternal         = New(ErrCodeInternal, "internal error")
)

// UserMessage returns a human-readable message for display.
func UserMessage(err error) string {
	if e, ok := err.(*DockkitError); ok {
		switch e.Code {
		case ErrCodeDockerNotRunning:
			return "Docker is not running. Start Docker Desktop or run: sudo systemctl start docker"
		case ErrCodePortConflict:
			if e.Detail != "" {
				return fmt.Sprintf("Port conflict: %s", e.Detail)
			}
			return "Port conflict detected. Another service is using the same port."
		case ErrCodeContainerExists:
			if e.Detail != "" {
				return fmt.Sprintf("Container already exists: %s", e.Detail)
			}
			return "A container with this name already exists."
		case ErrCodePermissionDenied:
			return "Permission denied. Try running with sudo or add your user to the docker group."
		case ErrCodeConfigCorrupted:
			return "Config file is corrupted. A backup has been created."
		case ErrCodeHubOffline:
			return "Docker Hub is unreachable. Using cached data."
		case ErrCodeHubRateLimited:
			return "Docker Hub rate limited. Using cached data."
		default:
			return e.Message
		}
	}
	return err.Error()
}
