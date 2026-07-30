package components

import (
	"charm.land/lipgloss/v2"
)

var (
	helpBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	helpTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	helpSection = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D")).
			Bold(true)

	helpKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Width(20)

	helpDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EAEAEA"))
)

// HelpOverlay displays help information.
type HelpOverlay struct {
	Visible bool
}

// NewHelpOverlay creates a new help overlay.
func NewHelpOverlay() HelpOverlay {
	return HelpOverlay{Visible: false}
}

// Render renders the help overlay.
func (h HelpOverlay) Render() string {
	if !h.Visible {
		return ""
	}

	title := helpTitle.Render("dockkit — Help")

	sections := []struct {
		title string
		items []struct{ key, desc string }
	}{
		{
			title: "Navigation",
			items: []struct{ key, desc string }{
				{"↑/↓ or j/k", "Navigate up/down"},
				{"←/→ or h/l", "Navigate left/right"},
				{"Enter", "Select / Open"},
				{"Esc / q", "Back / Quit"},
				{"Tab", "Next section"},
			},
		},
		{
			title: "Service Actions",
			items: []struct{ key, desc string }{
				{"s", "Start service"},
				{"x", "Stop service"},
				{"r", "Restart service"},
				{"l", "View logs"},
				{"c", "Edit config"},
				{"d", "Remove service"},
			},
		},
		{
			title: "Quick Actions",
			items: []struct{ key, desc string }{
				{"+", "Add new service"},
				{"S", "Start all services"},
				{"X", "Stop all services"},
				{"R", "Refresh status"},
				{"/", "Search / Filter"},
				{"?", "Toggle this help"},
			},
		},
	}

	content := title + "\n\n"

	for _, section := range sections {
		content += helpSection.Render(section.title) + "\n"
		for _, item := range section.items {
			content += helpKey.Render(item.key) + helpDesc.Render(item.desc) + "\n"
		}
		content += "\n"
	}

	return helpBorder.Render(content)
}

// Toggle toggles the help overlay visibility.
func (h *HelpOverlay) Toggle() {
	h.Visible = !h.Visible
}

// Show shows the help overlay.
func (h *HelpOverlay) Show() {
	h.Visible = true
}

// Hide hides the help overlay.
func (h *HelpOverlay) Hide() {
	h.Visible = false
}
