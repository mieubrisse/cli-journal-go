package filterable_list_item

import (
	"github.com/mieubrisse/cli-journal-go/components/framework"
)

type implementation[T framework.Component] struct {
	inner T

	value string

	width  int
	height int
}

func New[T framework.Component](value string, innerComponent T) FilterableListItemComponent {
	return &implementation[T]{
		inner:  innerComponent,
		value:  value,
		width:  0,
		height: 0,
	}
}

func (impl implementation[T]) View() string {
	return impl.inner.View()
}

func (impl *implementation[T]) Resize(width int, height int) {
	impl.inner.Resize(width, height)
	impl.width = width
	impl.height = height
}

func (impl implementation[T]) GetWidth() int {
	return impl.width
}

func (impl implementation[T]) GetHeight() int {
	return impl.height
}

func (impl implementation[T]) GetValue() string {
	return impl.value
}
