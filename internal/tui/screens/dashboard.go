package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// Local styles to avoid import cycle with tui package.
var (
	colorPrimary = lipgloss.Color("#7D56F4")
	colorMuted   = lipgloss.Color("#6C757D")
	colorHighlight = lipgloss.Color("#7D56F4")

	styleHeader = lipgloss.NewStyle().Bold(true).Foreground(colorPrimary).Padding(0, 1)
	styleMuted = lipgloss.NewStyle().Foreground(colorMuted)
	styleHighlight = lipgloss.NewStyle().Foreground(colorHighlight).Bold(true)
	styleBold = lipgloss.NewStyle().Bold(true)
)

// KeyMap for screens (duplicated to avoid cycle).
type screenKeys struct {
	Up, Down, Enter, Escape key.Binding
}

var defaultScreenKeys = screenKeys{
	Up:      key.NewBinding(key.WithKeys("up", "k")),
	Down:    key.NewBinding(key.WithKeys("down", "j")),
	Enter:   key.NewBinding(key.WithKeys("enter")),
	Escape:  key.NewBinding(key.WithKeys("esc")),
}

// ServiceEntry represents a service row in the dashboard.
type ServiceEntry struct {
	Name          string
	Version       string
	Port          int
	Status        string // running, stopped, healthy, unhealthy
	ContainerName string
}

// DashboardModel is the dashboard screen.
type DashboardModel struct {
	services []ServiceEntry
	cursor   int
	keys     screenKeys
	width    int
	height   int
}

// NewDashboardModel creates a new dashboard.
func NewDashboardModel() DashboardModel {
	return DashboardModel{
		services: []ServiceEntry{},
		cursor:   0,
		keys:     defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *DashboardModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// SetServices updates the service list.
func (m *DashboardModel) SetServices(services []ServiceEntry) {
	m.services = services
	if m.cursor >= len(services) && len(services) > 0 {
		m.cursor = len(services) - 1
	}
}

// Update handles messages.
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.services)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Enter):
			if len(m.services) > 0 {
				svc := m.services[m.cursor]
				return m, func() tea.Msg {
					return messages.NavigateToMsg{
						Screen: messages.ScreenServiceDetail,
						Data:   svc,
					}
				}
			}
		case msg.String() == "+":
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenServicePicker}
			}
		}
	}
	return m, nil
}

// View renders the dashboard.
func (m DashboardModel) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("🐳 dockkit v1.0.0"))
	b.WriteString("\n\n")

	b.WriteString(styleBold.Render("Services"))
	b.WriteString("\n")

	if len(m.services) == 0 {
		b.WriteString(styleMuted.Render("  No services configured. Press [+] to add."))
	} else {
		for i, svc := range m.services {
			row := renderServiceRow(svc, i == m.cursor)
			b.WriteString(row)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(styleMuted.Render("  [+] Add   [Enter] Detail   [?] Help   [ctrl+c] Quit"))

	return b.String()
}

// renderServiceRow renders a single service row.
func renderServiceRow(svc ServiceEntry, selected bool) string {
	icon := "○"
	if svc.Status == "running" || svc.Status == "healthy" {
		icon = "●"
	}

	portStr := "-"
	if svc.Port > 0 {
		portStr = fmt.Sprintf(":%d", svc.Port)
	}

	row := fmt.Sprintf("  %s %-20s  %-8s  %s", icon, svc.Name+" "+svc.Version, portStr, svc.ContainerName)

	if selected {
		return styleHighlight.Render(row)
	}
	return row
}
