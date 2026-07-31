package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// ConfigWizardModel handles service configuration.
type ConfigWizardModel struct {
	serviceName string
	version     string
	inputs      []textinput.Model
	focus       int
	keys        screenKeys
	width       int
	height      int
}

// NewConfigWizardModel creates a new config wizard.
func NewConfigWizardModel(serviceName, version string) ConfigWizardModel {
	inputs := make([]textinput.Model, 4)

	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = 64
	}

	inputs[0].Placeholder = "5432"
	inputs[0].Prompt = "Port:          "

	inputs[1].Placeholder = "postgres"
	inputs[1].Prompt = "Username:      "

	inputs[2].Placeholder = "postgres"
	inputs[2].Prompt = "Password:      "
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].EchoCharacter = '•'

	inputs[3].Placeholder = "postgres"
	inputs[3].Prompt = "Database:      "

	return ConfigWizardModel{
		serviceName: serviceName,
		version:     version,
		inputs:      inputs,
		focus:       0,
		keys:        defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *ConfigWizardModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Update handles messages.
func (m ConfigWizardModel) Update(msg tea.Msg) (ConfigWizardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.focus > 0 {
				m.inputs[m.focus].Blur()
				m.focus--
				m.inputs[m.focus].Focus()
			}
		case key.Matches(msg, m.keys.Down):
			if m.focus < len(m.inputs)-1 {
				m.inputs[m.focus].Blur()
				m.focus++
				m.inputs[m.focus].Focus()
			}
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenServicePicker}
			}
		case key.Matches(msg, m.keys.Enter):
			return m, func() tea.Msg {
				return messages.ToastMsg{
					Message: fmt.Sprintf("%s %s configured!", m.serviceName, m.version),
					Type:    0,
				}
			}
		}
	}

	if m.focus < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the config wizard.
func (m ConfigWizardModel) View() string {
	var b strings.Builder

	title := fmt.Sprintf("Configure %s %s", m.serviceName, m.version)
	b.WriteString(styleHeader.Render(title))
	b.WriteString("\n\n")

	b.WriteString(styleBold.Render("Service Configuration"))
	b.WriteString("\n\n")

	for i, input := range m.inputs {
		if i == m.focus {
			b.WriteString(styleHighlight.Render(input.View()))
		} else {
			b.WriteString(input.View())
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleMuted.Render("  [↑/↓] Navigate   [Enter] Install   [Esc] Cancel"))

	return b.String()
}
