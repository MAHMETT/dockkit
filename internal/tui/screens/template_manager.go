package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// TemplateManagerModel lists available templates.
type TemplateManagerModel struct {
	templates []TemplateEntry
	cursor    int
	keys      screenKeys
	width     int
	height    int
}

// NewTemplateManagerModel creates a new template manager.
func NewTemplateManagerModel(templates []TemplateEntry) TemplateManagerModel {
	return TemplateManagerModel{
		templates: templates,
		cursor:    0,
		keys:      defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *TemplateManagerModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Update handles messages.
func (m TemplateManagerModel) Update(msg tea.Msg) (TemplateManagerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.templates)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenDashboard}
			}
		case key.Matches(msg, m.keys.Enter):
			if len(m.templates) > 0 {
				tmpl := m.templates[m.cursor]
				return m, func() tea.Msg {
					return messages.NavigateToMsg{
						Screen: messages.ScreenConfigWizard,
						Data:   tmpl,
					}
				}
			}
		}
	}
	return m, nil
}

// View renders the template manager.
func (m TemplateManagerModel) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("Templates"))
	b.WriteString("\n\n")

	for i, tmpl := range m.templates {
		selected := i == m.cursor
		versions := fmt.Sprintf("%d versions", len(tmpl.Versions))
		row := fmt.Sprintf("  %s %-15s %s", tmpl.Icon, tmpl.Name, versions)
		if selected {
			b.WriteString(styleHighlight.Render(row))
		} else {
			b.WriteString(row)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleMuted.Render("  [↑/↓] Navigate   [Enter] Use   [Esc] Back"))
	return b.String()
}
