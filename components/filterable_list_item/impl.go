package filterable_list_item

import (
	"github.com/mieubrisse/cli-journal-go/components"
)

type implementation struct {
	inner components.Component

	value string

	width  int
	height int
}

func New(value string, innerComponent components.Component) Component {
	return &implementation{
		inner:  innerComponent,
		value:  value,
		width:  0,
		height: 0,
	}
}

func (impl implementation) View() string {
	return impl.inner.View()
}

func (impl *implementation) Resize(width int, height int) {
	impl.inner.Resize(width, height)
	impl.width = width
	impl.height = height
}

func (impl implementation) GetWidth() int {
	return impl.width
}

func (impl implementation) GetHeight() int {
	return impl.height
}

func (impl implementation) GetValue() string {
	return impl.value
}
