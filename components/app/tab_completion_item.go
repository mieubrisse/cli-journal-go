package app

type tabCompletionItem struct {
	completion string

	// TODO add matched char index
}

func (t tabCompletionItem) Render() string {
	return t.completion
}
