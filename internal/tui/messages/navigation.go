package messages

import "time"

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

type BackMsg struct{}

// Docker messages
type StatusUpdateMsg struct {
	// Will be populated with container states
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
type ConfigLoadedMsg struct {
	// Config loaded
}

type ConfigSavedMsg struct {
	Service string
}

type ConfigErrorMsg struct {
	Err error
}

// System messages
type TickMsg struct {
	Time time.Time
}

type ToastMsg struct {
	Message string
	Type    ToastType // success, error, info
}

type ToastType int

const (
	ToastSuccess ToastType = iota
	ToastError
	ToastInfo
)

type ErrorMsg struct {
	Err     error
	Message string
}

type LoadingMsg struct {
	Loading bool
	Message string
}

type QuitMsg struct{}
