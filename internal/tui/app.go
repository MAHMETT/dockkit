package tui

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"

	"github.com/MAHMETT/dockkit/internal/config"
	"github.com/MAHMETT/dockkit/internal/docker"
	"github.com/MAHMETT/dockkit/internal/templates"
	"github.com/MAHMETT/dockkit/internal/tui/messages"
	"github.com/MAHMETT/dockkit/internal/tui/components"
	"github.com/MAHMETT/dockkit/internal/tui/screens"
)

// Model is the root TUI model.
type Model struct {
	screen     messages.ScreenID
	width      int
	height     int
	ready      bool

	keys             KeyMap
	help             help.Model
	spinner          spinner.Model
	toast            components.Toast
	helpOverlay      components.HelpOverlay
	permissionDialog screens.PermissionDialogModel

	dashboard     screens.DashboardModel
	serviceDetail screens.ServiceDetailModel
	servicePicker screens.ServicePickerModel
	configWizard  screens.ConfigWizardModel
	logsViewer    screens.LogsViewerModel
	templateMgr   screens.TemplateManagerModel
	versionFetch  screens.VersionFetcherModel
	templateEdit  screens.TemplateEditorModel

	toastTimer int
	cfg        *config.Config
}

func NewModel() Model {
	keys := DefaultKeyMap()
	h := help.New()
	h.ShowAll = true
	s := spinner.New()
	s.Spinner = spinner.Dot

	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	return Model{
		screen:        messages.ScreenDashboard,
		keys:          keys,
		help:          h,
		spinner:       s,
		cfg:           cfg,
		dashboard:     screens.NewDashboardModel(),
		servicePicker: screens.NewServicePickerModel(defaultTemplates()),
		templateMgr:   screens.NewTemplateManagerModel(defaultTemplates()),
	}
}

func defaultTemplates() []screens.TemplateEntry {
	return []screens.TemplateEntry{
		{Name: "PostgreSQL", Icon: "🐘", Category: "database", Versions: []string{"15", "16", "17"}, DefaultPort: 5432, DefaultUser: "postgres", DefaultPassword: "postgres", DefaultDatabase: "postgres"},
		{Name: "MySQL", Icon: "🐬", Category: "database", Versions: []string{"8.0", "8.4", "9.0"}, DefaultPort: 3306, DefaultUser: "root", DefaultPassword: "mysql", DefaultDatabase: "mysql"},
		{Name: "MariaDB", Icon: "🐭", Category: "database", Versions: []string{"11"}, DefaultPort: 3306, DefaultUser: "root", DefaultPassword: "mariadb", DefaultDatabase: "mariadb"},
		{Name: "Redis", Icon: "⚡", Category: "cache", Versions: []string{"7"}, DefaultPort: 6379},
		{Name: "MongoDB", Icon: "🍃", Category: "database", Versions: []string{"7", "8"}, DefaultPort: 27017, DefaultUser: "admin", DefaultPassword: "mongo", DefaultDatabase: "admin"},
		{Name: "MinIO", Icon: "📦", Category: "storage", Versions: []string{"latest"}, DefaultPort: 9000, DefaultUser: "minioadmin", DefaultPassword: "minioadmin"},
		{Name: "Elasticsearch", Icon: "🔍", Category: "search", Versions: []string{"8"}, DefaultPort: 9200},
		{Name: "Memcached", Icon: "🔧", Category: "cache", Versions: []string{"1.6"}, DefaultPort: 11211},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.helpOverlay.SetTermHeight(msg.Height)
		m.dashboard.SetSize(msg.Width, msg.Height)
		m.servicePicker.SetSize(msg.Width, msg.Height)
		m.templateMgr.SetSize(msg.Width, msg.Height)
		m.permissionDialog.SetSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		// === HELP OVERLAY (highest priority) ===
		if m.helpOverlay.Visible {
			switch {
			case key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Escape):
				m.helpOverlay.Hide()
			case key.Matches(msg, m.keys.Up):
				m.helpOverlay.ScrollUp()
			case key.Matches(msg, m.keys.Down):
				m.helpOverlay.ScrollDown()
			}
			return m, nil
		}

		// === DISMISS ERROR TOAST ===
		if m.toast.Visible && m.toast.Type == components.ToastError {
			m.toast.Hide()
			return m, nil
		}

		// === GLOBAL KEYS ===
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.Help) {
			m.helpOverlay.Toggle()
			return m, nil
		}
		if key.Matches(msg, m.keys.Escape) {
			if m.screen != messages.ScreenDashboard {
				m.screen = messages.ScreenDashboard
				return m, nil
			}
			return m, tea.Quit
		}

	case messages.NavigateToMsg:
		m.screen = msg.Screen
		m.helpOverlay.Hide()

		switch msg.Screen {
		case messages.ScreenServiceDetail:
			if svc, ok := msg.Data.(screens.ServiceEntry); ok {
				m.serviceDetail = screens.NewServiceDetailModel(svc)
				m.serviceDetail.SetSize(m.width, m.height)
			}
		case messages.ScreenConfigWizard:
			if data, ok := msg.Data.(screens.ConfigWizardData); ok {
				m.configWizard = screens.NewConfigWizardModel(data)
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
		case messages.ScreenPermissionDialog:
			if data, ok := msg.Data.(screens.PermissionDialogData); ok {
				m.permissionDialog = screens.NewPermissionDialogModel(data.Service, data.Version)
				m.permissionDialog.SetSize(m.width, m.height)
			}
		}

	case messages.ConfigSaveMsg:
		cmd := m.handleConfigSave(msg)
		cmds = append(cmds, cmd)

	case messages.ConfigSavedMsg:
		cmd := m.startServiceCmd(msg.Service, msg.Version)
		cmds = append(cmds, cmd)

	case messages.ConfigErrorMsg:
		m.toast = components.NewToast(msg.Message, 1)
		m.toastTimer = -1 // error: no auto-dismiss

	case messages.PermissionDeniedMsg:
		m.screen = messages.ScreenPermissionDialog
		m.permissionDialog = screens.NewPermissionDialogModel(msg.Service, msg.Version)
		m.permissionDialog.SetSize(m.width, m.height)

	case messages.ToastMsg:
		m.toast = components.NewToast(msg.Message, msg.Type)
		switch msg.Type {
		case 0: // success
			m.toastTimer = 5
		case 1: // error
			m.toastTimer = -1 // never auto-dismiss
		case 2: // info
			m.toastTimer = 3
		}

	case messages.ErrorMsg:
		m.toast = components.NewToast(msg.Message, 1)
		m.toastTimer = -1

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

		// Auto-dismiss toast (only for positive timer values)
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
	case messages.ScreenPermissionDialog:
		var cmd tea.Cmd
		m.permissionDialog, cmd = m.permissionDialog.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleConfigSave(msg messages.ConfigSaveMsg) tea.Cmd {
	return func() tea.Msg {
		cfg := m.cfg
		if cfg == nil {
			cfg = config.DefaultConfig()
		}
		if cfg.Services == nil {
			cfg.Services = map[string]config.Service{}
		}

		svc, ok := cfg.Services[msg.ServiceName]
		if !ok {
			prefix := msg.ServiceName
			if len(prefix) > 3 {
				prefix = prefix[:3]
			}
			svc = config.Service{
				Prefix:   strings.ToUpper(prefix),
				Versions: map[string]config.ServiceVersion{},
			}
		}
		if svc.Versions == nil {
			svc.Versions = map[string]config.ServiceVersion{}
		}

		svc.Versions[msg.Version] = config.ServiceVersion{
			Enabled:       true,
			Port:          msg.Port,
			ContainerName: msg.ContainerName,
			Image:         msg.ServiceName + ":" + msg.Version,
			User:          msg.User,
			Password:      msg.Password,
			Database:      msg.Database,
		}
		cfg.Services[msg.ServiceName] = svc

		if err := config.Save(cfg); err != nil {
			return messages.ConfigErrorMsg{
				Message: fmt.Sprintf("Failed to save config: %v", err),
			}
		}

		serviceDir, err := config.ServiceDir(msg.ServiceName, msg.Version)
		if err != nil {
			return messages.ConfigErrorMsg{
				Message: fmt.Sprintf("Failed to create service directory: %v", err),
			}
		}
		if err := os.MkdirAll(serviceDir, 0700); err != nil {
			return messages.ConfigErrorMsg{
				Message: fmt.Sprintf("Failed to create service directory: %v", err),
			}
		}

		tmpl, err := templates.LoadBuiltin(msg.ServiceName)
		if err == nil {
			opts := templates.RenderOptions{
				ServiceName:   msg.ServiceName,
				Version:       msg.Version,
				Port:          msg.Port,
				User:          msg.User,
				Password:      msg.Password,
				Database:      msg.Database,
				ContainerName: msg.ContainerName,
				Timezone:      cfg.General.Timezone,
				Network:       cfg.General.DefaultNetwork,
			}
			composeYAML, err := templates.RenderToString(tmpl, opts)
			if err == nil {
				composePath := serviceDir + "/docker-compose.yml"
				os.WriteFile(composePath, []byte(composeYAML), 0644)
			}
		}

		m.cfg = cfg

		return messages.ConfigSavedMsg{
			Service: msg.ServiceName,
			Version: msg.Version,
			Message: fmt.Sprintf("%s %s configured!", msg.ServiceName, msg.Version),
		}
	}
}

func (m Model) startServiceCmd(service, version string) tea.Cmd {
	return func() tea.Msg {
		dir, err := config.ServiceDir(service, version)
		if err != nil {
			return messages.ConfigErrorMsg{
				Message: fmt.Sprintf("Cannot find service directory: %v", err),
			}
		}

		if _, statErr := os.Stat(dir + "/docker-compose.yml"); os.IsNotExist(statErr) {
			return messages.ConfigErrorMsg{
				Message: "docker-compose.yml not found. Run setup again.",
			}
		}

		if err := docker.ComposeUp(context.Background(), dir); err != nil {
			// Check for permission denied
			if strings.Contains(err.Error(), "permission denied") {
				return messages.PermissionDeniedMsg{
					Service: service,
					Version: version,
					Error:   err.Error(),
				}
			}
			return messages.ConfigErrorMsg{
				Message: fmt.Sprintf("Failed to start %s: %v", service, err),
			}
		}

		return messages.ToastMsg{
			Message: fmt.Sprintf("%s %s started successfully!", service, version),
			Type:    0,
		}
	}
}

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
	case messages.ScreenPermissionDialog:
		content = m.permissionDialog.View()
	default:
		content = "Screen not implemented"
	}

	if m.helpOverlay.Visible {
		content = m.helpOverlay.Render()
	}

	if m.toast.Visible {
		content += "\n" + m.toast.Render()
	}

	v.SetContent(content)
	return v
}
