package helpers

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func GetMaxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetMinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func OverlayString(background string, overlay string) string {
	backgroundHeight := lipgloss.Height(background)
	backgroundWidth := lipgloss.Width(background)

	overlayHeight := lipgloss.Height(overlay)
	overlayWidth := lipgloss.Width(overlay)

	firstOverlaidLineIdx := (backgroundHeight / 2) - (overlayHeight / 2) + 1

	// The index of the first character that will be replaced
	cutpointIdx := (backgroundWidth / 2) - (overlayWidth / 2)

	// The index of the character that switches back to being background again
	resumepointIdx := backgroundWidth - cutpointIdx

	overlayLines := strings.Split(overlay, "\n")
	backgroundLines := strings.Split(background, "\n")

	resultLines := []string{}
	for idx, backgroundLine := range backgroundLines {
		// No more overlaying to do
		if idx >= firstOverlaidLineIdx+overlayHeight {
			break
		}

		if idx < firstOverlaidLineIdx {
			// TODO make it faint?
			resultLines = append(resultLines, backgroundLine)
			continue
		}

		overlayLine := overlayLines[idx-firstOverlaidLineIdx]

		resultLine := backgroundLine[:cutpointIdx] + overlayLine + backgroundLine[resumepointIdx:]
		resultLines = append(resultLines, resultLine)
	}

	return strings.Join(resultLines, "\n")
}
