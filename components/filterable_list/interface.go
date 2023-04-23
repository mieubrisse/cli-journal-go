package filterable_list

import (
	"github.com/mieubrisse/cli-journal-go/components"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
)

type Component interface {
	components.FocusableComponent

	UpdateFilter(newFilter func(idx int, item filterable_list_item.Component) bool)
	SetItems(items []filterable_list_item.Component)
	Scroll(scrollOffset int)
	GetItems() []filterable_list_item.Component
	GetFilteredItemIndices() []int
	GetHighlightedItemIndex() int
}
