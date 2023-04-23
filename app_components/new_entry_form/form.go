package new_entry_form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/text_input"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"regexp"
)

const (
	horizontalPadding = 2
	verticalPadding   = 1

	title = "Create Content"
)

var acceptableNameRegex = regexp.MustCompile("^[a-zA-Z0-9.-]+$")

// TODO something about a border?

type Model struct {
	// TODO more fields
	nameInput     text_input.Model
	nameValidator func(text string) bool

	isFocused bool

	height int
	width  int
}

func New() Component {
	input := text_input.New("Name: ")
	validator := func(text string) bool {
		return acceptableNameRegex.MatchString(text)
	}
	impl := Model{
		nameInput:     input,
		nameValidator: validator,
	}
	impl.recalculateInputColors()
	return &impl
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) tea.Cmd {
	cmd := model.nameInput.Update(msg)
	model.recalculateInputColors()
	return cmd
}

func (model Model) View() string {
	// TODO some fancy nonsense to truncate strings that are too long for the form

	renderedTitle := lipgloss.NewStyle().
		Foreground(global_styles.White).
		Bold(true).
		Render(title)

	lines := lipgloss.JoinVertical(
		lipgloss.Center,
		renderedTitle,
		"",
		model.nameInput.View(),
	)

	return lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		Padding(verticalPadding, horizontalPadding, verticalPadding, horizontalPadding).
		Render(lines)
}

func (model Model) GetValue() string {
	return model.nameInput.GetValue()
}

func (model *Model) Clear() {
	model.nameInput.SetValue("")
	model.recalculateInputColors()
}

func (model Model) GetNameValue() string {
	return model.nameInput.GetValue()
}

func (model Model) SetNameValue(name string) {
	model.nameInput.SetValue(name)
}

func (model *Model) Focus() tea.Cmd {
	model.isFocused = true
	return model.nameInput.Focus()
}

func (model Model) Blur() tea.Cmd {
	model.isFocused = false
	return model.nameInput.Blur()
}

func (model Model) Focused() bool {
	return model.isFocused
}

func (model *Model) Resize(width int, height int) {
	model.width = width
	model.height = height

	inputHeight := 1
	inputWidth := width - 2*horizontalPadding
	model.nameInput.Resize(inputWidth, inputHeight)
}

func (model Model) GetHeight() int {
	return model.height
}

func (model Model) GetWidth() int {
	return model.width
}

// ====================================================================================================
//                                   Private Helper Functions
// ====================================================================================================

func (model *Model) recalculateInputColors() {
	isValid := model.nameValidator(model.nameInput.GetValue())
	if isValid {
		model.nameInput.SetForegroundColor(global_styles.White)
	} else {
		model.nameInput.SetForegroundColor(global_styles.Red)
	}
}
