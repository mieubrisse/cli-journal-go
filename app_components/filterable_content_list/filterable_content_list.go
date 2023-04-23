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

	filterPredicate func(item content_item.ContentItem) bool

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
	// First calculate the footer, so we can get its height later
	footerStr := ""
	numSelectedItems := len(model.selectedItemIndexes.GetIndices())
	if numSelectedItems > 0 {
		numberStr := fmt.Sprintf("%d", numSelectedItems)
		numberStr = lipgloss.NewStyle().Foreground(global_styles.Orange).Render(numberStr)

		textStr := lipgloss.NewStyle().Foreground(global_styles.White).Render(" items selected")

		footerStr = numberStr + textStr
	}
	style := lipgloss.NewStyle().
		Width(model.width).
		Align(lipgloss.Center)
	finalFooter := style.Render(footerStr)

	// Calculate content
	lines := []string{}
	if len(model.filteredContentIndices) == 0 {
		noItemsLine := lipgloss.NewStyle().
			Width(model.width).
			Faint(true).
			Align(lipgloss.Center).
			Render("No items")
		lines = []string{noItemsLine}
	}

	for idx, contentIdx := range model.filteredContentIndices {
		content := model.allContent[contentIdx]
		isContentHighlighted := idx == model.cursorIdx
		isContentSelected := model.selectedItemIndexes.Contains(contentIdx)

		line := model.renderContentLine(content, isContentHighlighted, isContentSelected)
		lines = append(lines, line)
	}

	contentLinesStr := strings.Join(lines, "\n")
	contentHeight := model.height - lipgloss.Height(footerStr)
	finalContent := lipgloss.NewStyle().Height(contentHeight).Render(contentLinesStr)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		finalContent,
		finalFooter,
	)
}

func (model Model) SetFilters(nameFilterLines []string, tagFilterLines []string) Model {
	nameFilterRegexes := transformTermsToFilterRegexes(nameFilterLines)
	tagFilterRegexes := transformTermsToFilterRegexes(tagFilterLines)

	// The way this predicate is structured is as a "gauntlet" - there are many opportunities for an item to be discarded,
	// and only if it passes those will it be in
	// I believe this to be the best way to structure predicates, because it makes it easier to think about
	predicate := func(item content_item.ContentItem) bool {
		// Filter out non-matching names
		for _, nameRegex := range nameFilterRegexes {
			if !nameRegex.MatchString(item.Name) {
				return false
			}
		}

		// If no tag filters are specified, skip this step of the gauntlet
		if len(tagFilterRegexes) > 0 {
			// If we have tag filters, we run a sub-gauntlet for tags, where at least one tag must make it through
			hasTagMatch := false
		tagLoop:
			for _, tag := range item.Tags {
				for _, tagRegex := range tagFilterRegexes {
					if !tagRegex.MatchString(tag) {
						continue tagLoop
					}
				}

				// A tag made it through the gauntlet; no need to process further tags
				hasTagMatch = true
				break
			}
			if !hasTagMatch {
				return false
			}
		}

		return true
	}

	model.filterPredicate = predicate

	model = model.recalculateView()
	return model
}

func (model Model) AddItem(content content_item.ContentItem) Model {
	model.allContent = append(
		[]content_item.ContentItem{content},
		model.allContent...,
	)
	model = model.recalculateView()
	return model
}

func (model Model) Focused() bool {
	return model.isFocused
}

// TODO return a model?
func (model Model) Focus() Model {
	model.isFocused = true
	return model
}

func (model Model) Blur() Model {
	model.isFocused = false
	return model
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
func (model Model) recalculateView() Model {
	filteredContentIndices := []int{}
	for idx, content := range model.allContent {
		if !model.filterPredicate(content) {
			continue
		}

		filteredContentIndices = append(filteredContentIndices, idx)
	}

	// In the special case where we have exactly the same results, we can keep the cursor index
	newCursorIdx := 0
	if len(filteredContentIndices) == len(model.filteredContentIndices) {
		keepCursorIdx := true
		for idx, contentIdx := range filteredContentIndices {
			if model.filteredContentIndices[idx] != contentIdx {
				keepCursorIdx = false
				break
			}
		}

		if keepCursorIdx {
			newCursorIdx = model.cursorIdx
		}
	}
	model.cursorIdx = newCursorIdx

	model.filteredContentIndices = filteredContentIndices

	return model
}

/*
func getNameMatchPredicate(nameFilterLines []string) func(item content_item.ContentItem) bool {
	allSubPredicates := make([]func(item content_item.ContentItem) bool, len(nameFilterLines))
	for idx, filterLine := range nameFilterLines {
		terms := strings.Fields(filterLine)

		// No terms == match everything
		if len(terms) == 0 {
			allSubPredicates[idx] = func(item content_item.ContentItem) bool {
				return true
			}
			continue
		}

		matcher := getFuzzyMatcherFromTerms(terms)
		allSubPredicates[idx] = func(item content_item.ContentItem) bool {
			return matcher.MatchString(item.Name)
		}
	}

	// Glue subpredicates together
	return func(item content_item.ContentItem) bool {
		for _, subPredicate := range allSubPredicates {
			if !subPredicate(item) {
				return false
			}
		}
		return true
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

	// There's at least one search term now, so there must be at least one tag that matches all predicates
	return func(item content_item.ContentItem) bool {
		for _, tag := range item.Tags {
			if matcher.MatchString(tag) {
				return true
			}
		}
		return false
	}
}

*/

// Each line that has text, will produce a regex filter that must match
func transformTermsToFilterRegexes(lines []string) []*regexp.Regexp {
	result := make([]*regexp.Regexp, 0)
	for _, line := range lines {
		terms := strings.Fields(line)

		// Don't bother creating a regex for an empty line
		if len(terms) == 0 {
			continue
		}

		// We have terms, so we need to escape them
		escapedTerms := make([]string, 0, len(terms))
		for _, term := range terms {
			escapedTerms = append(escapedTerms, regexp.QuoteMeta(term))
		}

		// The (?i) makes the search case-insensitive
		regexStr := "(?i)" + strings.Join(escapedTerms, ".*")

		// Okay to use MustCompile here because we quote the user's input so it should be safe
		result = append(result, regexp.MustCompile(regexStr))
	}

	return result
}
