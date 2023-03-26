package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/list"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
	isSelected  bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list          list.Model
	selectedItems map[int]bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

type itemDelegate struct{}

func (i itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {

}

func (i itemDelegate) Height() int {
	//TODO implement me
	panic("implement me")
}

func (i itemDelegate) Spacing() int {
	//TODO implement me
	panic("implement me")
}

func (i itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	//TODO implement me
	panic("implement me")
}

func main() {
	content := []contentItem{
		{
			name: "Foo",
			tags: nil,
		},
		{
			name: "Bar",
			tags: nil,
		},
		{
			name: "Bang",
			tags: nil,
		},
	}

	model := &appModel{
		mode:                0,
		filterInput:         textinput.New(),
		content:             content,
		cursorIdx:           0,
		selectedItemIndexes: make(map[int]bool),
		height:              0,
		width:               0,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
