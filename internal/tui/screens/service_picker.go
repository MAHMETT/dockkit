package screens

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// TemplateEntry represents a service template.
type TemplateEntry struct {
	Name        string
	Description string
	Category    string
	Icon        string
	Versions    []string
}

// ServicePickerModel browses available templates.
type ServicePickerModel struct {
	templates []TemplateEntry
	cursor    int
	keys      screenKeys
	width     int
	height    int
}

// NewServicePickerModel creates a new service picker.
func NewServicePickerModel(templates []TemplateEntry) ServicePickerModel {
	return ServicePickerModel{
		templates: templates,
		cursor:    0,
		keys:      defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *ServicePickerModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Update handles messages.
func (m ServicePickerModel) Update(msg tea.Msg) (ServicePickerModel, tea.Cmd) {
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

// View renders the service picker.
func (m ServicePickerModel) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("Add New Service"))
	b.WriteString("\n\n")

	// Group by category with sorted keys for deterministic order
	categories := make(map[string][]TemplateEntry)
	for _, tmpl := range m.templates {
		categories[tmpl.Category] = append(categories[tmpl.Category], tmpl)
	}

	catNames := make([]string, 0, len(categories))
	for cat := range categories {
		catNames = append(catNames, cat)
	}
	sort.Strings(catNames)

	cursor := 0
	for _, cat := range catNames {
		templates := categories[cat]
		b.WriteString(styleHighlight.Render(cat))
		b.WriteString("\n")
		for _, tmpl := range templates {
			selected := cursor == m.cursor
			versions := strings.Join(tmpl.Versions, ", ")
			row := fmt.Sprintf("  %s %-15s %s", tmpl.Icon, tmpl.Name, versions)
			if selected {
				b.WriteString(styleHighlight.Render(row))
			} else {
				b.WriteString(row)
			}
			b.WriteString("\n")
			cursor++
		}
		b.WriteString("\n")
	}

	b.WriteString(styleMuted.Render("  [↑/↓] Navigate   [Enter] Select   [Esc] Back"))
	return b.String()
}
