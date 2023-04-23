package text_input

import "github.com/mieubrisse/cli-journal-go/components"

type Component interface {
	components.InteractiveComponent

	GetValue() string
	SetValue(value string)
}
