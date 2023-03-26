package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/filter_input"
	"github.com/mieubrisse/cli-journal-go/components/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/helpers"
)

type Mode int

const (
	padding = 1
)

var appStyle = lipgloss.NewStyle().Padding(padding)

type Model struct {
	mode Mode

	// filterInput tea.Model
	filterInput filter_input.Model

	contentList filterable_content_list.Model

	height int
	width  int
}

func New(filterInput filter_input.Model, contentList filterable_content_list.Model) Model {
	return Model{
		filterInput: filterInput,
		contentList: contentList,
		height:      0,
		width:       0,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

// NOTE: This returns a model because BubbleTea expects models to be passed by-value, so the way to "update" the model
// is to return a new instance of it
func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return model, tea.Quit
		}

		if model.contentList.Focused() {
			switch msg.String() {
			case "/":
				// TODO handle a command coming out the other side?
				model.contentList.Blur()

				// This will tell the input that it should display the cursor
				cmd := model.filterInput.Focus()

				return model, cmd
			case "c":
				model.filterInput.SetValue("")
				model.contentList.UpdateNameFilterText(model.filterInput.Value())
			}

			var cmd tea.Cmd
			model.contentList, cmd = model.contentList.Update(msg)
			return model, cmd

		} else if model.filterInput.Focused() {
			// Back out of filter mode
			if msg.String() == "esc" || msg.String() == "enter" {
				model.filterInput.Blur()
				model.contentList.Focus()
				return model, nil
			}

			var cmd tea.Cmd
			model.filterInput, cmd = model.filterInput.Update(msg)

			// Make sure to tell the content list about the new filter update
			model.contentList.UpdateNameFilterText(model.filterInput.Value())

			return model, cmd
		}
	case tea.WindowSizeMsg:
		return model.Resize(msg.Width, msg.Height), nil
	}

	return model, nil
}

func (model Model) View() string {
	sections := []string{
		model.filterInput.View(),
		"",
		model.contentList.View(),
	}

	contents := lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)

	return appStyle.Render(contents)
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height

	model.filterInput = model.filterInput.Resize(helpers.GetMaxInt(0, width-2*padding))

	// TODO this is a mess; figure out a better way to do this (maybe the parent keeps track of the size of the children??)
	model.contentList = model.contentList.Resize(
		helpers.GetMaxInt(0, width-2*padding),
		helpers.GetMaxInt(0, height-4),
	)

	return model
}
