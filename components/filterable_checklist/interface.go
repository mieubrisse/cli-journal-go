package filterable_checklist

import (
	"github.com/mieubrisse/cli-journal-go/components"
	"github.com/mieubrisse/cli-journal-go/components/filterable_checklist_item"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list"
)

type Component interface {
	components.InteractiveComponent

	// Used for manipulations of the inner list (no need to reimplement all the functions)
	// The items in the original list will match the items from GetItems
	GetFilterableList() filterable_list.Component

	SetItems(items []filterable_checklist_item.Component)
	GetItems() []filterable_checklist_item.Component

	// GetSelectedItemOriginalIndices gets the indices within the current items list that are selected
	GetSelectedItemOriginalIndices() map[int]bool

	ToggleHighlightedItemSelection()

	// SetHighlightedItemSelection sets the selection for the currently-highlighted item
	SetHighlightedItemSelection(isSelected bool)

	// SetAllViewableItemsSelection sets the selection for all items that are currently shown (i.e. matching the filter)
	SetAllViewableItemsSelection(isSelected bool)

	// SetAllItemsSelection sets the selection on ALL items in the list (whether they're viewable or not)
	SetAllItemsSelection(isSelected bool)
}
