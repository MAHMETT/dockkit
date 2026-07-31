package screens

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// ConfigWizardData holds data passed to the config wizard.
type ConfigWizardData struct {
	Template        TemplateEntry
	SelectedVersion string
}

// ConfigWizardModel handles service configuration.
type ConfigWizardModel struct {
	serviceName string
	version     string
	defaults    ConfigWizardData
	inputs      []textinput.Model
	focus       int
	keys        screenKeys
	width       int
	height      int
}

// NewConfigWizardModel creates a new config wizard.
func NewConfigWizardModel(data ConfigWizardData) ConfigWizardModel {
	tmpl := data.Template

	inputs := make([]textinput.Model, 4)
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = 64
	}

	// Determine port — auto-suggest if default is taken
	port := tmpl.DefaultPort
	if isPortOccupied(port) {
		port = suggestNextPort(port)
	}

	inputs[0].Prompt = "Port:          "
	inputs[0].SetValue(strconv.Itoa(port))

	inputs[1].Prompt = "Username:      "
	inputs[1].SetValue(tmpl.DefaultUser)

	inputs[2].Prompt = "Password:      "
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].EchoCharacter = '•'
	inputs[2].SetValue(tmpl.DefaultPassword)

	inputs[3].Prompt = "Database:      "
	inputs[3].SetValue(tmpl.DefaultDatabase)

	// Focus first input
	inputs[0].Focus()

	return ConfigWizardModel{
		serviceName: tmpl.Name,
		version:     data.SelectedVersion,
		defaults:    data,
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
				return m, nil
			}
		case key.Matches(msg, m.keys.Down):
			if m.focus < len(m.inputs)-1 {
				m.inputs[m.focus].Blur()
				m.focus++
				m.inputs[m.focus].Focus()
				return m, nil
			}
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenServicePicker}
			}
		case key.Matches(msg, m.keys.Enter):
			return m.saveConfig()
		}
	}

	// Pass non-navigation keys to focused input
	if m.focus < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}

	return m, nil
}

// saveConfig validates inputs and sends ConfigSaveMsg.
func (m ConfigWizardModel) saveConfig() (ConfigWizardModel, tea.Cmd) {
	// Parse port
	portStr := m.inputs[0].Value()
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1024 || port > 65535 {
		return m, func() tea.Msg {
			return messages.ConfigErrorMsg{
				Message: fmt.Sprintf("Invalid port: %s (must be 1024-65535)", portStr),
			}
		}
	}

	// Check port availability
	if isPortOccupied(port) {
		suggested := suggestNextPort(port)
		return m, func() tea.Msg {
			msg := fmt.Sprintf("Port %d is already in use.", port)
			if suggested > 0 {
				msg += fmt.Sprintf(" Suggested: %d", suggested)
			}
			return messages.ConfigErrorMsg{Message: msg}
		}
	}

	user := m.inputs[1].Value()
	password := m.inputs[2].Value()
	database := m.inputs[3].Value()

	if user == "" {
		return m, func() tea.Msg {
			return messages.ConfigErrorMsg{Message: "Username is required"}
		}
	}

	// Generate container name
	containerName := fmt.Sprintf("dockkit-%s-%s",
		strings.ToLower(m.serviceName), m.version)

	return m, func() tea.Msg {
		return messages.ConfigSaveMsg{
			ServiceName:   strings.ToLower(m.serviceName),
			Version:       m.version,
			Port:          port,
			User:          user,
			Password:      password,
			Database:      database,
			ContainerName: containerName,
		}
	}
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

// isPortOccupied checks if a port is in use on localhost.
func isPortOccupied(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 100*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// suggestNextPort finds the next available port starting from port+1.
func suggestNextPort(port int) int {
	for offset := 1; offset <= 100; offset++ {
		candidate := port + offset
		if candidate > 65535 {
			return 0
		}
		if candidate < 1024 {
			continue
		}
		if !isPortOccupied(candidate) {
			return candidate
		}
	}
	return 0
}
