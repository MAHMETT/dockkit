package components

import (
	"charm.land/lipgloss/v2"
)

var (
	confirmBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFE66D")).
			Padding(1, 2)

	confirmTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D")).
			Bold(true)

	confirmYes = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			Padding(0, 2)

	confirmNo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Padding(0, 2)
)

// ConfirmDialog asks for user confirmation.
type ConfirmDialog struct {
	Message  string
	Selected int // 0 = Yes, 1 = No
	Visible  bool
	OnYes    func()
	OnNo     func()
}

// NewConfirmDialog creates a new confirmation dialog.
func NewConfirmDialog(message string, onYes, onNo func()) ConfirmDialog {
	return ConfirmDialog{
		Message:  message,
		Selected: 1, // Default to No
		Visible:  true,
		OnYes:    onYes,
		OnNo:     onNo,
	}
}

// Render renders the confirmation dialog.
func (d ConfirmDialog) Render() string {
	if !d.Visible {
		return ""
	}

	title := confirmTitle.Render("Confirm")

	yes := confirmYes.Render("[Yes]")
	no := confirmNo.Render("[No]")

	if d.Selected == 0 {
		yes = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			Background(lipgloss.Color("#3D3D5C")).
			Padding(0, 2).
			Render("[Yes]")
	} else {
		no = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Background(lipgloss.Color("#3D3D5C")).
			Padding(0, 2).
			Render("[No]")
	}

	content := title + "\n\n" + d.Message + "\n\n" + yes + "  " + no

	return confirmBorder.Render(content)
}

// Left moves selection to the left.
func (d *ConfirmDialog) Left() {
	d.Selected = 0
}

// Right moves selection to the right.
func (d *ConfirmDialog) Right() {
	d.Selected = 1
}

// Confirm confirms the current selection and hides the dialog.
func (d *ConfirmDialog) Confirm() func() {
	d.Visible = false
	if d.Selected == 0 {
		return d.OnYes
	}
	return d.OnNo
}
