package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/text_input"
	"github.com/mieubrisse/cli-journal-go/global_styles"
)

const (
	horizontalPadding = 2
	verticalPadding   = 1
)

// TODO something about a border?

type Model struct {
	title string

	// TODO more fields
	input     text_input.Model
	validator func(text string) bool

	isFocused bool

	height int
	width  int
}

func New(
	title string,
	input text_input.Model,
	validator func(text string) bool,
) Model {
	return Model{
		title:     title,
		input:     input,
		validator: validator,
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

	isValid := model.validator(model.input.Value())

	model.input = model.input.SetTextStyle()
	if isValid

	return model, cmd
}

func (model Model) View() string {
	// TODO some fancy nonsense to truncate strings that are too long for the form

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		model.title,
		"",
		lipgloss.NewStyle().Foreground(global_styles.Red).Render(model.input.View()),
	)

	return lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		Padding(verticalPadding, horizontalPadding, verticalPadding, horizontalPadding).
		Render(content)
}

func (model Model) Focus() (Model, tea.Cmd) {
	model.isFocused = true

	var cmd tea.Cmd
	model.input, cmd = model.input.Focus()
	return model, cmd
}

func (model Model) Blur() Model {
	model.isFocused = false
	model.input.Blur()
	return model
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

func (model Model) SetValue(value string) Model {
	model.input = model.input.SetValue(value)
	return model
}
