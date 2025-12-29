package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for the game
var (
	// Primary colors
	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
	ColorSecondary = lipgloss.Color("#10B981") // Green
	ColorAccent    = lipgloss.Color("#F59E0B") // Amber
	
	// Status colors
	ColorSuccess = lipgloss.Color("#22C55E") // Green
	ColorWarning = lipgloss.Color("#EAB308") // Yellow
	ColorError   = lipgloss.Color("#EF4444") // Red
	ColorInfo    = lipgloss.Color("#3B82F6") // Blue
	
	// Element colors
	ColorFire    = lipgloss.Color("#EF4444") // Red
	ColorIce     = lipgloss.Color("#06B6D4") // Cyan
	ColorThunder = lipgloss.Color("#FACC15") // Yellow
	ColorArcane  = lipgloss.Color("#A855F7") // Purple
	
	// UI colors
	ColorBorder    = lipgloss.Color("#6B7280") // Gray
	ColorBorderDim = lipgloss.Color("#374151") // Darker gray
	ColorText      = lipgloss.Color("#F9FAFB") // Almost white
	ColorTextDim   = lipgloss.Color("#9CA3AF") // Gray text
	ColorBg        = lipgloss.Color("#111827") // Dark background
	ColorBgAlt     = lipgloss.Color("#1F2937") // Slightly lighter bg
)

// Base styles
var (
	// Title style
	TitleStyle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		MarginBottom(1)

	// Subtitle style
	SubtitleStyle = lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	// Normal text
	TextStyle = lipgloss.NewStyle().
		Foreground(ColorText)

	// Dimmed text
	DimStyle = lipgloss.NewStyle().
		Foreground(ColorTextDim)

	// Highlight style
	HighlightStyle = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
		Foreground(ColorError)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
		Foreground(ColorSuccess)

	// Warning style
	WarningStyle = lipgloss.NewStyle().
		Foreground(ColorWarning)

	// Selected item style
	SelectedStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(ColorBgAlt).
		Bold(true)
)

// Box styles
var (
	// Main container
	ContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2)

	// Header box
	HeaderStyle = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 2).
		Align(lipgloss.Center)

	// Section style
	SectionStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(ColorBorderDim).
		Padding(0, 1)

	// Footer/help style
	FooterStyle = lipgloss.NewStyle().
		Foreground(ColorTextDim).
		MarginTop(1)
)

// Element-specific styles
var (
	FireStyle = lipgloss.NewStyle().
		Foreground(ColorFire)

	IceStyle = lipgloss.NewStyle().
		Foreground(ColorIce)

	ThunderStyle = lipgloss.NewStyle().
		Foreground(ColorThunder)

	ArcaneStyle = lipgloss.NewStyle().
		Foreground(ColorArcane)
)

// GetElementStyle returns the style for an element type.
func GetElementStyle(element string) lipgloss.Style {
	switch element {
	case "fire":
		return FireStyle
	case "ice":
		return IceStyle
	case "thunder":
		return ThunderStyle
	case "arcane":
		return ArcaneStyle
	default:
		return TextStyle
	}
}

// GetElementIcon returns an emoji/icon for an element type.
func GetElementIcon(element string) string {
	switch element {
	case "fire":
		return "🔥"
	case "ice":
		return "❄️"
	case "thunder":
		return "⚡"
	case "arcane":
		return "✨"
	default:
		return "•"
	}
}

// Progress bar styles
var (
	ProgressBarFilled = lipgloss.NewStyle().
		Foreground(ColorSuccess)

	ProgressBarEmpty = lipgloss.NewStyle().
		Foreground(ColorBorderDim)
)
