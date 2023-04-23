package filterable_item_list

type FilterableListItem interface {
	// Returns a string representation of the list item
	Render() string
}
