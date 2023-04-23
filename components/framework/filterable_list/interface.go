package filterable_list

import (
	"github.com/mieubrisse/cli-journal-go/components/framework"
	"github.com/mieubrisse/cli-journal-go/components/framework/filterable_list_item"
)

type FilterableListComponent[T filterable_list_item.FilterableListItemComponent] interface {
	framework.FocusableComponent

	UpdateFilter(newFilter func(idx int, item T) bool)
	SetItems(items []T)
	Scroll(scrollOffset int)
	GetItems() []T
	GetFilteredItemIndices() []int
	GetHighlightedItemIndex() int
}
