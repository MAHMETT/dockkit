package messages

// ScreenID identifies a TUI screen.
type ScreenID int

const (
	ScreenDashboard ScreenID = iota
	ScreenServiceDetail
	ScreenServicePicker
	ScreenConfigWizard
	ScreenVersionFetcher
	ScreenLogsViewer
	ScreenTemplateManager
	ScreenTemplateEditor
	ScreenPermissionDialog
	ScreenError
)

// Navigation messages
type NavigateToMsg struct {
	Screen ScreenID
	Data   interface{}
}

// Docker messages
type StatusUpdateMsg struct {
	Service string
	Running bool
	Health  string
}

type ContainerStartedMsg struct {
	Service string
}

type ContainerStoppedMsg struct {
	Service string
}

type ContainerRestartedMsg struct {
	Service string
}

type LogsMsg struct {
	Service string
	Lines   []string
}

// Config messages
type ConfigSaveMsg struct {
	ServiceName   string
	Version       string
	Port          int
	User          string
	Password      string
	Database      string
	ContainerName string
}

type ConfigSavedMsg struct {
	Service string
	Version string
	Message string
}

type ConfigErrorMsg struct {
	Err     error
	Message string
}

// Permission messages
type PermissionDeniedMsg struct {
	Service string
	Version string
	Error   string
}

type SudoExecMsg struct {
	Service    string
	Version    string
	Password   string
	ServiceDir string
}

// System messages
// ToastMsg type: 0=success, 1=error, 2=info
type ToastMsg struct {
	Message string
	Type    int
}

type ErrorMsg struct {
	Err     error
	Message string
}

type LoadingMsg struct {
	Active  bool
	Message string
}

type TickMsg struct{}
