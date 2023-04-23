package filterable_list_item

import (
	"github.com/mieubrisse/cli-journal-go/components/framework"
)

type FilterableListItemComponent interface {
	framework.Component

	GetValue() string
}
