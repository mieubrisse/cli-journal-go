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

	// First calculate the footer, so we know its height so we know how big to make the content
	footerStr := ""
	numSelectedItems := len(model.selectedItemIndexes.GetIndices())
	if len(model.selectedItemIndexes.GetIndices()) > 0 {
		footerStr = fmt.Sprintf(" %d items selected", numSelectedItems)
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#eed202"))
	renderedFooter := style.Render(footerStr)

	// Now calculate content
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

		line := fmt.Sprintf(
			"%s  %s     %s",
			maybeCheckmark,
			content.Name,
			strings.Join(content.Tags, " "),
		)
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

	// Finally, slam everything together
	renderedContentStr := lipgloss.NewStyle().
		// Height(model.height - lipgloss.Height(renderedFooter)).
		Height(model.height - 2).
		Width(model.width).
		Render(contentStr)

	return lipgloss.JoinVertical(lipgloss.Left, renderedContentStr, renderedFooter)
}

// Updates the name filter text that this model knows about, and does the appropriate recalculations on the cursor
// NOTE: We have to do this because there doesn't seem to be a way to share a single textinput.Model component between two models
func (model *Model) UpdateFilters(nameFilterText string, tagFilterText string) {
	nameMatchPredicate := getPredicateFromSearchTerms(nameFilterText)
	tagMatchPredicate := getPredicateFromSearchTerms(tagFilterText)

	filteredContentIndices := []int{}
	for idx, content := range model.allContent {
		if !nameMatchPredicate(content.Name) {
			continue
		}

		hasTagMatch := false
		for _, tag := range content.Tags {
			if tagMatchPredicate(tag) {
				hasTagMatch = true
				break
			}
		}
		if !hasTagMatch {
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
	model.height = height
	model.width = width
	return model
}

func getPredicateFromSearchTerms(termsText string) func(string) bool {
	terms := strings.Fields(termsText)

	// No terms matches everything
	if len(terms) == 0 {
		return func(string) bool {
			return true
		}
	}

	// We have terms, so we need to escape them
	escapedTerms := make([]string, 0, len(terms))
	for _, term := range terms {
		escapedTerms = append(escapedTerms, regexp.QuoteMeta(term))
	}

	// The (?i) makes the search case-insensitive
	regexStr := "(?i)" + strings.Join(escapedTerms, ".*")

	// Okay to use MustCompile here because we quote the user's input so it should be safe
	matcher := regexp.MustCompile(regexStr)

	return func(str string) bool {
		return matcher.MatchString(str)
	}
}
