package components

import (
	"fmt"

	"github.com/Ltorre/ManaTTY/models"
	"github.com/Ltorre/ManaTTY/ui"
	"github.com/Ltorre/ManaTTY/utils"
	"github.com/charmbracelet/lipgloss"
)

// Spell list styles
var (
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F9FAFB")).
			Background(lipgloss.Color("#1F2937"))

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9FAFB"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	readyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22C55E")).
			Bold(true)

	cooldownStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EAB308"))
)

// SpellListItem renders a single spell in a list.
func SpellListItem(spell *models.Spell, selected bool, showDetails bool) string {
	style := normalStyle
	if selected {
		style = selectedStyle
	}

	// Element icon
	icon := GetElementIcon(string(spell.Element))

	// Cooldown status
	var status string
	if spell.IsReady() {
		status = readyStyle.Render("Ready!")
	} else {
		status = cooldownStyle.Render(utils.FormatCooldown(spell.CooldownRemainingMs))
	}

	// Build the line
	prefix := "  "
	if selected {
		prefix = "> "
	}

	line := fmt.Sprintf("%s%s %s (Lv%d) - %s",
		prefix, icon, spell.Name, spell.Level, status)

	if showDetails {
		details := dimStyle.Render(fmt.Sprintf(
			"      Element: %s | Cooldown: %s | Casts: %d",
			spell.Element,
			utils.FormatMilliseconds(spell.BaseCooldownMs),
			spell.CastCount,
		))
		return style.Render(line) + "\n" + details
	}

	return style.Render(line)
}

// SpellList renders a list of spells.
func SpellList(spells []*models.Spell, selectedIndex int, showDetails bool) string {
	if len(spells) == 0 {
		return dimStyle.Render("  No spells unlocked")
	}

	var lines []string
	for i, spell := range spells {
		lines = append(lines, SpellListItem(spell, i == selectedIndex, showDetails))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// GetElementIcon returns an icon for an element.
// Uses ASCII fallbacks on Windows for compatibility.
func GetElementIcon(element string) string {
	sym := ui.GetSymbols()
	switch element {
	case "fire":
		return sym.Fire
	case "ice":
		return sym.Ice
	case "thunder":
		return sym.Thunder
	case "arcane":
		return sym.Arcane
	default:
		return sym.Default
	}
}

// GetElementColor returns a color for an element.
func GetElementColor(element string) lipgloss.Color {
	switch element {
	case "fire":
		return lipgloss.Color("#EF4444")
	case "ice":
		return lipgloss.Color("#06B6D4")
	case "thunder":
		return lipgloss.Color("#FACC15")
	case "arcane":
		return lipgloss.Color("#A855F7")
	default:
		return lipgloss.Color("#F9FAFB")
	}
}
