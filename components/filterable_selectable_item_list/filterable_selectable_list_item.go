package filterable_selectable_item_list

import "github.com/mieubrisse/cli-journal-go/components/filterable_item_list"

type FilterableSelectableListItem interface {
	filterable_item_list.FilterableListItem

	SetSelection(isSelected bool)
}
