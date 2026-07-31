package components

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

// StatusBadge renders a status indicator.
type StatusBadge struct {
	Status string // running, stopped, healthy, unhealthy, starting
}

// NewStatusBadge creates a new status badge.
func NewStatusBadge(status string) StatusBadge {
	return StatusBadge{Status: status}
}

// Render renders the status badge using cached styles.
func (s StatusBadge) Render() string {
	var style lipgloss.Style
	var icon string

	switch s.Status {
	case "running":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true)
		icon = "●"
	case "stopped":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C757D"))
		icon = "○"
	case "healthy":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4"))
		icon = "✓"
	case "unhealthy":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)
		icon = "✗"
	case "starting":
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D"))
		icon = "⟳"
	default:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C757D"))
		icon = "-"
	}

	return style.Render(icon + " " + s.Status)
}

// PortBadge renders a port number.
type PortBadge struct {
	Port int
}

// NewPortBadge creates a new port badge.
func NewPortBadge(port int) PortBadge {
	return PortBadge{Port: port}
}

// Render renders the port badge.
func (p PortBadge) Render() string {
	if p.Port == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C757D")).
			Render("-")
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Render(fmt.Sprintf(":%d", p.Port))
}
