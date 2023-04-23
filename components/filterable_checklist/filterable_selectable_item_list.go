package filterable_checklist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/components/filterable_checklist_item"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list"
)

type implementation[T filterable_checklist_item.Component] struct {
	filterableList filterable_list.FilterableListComponent[T]

	items []T

	// Indices of selected items within the *unfiltered* list
	selectedItemIndices map[int]bool

	isFocused bool
	width     int
	height    int
}

func New() {

}

func (impl implementation[T]) Resize(width int, height int) implementation[T] {
	impl.width = width
	impl.height = height

	impl.filterableList.Resize(width, height)
	return impl
}

func (impl implementation[T]) GetHeight() int {
	return impl.height
}

func (impl implementation[T]) GetWidth() int {
	return impl.width
}

func (impl *implementation[T]) Focus() tea.Cmd {
	impl.isFocused = true

	return impl.filterableList.Focus()
}

func (impl *implementation[T]) Blur() tea.Cmd {
	impl.isFocused = false
	return impl.filterableList.Blur()
}

func (impl implementation[T]) Focused() bool {
	return impl.isFocused
}

func (impl *implementation[T]) UpdateFilter(newFilter func(int, T) bool) {
	impl.filterableList.UpdateFilter(newFilter)
}

func (impl *implementation[T]) SetItems(items []T) {
	impl.filterableList.SetItems(items)
}

func (impl *implementation[T]) Scroll(scrollOffset int) {
	impl.filterableList.Scroll(scrollOffset)
}

func (impl implementation[T]) GetItems() []T {
	return impl.items
}

func (impl implementation[T]) GetFilteredItemIndices() []int {
	return impl.filterableList.GetFilteredItemIndices()
}

func (impl implementation[T]) GetHighlightedItemIndex() int {
	return impl.filterableList.GetHighlightedItemIndex()
}

func (impl *implementation[T]) SetHighlightedItemSelection(isSelected bool) {
	filteredItemIndices := impl.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return
	}

	highlightedItemIdxInFilteredList := impl.GetHighlightedItemIndex()
	highlightedItemIdxInOriginalList := filteredItemIndices[highlightedItemIdxInFilteredList]

	impl.setItemSelection(highlightedItemIdxInOriginalList, isSelected)
}

func (impl *implementation[T]) SetAllViewableItemsSelection(isSelected bool) {
	filteredItemIndices := impl.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return
	}

	for _, originalItemIdx := range filteredItemIndices {
		impl.setItemSelection(originalItemIdx, isSelected)
	}
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
func (impl implementation[T]) setItemSelection(itemIdx int, isSelected bool) implementation[T] {
	item := impl.items[itemIdx]
	item.SetSelection(isSelected)

	if isSelected {
		impl.selectedItemIndices[itemIdx] = true
	} else {
		delete(impl.selectedItemIndices, itemIdx)
	}
	return impl
}
