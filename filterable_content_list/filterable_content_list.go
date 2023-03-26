package filterable_content_list

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go"
)

// This
type Model struct {
	filterInput main.filterTextInput

	content             []main.contentItem
	cursorIdx           int
	selectedItemIndexes map[int]bool
}

// TODO replace content with contentProvider
func New(content []main.contentItem, filterInput textinput.Model) Model {
	return Model{
		filterInput:         filterInput,
		content:             content,
		cursorIdx:           0,
		selectedItemIndexes: make(map[int]bool),
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			newCursorIdx := model.cursorIdx + 1
			if newCursorIdx < len(model.content) {
				model.cursorIdx = newCursorIdx
			}
		case "k":
			newCursorIdx := model.cursorIdx - 1
			if newCursorIdx >= 0 {
				model.cursorIdx = newCursorIdx
			}
		case "x":
			_, found := model.selectedItemIndexes[model.cursorIdx]
			if found {
				delete(model.selectedItemIndexes, model.cursorIdx)
			} else {
				model.selectedItemIndexes[model.cursorIdx] = true
			}
		case "a":
			if len(model.selectedItemIndexes) < len(model.content) {
				for idx := range model.content {
					model.selectedItemIndexes[idx] = true
				}
			} else {
				model.selectedItemIndexes = make(map[int]bool)
			}
		case "c":
			// Clear the filter
			model.filterInput.SetValue("")
		case "/":
			model.mode = main.filterMode

			// This will tell the input that it should display the cursor
			cmd := model.filterInput.Focus()
			return model, cmd
		}
	}

	return model, nil
}

func (model Model) View() string {
	//TODO implement me
	panic("implement me")
}
