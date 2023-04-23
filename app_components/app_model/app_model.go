package app_model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mieubrisse/cli-journal-go/app_components/entry_item"
	"github.com/mieubrisse/cli-journal-go/app_components/entry_list"
	"github.com/mieubrisse/cli-journal-go/app_components/filter_pane"
	"github.com/mieubrisse/cli-journal-go/app_components/new_entry_form"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list"
	"github.com/mieubrisse/cli-journal-go/components/filterable_list_item"
	"github.com/mieubrisse/cli-journal-go/global_styles"
	"github.com/mieubrisse/cli-journal-go/helpers"
	"github.com/mieubrisse/vim-bubble/vim"
	"github.com/sahilm/fuzzy"
	"sort"
)

const (
	maxCreateContentModalWidth  = 50
	maxCreateContentModalHeight = 3

	filterPaneHeight = 6
)

// "Constants"
var horizontalPadThresholdsByTerminalWidth = map[int]int{
	0:   0,
	60:  1,
	120: 2,
}
var verticalPadThresholdsByTerminalHeight = map[int]int{
	0:  0,
	40: 1,
}
var filtersLabelLine = lipgloss.NewStyle().
	Foreground(global_styles.Cyan).
	Bold(true).
	Render("FILTERS")

type Model struct {
	createContentForm new_entry_form.Component

	filterPane filter_pane.Model

	filterTabCompletionPane filterable_list.Component[filterable_list_item.Component]

	contentList entry_list.Model

	tags []string

	height int
	width  int
}

func New(
	content []entry_item.Component,
) Model {
	createContentForm := new_entry_form.New()

	contentList := entry_list.New(content)
	contentList.Focus()

	filterPane := filter_pane.New()

	deduplicatedTags := make(map[string]bool, 0)
	for _, contentItem := range content {
		for _, tag := range contentItem.GetTags() {
			deduplicatedTags[tag] = true
		}
	}

	sortedTags := make([]string, 0, len(deduplicatedTags))
	for tag := range deduplicatedTags {
		sortedTags = append(sortedTags, tag)
	}
	sort.Strings(sortedTags)

	completionPane := filterable_list.New[filterable_list_item.Component]()

	return Model{
		createContentForm:       createContentForm,
		filterPane:              filterPane,
		filterTabCompletionPane: completionPane,
		contentList:             contentList,
		height:                  0,
		width:                   0,
		tags:                    sortedTags,
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
			case "\\":
				cmds := make([]tea.Cmd, 0)
				cmds = append(cmds, model.contentList.Blur())

				cmds = append(cmds, model.filterPane.Focus())
				cmds = append(cmds, model.filterTabCompletionPane.Focus())
				model.filterPane.SetMode(vim.InsertMode)

				return model, tea.Batch(cmds...)
			case "c":
				// Clear all filters
				model.filterPane.Clear()
				nameFilterLines, tagFilterLines := model.filterPane.GetFilterLines()
				model.contentList.SetFilters(nameFilterLines, tagFilterLines)
				model.filterTabCompletionPane.SetItems([]filterable_list_item.Component{})

				return model, nil
			case "n":
				cmds := make([]tea.Cmd, 0)
				cmds = append(cmds, model.contentList.Blur())
				cmds = append(cmds, model.createContentForm.Focus())
				return model, tea.Batch(cmds...)
			}

			cmd := model.contentList.Update(msg)
			return model, cmd

		} else if model.filterPane.Focused() {
			switch msg.String() {
			case "\\":
				// TODO switch to by-value
				model.filterPane.Blur()
				model.filterTabCompletionPane.Blur()

				cmd := model.contentList.Focus()
				return model, cmd
			case "ctrl+j":
				model.filterTabCompletionPane.Scroll(1)
				return model, nil
			case "ctrl+k":
				model.filterTabCompletionPane.Scroll(-1)
				return model, nil
			}

			var cmd tea.Cmd
			if msg.String() == "tab" {
				// The user is tab-completing
				filteredItemIndices := model.filterTabCompletionPane.GetFilteredItemIndices()
				if len(filteredItemIndices) > 0 {
					highlightedCompletionIdxInFilteredList := model.filterTabCompletionPane.GetHighlightedItemIndex()
					highlightedCompletionIdxInOriginalList := filteredItemIndices[highlightedCompletionIdxInFilteredList]
					selectedCompletion := model.filterTabCompletionPane.GetItems()[highlightedCompletionIdxInOriginalList]
					model.filterPane.ReplaceCurrentFilter(selectedCompletion.GetValue(), true)
				}
			} else {
				cmd = model.filterPane.Update(msg)
			}

			// Make sure to let the content list know about the changes
			nameFilterList, tagFilterList := model.filterPane.GetFilterLines()
			model.contentList.SetFilters(nameFilterList, tagFilterList)

			// Update the tab-contents pane with changes, displaying nothing if the line isn't a tag filter line
			filterText, isTagFilter := model.filterPane.GetCurrentFilter()
			tabCompletionItems := make([]filterable_list_item.Component, 0)
			if isTagFilter {
				if len(filterText) > 0 {
					matches := fuzzy.Find(filterText, model.tags)

					tabCompletionItems = make([]filterable_list_item.Component, len(matches))
					for idx, match := range matches {
						tabCompletionItems[idx] = filterable_list_item.New(model.tags[match.Index])
					}
				} else {
					tabCompletionItems = make([]filterable_list_item.Component, len(model.tags))
					for idx, tag := range model.tags {
						tabCompletionItems[idx] = filterable_list_item.New(tag)
					}
				}
			}
			model.filterTabCompletionPane.SetItems(tabCompletionItems)

			return model, cmd
		} else if model.createContentForm.Focused() {
			switch msg.String() {
			case "esc":
				// Back out of the create content modal
				model.createContentForm.Clear()

				cmds := make([]tea.Cmd, 0)
				cmds = append(cmds, model.createContentForm.Blur())
				cmds = append(cmds, model.contentList.Focus())
				return model, tea.Batch(cmds...)
			case "enter":
				// TODO reenable
				/*
					// Create the new content piece
					content := content_item.ContentItem{
						Timestamp: time.Now(),
						Name:      model.createContentForm.GetNameValue(),
						Tags:      []string{},
					}
					model.contentList.AddItem(content)
				*/

				model.createContentForm.Clear()

				cmds := make([]tea.Cmd, 0)
				cmds = append(cmds, model.createContentForm.Blur())
				cmds = append(cmds, model.contentList.Focus())
				return model, tea.Batch(cmds...)
			}

			cmd := model.createContentForm.Update(msg)
			return model, cmd
		}
	case tea.WindowSizeMsg:
		return model.Resize(msg.Width, msg.Height), nil
	}

	return model, nil
}

func (model Model) View() string {
	filterView := lipgloss.JoinHorizontal(
		lipgloss.Center,
		model.filterPane.View(),
		model.filterTabCompletionPane.View(),
	)

	sections := []string{
		model.contentList.View(),
		filtersLabelLine,
		filterView,
	}

	contents := lipgloss.JoinVertical(lipgloss.Left, sections...)

	horizontalPad, verticalPad := getPadsForSize(model.width, model.height)
	result := lipgloss.NewStyle().
		Width(model.width).
		Height(model.height).
		MaxWidth(model.width).
		MaxHeight(model.height).
		Padding(verticalPad, horizontalPad, verticalPad, horizontalPad).
		Render(contents)

	if model.createContentForm.Focused() {
		createContentFormStr := model.createContentForm.View()

		createContentFormStr = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			Render(createContentFormStr)

		result = helpers.OverlayString(result, createContentFormStr)
	}

	return result
}

func (model Model) Resize(width int, height int) Model {
	model.width = width
	model.height = height

	horizontalPad, verticalPad := getPadsForSize(model.width, model.height)
	displaySpaceWidth := helpers.GetMaxInt(0, model.width-2*horizontalPad)
	displaySpaceHeight := helpers.GetMaxInt(0, model.height-2*verticalPad)

	filterPaneWidth := int(0.5 * float64(displaySpaceWidth))
	model.filterPane.Resize(filterPaneWidth, filterPaneHeight)

	completionPaneWidth := displaySpaceWidth - filterPaneWidth
	model.filterTabCompletionPane.Resize(completionPaneWidth, filterPaneHeight)

	// Leave one blank line for filters label
	contentListHeight := helpers.GetMaxInt(0, displaySpaceHeight-filterPaneHeight-1)

	model.contentList.Resize(displaySpaceWidth, contentListHeight)

	createContentModalWidth := helpers.GetMinInt(model.width, maxCreateContentModalWidth)
	createContentModalHeight := helpers.GetMinInt(model.height, maxCreateContentModalHeight)
	model.createContentForm.Resize(createContentModalWidth, createContentModalHeight)

	return model
}

// =================================== Private Helper Functions ===================================
func getPadsForSize(width int, height int) (int, int) {
	actualHorizontalPad := 0
	for threshold, trialHorizontalPad := range horizontalPadThresholdsByTerminalWidth {
		if width > threshold && actualHorizontalPad < trialHorizontalPad {
			actualHorizontalPad = trialHorizontalPad
		}

	}

	actualVerticalPad := 0
	for threshold, trialVerticalPad := range verticalPadThresholdsByTerminalHeight {
		if height > threshold && actualVerticalPad < trialVerticalPad {
			actualVerticalPad = trialVerticalPad
		}
	}

	return actualHorizontalPad, actualVerticalPad
}

func clampInt(value int, min int, max int) int {
	if max < min {
		max, min = min, max
	}

	return helpers.GetMaxInt(min, helpers.GetMinInt(max, value))
}
