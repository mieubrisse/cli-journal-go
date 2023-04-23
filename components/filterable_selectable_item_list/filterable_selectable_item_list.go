package filterable_selectable_item_list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/components/filterable_item_list"
)

type Model[T FilterableSelectableListItem] struct {
	filterableList filterable_item_list.Model[T]

	items []T

	// Indices of selected items within the *unfiltered* list
	selectedItemIndices map[int]bool

	isFocused bool
	width     int
	height    int
}

func New() {

}

func (model Model[T]) Resize(width int, height int) Model[T] {
	model.width = width
	model.height = height

	model.filterableList = model.filterableList.Resize(width, height)
	return model
}

func (model Model[T]) GetHeight() int {
	return model.height
}

func (model Model[T]) GetWidth() int {
	return model.width
}

func (model *Model[T]) Focus() tea.Cmd {
	model.isFocused = true

	return model.filterableList.Focus()
}

func (model *Model[T]) Blur() tea.Cmd {
	model.isFocused = false
	return model.filterableList.Blur()
}

func (model Model[T]) Focused() bool {
	return model.isFocused
}

func (model Model[T]) UpdateFilter(newFilter func(int, T) bool) Model[T] {
	model.filterableList = model.filterableList.UpdateFilter(newFilter)
	return model
}

func (model Model[T]) SetItems(items []T) Model[T] {
	// TDOO more stuff
	model.filterableList = model.filterableList.SetItems(items)
	return model
}

func (model Model[T]) Scroll(scrollOffset int) Model[T] {
	model.filterableList = model.filterableList.Scroll(scrollOffset)
	return model
}

func (model Model[T]) GetItems() []T {
	return model.items
}

func (model Model[T]) GetFilteredItemIndices() []int {
	return model.filterableList.GetFilteredItemIndices()
}

func (model Model[T]) GetHighlightedItemIndex() int {
	return model.filterableList.GetHighlightedItemIndex()
}

func (model Model[T]) SetHighlightedItemSelection(isSelected bool) Model[T] {
	filteredItemIndices := model.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return model
	}

	highlightedItemIdxInFilteredList := model.GetHighlightedItemIndex()
	highlightedItemIdxInOriginalList := filteredItemIndices[highlightedItemIdxInFilteredList]

	model = model.setItemSelection(highlightedItemIdxInOriginalList, isSelected)

	return model
}

func (model Model[T]) SetAllViewableItemsSelection(isSelected bool) Model[T] {
	filteredItemIndices := model.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return model
	}

	for _, originalItemIdx := range filteredItemIndices {
		model.setItemSelection(originalItemIdx, isSelected)
	}

	return model
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
func (model Model[T]) setItemSelection(itemIdx int, isSelected bool) Model[T] {
	item := model.items[itemIdx]
	item.SetSelection(isSelected)

	if isSelected {
		model.selectedItemIndices[itemIdx] = true
	} else {
		delete(model.selectedItemIndices, itemIdx)
	}
	return model
}
