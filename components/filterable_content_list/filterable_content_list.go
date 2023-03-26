package filterable_content_list

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
	"github.com/mieubrisse/cli-journal-go/data_structures/selected_item_index_set"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"regexp"
	"strings"
)

// This
type Model struct {
	// Whether to highlight the cursor line or not
	isFocused bool

	// All content that exists
	allContent []content_item.ContentItem

	// A list of indexes from the allContent list that match the given filter
	filteredContentIndices []int

	// Cursor index _within the filtered list_
	cursorIdx int

	selectedItemIndexes *selected_item_index_set.SelectedItemIndexSet

	height int
	width  int
}

// TODO replace content with contentProvider
func New(content []content_item.ContentItem) Model {
	filteredContentIndices := []int{}
	for idx := range content {
		filteredContentIndices = append(filteredContentIndices, idx)
	}

	return Model{
		// TODO this needs to be dynamic from the start args of the program - not sure how to do this!
		isFocused:              true,
		allContent:             content,
		filteredContentIndices: filteredContentIndices,
		cursorIdx:              0,
		selectedItemIndexes:    selected_item_index_set.New(),
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			newCursorIdx := model.cursorIdx + 1
			if newCursorIdx >= len(model.filteredContentIndices) {
				return model, nil
			}
			model.cursorIdx = newCursorIdx
			return model, nil
		case "k":
			newCursorIdx := model.cursorIdx - 1
			if newCursorIdx < 0 {
				return model, nil
			}
			model.cursorIdx = newCursorIdx
			return model, nil
		case "x":
			if model.selectedItemIndexes.Contains(model.cursorIdx) {
				model.selectedItemIndexes.Remove(model.cursorIdx)
			} else {
				model.selectedItemIndexes.Add(model.cursorIdx)
			}
			return model, nil
		}
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
	}

	return model, nil
}

func (model Model) View() string {
	componentStyle := lipgloss.NewStyle()
	if model.isFocused {
		componentStyle = componentStyle.Background(global_styles.FocusedComponentBackgroundColor)
	}

	contentLines := []string{}
	if len(model.filteredContentIndices) == 0 {
		contentLines = append(contentLines, "<no items>")
	}

	for idx, contentIdx := range model.filteredContentIndices {
		content := model.allContent[contentIdx]

		maybeCheckmark := " "
		if model.selectedItemIndexes.Contains(idx) {
			maybeCheckmark = "✔️"
		}

		line := fmt.Sprintf("%s  %s", maybeCheckmark, content.Name)
		lineStyle := lipgloss.NewStyle().Width(model.width)
		if model.isFocused && idx == model.cursorIdx {
			lineStyle = lineStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
		}
		renderedLine := lineStyle.Render(line)
		contentLines = append(contentLines, renderedLine)
	}
	contentStr := lipgloss.JoinVertical(
		lipgloss.Left,
		contentLines...,
	)

	footerStr := ""
	numSelectedItems := len(model.selectedItemIndexes.GetIndices())
	if len(model.selectedItemIndexes.GetIndices()) > 0 {
		footerStr = fmt.Sprintf(" %d items selected", numSelectedItems)
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#eed202"))
	renderedFooter := style.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left, contentStr, "", renderedFooter)
}

func (model *Model) Focused() bool {
	return model.isFocused
}

func (model *Model) Focus() {
	model.isFocused = true
}

func (model *Model) Blur() {
	model.isFocused = false
}

// Updates the name filter text that this model knows about, and does the appropriate recalculations on the cursor
// NOTE: We have to do this because there doesn't seem to be a way to share a single textinput.Model component between two models
func (model *Model) UpdateNameFilterText(filterText string) {
	escapedFilterText := regexp.QuoteMeta(filterText)
	nameSearchTerms := strings.Fields(escapedFilterText)
	// The (?i) makes the search case-insensitive
	regexStr := "(?i)" + strings.Join(nameSearchTerms, ".*")
	// Okay to use MustCompile here because we quote the user's input so it should be safe
	matcher := regexp.MustCompile(regexStr)

	filteredContentIndices := []int{}
	for idx, content := range model.allContent {
		if !matcher.MatchString(content.Name) {
			continue
		}
		filteredContentIndices = append(filteredContentIndices, idx)
	}

	model.filteredContentIndices = filteredContentIndices
	model.cursorIdx = 0
}

func (model Model) Resize(width int, height int) Model {
	model.height = height
	model.width = width
	return model
}