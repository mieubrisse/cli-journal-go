package helpers

import (
	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var overlayBackgroundStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))

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

// TODO FIX BUG IN THIS - we can't just overlay using string length, because non-printing chars mess this up (e.g. color chars)
func OverlayString(background string, overlay string, shouldDimBackground bool) string {
	backgroundHeight := lipgloss.Height(background)
	backgroundWidth := lipgloss.Width(background)

	overlayHeight := lipgloss.Height(overlay)
	overlayWidth := lipgloss.Width(overlay)

	// The index of the first line that will suffer replacement
	firstOverlaidLineIdx := (backgroundHeight / 2) - (overlayHeight / 2)

	// The index of the line that switches back to being background again
	resumeLineIdx := firstOverlaidLineIdx + overlayHeight

	// The index of the first character that will be replaced
	cutpointIdx := (backgroundWidth / 2) - (overlayWidth / 2)

	// The index of the character that switches back to being background again
	resumepointIdx := backgroundWidth - cutpointIdx

	overlayLines := strings.Split(overlay, "\n")
	backgroundLines := strings.Split(background, "\n")

	resultLines := []string{}
	for idx, backgroundLine := range backgroundLines {
		if idx < firstOverlaidLineIdx || idx >= resumeLineIdx {
			resultLines = append(resultLines, backgroundLine)
			continue
		}

		overlayLine := overlayLines[idx-firstOverlaidLineIdx]

		resultLine := backgroundLine[:cutpointIdx] + overlayLine + backgroundLine[resumepointIdx:]
		resultLines = append(resultLines, resultLine)
	}

	return strings.Join(resultLines, "\n")
}

// TODO use this, somehow, to do background-dimming
// Removes all non-graphic characters
func ClearFormatting(input string) string {
	return stripansi.Strip(input)
}
