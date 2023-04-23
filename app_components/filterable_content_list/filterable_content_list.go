package filterable_content_list

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/app_components/entry_item"
	"github.com/mieubrisse/cli-journal-go/components/filterable_checklist"
	"github.com/mieubrisse/cli-journal-go/components/filterable_checklist_item"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"regexp"
	"strings"
)

// This
type Model struct {
	checklist filterable_checklist.Component

	items []entry_item.Component

	// Whether to highlight the cursor line or not
	isFocused bool

	height int
	width  int
}

// TODO replace content with contentProvider
func New(content []entry_item.Component) Model {
	castedContent := make([]filterable_checklist_item.Component, len(content))
	for idx, item := range content {
		castedContent[idx] = filterable_checklist_item.Component(item)
	}
	checklist := filterable_checklist.New(castedContent)

	return Model{
		checklist: nil,
		items:     nil,
		isFocused: false,
		height:    0,
		width:     0,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) tea.Cmd {
	// Do nothing on non-Keymsgs
	switch msg.(type) {
	case tea.KeyMsg:
		// Proceed to rest of function
	default:
		return nil
	}

	if !model.isFocused {
		return nil
	}

	return model.checklist.Update(msg)
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

func (model *Model) Focus() {
	model.isFocused = true
	model.checklist.Focus()
}

func (model *Model) Blur() Model {
	model.isFocused = false
	model.checklist.Blur()
}

func (model *Model) Resize(width int, height int) {
	model.width = width
	model.height = height
	model.checklist.Resize(width, height)
}

// ====================================================================================================
//
//	PRIVATE HELPER FUNCTIONS
//
// ====================================================================================================
/*
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
