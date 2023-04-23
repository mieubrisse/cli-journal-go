package checklist_item

import "github.com/mieubrisse/cli-journal-go/components"

type ChecklistItemComponent interface {
	components.Component

	SetSelection(isSelected bool)
}
