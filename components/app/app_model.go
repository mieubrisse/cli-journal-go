package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/components/filterable_content_list"
	"github.com/mieubrisse/cli-journal-go/components/form"
	"github.com/mieubrisse/cli-journal-go/components/text_input"
	"github.com/mieubrisse/cli-journal-go/data_structures/content_item"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"time"
)

const (
	maxCreateContentModalWidth  = 50
	maxCreateContentModalHeight = 3

	shouldDimBackgroundForModals = false
)

// "Constants"
var horizontalPadThresholds = map[int]int{
	0:   0,
	60:  1,
	120: 2,
}
var verticalPadThresholds = map[int]int{
	0:  0,
	40: 1,
}

type Model struct {
	createContentForm form.Model

	nameFilterInput text_input.Model

	tagFilterInput text_input.Model

	contentList filterable_content_list.Model

	height int
	width  int
}

func New(
	createContentForm form.Model,
	nameFilterInput text_input.Model,
	tagFilterInput text_input.Model,
	contentList filterable_content_list.Model,
) Model {
	return Model{
		createContentForm: createContentForm,
		nameFilterInput:   nameFilterInput,
		tagFilterInput:    tagFilterInput,
		contentList:       contentList,
		height:            0,
		width:             0,
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
			case "esc":
				// Clear all filters
				model.nameFilterInput.SetValue("")
				model.tagFilterInput.SetValue("")

				// Need to tell the content list that we changed
				model.contentList.SetFilterTexts(
					model.nameFilterInput.Value(),
					model.tagFilterInput.Value(),
				)
			case "c":
				model.contentList.Blur()
				model.createContentForm.Focus()
			}

			var cmd tea.Cmd
			model.contentList, cmd = model.contentList.Update(msg)
			return model, cmd

		} else if model.nameFilterInput.Focused() {
			switch msg.String() {
			case "esc":
				// User is clearing the filter

				// Revert to the previous filter value, and update the content list
				model.nameFilterInput.SetValue("")
				model.contentList.SetFilterTexts(
					model.nameFilterInput.Value(),
					model.tagFilterInput.Value(),
				)

				model.nameFilterInput.Blur()
				model.contentList.Focus()
				return model, nil
			case "enter":
				// User is exiting the filtering mode, persisting their changes
				model.nameFilterInput.Blur()
				model.contentList.Focus()
				return model, nil
			}

			var cmd tea.Cmd
			model.nameFilterInput, cmd = model.nameFilterInput.Update(msg)

			// Make sure to tell the content list about the new filter update
			model.contentList.SetFilterTexts(
				model.nameFilterInput.Value(),
				model.tagFilterInput.Value(),
			)

			return model, cmd
		} else if model.tagFilterInput.Focused() {
			switch msg.String() {
			case "esc":
				// User is clearing the filter

				// Revert to the previous filter value, and update the content list
				model.tagFilterInput.SetValue("")
				model.contentList.SetFilterTexts(
					model.tagFilterInput.Value(),
					model.tagFilterInput.Value(),
				)

				model.tagFilterInput.Blur()
				model.contentList.Focus()
				return model, nil
			case "enter":
				// User is exiting the filtering mode, persisting their changes
				model.tagFilterInput.Blur()
				model.contentList.Focus()
				return model, nil
			}

			var cmd tea.Cmd
			model.tagFilterInput, cmd = model.tagFilterInput.Update(msg)

			// Make sure to tell the content list about the new filter update
			model.contentList.SetFilterTexts(
				model.nameFilterInput.Value(),
				model.tagFilterInput.Value(),
			)

			return model, cmd
		} else if model.createContentForm.Focused() {
			switch msg.String() {
			case "esc":
				// Back out of the create content modal
				model.createContentForm.SetValue("")
				model.createContentForm.Blur()
				model.contentList.Focus()
				return model, nil
			case "enter":
				// Create the new content piece
				content := content_item.ContentItem{
					Timestamp: time.Now(),
					Name:      model.createContentForm.GetValue(),
					Tags:      []string{},
				}
				model.contentList = model.contentList.AddItem(content)

				model.createContentForm.SetValue("")
				model.createContentForm.Blur()
				model.contentList.Focus()

				return model, nil
			}

			var cmd tea.Cmd
			model.createContentForm, cmd = model.createContentForm.Update(msg)
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

	contents := lipgloss.JoinVertical(lipgloss.Left, sections...)

	horizontalPad, verticalPad := getPadsForSize(model.width, model.height)
	result := lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		Padding(verticalPad, horizontalPad, verticalPad, horizontalPad).
		Render(contents)

	if model.createContentForm.Focused() {
		createContentFormStr := model.createContentForm.View()
		/*
			createContentModalStr = lipgloss.NewStyle().
				Border()
		*/

		createContentFormStr = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Render(createContentFormStr)

		result = helpers.OverlayString(result, createContentFormStr, shouldDimBackgroundForModals)
	}

	return result
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height

	horizontalPad, verticalPad := getPadsForSize(model.width, model.height)
	componentSpaceWidth := helpers.GetMaxInt(0, model.width-2*horizontalPad)
	componentSpaceHeight := helpers.GetMaxInt(0, model.height-2*verticalPad)

	model.nameFilterInput = model.nameFilterInput.Resize(componentSpaceWidth, 1)
	model.tagFilterInput = model.tagFilterInput.Resize(componentSpaceWidth, 1)

	contentListHeight := helpers.GetMaxInt(
		0,
		componentSpaceHeight-model.nameFilterInput.GetHeight()-model.tagFilterInput.GetHeight(),
	)

	model.contentList = model.contentList.Resize(componentSpaceWidth, contentListHeight)

	createContentModalWidth := helpers.GetMinInt(model.width, maxCreateContentModalWidth)
	createContentModalHeight := helpers.GetMinInt(model.height, maxCreateContentModalHeight)
	model.createContentForm = model.createContentForm.Resize(createContentModalWidth, createContentModalHeight)

	return model
}

func getPadsForSize(width int, height int) (int, int) {
	actualHorizontalPad := 0
	for threshold, trialHorizontalPad := range horizontalPadThresholds {
		if width > threshold && actualHorizontalPad < trialHorizontalPad {
			actualHorizontalPad = trialHorizontalPad
		}

	}

	actualVerticalPad := 0
	for threshold, trialVerticalPad := range verticalPadThresholds {
		if height > threshold && actualVerticalPad < trialVerticalPad {
			actualVerticalPad = trialVerticalPad
		}
	}

	return actualHorizontalPad, actualVerticalPad
}
