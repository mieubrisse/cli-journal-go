package list_item

import "github.com/mieubrisse/cli-journal-go/components"

type ListItemComponent interface {
	components.Component

	GetValue() string
}
