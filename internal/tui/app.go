package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
	"github.com/MAHMETT/dockkit/internal/tui/components"
)

// Model is the root TUI model.
type Model struct {
	// State
	screen     messages.ScreenID
	width      int
	height     int
	ready      bool

	// Components
	keys        KeyMap
	help        help.Model
	spinner     spinner.Model
	toast       components.Toast
	helpOverlay components.HelpOverlay
	loading     components.LoadingSpinner

	// Data
	err error
}

// NewModel creates a new root TUI model.
func NewModel() Model {
	keys := DefaultKeyMap()

	h := help.New()
	h.ShowAll = true

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		screen:  messages.ScreenDashboard,
		keys:    keys,
		help:    h,
		spinner: s,
		loading: components.NewLoadingSpinner("Loading..."),
	}
}

// Init initializes the TUI.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tea.KeyPressMsg:
		// Global keys
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.Help) {
			m.helpOverlay.Toggle()
			return m, nil
		}
		if key.Matches(msg, m.keys.Escape) {
			if m.helpOverlay.Visible {
				m.helpOverlay.Hide()
				return m, nil
			}
			if m.screen != messages.ScreenDashboard {
				m.screen = messages.ScreenDashboard
				return m, nil
			}
			return m, tea.Quit
		}

	case messages.NavigateToMsg:
		m.screen = msg.Screen
		m.helpOverlay.Hide()

	case messages.ToastMsg:
		m.toast = components.NewToast(msg.Message, components.ToastType(msg.Type))

	case messages.ErrorMsg:
		m.err = msg.Err
		m.toast = components.NewToast(msg.Message, components.ToastError)

	case messages.LoadingMsg:
		m.loading.SetMessage(msg.Message)
		if msg.Loading {
			m.loading.Show()
		} else {
			m.loading.Hide()
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the TUI.
func (m Model) View() tea.View {
	v := tea.NewView("")

	if !m.ready {
		v.SetContent("Initializing...")
		return v
	}

	var content string

	switch m.screen {
	case messages.ScreenDashboard:
		content = m.viewDashboard()
	default:
		content = "Screen not implemented"
	}

	// Layer help overlay on top
	if m.helpOverlay.Visible {
		content = m.helpOverlay.Render()
	}

	// Layer toast on top
	if m.toast.Visible {
		content += "\n" + m.toast.Render()
	}

	v.SetContent(content)
	return v
}

// viewDashboard renders the dashboard screen.
func (m Model) viewDashboard() string {
	header := Styles.Header.Render("🐳 dockkit v1.0.0")
	footer := Styles.Footer.Render("[?] Help [q] Quit")

	status := Styles.Muted.Render(fmt.Sprintf("Width: %d | Height: %d", m.width, m.height))

	return header + "\n\n" + status + "\n\n" + footer
}
