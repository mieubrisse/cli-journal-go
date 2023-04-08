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
	model := Model{
		title:     title,
		input:     input,
		validator: validator,
	}
	return model.recalculateInputColors()
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

	model = model.recalculateInputColors()

	return model, cmd
}

func (model Model) View() string {
	// TODO some fancy nonsense to truncate strings that are too long for the form

	renderedTitle := lipgloss.NewStyle().Foreground(global_styles.White).Bold(true).Render(model.title)

	lines := lipgloss.JoinVertical(
		lipgloss.Center,
		renderedTitle,
		"",
		model.input.View(),
	)

	return lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		Padding(verticalPadding, horizontalPadding, verticalPadding, horizontalPadding).
		Render(lines)
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

func (model Model) Clear() Model {
	model.input = model.input.SetValue("")
	model = model.recalculateInputColors()
	return model
}

func (model Model) recalculateInputColors() Model {
	isValid := model.validator(model.input.Value())
	if isValid {
		model.input = model.input.SetForegroundColor(global_styles.White)
	} else {
		model.input = model.input.SetForegroundColor(global_styles.Red)
	}
	return model
}
