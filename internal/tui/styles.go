package tui

import (
	"charm.land/lipgloss/v2"
)

// Color palette
var (
	ColorPrimary    = lipgloss.Color("#7D56F4")
	ColorSecondary  = lipgloss.Color("#FF6B6B")
	ColorSuccess    = lipgloss.Color("#4ECDC4")
	ColorWarning    = lipgloss.Color("#FFE66D")
	ColorError      = lipgloss.Color("#FF6B6B")
	ColorMuted      = lipgloss.Color("#6C757D")
	ColorBackground = lipgloss.Color("#1A1A2E")
	ColorForeground = lipgloss.Color("#EAEAEA")
	ColorBorder     = lipgloss.Color("#3D3D5C")
	ColorHighlight  = lipgloss.Color("#7D56F4")
)

// Styles holds all reusable styles
var Styles = struct {
	// App
	Header lipgloss.Style
	Footer lipgloss.Style
	Title  lipgloss.Style

	// Status
	StatusRunning  lipgloss.Style
	StatusStopped  lipgloss.Style
	StatusHealthy  lipgloss.Style
	StatusUnhealthy lipgloss.Style
	StatusStarting lipgloss.Style

	// Text
	Bold      lipgloss.Style
	Italic    lipgloss.Style
	Muted     lipgloss.Style
	Highlight lipgloss.Style

	// Layout
	Box       lipgloss.Style
	BoxActive lipgloss.Style
	Divider   lipgloss.Style

	// List
 ListItem       lipgloss.Style
	ListItemSelected lipgloss.Style

	// Key
	KeyHint lipgloss.Style
	Key     lipgloss.Style

	// Toast
	ToastSuccess lipgloss.Style
	ToastError   lipgloss.Style
	ToastInfo    lipgloss.Style
}{
	// App
	Header: lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Padding(0, 1),

	Footer: lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 1),

	Title: lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorForeground).
		Padding(1, 0),

	// Status
	StatusRunning: lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true),

	StatusStopped: lipgloss.NewStyle().
		Foreground(ColorMuted),

	StatusHealthy: lipgloss.NewStyle().
		Foreground(ColorSuccess),

	StatusUnhealthy: lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true),

	StatusStarting: lipgloss.NewStyle().
		Foreground(ColorWarning),

	// Text
	Bold: lipgloss.NewStyle().
		Bold(true),

	Italic: lipgloss.NewStyle().
		Italic(true),

	Muted: lipgloss.NewStyle().
		Foreground(ColorMuted),

	Highlight: lipgloss.NewStyle().
		Foreground(ColorHighlight).
		Bold(true),

	// Layout
	Box: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2),

	BoxActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2),

	Divider: lipgloss.NewStyle().
		Foreground(ColorBorder).
		Padding(0, 0),

	// List
	ListItem: lipgloss.NewStyle().
		Padding(0, 1),

	ListItemSelected: lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1),

	// Key
	KeyHint: lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 1),

	Key: lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1),

	// Toast
	ToastSuccess: lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true).
		Padding(0, 1),

	ToastError: lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true).
		Padding(0, 1),

	ToastInfo: lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Padding(0, 1),
}

// StatusIcon returns a status icon with style.
func StatusIcon(running bool) string {
	if running {
		return Styles.StatusRunning.Render("●")
	}
	return Styles.StatusStopped.Render("○")
}

// HealthIcon returns a health icon with style.
func HealthIcon(health string) string {
	switch health {
	case "healthy":
		return Styles.StatusHealthy.Render("✓")
	case "unhealthy":
		return Styles.StatusUnhealthy.Render("✗")
	case "starting":
		return Styles.StatusStarting.Render("⟳")
	default:
		return Styles.Muted.Render("-")
	}
}
