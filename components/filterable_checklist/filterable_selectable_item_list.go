package filterable_checklist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/components/filterable_checklist_item"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
)

type implementation struct {
	filterableList filterable_list.Component

	items []filterable_checklist_item.Component

	// Indices of selected items within the *unfiltered* list
	selectedItemIndices map[int]bool

	isFocused bool
	width     int
	height    int
}

func New() {

}

func (impl *implementation) Resize(width int, height int) {
	impl.width = width
	impl.height = height

	impl.filterableList.Resize(width, height)
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl *implementation) Focus() tea.Cmd {
	impl.isFocused = true

	return impl.filterableList.Focus()
}

func (impl *implementation) Blur() tea.Cmd {
	impl.isFocused = false
	return impl.filterableList.Blur()
}

func (impl implementation) Focused() bool {
	return impl.isFocused
}

func (impl *implementation) UpdateFilter(newFilter func(int, filterable_list_item.Component) bool) {
	impl.filterableList.UpdateFilter(newFilter)
}

func (impl *implementation) SetItems(items []filterable_checklist_item.Component) {
	// TODO something about preserving the selected item indices when the list changes??
	impl.items = items
	impl.selectedItemIndices = make(map[int]bool, 0)

	castedItems := make([]filterable_list_item.Component, len(items))
	for idx, item := range items {
		castedItems[idx] = filterable_list_item.Component(item)
	}
	impl.filterableList.SetItems(castedItems)
}

func (impl *implementation) Scroll(scrollOffset int) {
	impl.filterableList.Scroll(scrollOffset)
}

func (impl implementation) GetItems() []filterable_checklist_item.Component {
	return impl.items
}

func (impl implementation) GetFilteredItemIndices() []int {
	return impl.filterableList.GetFilteredItemIndices()
}

func (impl implementation) GetHighlightedItemIndex() int {
	return impl.filterableList.GetHighlightedItemIndex()
}

func (impl *implementation) SetHighlightedItemSelection(isSelected bool) {
	filteredItemIndices := impl.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return
	}

	highlightedItemIdxInFilteredList := impl.GetHighlightedItemIndex()
	highlightedItemIdxInOriginalList := filteredItemIndices[highlightedItemIdxInFilteredList]

	impl.setItemSelection(highlightedItemIdxInOriginalList, isSelected)
}

func (impl *implementation) SetAllViewableItemsSelection(isSelected bool) {
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
func (impl *implementation) setItemSelection(itemIdx int, isSelected bool) {
	item := impl.items[itemIdx]
	item.SetSelection(isSelected)

	if isSelected {
		impl.selectedItemIndices[itemIdx] = true
	} else {
		delete(impl.selectedItemIndices, itemIdx)
	}
}
