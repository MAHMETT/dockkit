package errors

import (
	"errors"
	"testing"
)

func TestDockkitError_Error(t *testing.T) {
	e := New(ErrCodePortConflict, "port conflict")
	if e.Error() != "[4001] port conflict" {
		t.Errorf("Error() = %q", e.Error())
	}
}

func TestDockkitError_WithDetail(t *testing.T) {
	e := New(ErrCodePortConflict, "port conflict").WithDetail("port 5432 used by pg16")
	if e.Detail != "port 5432 used by pg16" {
		t.Errorf("Detail = %q", e.Detail)
	}
	if e.Error() != "[4001] port conflict: port 5432 used by pg16" {
		t.Errorf("Error() = %q", e.Error())
	}
}

func TestDockkitError_Unwrap(t *testing.T) {
	inner := errors.New("original error")
	e := Wrap(ErrCodeInternal, "wrapped", inner)
	if e.Unwrap() != inner {
		t.Error("Unwrap() did not return inner error")
	}
}

func TestDockkitError_Is(t *testing.T) {
	e1 := New(ErrCodePortConflict, "test")
	e2 := New(ErrCodePortConflict, "other")
	e3 := New(ErrCodeDockerNotRunning, "test")

	if !errors.Is(e1, e2) {
		t.Error("errors.Is(same code) = false, want true")
	}
	if errors.Is(e1, e3) {
		t.Error("errors.Is(different code) = true, want false")
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *DockkitError
		code ErrorCode
	}{
		{"ErrDockerNotRunning", ErrDockerNotRunning, ErrCodeDockerNotRunning},
		{"ErrPortConflict", ErrPortConflict, ErrCodePortConflict},
		{"ErrContainerExists", ErrContainerExists, ErrCodeContainerExists},
		{"ErrConfigCorrupted", ErrConfigCorrupted, ErrCodeConfigCorrupted},
		{"ErrHubRateLimited", ErrHubRateLimited, ErrCodeHubRateLimited},
		{"ErrHubOffline", ErrHubOffline, ErrCodeHubOffline},
		{"ErrPermissionDenied", ErrPermissionDenied, ErrCodePermissionDenied},
		{"ErrTemplateNotFound", ErrTemplateNotFound, ErrCodeTemplateNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Code = %d, want %d", tt.err.Code, tt.code)
			}
			if tt.err.Error() == "" {
				t.Error("Error() returned empty string")
			}
		})
	}
}

func TestUserMessage(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"docker not running", ErrDockerNotRunning, "Docker is not running"},
		{"port conflict", ErrPortConflict, "Port conflict detected"},
		{"permission denied", ErrPermissionDenied, "Permission denied"},
		{"hub offline", ErrHubOffline, "Docker Hub is unreachable"},
		{"generic error", errors.New("something"), "something"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := UserMessage(tt.err)
			if msg == "" {
				t.Error("UserMessage() returned empty string")
			}
		})
	}
}

func TestUserMessage_WithDetail(t *testing.T) {
	e := ErrPortConflict.WithDetail("port 5432 used by pg16")
	msg := UserMessage(e)
	if msg != "Port conflict: port 5432 used by pg16" {
		t.Errorf("UserMessage() = %q", msg)
	}
}
