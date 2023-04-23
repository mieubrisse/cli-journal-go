package filterable_list_item

import (
	"github.com/mieubrisse/cli-journal-go/components"
)

type Component interface {
	components.Component

	GetValue() string
}
