package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Colors for progress bar
var (
	filledColor = lipgloss.Color("#22C55E")
	emptyColor  = lipgloss.Color("#374151")
)

// ProgressBar creates a text-based progress bar.
func ProgressBar(width int, progress float64) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filled := int(progress * float64(width))
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(filledColor)
	emptyStyle := lipgloss.NewStyle().Foreground(emptyColor)

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	return bar
}

// ProgressBarWithLabel creates a progress bar with a label.
func ProgressBarWithLabel(width int, progress float64, label string) string {
	bar := ProgressBar(width, progress)
	return lipgloss.JoinHorizontal(lipgloss.Center, "[", bar, "] ", label)
}

// ProgressBarColored creates a colored progress bar.
func ProgressBarColored(width int, progress float64, color lipgloss.Color) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filled := int(progress * float64(width))
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(emptyColor)

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	return bar
}
