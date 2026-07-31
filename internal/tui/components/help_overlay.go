package components

import (
	"charm.land/bubbles/v2/key"
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

	helpMuted = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C757D"))
)

// KeyMapForHelp defines the keybindings for the help overlay.
type KeyMapForHelp struct {
	Navigation  []key.Binding
	Actions     []key.Binding
	Quick       []key.Binding
}

// DefaultKeyMapForHelp returns the default help keybindings.
func DefaultKeyMapForHelp() KeyMapForHelp {
	return KeyMapForHelp{
		Navigation: []key.Binding{
			key.NewBinding(key.WithKeys("up/k"), key.WithHelp("↑/k", "up")),
			key.NewBinding(key.WithKeys("down/j"), key.WithHelp("↓/j", "down")),
			key.NewBinding(key.WithKeys("left/h"), key.WithHelp("←/h", "left")),
			key.NewBinding(key.WithKeys("right/l"), key.WithHelp("→/l", "right")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		},
		Actions: []key.Binding{
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start")),
			key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "stop")),
			key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "restart")),
			key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logs")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "config")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "remove")),
		},
		Quick: []key.Binding{
			key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "add service")),
			key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "start all")),
			key.NewBinding(key.WithKeys("X"), key.WithHelp("X", "stop all")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
			key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
			key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		},
	}
}

// HelpOverlay displays help information.
type HelpOverlay struct {
	Visible bool
	keyMap  KeyMapForHelp
}

// NewHelpOverlay creates a new help overlay.
func NewHelpOverlay() HelpOverlay {
	return HelpOverlay{
		Visible: false,
		keyMap:  DefaultKeyMapForHelp(),
	}
}

// Render renders the help overlay.
func (h HelpOverlay) Render() string {
	if !h.Visible {
		return ""
	}

	title := helpTitle.Render("dockkit — Help")

	sections := []struct {
		title string
		items []key.Binding
	}{
		{"Navigation", h.keyMap.Navigation},
		{"Service Actions", h.keyMap.Actions},
		{"Quick Actions", h.keyMap.Quick},
	}

	content := title + "\n\n"

	for _, section := range sections {
		content += helpSection.Render(section.title) + "\n"
		for _, item := range section.items {
			help := item.Help()
			content += helpKey.Render(help.Key) + helpDesc.Render(help.Desc) + "\n"
		}
		content += "\n"
	}

	content += helpMuted.Render("Press ? or Esc to close")

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
