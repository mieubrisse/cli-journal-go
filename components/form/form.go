package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/text_input"
)

const (
	horizontalPadding = 2
	verticalPadding   = 1
)

// TODO something about a border?

type Model struct {
	title string

	// TODO more fields
	input text_input.Model

	isFocused bool

	height int
	width  int
}

func New(title string, input text_input.Model) Model {
	return Model{
		title: title,
		input: input,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !model.isFocused {
		return model, nil
	}

	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)
	return model, cmd
}

func (model Model) View() string {
	// TODO some fancy nonsense to truncate strings that are too long for the form

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		model.title,
		"",
		model.input.View(),
	)

	return lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		Padding(verticalPadding, horizontalPadding, verticalPadding, horizontalPadding).
		Render(content)
}

func (model *Model) Focus() tea.Cmd {
	model.isFocused = true
	return model.input.Focus()
}

func (model *Model) Blur() {
	model.input.Blur()
	model.isFocused = false
}

func (model Model) Focused() bool {
	return model.isFocused
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height

	inputHeight := 1
	inputWidth := width - 2*horizontalPadding
	model.input = model.input.Resize(inputWidth, inputHeight)

	return model
}

func (model Model) GetHeight() int {
	return model.height
}

func (model Model) GetWidth() int {
	return model.width
}

func (model Model) GetValue() string {
	return model.input.Value()
}

func (model *Model) SetValue(value string) {
	model.input.SetValue(value)
}
