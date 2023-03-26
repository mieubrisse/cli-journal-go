package components

import tea "github.com/charmbracelet/bubbletea"

type FocusableComponent interface {
	Focus() tea.Cmd
	Blur()
	Focused() bool
}
