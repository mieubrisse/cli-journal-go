package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/components/text_filter_input"
	"github.com/mieubrisse/cli-journal-go/helpers"
)

type Mode int

const (
	padding = 2
)

var appStyle = lipgloss.NewStyle().Padding(padding)

type Model struct {
	mode Mode

	nameFilterInput text_filter_input.Model

	tagFilterInput text_filter_input.Model

	contentList filterable_content_list.Model

	height int
	width  int
}

func New(
	nameFilterInput text_filter_input.Model,
	tagFilterInput text_filter_input.Model,
	contentList filterable_content_list.Model,
) Model {
	return Model{
		nameFilterInput: nameFilterInput,
		tagFilterInput:  tagFilterInput,
		contentList:     contentList,
		height:          0,
		width:           0,
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
				model.contentList.Blur()

				// This will tell the input that it should display the cursor
				cmd := model.nameFilterInput.Focus()

				return model, cmd
			case "#":
				model.contentList.Blur()

				// This will tell the input that it should display the cursor
				cmd := model.tagFilterInput.Focus()

				return model, cmd
			case "c":
				model.nameFilterInput.SetValue("")
				model.tagFilterInput.SetValue("")

				// Need to tell the content list that we changed
				model.contentList.UpdateFilters(
					model.nameFilterInput.Value(),
					model.tagFilterInput.Value(),
				)
			}

			var cmd tea.Cmd
			model.contentList, cmd = model.contentList.Update(msg)
			return model, cmd

		} else if model.nameFilterInput.Focused() {
			// Back out of filter mode
			if msg.String() == "esc" || msg.String() == "enter" {
				model.nameFilterInput.Blur()
				model.contentList.Focus()
				return model, nil
			}

			var cmd tea.Cmd
			model.nameFilterInput, cmd = model.nameFilterInput.Update(msg)

			// Make sure to tell the content list about the new filter update
			model.contentList.UpdateFilters(
				model.nameFilterInput.Value(),
				model.tagFilterInput.Value(),
			)

			return model, cmd
		} else if model.tagFilterInput.Focused() {
			// Back out of filter mode
			if msg.String() == "esc" || msg.String() == "enter" {
				model.tagFilterInput.Blur()
				model.contentList.Focus()
				return model, nil
			}

			var cmd tea.Cmd
			model.tagFilterInput, cmd = model.tagFilterInput.Update(msg)

			// Make sure to tell the content list about the new filter update
			model.contentList.UpdateFilters(
				model.nameFilterInput.Value(),
				model.tagFilterInput.Value(),
			)

			return model, cmd
		}
	case tea.WindowSizeMsg:
		return model.Resize(msg.Width, msg.Height), nil
	}

	return model, nil
}

func (model Model) View() string {
	sections := []string{
		model.contentList.View(),
		model.nameFilterInput.View(),
		model.tagFilterInput.View(),
	}

	contents := lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)

	return appStyle.Copy().
		Width(model.width).
		Height(model.height).
		Render(contents)
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height

	componentSpaceWidth := helpers.GetMaxInt(0, model.width-2*padding)
	componentSpaceHeight := helpers.GetMaxInt(0, model.height-2*padding)

	model.nameFilterInput = model.nameFilterInput.Resize(componentSpaceWidth, 1)
	model.tagFilterInput = model.tagFilterInput.Resize(componentSpaceWidth, 1)

	contentListHeight := helpers.GetMaxInt(0, componentSpaceHeight-model.nameFilterInput.GetHeight()-model.tagFilterInput.GetHeight())

	model.contentList = model.contentList.Resize(componentSpaceWidth, contentListHeight)

	return model
}
