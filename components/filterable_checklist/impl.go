package filterable_checklist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/components/filterable_checklist_item"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
)

type implementation struct {
	innerList filterable_list.Component

	items []filterable_checklist_item.Component

	// Indices of selected items within the *unfiltered* list
	selectedItemIndices map[int]bool

	isFocused bool
	width     int
	height    int
}

// TODO get rid of items in constructor
func New(items []filterable_checklist_item.Component) Component {
	castedItems := make([]filterable_list_item.Component, len(items))
	for idx, item := range items {
		castedItems[idx] = filterable_list_item.Component(item)
	}
	inner := filterable_list.New(castedItems)
	return &implementation{
		innerList:           inner,
		items:               items,
		selectedItemIndices: make(map[int]bool, 0),
		isFocused:           false,
		width:               0,
		height:              0,
	}
}

func (impl *implementation) View() string {
	return impl.innerList.View()
}

func (impl *implementation) Update(msg tea.Msg) tea.Cmd {
	// Do nothing on non-Keymsgs
	switch msg.(type) {
	case tea.KeyMsg:
		// Proceed to rest of function
	default:
		return nil
	}

	// TODO allow for KeyMap overrides here?
	var returnCmd tea.Cmd
	castedMsg := msg.(tea.KeyMsg)
	switch castedMsg.String() {
	case "x":
		filteredItemOriginalIndicies := impl.innerList.GetFilteredItemIndices()
		if len(filteredItemOriginalIndicies) == 0 {
			break
		}
		itemOriginalIdx := filteredItemOriginalIndicies[impl.innerList.GetHighlightedItemIndex()]
		item := impl.items[itemOriginalIdx]
		isSelected := item.IsSelected()
		impl.setItemSelection(itemOriginalIdx, !isSelected)
	case "S":
		for _, originalIdx := range impl.innerList.GetFilteredItemIndices() {
			impl.setItemSelection(originalIdx, true)
		}
	case "D":
		for _, originalIdx := range impl.innerList.GetFilteredItemIndices() {
			impl.setItemSelection(originalIdx, false)
		}
	default:
		returnCmd = impl.innerList.Update(msg)
	}

	return returnCmd
}

func (impl implementation) GetItems() []filterable_checklist_item.Component {
	return impl.items
}

func (impl *implementation) SetItems(items []filterable_checklist_item.Component) {
	// TODO something about preserving the selected item indices when the list changes??
	impl.items = items
	impl.selectedItemIndices = make(map[int]bool, 0)

	castedItems := make([]filterable_list_item.Component, len(items))
	for idx, item := range items {
		castedItems[idx] = filterable_list_item.Component(item)
	}
	impl.innerList.SetItems(castedItems)
}

func (impl implementation) GetFilterableList() filterable_list.Component {
	return impl.innerList
}

func (impl implementation) GetSelectedItemOriginalIndices() map[int]bool {
	return impl.selectedItemIndices
}

func (impl *implementation) SetHighlightedItemSelection(isSelected bool) {
	filteredItemIndices := impl.innerList.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return
	}

	highlightedItemIdxInFilteredList := impl.innerList.GetHighlightedItemIndex()
	highlightedItemIdxInOriginalList := filteredItemIndices[highlightedItemIdxInFilteredList]

	impl.setItemSelection(highlightedItemIdxInOriginalList, isSelected)
}

func (impl *implementation) SetAllViewableItemsSelection(isSelected bool) {
	filteredItemIndices := impl.innerList.GetFilteredItemIndices()
	if len(filteredItemIndices) == 0 {
		return
	}

	for _, originalItemIdx := range filteredItemIndices {
		impl.setItemSelection(originalItemIdx, isSelected)
	}
}

func (impl *implementation) Resize(width int, height int) {
	impl.width = width
	impl.height = height

	impl.innerList.Resize(width, height)
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl *implementation) Focus() tea.Cmd {
	impl.isFocused = true

	return impl.innerList.Focus()
}

func (impl *implementation) Blur() tea.Cmd {
	impl.isFocused = false
	return impl.innerList.Blur()
}

func (impl implementation) Focused() bool {
	return impl.isFocused
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
// setItemSelection sets the selection of the given item, and does the appropriate bookkeeping
func (impl *implementation) setItemSelection(itemIdx int, isSelected bool) {
	item := impl.items[itemIdx]
	item.SetSelection(isSelected)

	if isSelected {
		impl.selectedItemIndices[itemIdx] = true
	} else {
		delete(impl.selectedItemIndices, itemIdx)
	}
}
