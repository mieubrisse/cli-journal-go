package text_filter_input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/global_styles"
)

type Model struct {
	input textinput.Model

	width  int
	height int
}

func New(promptText string) Model {
	input := textinput.New()
	input.Prompt = promptText
	return Model{
		input: input,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	baseStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height)
	if m.input.Focused() {
		baseStyle = baseStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
	}

	// It'd be really nice if textinput.Model allowed us to control the style based on focus, but alas it does not
	return baseStyle.Render(m.input.View())
}

func (m *Model) Focus() tea.Cmd {
	return m.input.Focus()
}

func (m *Model) Blur() {
	m.input.Blur()
}

func (m Model) Focused() bool {
	return m.input.Focused()
}

func (m *Model) SetValue(newValue string) {
	m.input.SetValue(newValue)
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m Model) Resize(width int, height int) Model {
	m.width = width
	m.input.Width = width
	return m
}

func (m Model) GetHeight() int {
	return m.height
}

func (m Model) GetWidth() int {
	return m.width
}