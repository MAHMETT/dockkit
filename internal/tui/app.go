package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
	"github.com/MAHMETT/dockkit/internal/tui/components"
	"github.com/MAHMETT/dockkit/internal/tui/screens"
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

	// Screens
	dashboard     screens.DashboardModel
	serviceDetail screens.ServiceDetailModel
	servicePicker screens.ServicePickerModel
	configWizard  screens.ConfigWizardModel
	logsViewer    screens.LogsViewerModel
	templateMgr   screens.TemplateManagerModel
	versionFetch  screens.VersionFetcherModel
	templateEdit  screens.TemplateEditorModel

	// State
	toastTimer int
}

// NewModel creates a new root TUI model.
func NewModel() Model {
	keys := DefaultKeyMap()

	h := help.New()
	h.ShowAll = true

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		screen:      messages.ScreenDashboard,
		keys:        keys,
		help:        h,
		spinner:     s,
		dashboard:   screens.NewDashboardModel(),
		servicePicker: screens.NewServicePickerModel(defaultTemplates()),
		templateMgr: screens.NewTemplateManagerModel(defaultTemplates()),
	}
}

// defaultTemplates returns the built-in service templates.
func defaultTemplates() []screens.TemplateEntry {
	return []screens.TemplateEntry{
		{Name: "PostgreSQL", Icon: "🐘", Category: "database", Versions: []string{"15", "16", "17"}},
		{Name: "MySQL", Icon: "🐬", Category: "database", Versions: []string{"8.0", "8.4", "9.0"}},
		{Name: "MariaDB", Icon: "🐭", Category: "database", Versions: []string{"11"}},
		{Name: "Redis", Icon: "⚡", Category: "cache", Versions: []string{"7"}},
		{Name: "MongoDB", Icon: "🍃", Category: "database", Versions: []string{"7", "8"}},
		{Name: "MinIO", Icon: "📦", Category: "storage", Versions: []string{"latest"}},
		{Name: "Elasticsearch", Icon: "🔍", Category: "search", Versions: []string{"8"}},
		{Name: "Memcached", Icon: "🔧", Category: "cache", Versions: []string{"1.6"}},
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

		// Set sizes on all screens
		m.dashboard.SetSize(msg.Width, msg.Height)
		m.servicePicker.SetSize(msg.Width, msg.Height)
		m.templateMgr.SetSize(msg.Width, msg.Height)

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

		// Initialize screen-specific data
		switch msg.Screen {
		case messages.ScreenServiceDetail:
			if svc, ok := msg.Data.(screens.ServiceEntry); ok {
				m.serviceDetail = screens.NewServiceDetailModel(svc)
				m.serviceDetail.SetSize(m.width, m.height)
			}
		case messages.ScreenConfigWizard:
			if tmpl, ok := msg.Data.(screens.TemplateEntry); ok {
				m.configWizard = screens.NewConfigWizardModel(tmpl.Name, "latest")
				m.configWizard.SetSize(m.width, m.height)
			} else if svc, ok := msg.Data.(screens.ServiceEntry); ok {
				m.configWizard = screens.NewConfigWizardModel(svc.Name, svc.Version)
				m.configWizard.SetSize(m.width, m.height)
			}
		case messages.ScreenLogsViewer:
			if svc, ok := msg.Data.(screens.ServiceEntry); ok {
				m.logsViewer = screens.NewLogsViewerModel(svc)
				m.logsViewer.SetSize(m.width, m.height)
			}
		case messages.ScreenVersionFetcher:
			if image, ok := msg.Data.(string); ok {
				m.versionFetch = screens.NewVersionFetcherModel(image)
				m.versionFetch.SetSize(m.width, m.height)
			}
		case messages.ScreenTemplateEditor:
			if name, ok := msg.Data.(string); ok {
				m.templateEdit = screens.NewTemplateEditorModel(name, "# Template: "+name+"\n")
				m.templateEdit.SetSize(m.width, m.height)
			}
		}

	case messages.ToastMsg:
		m.toast = components.NewToast(msg.Message, msg.Type)
		m.toastTimer = 5

	case messages.ErrorMsg:
		m.toast = components.NewToast(msg.Message, 1)
		m.toastTimer = 8

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

		if m.toast.Visible && m.toastTimer > 0 {
			m.toastTimer--
			if m.toastTimer == 0 {
				m.toast.Hide()
			}
		}
	}

	// Route to screen update
	switch m.screen {
	case messages.ScreenDashboard:
		var cmd tea.Cmd
		m.dashboard, cmd = m.dashboard.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ScreenServiceDetail:
		var cmd tea.Cmd
		m.serviceDetail, cmd = m.serviceDetail.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ScreenServicePicker:
		var cmd tea.Cmd
		m.servicePicker, cmd = m.servicePicker.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ScreenConfigWizard:
		var cmd tea.Cmd
		m.configWizard, cmd = m.configWizard.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ScreenLogsViewer:
		var cmd tea.Cmd
		m.logsViewer, cmd = m.logsViewer.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ScreenTemplateManager:
		var cmd tea.Cmd
		m.templateMgr, cmd = m.templateMgr.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ScreenVersionFetcher:
		var cmd tea.Cmd
		m.versionFetch, cmd = m.versionFetch.Update(msg)
		cmds = append(cmds, cmd)
	case messages.ScreenTemplateEditor:
		var cmd tea.Cmd
		m.templateEdit, cmd = m.templateEdit.Update(msg)
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
		content = m.dashboard.View()
	case messages.ScreenServiceDetail:
		content = m.serviceDetail.View()
	case messages.ScreenServicePicker:
		content = m.servicePicker.View()
	case messages.ScreenConfigWizard:
		content = m.configWizard.View()
	case messages.ScreenLogsViewer:
		content = m.logsViewer.View()
	case messages.ScreenTemplateManager:
		content = m.templateMgr.View()
	case messages.ScreenVersionFetcher:
		content = m.versionFetch.View()
	case messages.ScreenTemplateEditor:
		content = m.templateEdit.View()
	default:
		content = "Screen not implemented"
	}

	// Layer help overlay
	if m.helpOverlay.Visible {
		content = m.helpOverlay.Render()
	}

	// Layer toast
	if m.toast.Visible {
		content += "\n" + m.toast.Render()
	}

	v.SetContent(content)
	return v
}
