package text_input

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

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)
	return model, cmd
}

func (model Model) View() string {
	baseStyle := lipgloss.NewStyle().
		Width(model.width).
		Height(model.height)
	if model.input.Focused() {
		baseStyle = baseStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
	}

	// It'd be really nice if textinput.Model allowed us to control the style based on focus, but alas it does not
	return baseStyle.Render(model.input.View())
}

func (model Model) Focus() (Model, tea.Cmd) {
	return model, model.input.Focus()
}

func (model Model) Blur() Model {
	model.input.Blur()
	return model
}

func (model Model) Focused() bool {
	return model.input.Focused()
}

func (model Model) SetValue(newValue string) Model {
	model.input.SetValue(newValue)
	return model
}

func (model Model) Value() string {
	return model.input.Value()
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.input.Width = width
	model.height = height
	return model
}

func (model Model) GetHeight() int {
	return model.height
}

func (model Model) GetWidth() int {
	return model.width
}
