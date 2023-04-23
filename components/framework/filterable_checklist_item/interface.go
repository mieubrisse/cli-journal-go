package filterable_checklist_item

import (
	"github.com/mieubrisse/cli-journal-go/components/framework/filterable_list_item"
)

type Component interface {
	filterable_list_item.FilterableListItemComponent

	SetSelection(isSelected bool)
}
