package components

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/lipgloss/v2"
)

var (
	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))
)

// LoadingSpinner shows a loading indicator.
type LoadingSpinner struct {
	spinner spinner.Model
	Message string
	Visible bool
}

// NewLoadingSpinner creates a new loading spinner.
func NewLoadingSpinner(message string) LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = loadingStyle

	return LoadingSpinner{
		spinner: s,
		Message: message,
		Visible: true,
	}
}

// Tick advances the spinner animation.
func (l *LoadingSpinner) Tick() {
	l.spinner.Tick()
}

// Render renders the loading spinner.
func (l LoadingSpinner) Render() string {
	if !l.Visible {
		return ""
	}
	return l.spinner.View() + " " + l.Message
}

// Hide hides the spinner.
func (l *LoadingSpinner) Hide() {
	l.Visible = false
}

// Show shows the spinner.
func (l *LoadingSpinner) Show() {
	l.Visible = true
}

// SetMessage updates the loading message.
func (l *LoadingSpinner) SetMessage(msg string) {
	l.Message = msg
}
