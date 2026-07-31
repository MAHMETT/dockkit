package components

import (
	"fmt"

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

	helpIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C757D")).
			Faint(true)
)

// helpItem is a single help entry.
type helpItem struct {
	section string
	key     string
	desc    string
}

// KeyMapForHelp defines the keybindings for the help overlay.
type KeyMapForHelp struct {
	Navigation []key.Binding
	Actions    []key.Binding
	Quick      []key.Binding
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
	Visible      bool
	scrollOffset int
	keyMap       KeyMapForHelp
	termHeight   int
}

// NewHelpOverlay creates a new help overlay.
func NewHelpOverlay() HelpOverlay {
	return HelpOverlay{
		Visible: false,
		keyMap:  DefaultKeyMapForHelp(),
	}
}

// SetTermHeight sets the terminal height for scroll calculation.
func (h *HelpOverlay) SetTermHeight(height int) {
	h.termHeight = height
}

// ScrollUp scrolls the help overlay up.
func (h *HelpOverlay) ScrollUp() {
	if h.scrollOffset > 0 {
		h.scrollOffset--
	}
}

// ScrollDown scrolls the help overlay down.
func (h *HelpOverlay) ScrollDown() {
	items := h.allItems()
	visible := h.visibleHeight()
	if h.scrollOffset < len(items)-visible {
		h.scrollOffset++
	}
}

// visibleHeight calculates how many help items fit on screen.
func (h *HelpOverlay) visibleHeight() int {
	// Title (2 lines) + border (2 lines) + footer (1 line) = 5 lines overhead
	height := h.termHeight - 5
	if height < 3 {
		height = 3
	}
	return height
}

// allItems flattens all help sections into a single list.
func (h *HelpOverlay) allItems() []helpItem {
	var items []helpItem

	sections := []struct {
		title string
		items []key.Binding
	}{
		{"Navigation", h.keyMap.Navigation},
		{"Service Actions", h.keyMap.Actions},
		{"Quick Actions", h.keyMap.Quick},
	}

	for _, section := range sections {
		for _, item := range section.items {
			help := item.Help()
			items = append(items, helpItem{
				section: section.title,
				key:     help.Key,
				desc:    help.Desc,
			})
		}
	}

	return items
}

// Render renders the help overlay.
func (h HelpOverlay) Render() string {
	if !h.Visible {
		return ""
	}

	title := helpTitle.Render("dockkit — Help")

	items := h.allItems()
	visible := h.visibleHeight()

	// Apply scroll bounds
	if h.scrollOffset > len(items)-visible {
		h.scrollOffset = len(items) - visible
	}
	if h.scrollOffset < 0 {
		h.scrollOffset = 0
	}

	// Slice visible items
	start := h.scrollOffset
	end := start + visible
	if end > len(items) {
		end = len(items)
	}
	visibleItems := items[start:end]

	content := title + "\n\n"

	// Render visible items
	lastSection := ""
	for _, item := range visibleItems {
		if item.section != lastSection {
			if lastSection != "" {
				content += "\n"
			}
			content += helpSection.Render(item.section) + "\n"
			lastSection = item.section
		}
		content += helpKey.Render(item.key) + helpDesc.Render(item.desc) + "\n"
	}

	// Scroll indicator
	if len(items) > visible {
		content += "\n"
		content += helpIndicator.Render(fmt.Sprintf(
			"↑↓ to scroll (%d/%d)   ? or Esc to close",
			start+1, len(items),
		))
	} else {
		content += "\n" + helpMuted.Render("Press ? or Esc to close")
	}

	return helpBorder.Render(content)
}

// Toggle toggles the help overlay visibility.
func (h *HelpOverlay) Toggle() {
	h.Visible = !h.Visible
	if !h.Visible {
		h.scrollOffset = 0
	}
}

// Show shows the help overlay.
func (h *HelpOverlay) Show() {
	h.Visible = true
	h.scrollOffset = 0
}

// Hide hides the help overlay.
func (h *HelpOverlay) Hide() {
	h.Visible = false
	h.scrollOffset = 0
}
