package components

import (
	"charm.land/lipgloss/v2"
)

var (
	toastSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			Padding(0, 1)

	toastError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Padding(0, 1)

	toastInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	toastBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3D3D5C")).
			Padding(0, 1)
)

// ToastType represents the type of toast message.
type ToastType int

const (
	ToastSuccess ToastType = iota
	ToastError
	ToastInfo
)

// Toast displays a temporary message.
type Toast struct {
	Message string
	Type    ToastType
	Visible bool
}

// NewToast creates a new toast.
func NewToast(message string, toastType ToastType) Toast {
	return Toast{
		Message: message,
		Type:    toastType,
		Visible: true,
	}
}

// Render renders the toast.
func (t Toast) Render() string {
	if !t.Visible {
		return ""
	}

	var styled string
	switch t.Type {
	case ToastSuccess:
		styled = toastSuccess.Render("✓ " + t.Message)
	case ToastError:
		styled = toastError.Render("✗ " + t.Message)
	case ToastInfo:
		styled = toastInfo.Render("ℹ " + t.Message)
	}

	return toastBorder.Render(styled)
}

// Hide hides the toast.
func (t *Toast) Hide() {
	t.Visible = false
}

// Show shows the toast.
func (t *Toast) Show() {
	t.Visible = true
}
