package filterable_list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/list_item"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"strings"
)

/*
Component for displaying a scrollable, filterable list of items
*/
type Model[T list_item.ListItemComponent] struct {
	unfilteredItems []T

	// The indices of the filtered items within the unfiltered items list
	filteredItemsOriginalIndices []int

	// The index of the highlighted item within the *filtered list*
	highlightedItemIdx int

	isFocused bool
	width     int
	height    int
}

func New[T list_item.ListItemComponent](items []T) Model[T] {
	filteredIndices := []int{}
	for idx := range items {
		filteredIndices = append(filteredIndices, idx)
	}

	return Model[T]{
		unfilteredItems:              items,
		filteredItemsOriginalIndices: filteredIndices,
		highlightedItemIdx:           0,
		width:                        0,
		height:                       0,
	}
}

func (model Model[T]) View() string {
	baseLineStyle := lipgloss.NewStyle().
		Width(model.width)

	// As aesthetic choices, when there are more item lines than display lines:
	// 1. We want the entire list to scroll around the cursor if it's in the center of the screen, rather than
	//    the user needing to scroll to top or bottom to get the list to move. This helps the user see more
	//    relevant information at once
	// 2. When the cursor is near the top or bottom of the list, scroll the cursor rather than the entire list
	//    so that we don't get blank space
	// The easiest way to accomplish this is to calculate the range of acceptable first-line indexes of the view,
	//   which will range from [0, num_items - num_display_lines], and when the user is in the middle of the list
	//   the view will have the cursor line in the center
	halfHeight := model.height / 2

	// Ensure that, when near the bottom of the list, the cursor is no longer centered and scrolls to the bottom
	firstDisplayedLineIdxInclusive := helpers.GetMinInt(
		model.highlightedItemIdx-halfHeight,
		len(model.filteredItemsOriginalIndices)-model.height,
	)

	// Ensure that, when near the top of the list, the cursor is no longer centered and scrolls to the top
	firstDisplayedLineIdxInclusive = helpers.GetMaxInt(
		firstDisplayedLineIdxInclusive,
		0,
	)

	lastDisplayedLineIdxExclusive := helpers.GetMinInt(
		len(model.filteredItemsOriginalIndices),
		firstDisplayedLineIdxInclusive+model.height,
	)

	if firstDisplayedLineIdxInclusive == lastDisplayedLineIdxExclusive {
		return ""
	}

	displayedItems := model.filteredItemsOriginalIndices[firstDisplayedLineIdxInclusive:lastDisplayedLineIdxExclusive]

	viewableLinesHighlightedItemIdx := model.highlightedItemIdx - firstDisplayedLineIdxInclusive

	resultLines := []string{}
	for idx, originalItemIdx := range displayedItems {
		item := model.unfilteredItems[originalItemIdx]

		lineStyle := baseLineStyle
		if model.isFocused && idx == viewableLinesHighlightedItemIdx {
			lineStyle = baseLineStyle.Copy().Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
		}
		renderedLine := lineStyle.Render(item.View())

		resultLines = append(resultLines, renderedLine)
	}

	result := strings.Join(resultLines, "\n")

	// TODO truncating long lines by printable char

	return lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		MaxWidth(model.width).
		MaxHeight(model.height).
		Render(result)
}

func (model *Model[T]) UpdateFilter(newFilter func(int, T) bool) {
	highlightedOriginalItemIdx := model.filteredItemsOriginalIndices[model.highlightedItemIdx]

	// By default, assume that the highlighted item in the pre-filter list doesn't exist in the
	// post-filter list (but we'll fix this below if the assumption is false)
	model.highlightedItemIdx = 0

	newFilteredItemOriginalIndices := []int{}
	for idx, item := range model.unfilteredItems {
		if newFilter(idx, item) {
			newFilteredItemOriginalIndices = append(newFilteredItemOriginalIndices, idx)

			// If the highlighted item in the pre-filter list also exists in the post-filter list,
			// leave it highlighted
			if idx == highlightedOriginalItemIdx {
				model.highlightedItemIdx = len(newFilteredItemOriginalIndices) - 1
			}
		}
	}
	model.filteredItemsOriginalIndices = newFilteredItemOriginalIndices
}

func (model *Model[T]) SetItems(items []T) {
	filteredIndices := []int{}
	for idx := range items {
		filteredIndices = append(filteredIndices, idx)
	}

	model.unfilteredItems = items
	model.filteredItemsOriginalIndices = filteredIndices
	model.highlightedItemIdx = 0
}

// Scrolls the highlighted selection down or up by the specified number of items, with safeguards to
// prevent scrolling off the ends of the list
func (model *Model[T]) Scroll(scrollOffset int) {
	newHighlightedItemIdx := model.highlightedItemIdx + scrollOffset
	if newHighlightedItemIdx < 0 {
		newHighlightedItemIdx = 0
	}
	if newHighlightedItemIdx >= len(model.filteredItemsOriginalIndices) {
		newHighlightedItemIdx = len(model.filteredItemsOriginalIndices) - 1
	}
	model.highlightedItemIdx = newHighlightedItemIdx
}

func (model Model[T]) GetItems() []T {
	return model.unfilteredItems
}

// Gets the indices (within the original list) of the items currently being displayed
func (model Model[T]) GetFilteredItemIndices() []int {
	return model.filteredItemsOriginalIndices
}

// GetHighlightedItemIndex returns the index *within the filtered list* of the highlighted item
func (model Model[T]) GetHighlightedItemIndex() int {
	return model.highlightedItemIdx
}

func (model *Model[T]) Resize(width int, height int) {
	model.width = width
	model.height = height

	for _, item := range model.unfilteredItems {
		// TODO Allow items to wrap (but requires a whole viewing framework)
		item.Resize(width, 1)
	}
}

func (model Model[T]) GetHeight() int {
	return model.height
}

func (model Model[T]) GetWidth() int {
	return model.width
}

func (model *Model[T]) Focus() tea.Cmd {
	model.isFocused = true
	return nil
}

func (model *Model[T]) Blur() tea.Cmd {
	model.isFocused = false
	return nil
}

func (model Model[T]) Focused() bool {
	return model.isFocused
}
