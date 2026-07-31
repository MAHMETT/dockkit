package screens

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// TemplateEditorModel edits a template YAML.
type TemplateEditorModel struct {
	name     string
	editor   textarea.Model
	keys     screenKeys
	width    int
	height   int
}

// NewTemplateEditorModel creates a new template editor.
func NewTemplateEditorModel(name, content string) TemplateEditorModel {
	ta := textarea.New()
	ta.SetValue(content)
	ta.SetWidth(80)
	ta.SetHeight(20)
	ta.Focus()

	return TemplateEditorModel{
		name:   name,
		editor: ta,
		keys:   defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *TemplateEditorModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.editor.SetWidth(w - 4)
	m.editor.SetHeight(h - 8)
}

// Update handles messages.
func (m TemplateEditorModel) Update(msg tea.Msg) (TemplateEditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenTemplateManager}
			}
		case key.Matches(msg, m.keys.Enter):
			return m, func() tea.Msg {
				return messages.ToastMsg{Message: "Template saved!", Type: 0}
			}
		}
	}

	var cmd tea.Cmd
	m.editor, cmd = m.editor.Update(msg)
	return m, cmd
}

// View renders the template editor.
func (m TemplateEditorModel) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("Edit Template: " + m.name))
	b.WriteString("\n\n")

	b.WriteString(m.editor.View())

	b.WriteString("\n")
	b.WriteString(styleMuted.Render("  [Ctrl+S] Save   [Esc] Cancel"))

	return b.String()
}
