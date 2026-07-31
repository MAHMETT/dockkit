package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// LogsViewerModel displays container logs.
type LogsViewerModel struct {
	service  ServiceEntry
	viewport viewport.Model
	lines    []string
	follow   bool
	keys     screenKeys
	width    int
	height   int
}

// NewLogsViewerModel creates a new logs viewer.
func NewLogsViewerModel(service ServiceEntry) LogsViewerModel {
	vp := viewport.New(
		viewport.WithWidth(80),
		viewport.WithHeight(20),
	)

	return LogsViewerModel{
		service:  service,
		viewport: vp,
		lines:    []string{"Loading logs..."},
		follow:   true,
		keys:     defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *LogsViewerModel) SetSize(w, h int) {
	m.width = w
	m.height = h - 4
	if m.height < 1 {
		m.height = 1
	}
	m.viewport.SetWidth(w - 4)
	m.viewport.SetHeight(m.height)
}

// Update handles messages.
func (m LogsViewerModel) Update(msg tea.Msg) (LogsViewerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{
					Screen: messages.ScreenServiceDetail,
					Data:   m.service,
				}
			}
		case key.Matches(msg, m.keys.Up):
			m.viewport.ScrollUp(1)
		case key.Matches(msg, m.keys.Down):
			m.viewport.ScrollDown(1)
		case msg.String() == "f" || msg.String() == "F":
			m.follow = !m.follow
		}
	case logsUpdateMsg:
		m.lines = msg.lines
		content := strings.Join(m.lines, "\n")
		m.viewport.SetContent(content)
		if m.follow {
			m.viewport.GotoBottom()
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the logs viewer.
func (m LogsViewerModel) View() string {
	var b strings.Builder

	title := fmt.Sprintf("Logs: %s %s", m.service.Name, m.service.Version)
	b.WriteString(styleHeader.Render(title))
	b.WriteString("\n")

	followStatus := "OFF"
	if m.follow {
		followStatus = "ON"
	}
	b.WriteString(styleMuted.Render(fmt.Sprintf("  [F] Follow: %s   [Esc] Back", followStatus)))
	b.WriteString("\n")

	b.WriteString(m.viewport.View())

	return b.String()
}

type logsUpdateMsg struct {
	lines []string
}
