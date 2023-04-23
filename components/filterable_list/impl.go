package filterable_list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"strings"
)

/*
Component for displaying a scrollable, filterable list of items
*/
type implementation struct {
	unfilteredItems []filterable_list_item.Component

	// The indices of the filtered items within the unfiltered items list
	filteredItemsOriginalIndices []int

	// The index of the highlighted item within the *filtered list*
	highlightedItemIdx int

	isFocused bool
	width     int
	height    int
}

func New(items []filterable_list_item.Component) Component {
	filteredIndices := []int{}
	for idx := range items {
		filteredIndices = append(filteredIndices, idx)
	}

	return &implementation{
		unfilteredItems:              items,
		filteredItemsOriginalIndices: filteredIndices,
		highlightedItemIdx:           0,
		width:                        0,
		height:                       0,
	}
}

func (impl implementation) View() string {
	baseLineStyle := lipgloss.NewStyle().
		Width(impl.width)

	// As aesthetic choices, when there are more item lines than display lines:
	// 1. We want the entire list to scroll around the cursor if it's in the center of the screen, rather than
	//    the user needing to scroll to top or bottom to get the list to move. This helps the user see more
	//    relevant information at once
	// 2. When the cursor is near the top or bottom of the list, scroll the cursor rather than the entire list
	//    so that we don't get blank space
	// The easiest way to accomplish this is to calculate the range of acceptable first-line indexes of the view,
	//   which will range from [0, num_items - num_display_lines], and when the user is in the middle of the list
	//   the view will have the cursor line in the center
	halfHeight := impl.height / 2

	// Ensure that, when near the bottom of the list, the cursor is no longer centered and scrolls to the bottom
	firstDisplayedLineIdxInclusive := helpers.GetMinInt(
		impl.highlightedItemIdx-halfHeight,
		len(impl.filteredItemsOriginalIndices)-impl.height,
	)

	// Ensure that, when near the top of the list, the cursor is no longer centered and scrolls to the top
	firstDisplayedLineIdxInclusive = helpers.GetMaxInt(
		firstDisplayedLineIdxInclusive,
		0,
	)

	lastDisplayedLineIdxExclusive := helpers.GetMinInt(
		len(impl.filteredItemsOriginalIndices),
		firstDisplayedLineIdxInclusive+impl.height,
	)

	if firstDisplayedLineIdxInclusive == lastDisplayedLineIdxExclusive {
		return ""
	}

	displayedItems := impl.filteredItemsOriginalIndices[firstDisplayedLineIdxInclusive:lastDisplayedLineIdxExclusive]

	viewableLinesHighlightedItemIdx := impl.highlightedItemIdx - firstDisplayedLineIdxInclusive

	resultLines := []string{}
	for idx, originalItemIdx := range displayedItems {
		item := impl.unfilteredItems[originalItemIdx]

		lineStyle := baseLineStyle
		if impl.isFocused && idx == viewableLinesHighlightedItemIdx {
			// NOTE: this _may_ mess up the styling of the inner stuff
			lineStyle = baseLineStyle.Copy().Background(global_styles.FocusedComponentBackgroundColor).Bold(true)
		}
		renderedLine := lineStyle.Render(item.View())

		resultLines = append(resultLines, renderedLine)
	}

	result := strings.Join(resultLines, "\n")

	// TODO truncating long lines by printable char

	return lipgloss.NewStyle().
		Width(impl.width).
		Height(impl.height).
		MaxWidth(impl.width).
		MaxHeight(impl.height).
		Render(result)
}

func (impl *implementation) Update(msg tea.Msg) tea.Cmd {
	// Do nothing on non-Keymsgs
	switch msg.(type) {
	case tea.KeyMsg:
		// Proceed to rest of function
	default:
		return nil
	}

	if !impl.isFocused {
		return nil
	}

	// TODO allow for KeyMap overrides here?
	castedMsg := msg.(tea.KeyMsg)
	switch castedMsg.String() {
	case "j":
		impl.Scroll(1)
	case "k":
		impl.Scroll(-1)
	case "J":
		impl.Scroll(impl.height)
	case "K":
		impl.Scroll(-impl.height)
	}
	return nil
}

func (impl *implementation) UpdateFilter(newFilter func(int, filterable_list_item.Component) bool) {
	// This is a hack to indicate "the filtered list was empty, so there's no highlighted item original idx"
	oldHighlightedItemOriginalIdx := -1

	// If there are items being displayed, de-highlight the current item (if there are
	if len(impl.filteredItemsOriginalIndices) > 0 {
		oldHighlightedItemOriginalIdx = impl.filteredItemsOriginalIndices[impl.highlightedItemIdx]
		oldHighlightedItem := impl.unfilteredItems[oldHighlightedItemOriginalIdx]
		oldHighlightedItem.SetHighlighted(false)
	}

	// By default, assume that the highlighted item in the pre-filter list doesn't exist in the
	// post-filter list (but we'll fix this below if the assumption is false)
	newHighlightedItemIdx := 0

	newFilteredItemOriginalIndices := []int{}
	for idx, item := range impl.unfilteredItems {
		if newFilter(idx, item) {
			newFilteredItemOriginalIndices = append(newFilteredItemOriginalIndices, idx)

			// If the previously-highlighted item also exists in the post-filter list,
			// leave it highlighted
			if idx == oldHighlightedItemOriginalIdx {
				newHighlightedItemIdx = len(newFilteredItemOriginalIndices) - 1
			}
		}
	}

	impl.filteredItemsOriginalIndices = newFilteredItemOriginalIndices
	impl.highlightedItemIdx = newHighlightedItemIdx

	// Highlight the new item (if possible)
	if len(impl.filteredItemsOriginalIndices) > 0 {
		item := impl.unfilteredItems[impl.highlightedItemIdx]
		item.SetHighlighted(true)
	}
}

func (impl *implementation) SetItems(items []filterable_list_item.Component) {
	filteredIndices := []int{}
	for idx := range items {
		filteredIndices = append(filteredIndices, idx)
	}

	impl.unfilteredItems = items
	impl.filteredItemsOriginalIndices = filteredIndices
	impl.highlightedItemIdx = 0

	if len(impl.filteredItemsOriginalIndices) > 0 {
		highlightedItemOriginalIdx := impl.filteredItemsOriginalIndices[impl.highlightedItemIdx]
		item := impl.unfilteredItems[highlightedItemOriginalIdx]
		item.SetHighlighted(true)
	}
}

// Scrolls the highlighted selection down or up by the specified number of items, with safeguards to
// prevent scrolling off the ends of the list
func (impl *implementation) Scroll(scrollOffset int) {
	newHighlightedItemIdx := impl.highlightedItemIdx + scrollOffset
	if newHighlightedItemIdx < 0 {
		newHighlightedItemIdx = 0
	}
	if newHighlightedItemIdx >= len(impl.filteredItemsOriginalIndices) {
		newHighlightedItemIdx = len(impl.filteredItemsOriginalIndices) - 1
	}

	if newHighlightedItemIdx == impl.highlightedItemIdx {
		return
	}

	// De-highlight the previous item
	oldHighlightOriginalIdx := impl.filteredItemsOriginalIndices[impl.highlightedItemIdx]
	oldItem := impl.unfilteredItems[oldHighlightOriginalIdx]
	oldItem.SetHighlighted(false)

	// Highlight the new item
	newHighlightOriginalIdx := impl.filteredItemsOriginalIndices[newHighlightedItemIdx]
	newItem := impl.unfilteredItems[newHighlightOriginalIdx]
	newItem.SetHighlighted(true)

	impl.highlightedItemIdx = newHighlightedItemIdx
}

func (impl implementation) GetItems() []filterable_list_item.Component {
	return impl.unfilteredItems
}

// Gets the indices (within the original list) of the items currently being displayed
func (impl implementation) GetFilteredItemIndices() []int {
	return impl.filteredItemsOriginalIndices
}

// GetHighlightedItemIndex returns the index *within the filtered list* of the highlighted item
func (impl implementation) GetHighlightedItemIndex() int {
	return impl.highlightedItemIdx
}

func (impl *implementation) Resize(width int, height int) {
	impl.width = width
	impl.height = height

	for _, item := range impl.unfilteredItems {
		// TODO Allow items to wrap (but requires a whole viewing framework)
		item.Resize(width, 1)
	}
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl *implementation) Focus() tea.Cmd {
	impl.isFocused = true
	return nil
}

func (impl *implementation) Blur() tea.Cmd {
	impl.isFocused = false
	return nil
}

func (impl implementation) Focused() bool {
	return impl.isFocused
}
