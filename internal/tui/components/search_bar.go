package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
)

var (
	searchBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3D3D5C")).
			Padding(0, 1)

	searchActive = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
)

// SearchBar is a text input for filtering/searching.
type SearchBar struct {
	input  textinput.Model
	Active bool
}

// NewSearchBar creates a new search bar.
func NewSearchBar(placeholder string) SearchBar {
	input := textinput.New()
	input.Placeholder = placeholder
	input.SetWidth(40)

	return SearchBar{
		input:  input,
		Active: false,
	}
}

// Focus focuses the search bar.
func (s *SearchBar) Focus() {
	s.input.Focus()
	s.Active = true
}

// Blur blurs the search bar.
func (s *SearchBar) Blur() {
	s.input.Blur()
	s.Active = false
}

// Value returns the current search value.
func (s SearchBar) Value() string {
	return s.input.Value()
}

// SetValue sets the search value.
func (s *SearchBar) SetValue(val string) {
	s.input.SetValue(val)
}

// Reset resets the search bar.
func (s *SearchBar) Reset() {
	s.input.Reset()
}

// Render renders the search bar.
func (s SearchBar) Render() string {
	style := searchBorder
	if s.Active {
		style = searchActive
	}
	return style.Render(s.input.View())
}

// Update updates the search bar with a tea.Msg.
func (s *SearchBar) Update(msg tea.Msg) {
	s.input.Update(msg)
}
