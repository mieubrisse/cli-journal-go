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

type implementation struct {
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
	impl := implementation{
		nameInput:     input,
		nameValidator: validator,
	}
	impl.recalculateInputColors()
	return &impl
}

func (impl implementation) Init() tea.Cmd {
	return nil
}

func (impl implementation) Update(msg tea.Msg) tea.Cmd {
	cmd := impl.nameInput.Update(msg)
	impl.recalculateInputColors()
	return cmd
}

func (impl implementation) View() string {
	// TODO some fancy nonsense to truncate strings that are too long for the form

	renderedTitle := lipgloss.NewStyle().
		Foreground(global_styles.White).
		Bold(true).
		Render(title)

	lines := lipgloss.JoinVertical(
		lipgloss.Center,
		renderedTitle,
		"",
		impl.nameInput.View(),
	)

	return lipgloss.NewStyle().
		Width(impl.width).
		Height(impl.height).
		Padding(verticalPadding, horizontalPadding, verticalPadding, horizontalPadding).
		Render(lines)
}

func (impl implementation) GetValue() string {
	return impl.nameInput.GetValue()
}

func (impl *implementation) Clear() {
	impl.nameInput.SetValue("")
	impl.recalculateInputColors()
}

func (impl implementation) GetNameValue() string {
	return impl.nameInput.GetValue()
}

func (impl implementation) SetNameValue(name string) {
	impl.nameInput.SetValue(name)
}

func (impl *implementation) Focus() tea.Cmd {
	impl.isFocused = true
	return impl.nameInput.Focus()
}

func (impl implementation) Blur() tea.Cmd {
	impl.isFocused = false
	return impl.nameInput.Blur()
}

func (impl implementation) Focused() bool {
	return impl.isFocused
}

func (impl *implementation) Resize(width int, height int) {
	impl.width = width
	impl.height = height

	inputHeight := 1
	inputWidth := width - 2*horizontalPadding
	impl.nameInput.Resize(inputWidth, inputHeight)
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetWidth() int {
	return impl.width
}

// ====================================================================================================
//                                   Private Helper Functions
// ====================================================================================================

func (impl *implementation) recalculateInputColors() {
	isValid := impl.nameValidator(impl.nameInput.GetValue())
	if isValid {
		impl.nameInput.SetForegroundColor(global_styles.White)
	} else {
		impl.nameInput.SetForegroundColor(global_styles.Red)
	}
}
