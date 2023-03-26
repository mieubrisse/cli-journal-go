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
	baseStyle := lipgloss.NewStyle().Width(model.width)

	// First calculate the footer, so we can get its height later
	footerStr := ""
	numSelectedItems := len(model.selectedItemIndexes.GetIndices())
	if len(model.selectedItemIndexes.GetIndices()) > 0 {
		footerStr = fmt.Sprintf(" %d items selected", numSelectedItems)
	}
	style := baseStyle.Copy().Foreground(lipgloss.Color("#eed202")).Align(lipgloss.Center)
	finalFooter := style.Render(footerStr)

	// Calculate content
	lines := []string{}
	if len(model.filteredContentIndices) == 0 {
		noItemsLine := baseStyle.Copy().Faint(true).Align(lipgloss.Center).Render("No items")
		lines = []string{noItemsLine}

	}

	for idx, contentIdx := range model.filteredContentIndices {
		content := model.allContent[contentIdx]

		maybeCheckmark := " "
		if model.selectedItemIndexes.Contains(idx) {
			maybeCheckmark = "✔️"
		}

		line := fmt.Sprintf(
			"%s  %s     %s",
			maybeCheckmark,
			content.Name,
			strings.Join(content.Tags, " "),
		)
		lineStyle := baseStyle.Copy()
		if model.isFocused && idx == model.cursorIdx {
			lineStyle = lineStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
		}
		renderedLine := lineStyle.Render(line)
		lines = append(lines, renderedLine)
	}

	contentLinesStr := lipgloss.JoinVertical(lipgloss.Left, lines...)
	contentHeight := model.height - lipgloss.Height(footerStr)
	finalContent := baseStyle.Copy().Height(contentHeight).Render(contentLinesStr)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		finalContent,
		finalFooter,
	)
}

// Updates the name filter text that this model knows about, and does the appropriate recalculations on the cursor
// NOTE: We have to do this because there doesn't seem to be a way to share a single textinput.Model component between two models
func (model *Model) UpdateFilters(nameFilterText string, tagFilterText string) {
	nameMatchPredicate := getNameMatchPredicate(nameFilterText)
	tagMatchPredicate := getTagMatchPredicate(tagFilterText)

	filteredContentIndices := []int{}
	for idx, content := range model.allContent {
		if !nameMatchPredicate(content) {
			continue
		}

		if !tagMatchPredicate(content) {
			continue
		}

		filteredContentIndices = append(filteredContentIndices, idx)
	}

	model.filteredContentIndices = filteredContentIndices
	model.cursorIdx = 0
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

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height
	return model
}

// ====================================================================================================
//
//	PRIVATE HELPER FUNCTIONS
//
// ====================================================================================================
func getNameMatchPredicate(filterText string) func(item content_item.ContentItem) bool {
	terms := strings.Fields(filterText)

	// No terms == match everything
	if len(terms) == 0 {
		return func(item content_item.ContentItem) bool {
			return true
		}
	}

	matcher := getFuzzyMatcherFromTerms(terms)

	return func(item content_item.ContentItem) bool {
		return matcher.MatchString(item.Name)
	}
}

func getTagMatchPredicate(filterText string) func(item content_item.ContentItem) bool {
	terms := strings.Fields(filterText)

	// No terms == match everything
	if len(terms) == 0 {
		return func(item content_item.ContentItem) bool {
			return true
		}
	}

	matcher := getFuzzyMatcherFromTerms(terms)

	// There's at least one search term now, so there must be at least one matching tag
	return func(item content_item.ContentItem) bool {
		for _, tag := range item.Tags {
			if matcher.MatchString(tag) {
				return true
			}
		}
		return false
	}
}

func getFuzzyMatcherFromTerms(terms []string) *regexp.Regexp {
	// We have terms, so we need to escape them
	escapedTerms := make([]string, 0, len(terms))
	for _, term := range terms {
		escapedTerms = append(escapedTerms, regexp.QuoteMeta(term))
	}

	// The (?i) makes the search case-insensitive
	regexStr := "(?i)" + strings.Join(escapedTerms, ".*")

	// Okay to use MustCompile here because we quote the user's input so it should be safe
	return regexp.MustCompile(regexStr)
}
