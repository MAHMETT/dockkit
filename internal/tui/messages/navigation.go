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
	ScreenError
)

// Navigation messages
type NavigateToMsg struct {
	Screen ScreenID
	Data   interface{} // optional data to pass to screen
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
type ConfigSavedMsg struct {
	Service string
}

type ConfigErrorMsg struct {
	Err error
}

// System messages
// ToastMsg uses int for Type to avoid circular dependency with components.
// Convert to components.ToastType at the call site.
type ToastMsg struct {
	Message string
	Type    int // 0=success, 1=error, 2=info (matches components.ToastType)
}

type ErrorMsg struct {
	Err     error
	Message string
}

type LoadingMsg struct {
	Active  bool // true to show, false to hide
	Message string
}

type TickMsg struct {
	// Used for periodic updates
}
