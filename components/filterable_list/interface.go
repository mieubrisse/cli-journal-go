package filterable_list

import (
	"github.com/mieubrisse/cli-journal-go/components"
	"github.com/mieubrisse/cli-journal-go/components/list_item"
)

type FilterableListComponent[T list_item.ListItemComponent] interface {
	components.FocusableComponent

	UpdateFilter(newFilter func(int, T) bool)
	SetItems(items []T)
	Scroll(scrollOffset int)
	GetItems() []T
	GetFilteredItemIndices() []int
	GetHighlightedItemIndex() int
}
