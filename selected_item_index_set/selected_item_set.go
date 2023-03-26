package selected_item_index_set

// Mathematical set to track which content item indexes are selected
type SelectedItemIndexSet struct {
	selected map[int]bool
}

func New(indices ...int) *SelectedItemIndexSet {
	selected := make(map[int]bool)
	for _, index := range indices {
		selected[index] = true
	}

	return &SelectedItemIndexSet{selected: selected}
}

func (set *SelectedItemIndexSet) Contains(index int) bool {
	_, found := set.selected[index]
	return found
}

func (set *SelectedItemIndexSet) Add(indices ...int) {
	for _, index := range indices {
		set.selected[index] = true
	}
}

func (set *SelectedItemIndexSet) Remove(indices ...int) {
	for _, index := range indices {
		delete(set.selected, index)
	}
}

func (set *SelectedItemIndexSet) Clear() {
	set.selected = make(map[int]bool)
}

func (set *SelectedItemIndexSet) GetIndices() map[int]bool {
	return set.selected
}
