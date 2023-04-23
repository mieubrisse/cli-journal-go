package framework

import tea "github.com/charmbracelet/bubbletea"

// TODO shift these signatures to be by-value??? But then have to figure out generic interfaces
type FocusableComponent interface {
	Component

	Focus() tea.Cmd
	Blur() tea.Cmd
	Focused() bool
}
