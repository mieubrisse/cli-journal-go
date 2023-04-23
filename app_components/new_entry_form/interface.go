package new_entry_form

import "github.com/mieubrisse/cli-journal-go/components"

type Component interface {
	components.InteractiveComponent

	GetNameValue() string
	SetNameValue(name string)
	Clear()
}
