package filterable_content_list

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/content_item"
	"github.com/mieubrisse/cli-journal-go/selected_item_index_set"
	"regexp"
	"strings"
)

// This
type Model struct {
	filterInput textinput.Model

	// Whether to highlight the cursor line or not
	isFocused bool

	content   []content_item.Model
	cursorIdx int

	selectedItemIndexes *selected_item_index_set.SelectedItemIndexSet
}

// TODO replace content with contentProvider
func New(content []content_item.Model, filterInput textinput.Model) Model {
	return Model{
		filterInput:         filterInput,
		content:             content,
		cursorIdx:           0,
		selectedItemIndexes: selected_item_index_set.New(),
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
			if newCursorIdx >= len(model.content) {
				return model, nil
			}

			model.content[model.cursorIdx].SetHighlighted(false)
			model.content[newCursorIdx].SetHighlighted(true)
			model.cursorIdx = newCursorIdx
			return model, nil
		case "k":
			newCursorIdx := model.cursorIdx + 1
			if newCursorIdx < 0 {
				return model, nil
			}

			model.content[model.cursorIdx].SetHighlighted(false)
			model.content[newCursorIdx].SetHighlighted(true)
			model.cursorIdx = newCursorIdx
			return model, nil
		case "x":
			cursorItem := model.content[model.cursorIdx]
			cursorItem.SetSelected(!cursorItem.IsSelected())
			return model, nil
		case "s":
			// Select all
			for _, item := range model.content {
				item.SetSelected(true)
			}
		case "d":
			// Deselect all
			for _, item := range model.content {
				item.SetSelected(false)
			}
		case "c":
			// Clear the filter
			model.filterInput.SetValue("")
			return model, nil
		}
	}

	return model, nil
}

func (model Model) View() string {
	resultBuilder := strings.Builder{}

	escapedFilterText := regexp.QuoteMeta(model.filterInput.Value())
	nameSearchTerms := strings.Fields(escapedFilterText)
	// The (?i) makes the search case-insensitive
	regexStr := "(?i)" + strings.Join(nameSearchTerms, ".*")
	// Okay to use MustCompile here because we quote the user's input so it should be safe
	matcher := regexp.MustCompile(regexStr)
	for _, item := range model.content {
		if !matcher.MatchString(item.GetName()) {
			continue
		}

		renderedItem := item.View()

		resultBuilder.WriteString(" " + renderedItem + "\n")
	}
	resultBuilder.WriteString("\n")

	return resultBuilder.String()
}

func (model *Model) Focus() {
	model.isFocused = true
}

func (model *Model) Blur() {
	model.isFocused = false
}
