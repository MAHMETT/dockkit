package components

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

// Cached styles for status badges (avoid allocation per render).
var (
	styleRunning   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4")).Bold(true)
	styleStopped   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C757D"))
	styleHealthy   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))
	styleUnhealthy = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true)
	styleStarting  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFE66D"))
	styleDefault   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C757D"))
	stylePort      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	stylePortNone  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C757D"))
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
		style = styleRunning
		icon = "●"
	case "stopped":
		style = styleStopped
		icon = "○"
	case "healthy":
		style = styleHealthy
		icon = "✓"
	case "unhealthy":
		style = styleUnhealthy
		icon = "✗"
	case "starting":
		style = styleStarting
		icon = "⟳"
	default:
		style = styleDefault
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
		return stylePortNone.Render("-")
	}
	return stylePort.Render(fmt.Sprintf(":%d", p.Port))
}
