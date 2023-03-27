package filterable_content_list

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
	"github.com/mieubrisse/cli-journal-go/data_structures/selected_item_index_set"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"regexp"
	"strings"
)

type componentSize int

const (
	contentTimestampFormat = "2006-01-02 15:04:05"

	checkmarkChar = '•'

	// Used when a line is too small
	continuationChar = '…'

	maxNameWidth = 45

	minimumNameAndTagWidth = 5

	wide componentSize = iota
	medium
	narrow
	sliver
)

// Minimum width, in characters, for the component to be classed as each size
var componentSizeThresholds = map[componentSize]int{
	wide:   150,
	medium: 120,
	narrow: 80,
	sliver: 0,
}
var checkmarkWidthsByComponentSize = map[componentSize]int{
	wide:   5,
	medium: 4,
	narrow: 3,
	sliver: 2,
}
var timestampWidthsByComponentSize = map[componentSize]int{
	wide:   len(contentTimestampFormat) + 4,
	medium: len(contentTimestampFormat) + 2,
	narrow: 0,
	sliver: 0,
}

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

// Updates the name filter text that this model knows about, and does the appropriate recalculations on the cursor
// NOTE: We have to do this because there doesn't seem to be a way to share a single textinput.Model component between two models
func (model *Model) SetFilterTexts(nameFilterText string, tagFilterText string) {
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
}

func (model Model) AddItem(content content_item.ContentItem) Model {
	model.allContent = append(model.allContent, content)
	return model
}

func (model *Model) Focused() bool {
	return model.isFocused
}

// TODO return a model?
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
func recalculateView() {

}

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

func (model Model) renderContentLine(content content_item.ContentItem, isContentHighlighted bool, isContentSelected bool) string {
	baseLineStyle := lipgloss.NewStyle()
	if model.isFocused && isContentHighlighted {
		baseLineStyle = baseLineStyle.Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
	}

	// Calculate the widths for the various components
	biggestThresholdPassed := sliver
	for trialComponentSize, threshold := range componentSizeThresholds {
		if model.width > threshold && threshold > componentSizeThresholds[biggestThresholdPassed] {
			biggestThresholdPassed = trialComponentSize
		}
	}
	actualComponentSize := biggestThresholdPassed
	checkmarkWidth, found := checkmarkWidthsByComponentSize[actualComponentSize]
	if !found {
		panic("No checkmark width for terminal size")
	}
	timestampWidth, found := timestampWidthsByComponentSize[actualComponentSize]
	if !found {
		panic("No timestamp width for terminal size")
	}

	widthRemaining := helpers.GetMaxInt(0, model.width-checkmarkWidth-timestampWidth)
	// Safety valve: if we don't have at least 10 characters, don't even bother

	nameWidth := helpers.GetMinInt(
		maxNameWidth,
		int(0.6*float64(widthRemaining)),
	)
	tagsWidth := helpers.GetMaxInt(0, widthRemaining-nameWidth)

	// Checkmark string
	checkmarkStr := ""
	if isContentSelected {
		checkmarkStr = string(checkmarkChar)
	}
	checkmarkStr = baseLineStyle.Copy().
		Foreground(global_styles.Orange).
		Width(checkmarkWidth).
		AlignHorizontal(lipgloss.Center).
		Render(checkmarkStr)

	// Timestamp (disabled if too small)
	timestampStr := ""
	if timestampWidth > 0 {
		timestampStr = content.Timestamp.Format(contentTimestampFormat)
		timestampStr = baseLineStyle.Copy().
			Foreground(global_styles.Cyan).
			Width(timestampWidth).
			AlignHorizontal(lipgloss.Left).
			Render(timestampStr)
	}

	// Name
	nameStr := ""
	if nameWidth > minimumNameAndTagWidth {
		nameStr = content.Name
		nameLen := len(nameStr)
		if nameLen > nameWidth-1 {
			nameStr = nameStr[:nameWidth-2] + string(continuationChar)
		}
		nameStr = baseLineStyle.Copy().
			Foreground(global_styles.White).
			Width(nameWidth).
			AlignHorizontal(lipgloss.Left).
			Render(nameStr)
	}

	tagsStr := ""
	if tagsWidth > minimumNameAndTagWidth {
		tagsStr = strings.Join(content.Tags, " ")
		tagsLen := len(tagsStr)
		if tagsLen > tagsWidth-1 {
			tagsStr = tagsStr[:tagsWidth-2] + string(continuationChar)
		}
		tagsStr = baseLineStyle.Copy().
			Foreground(global_styles.Red).
			Width(tagsWidth).
			AlignHorizontal(lipgloss.Left).
			Render(tagsStr)
	}

	line := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checkmarkStr,
		timestampStr,
		nameStr,
		tagsStr,
	)
	/*
		line := fmt.Sprintf(
			"%s  %s %s     %s",
			maybeCheckmark,
			timestampStr,
			content.Name,
			joinedTags,
		)

	*/

	return baseLineStyle.Copy().Width(model.width).Render(line)
}
