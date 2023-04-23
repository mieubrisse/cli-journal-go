package filterable_list

import (
	"github.com/mieubrisse/cli-journal-go/components"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
)

type Component[T filterable_list_item.Component] interface {
	components.InteractiveComponent

	UpdateFilter(newFilter func(idx int, item T) bool)
	SetItems(items []T)
	Scroll(scrollOffset int)
	GetItems() []T
	GetFilteredItemIndices() []int
	GetHighlightedItemIndex() int
}
