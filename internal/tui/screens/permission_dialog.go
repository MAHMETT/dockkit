package screens

import (
	"fmt"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// SetupInfo holds OS-specific setup instructions.
type SetupInfo struct {
	OS    string
	Steps []string
	Link  string
}

var (
	permHelpMuted = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C757D"))
)

// PermissionDialogData holds data passed to the permission dialog.
type PermissionDialogData struct {
	Service string
	Version string
}

// PermissionDialogModel handles Docker permission denied.
type PermissionDialogModel struct {
	service      string
	version      string
	option       int // 0=Run with sudo, 1=Setup without sudo, 2=Cancel
	showSudo     bool
	sudoInput    textinput.Model
	sudoFocused  bool
	keys         screenKeys
	width        int
	height       int
}

// NewPermissionDialogModel creates a new permission dialog.
func NewPermissionDialogModel(service, version string) PermissionDialogModel {
	sudoInput := textinput.New()
	sudoInput.Placeholder = "Enter sudo password"
	sudoInput.EchoMode = textinput.EchoPassword
	sudoInput.EchoCharacter = '•'
	sudoInput.SetWidth(40)

	return PermissionDialogModel{
		service:   service,
		version:   version,
		option:    1, // Default to "Setup without sudo"
		sudoInput: sudoInput,
		keys:      defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *PermissionDialogModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Update handles messages.
func (m PermissionDialogModel) Update(msg tea.Msg) (PermissionDialogModel, tea.Cmd) {
	// If showing sudo input, handle text input
	if m.showSudo {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, m.keys.Escape):
				m.showSudo = false
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				password := m.sudoInput.Value()
				if password == "" {
					return m, nil
				}
				// Execute with sudo
				return m, m.execWithSudo(password)
			}
		}
		var cmd tea.Cmd
		m.sudoInput, cmd = m.sudoInput.Update(msg)
		return m, cmd
	}

	// Main dialog navigation
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenDashboard}
			}
		case key.Matches(msg, m.keys.Up):
			if m.option > 0 {
				m.option--
			}
		case key.Matches(msg, m.keys.Down):
			if m.option < 2 {
				m.option++
			}
		case key.Matches(msg, m.keys.Enter):
			return m.handleOption()
		}
	}
	return m, nil
}

// handleOption processes the selected option.
func (m PermissionDialogModel) handleOption() (PermissionDialogModel, tea.Cmd) {
	switch m.option {
	case 0: // Run with sudo
		m.showSudo = true
		m.sudoInput.Focus()
		return m, nil
	case 1: // Setup without sudo
		// Already showing instructions, just stay on this option
		return m, nil
	case 2: // Cancel
		return m, func() tea.Msg {
			return messages.NavigateToMsg{Screen: messages.ScreenDashboard}
		}
	}
	return m, nil
}

// execWithSudo runs docker compose with sudo.
func (m PermissionDialogModel) execWithSudo(password string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement sudo execution
		// For now, show instructions
		return messages.ToastMsg{
			Message: "Sudo execution not yet implemented. Use terminal to run:\nsudo docker compose up -d",
			Type:    1, // error type — stays until dismissed
		}
	}
}

// getSetupInfo returns OS-specific setup instructions.
func getSetupInfo() SetupInfo {
	switch runtime.GOOS {
	case "linux":
		return SetupInfo{
			OS: "Linux",
			Steps: []string{
				"1. Add your user to the docker group:",
				"   sudo usermod -aG docker $USER",
				"",
				"2. Log out and log back in",
				"   (or run: newgrp docker)",
				"",
				"3. Verify: docker ps",
			},
			Link: "https://docs.docker.com/engine/install/linux-postinstall/",
		}
	case "darwin":
		return SetupInfo{
			OS: "macOS",
			Steps: []string{
				"1. Open Docker Desktop app",
				"2. Docker Desktop handles permissions",
				"   automatically. No sudo required.",
			},
			Link: "https://docs.docker.com/desktop/install/mac-install/",
		}
	case "windows":
		return SetupInfo{
			OS: "Windows (WSL2)",
			Steps: []string{
				"1. Ensure your user is in the",
				"   docker-users group",
				"2. Or run from WSL2 terminal",
				"3. Restart terminal after change",
			},
			Link: "https://docs.docker.com/desktop/install/windows-install/",
		}
	default:
		return SetupInfo{
			OS: "Unknown",
			Steps: []string{
				"Please ensure your user has Docker access.",
				"Check your Docker installation docs.",
			},
			Link: "https://docs.docker.com/engine/install/",
		}
	}
}

// View renders the permission dialog.
func (m PermissionDialogModel) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("Docker Permission Denied"))
	b.WriteString("\n\n")

	b.WriteString("dockkit needs Docker access but got permission denied.")
	b.WriteString("\n\n")

	b.WriteString(styleBold.Render("Choose an option:"))
	b.WriteString("\n\n")

	options := []string{
		"Run with sudo",
		"Setup Docker without sudo (recommended)",
		"Cancel",
	}

	for i, opt := range options {
		prefix := "  "
		if i == m.option {
			prefix = styleHighlight.Render("> ")
		}
		b.WriteString(fmt.Sprintf("%s%d. %s\n", prefix, i+1, opt))
	}

	// Show sudo input if selected
	if m.showSudo {
		b.WriteString("\n")
		b.WriteString(styleBold.Render("Enter sudo password:"))
		b.WriteString("\n")
		b.WriteString(m.sudoInput.View())
		b.WriteString("\n")
		b.WriteString(permHelpMuted.Render("Press Enter to confirm, Esc to cancel"))
	}

	// Show setup instructions
	if m.option == 1 && !m.showSudo {
		info := getSetupInfo()
		b.WriteString("\n")
		b.WriteString(styleBold.Render(fmt.Sprintf("Setup Instructions (%s):", info.OS)))
		b.WriteString("\n")
		for _, step := range info.Steps {
			b.WriteString("  " + step + "\n")
		}
		b.WriteString("\n")
		b.WriteString(permHelpMuted.Render("Reference:"))
		b.WriteString("\n")
		b.WriteString(permHelpMuted.Render("  "+info.Link))
	}

	b.WriteString("\n")
	b.WriteString(permHelpMuted.Render("  [↑/↓] Select   [Enter] Confirm   [Esc] Cancel"))

	return b.String()
}
