package components

import (
	"github.com/charmbracelet/lipgloss"
)

// Loader styles
var (
	spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	loaderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
)

// Loader represents a loading spinner.
type Loader struct {
	frame int
	text  string
}

// NewLoader creates a new loader.
func NewLoader(text string) *Loader {
	return &Loader{
		frame: 0,
		text:  text,
	}
}

// Tick advances the spinner to the next frame.
func (l *Loader) Tick() {
	l.frame = (l.frame + 1) % len(spinnerFrames)
}

// Render renders the loader.
func (l *Loader) Render() string {
	spinner := loaderStyle.Render(spinnerFrames[l.frame])
	return spinner + " " + l.text
}

// RenderLoader renders a static loading message.
func RenderLoader(text string) string {
	spinner := loaderStyle.Render(spinnerFrames[0])
	return spinner + " " + text
}
