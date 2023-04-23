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
