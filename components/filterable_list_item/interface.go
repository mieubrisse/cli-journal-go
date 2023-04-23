package filterable_list_item

import "github.com/mieubrisse/cli-journal-go/components"

type FilterableListItemComponent interface {
	components.Component

	GetValue() string
}
