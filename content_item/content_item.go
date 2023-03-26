package content_item

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	name string
	tags []string

	// Means that a checkbox will be displayed next to the item
	isSelected bool

	// Means that the item will be displayed
	isHighlighted bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	maybeCheckmark := " "
	if m.isSelected {
		maybeCheckmark = "✔️"
	}

	style := lipgloss.NewStyle()
	if m.isHighlighted {
		// TODO make color a constant
		style = style.Bold(true).Foreground(lipgloss.Color("#4e9a06"))
	}

	value := fmt.Sprintf("%s  %s", maybeCheckmark, m.name)
	return style.Render(value)
}

func (m *Model) GetName() string {
	return m.name
}

func (m *Model) IsSelected() bool {
	return m.isSelected
}

func (m *Model) SetSelected(newValue bool) {
	m.isSelected = newValue
}

func (m *Model) IsHighlighted() bool {
	return m.isHighlighted
}

func (m *Model) SetHighlighted(newValue bool) {
	m.isHighlighted = newValue
}
