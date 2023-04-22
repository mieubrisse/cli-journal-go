package components

import tea "github.com/charmbracelet/bubbletea"

// TODO shift these signatures to be by-value??? But then have to figure out generic interfaces
type FocusableComponent interface {
	Focus() tea.Cmd
	Blur() tea.Cmd
	Focused() bool
}
