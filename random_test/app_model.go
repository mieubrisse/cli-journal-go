package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mieubrisse/cli-journal-go/filterable_item_list"
)

type item struct {
	number int
}

func (i item) Render() string {
	return fmt.Sprintf("%d", i.number)
}

type appModel struct {
	list filterable_item_list.Model[item]

	width  int
	height int
}

func New() appModel {
	items := []item{}
	for i := 0; i < 70; i++ {
		items = append(items, item{number: i})
	}
	itemList := filterable_item_list.New[item](items)
	return appModel{
		list:   itemList,
		width:  0,
		height: 0,
	}
}

func (model appModel) Init() tea.Cmd {
	return nil
}

func (model appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return model, tea.Quit
		case "k":
			model.list = model.list.Scroll(-1)
		case "j":
			model.list = model.list.Scroll(1)
		case "K":
			model.list = model.list.Scroll(-20)
		case "J":
			model.list = model.list.Scroll(20)
		}
	case tea.WindowSizeMsg:
		model = model.Resize(msg.Width, msg.Height)
	}

	return model, nil
}

func (model appModel) View() string {
	return model.list.View()
}

func (model appModel) Resize(width int, height int) appModel {
	model.width = width
	model.height = height

	model.list = model.list.Resize(width, height)

	return model
}

func (model appModel) GetHeight() int {
	return model.height
}

func (model appModel) GetWidth() int {
	return model.width
}
