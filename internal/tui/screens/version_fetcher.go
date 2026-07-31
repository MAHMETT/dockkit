package screens

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"

	"github.com/MAHMETT/dockkit/internal/tui/messages"
)

// TagEntry represents a Docker image tag.
type TagEntry struct {
	Name string
}

// VersionFetcherModel fetches and displays Docker Hub tags.
type VersionFetcherModel struct {
	image    string
	tags     []TagEntry
	cursor   int
	loading  bool
	keys     screenKeys
	width    int
	height   int
}

// NewVersionFetcherModel creates a new version fetcher.
func NewVersionFetcherModel(image string) VersionFetcherModel {
	return VersionFetcherModel{
		image:   image,
		tags:    []TagEntry{},
		loading: true,
		keys:    defaultScreenKeys,
	}
}

// SetSize sets the terminal dimensions.
func (m *VersionFetcherModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Update handles messages.
func (m VersionFetcherModel) Update(msg tea.Msg) (VersionFetcherModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Escape):
			return m, func() tea.Msg {
				return messages.NavigateToMsg{Screen: messages.ScreenServicePicker}
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.tags)-1 {
				m.cursor++
			}
		}
	case tagsLoadedMsg:
		m.tags = msg.tags
		m.loading = false
	}
	return m, nil
}

// View renders the version fetcher.
func (m VersionFetcherModel) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render(fmt.Sprintf("Versions: %s", m.image)))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(styleMuted.Render("  Fetching tags from Docker Hub..."))
	} else if len(m.tags) == 0 {
		b.WriteString(styleMuted.Render("  No tags found."))
	} else {
		for i, tag := range m.tags {
			selected := i == m.cursor
			row := fmt.Sprintf("  %s", tag.Name)
			if selected {
				b.WriteString(styleHighlight.Render(row))
			} else {
				b.WriteString(row)
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(styleMuted.Render("  [↑/↓] Navigate   [Enter] Select   [Esc] Back"))
	return b.String()
}

type tagsLoadedMsg struct {
	tags []TagEntry
}
