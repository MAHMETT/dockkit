package screens

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// TemplateEntry represents a service template with default config.
type TemplateEntry struct {
	Name            string
	Description     string
	Category        string
	Icon            string
	Versions        []string
	DefaultPort     int
	DefaultUser     string
	DefaultPassword string
	DefaultDatabase string
}

// ServicePickerModel browses available templates.
type ServicePickerModel struct {
	templates    []TemplateEntry
	cursor       int
	versionCursor int
	selectingVersion bool // true when selecting version sub-item
	keys         screenKeys
	width        int
	height       int
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
		if m.selectingVersion {
			return m.handleVersionKeys(msg)
		}
		return m.handleServiceKeys(msg)
	}
	return m, nil
}

func (m ServicePickerModel) handleServiceKeys(msg tea.KeyPressMsg) (ServicePickerModel, tea.Cmd) {
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
			if len(tmpl.Versions) == 1 {
				// Only one version, go directly to config
				return m, func() tea.Msg {
					return messages.NavigateToMsg{
						Screen: messages.ScreenConfigWizard,
						Data: ConfigWizardData{
							Template:      tmpl,
							SelectedVersion: tmpl.Versions[0],
						},
					}
				}
			}
			// Multiple versions, enter version selection mode
			m.selectingVersion = true
			m.versionCursor = 0
		}
	case msg.String() == "/":
		// TODO: Implement search/filter
	}
	return m, nil
}

func (m ServicePickerModel) handleVersionKeys(msg tea.KeyPressMsg) (ServicePickerModel, tea.Cmd) {
	tmpl := m.templates[m.cursor]

	switch {
	case key.Matches(msg, m.keys.Up):
		if m.versionCursor > 0 {
			m.versionCursor--
		}
	case key.Matches(msg, m.keys.Down):
		if m.versionCursor < len(tmpl.Versions)-1 {
			m.versionCursor++
		}
	case key.Matches(msg, m.keys.Escape):
		m.selectingVersion = false
		return m, nil
	case key.Matches(msg, m.keys.Enter):
		if m.versionCursor < len(tmpl.Versions) {
			selectedVersion := tmpl.Versions[m.versionCursor]
			return m, func() tea.Msg {
				return messages.NavigateToMsg{
					Screen: messages.ScreenConfigWizard,
					Data: ConfigWizardData{
						Template:        tmpl,
						SelectedVersion: selectedVersion,
					},
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

	// Group by category with sorted keys
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
			indent := "  "
			if selected {
				indent = "▶ "
			}

			row := fmt.Sprintf("%s%s %s", indent, tmpl.Icon, tmpl.Name)
			b.WriteString(row)
			b.WriteString("\n")

			// Show versions if this template is selected
			if selected {
				for vi, ver := range tmpl.Versions {
					versionSelected := m.selectingVersion && vi == m.versionCursor
					prefix := "    "
					if versionSelected {
						prefix = "    ▶ "
					}
					b.WriteString(fmt.Sprintf("%sv%s", prefix, ver))
					b.WriteString("\n")
				}
			}

			cursor++
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if m.selectingVersion {
		b.WriteString(styleMuted.Render("  [↑/↓] Select version   [Enter] Confirm   [Esc] Back"))
	} else {
		b.WriteString(styleMuted.Render("  [↑/↓] Navigate   [Enter] Select   [Esc] Back"))
	}
	return b.String()
}
