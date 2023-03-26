package components

type ResizableComponent[T any] interface {
	Resize(width int, height int) T
	GetHeight() int
	GetWidth() int
}
