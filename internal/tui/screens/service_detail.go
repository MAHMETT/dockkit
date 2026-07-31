package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
	"github.com/MAHMETT/dockkit/internal/tui/components"
)

// ServiceDetailModel shows details for a single service.
type ServiceDetailModel struct {
	service  ServiceEntry
	action   int
	keys     screenKeys
	width    int
	height   int
}

// NewServiceDetailModel creates a new service detail screen.
func NewServiceDetailModel(service ServiceEntry) ServiceDetailModel {
	return ServiceDetailModel{
		service: service,
		action:  0,
		keys:    defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *ServiceDetailModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Update handles messages.
func (m ServiceDetailModel) Update(msg tea.Msg) (ServiceDetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.action > 0 {
				m.action--
			}
		case key.Matches(msg, m.keys.Down):
			if m.action < 5 {
				m.action++
			}
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenDashboard}
			}
		case key.Matches(msg, m.keys.Enter):
			return m.handleAction()
		}
	}
	return m, nil
}

func (m ServiceDetailModel) handleAction() (ServiceDetailModel, tea.Cmd) {
	svc := m.service
	switch m.action {
	case 0:
		return m, func() tea.Msg {
			return messages.ToastMsg{Message: fmt.Sprintf("Starting %s...", svc.Name), Type: 0}
		}
	case 1:
		return m, func() tea.Msg {
			return messages.ToastMsg{Message: fmt.Sprintf("Stopping %s...", svc.Name), Type: 0}
		}
	case 2:
		return m, func() tea.Msg {
			return messages.ToastMsg{Message: fmt.Sprintf("Restarting %s...", svc.Name), Type: 0}
		}
	case 3:
		return m, func() tea.Msg {
			return messages.NavigateToMsg{Screen: messages.ScreenLogsViewer, Data: svc}
		}
	case 4:
		return m, func() tea.Msg {
			return messages.NavigateToMsg{Screen: messages.ScreenConfigWizard, Data: svc}
		}
	case 5:
		return m, func() tea.Msg {
			return messages.NavigateToMsg{Screen: messages.ScreenDashboard}
		}
	}
	return m, nil
}

// View renders the service detail screen.
func (m ServiceDetailModel) View() string {
	var b strings.Builder

	title := fmt.Sprintf("%s %s", m.service.Name, m.service.Version)
	b.WriteString(styleHeader.Render(title))
	b.WriteString("\n\n")

	statusBadge := components.NewStatusBadge(m.service.Status)
	b.WriteString("Status: ")
	b.WriteString(statusBadge.Render())
	b.WriteString("\n\n")

	b.WriteString(styleBold.Render("Details"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Container: %s\n", m.service.ContainerName))
	b.WriteString(fmt.Sprintf("  Port:      %d\n", m.service.Port))
	b.WriteString("\n")

	b.WriteString(styleBold.Render("Actions"))
	b.WriteString("\n")

	actions := []string{"Start", "Stop", "Restart", "Logs", "Config", "Back"}
	for i, action := range actions {
		if i == m.action {
			b.WriteString(styleHighlight.Render("  > " + action))
		} else {
			b.WriteString("  " + action)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleMuted.Render("  [↑/↓] Navigate   [Enter] Select   [Esc] Back"))

	return b.String()
}
