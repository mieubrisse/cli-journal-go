package filter_pane

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"github.com/mieubrisse/vim-bubble/vim"
	"strings"
)

const (
	tagFilterLineLeader = "#"
)

type Model struct {
	input vim.Model

	isFocused bool
	width     int
	height    int
}

// TODO allow initializing with a state
func New() Model {
	return Model{
		input:     vim.New(),
		isFocused: false,
		width:     0,
		height:    0,
	}
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)
	return model, cmd
}

func (model Model) View() string {
	// TODO handle the Really Small case
	lines := []string{
		// TODO Replace with a better border
		"=========",
		model.input.View(),
	}
	return strings.Join(lines, "\n")
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height

	// Leave room for border
	model.input.Resize(width, helpers.GetMaxInt(height-1, 0))

	return model
}

func (model Model) GetHeight() int {
	return model.height
}

func (model Model) GetWidth() int {
	return model.width
}

func (model *Model) Focus() tea.Cmd {
	model.isFocused = true
	model.input.Focus()
	return nil
}

func (model *Model) Blur() tea.Cmd {
	model.isFocused = false
	model.input.Blur()
	return nil
}

func (model Model) Focused() bool {
	return model.isFocused
}

func (model Model) GetValue() string {
	return model.input.GetValue()
}

func (model Model) Clear() Model {
	model.input.SetValue("")
	return model
}

// Returns nameFilterLines, tagFilterLInes
func (model Model) GetFilterLines() ([]string, []string) {
	rawLines := strings.Split(model.input.GetValue(), "\n")

	nameFilterLines := make([]string, 0)
	tagFilterLines := make([]string, 0)
	for _, rawLine := range rawLines {
		filter := strings.TrimSpace(rawLine)

		isTagFilter := false
		if strings.HasPrefix(filter, tagFilterLineLeader) {
			isTagFilter = true
			filter = filter[1:]
		}

		if len(filter) == 0 {
			continue
		}

		if isTagFilter {
			tagFilterLines = append(tagFilterLines, filter)
		} else {
			nameFilterLines = append(nameFilterLines, filter)
		}
	}

	return nameFilterLines, tagFilterLines
}

func (model Model) IsInNormalMode() bool {
	return model.input.GetMode() == vim.NormalMode
}
