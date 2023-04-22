package app

type tabCompletionItem struct {
	completion string
}

func (t tabCompletionItem) Render() string {
	return t.completion
}
