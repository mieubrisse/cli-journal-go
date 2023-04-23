package filterable_checklist_item

import (
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
)

type Component interface {
	filterable_list_item.Component

	IsSelected() bool
	SetSelection(isSelected bool)
}
