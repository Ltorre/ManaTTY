package ui

import (
	"fmt"
	"strings"

	"github.com/Ltorre/ManaTTY/utils"
	"github.com/charmbracelet/lipgloss"
)

// View renders the current view.
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Build view based on current screen
	var content string
	switch m.currentView {
	case ViewTower:
		content = m.viewTower()
	case ViewSpells:
		content = m.viewSpells()
	case ViewRituals:
		content = m.viewRituals()
	case ViewStats:
		content = m.viewStats()
	case ViewPrestige:
		content = m.viewPrestige()
	case ViewMenu:
		content = m.viewMenu()
	default:
		content = m.viewTower()
	}

	// Add notification if present
	if m.notification != "" {
		notification := lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Render("📢 " + m.notification)
		content = lipgloss.JoinVertical(lipgloss.Top, content, "", notification)
	}

	// Add confirmation dialog if active
	if m.confirming {
		confirm := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorWarning).
			Padding(1, 2).
			Render(m.confirmText)
		content = lipgloss.JoinVertical(lipgloss.Top, content, "", confirm)
	}

	return content
}

// viewTower renders the main tower view.
func (m Model) viewTower() string {
	if m.gameState == nil {
		return "No game loaded"
	}

	gs := m.gameState
	var lines []string

	// Header
	era := ""
	if gs.PrestigeData.CurrentEra > 0 {
		era = fmt.Sprintf(" - ERA %d", gs.PrestigeData.CurrentEra)
	}
	header := HeaderStyle.Width(60).Render(
		TitleStyle.Render("🏰 MAGE TOWER ASCENSION" + era),
	)
	lines = append(lines, header)
	lines = append(lines, "")

	// Floor display
	floorStr := fmt.Sprintf("Current Floor: %d", gs.Tower.CurrentFloor)
	if gs.Tower.MaxFloorReached > gs.Tower.CurrentFloor {
		floorStr += fmt.Sprintf(" (Max: %d)", gs.Tower.MaxFloorReached)
	}
	lines = append(lines, SubtitleStyle.Render(floorStr))
	lines = append(lines, "")

	// Mana progress bar
	progress := gs.Tower.GetFloorProgress()
	barWidth := 40
	filled := int(progress * float64(barWidth))
	bar := ProgressBarFilled.Render(strings.Repeat("█", filled)) +
		ProgressBarEmpty.Render(strings.Repeat("░", barWidth-filled))

	percentage := int(progress * 100)
	manaStr := fmt.Sprintf("[%s] %d%% (%s / %s)",
		bar,
		percentage,
		utils.FormatNumber(gs.Tower.CurrentMana),
		utils.FormatNumber(gs.Tower.MaxMana),
	)
	lines = append(lines, manaStr)
	lines = append(lines, "")

	// Stats section
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	lines = append(lines, SubtitleStyle.Render("📊 STATS"))

	if m.engine != nil {
		manaPerSec := m.engine.CalculateManaPerSecond(gs)
		multiplier := m.engine.GetTotalMultiplier(gs)

		lines = append(lines, fmt.Sprintf("  Mana/sec:    %s (%sx multiplier)",
			HighlightStyle.Render(utils.FormatNumber(manaPerSec)),
			utils.FormatMultiplier(multiplier),
		))
	}

	lines = append(lines, fmt.Sprintf("  Total earned: %s",
		utils.FormatNumber(gs.Tower.LifetimeManaEarned),
	))
	lines = append(lines, fmt.Sprintf("  Spells: %d | Rituals: %d/%d",
		len(gs.Spells),
		len(gs.GetActiveRituals()),
		gs.PrestigeData.RitualCapacity,
	))

	// Auto-cast status
	autoCastStatus := "ON"
	if !gs.Session.AutoCastEnabled {
		autoCastStatus = "OFF"
	}
	lines = append(lines, fmt.Sprintf("  Auto-cast: %s", autoCastStatus))

	// Element synergy status (use GameState methods via m.gameState)
	if m.gameState.HasActiveSynergy() {
		element := m.gameState.GetActiveSynergy()
		remaining := m.gameState.GetSynergyTimeRemaining() / 1000 // convert to seconds
		synergyStr := fmt.Sprintf("  %s Synergy: %ds remaining (20%% bonus)",
			string(element), remaining)
		lines = append(lines, HighlightStyle.Render(synergyStr))
	} else if len(gs.Session.LastCastElements) > 0 {
		// Show element streak progress
		streakLen := len(gs.Session.LastCastElements)
		if streakLen > 0 {
			lastElement := gs.Session.LastCastElements[streakLen-1]
			lines = append(lines, DimStyle.Render(fmt.Sprintf("  Element streak: %s ×%d/3",
				string(lastElement), streakLen)))
		}
	}
	lines = append(lines, "")

	// Active rituals
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	activeRituals := gs.GetActiveRituals()
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf("🔥 ACTIVE RITUALS (%d/%d)",
		len(activeRituals), gs.PrestigeData.RitualCapacity)))

	if len(activeRituals) == 0 {
		lines = append(lines, DimStyle.Render("  No active rituals"))
	} else {
		for i, ritual := range activeRituals {
			cooldownStr := "Ready"
			if ritual.CooldownRemaining > 0 {
				cooldownStr = utils.FormatCooldown(ritual.CooldownRemaining)
			}
			lines = append(lines, fmt.Sprintf("  [%d] %s (%s)",
				i+1, ritual.Name, cooldownStr))
		}
	}
	lines = append(lines, "")

	// Footer
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	footer := FooterStyle.Render("[S] Spells  [R] Rituals  [T] Stats  [P] Prestige  [A] Auto-cast  [Q] Quit")
	lines = append(lines, footer)

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewMenu renders the menu view.
func (m Model) viewMenu() string {
	var lines []string

	header := HeaderStyle.Width(50).Render(
		TitleStyle.Render("📋 MENU"),
	)
	lines = append(lines, header)
	lines = append(lines, "")

	lines = append(lines, TextStyle.Render("  [S] Save Game"))
	lines = append(lines, TextStyle.Render("  [Q] Save & Quit"))
	lines = append(lines, "")
	lines = append(lines, FooterStyle.Render("[B/Esc] Back"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewSpells renders the spells view.
func (m Model) viewSpells() string {
	if m.gameState == nil {
		return "No game loaded"
	}

	var lines []string

	header := HeaderStyle.Width(70).Render(
		TitleStyle.Render("📜 SPELLS"),
	)
	lines = append(lines, header)
	lines = append(lines, "")

	// Element Synergy Status
	if m.gameState.HasActiveSynergy() {
		synergy := m.gameState.GetActiveSynergy()
		remaining := m.gameState.GetSynergyTimeRemaining() / 1000
		icon := GetElementIcon(string(synergy))
		lines = append(lines, SuccessStyle.Render(fmt.Sprintf("🔥 %s SYNERGY ACTIVE! +20%% bonus (%ds remaining)", icon, remaining)))
		lines = append(lines, "")
	}

	// Auto-Cast Loadout Panel
	usedSlots := len(m.gameState.Session.AutoCastSlots)
	maxSlots := m.gameState.GetAutoCastSlotCount()
	autoCastStatus := "OFF"
	if m.gameState.Session.AutoCastEnabled {
		autoCastStatus = "ON"
	}
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf("⚡ Auto-Cast Loadout [%s] (%d/%d slots)", autoCastStatus, usedSlots, maxSlots)))
	
	if len(m.gameState.Session.AutoCastSlots) == 0 {
		lines = append(lines, DimStyle.Render("  (empty - press Space on a spell to add)"))
	} else {
		for i, spellID := range m.gameState.Session.AutoCastSlots {
			spell := m.gameState.GetSpellByID(spellID)
			if spell != nil {
				icon := GetElementIcon(string(spell.Element))
				lines = append(lines, TextStyle.Render(fmt.Sprintf("  %d. %s %s (Lv%d)", i+1, icon, spell.Name, spell.Level)))
			}
		}
	}
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	lines = append(lines, "")

	// Spell list
	lines = append(lines, SubtitleStyle.Render("All Spells"))
	for i, spell := range m.gameState.Spells {
		selected := i == m.selectedIndex
		icon := GetElementIcon(string(spell.Element))

		prefix := "  "
		if selected {
			prefix = "> "
		}

		// Auto-cast indicator with position
		autoIndicator := "   "
		for idx, slotID := range m.gameState.Session.AutoCastSlots {
			if slotID == spell.ID {
				autoIndicator = fmt.Sprintf("[%d]", idx+1)
				break
			}
		}

		// Status
		status := SuccessStyle.Render("Ready")
		if !spell.IsReady() {
			status = WarningStyle.Render(utils.FormatCooldown(spell.CooldownRemainingMs))
		}

		// Level indicator with max check
		levelStr := fmt.Sprintf("Lv%d", spell.Level)
		if spell.Level >= 10 {
			levelStr = "MAX"
		}

		line := fmt.Sprintf("%s%s %s %s (%s) - %s",
			prefix, autoIndicator, icon, spell.Name, levelStr, status)

		if selected {
			lines = append(lines, SelectedStyle.Render(line))
		} else if m.gameState.IsSpellInAutoCast(spell.ID) {
			lines = append(lines, HighlightStyle.Render(line))
		} else {
			lines = append(lines, TextStyle.Render(line))
		}

		// Show details for selected spell
		if selected && m.engine != nil {
			stats := m.engine.GetSpellEffectiveStats(m.gameState, spell)
			upgradeStr := fmt.Sprintf("Upgrade: %.0f mana", stats.UpgradeCost)
			if !stats.CanUpgrade {
				upgradeStr = "MAX LEVEL"
			}
			details := DimStyle.Render(fmt.Sprintf(
				"      %s | CD: %s | Cost: %.0f | %s",
				spell.Element,
				utils.FormatMilliseconds(stats.CooldownMs),
				stats.ManaCost,
				upgradeStr,
			))
			lines = append(lines, details)
		}
	}

	lines = append(lines, "")

	// Footer with contextual help
	lines = append(lines, FooterStyle.Render("[↑↓] Select  [Enter] Cast  [U] Upgrade  [Space] Toggle Slot  [<>] Reorder  [A] Auto  [B] Back"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewRituals renders the rituals view.
func (m Model) viewRituals() string {
	if m.gameState == nil {
		return "No game loaded"
	}

	var lines []string

	header := HeaderStyle.Width(60).Render(
		TitleStyle.Render("⚡ RITUALS"),
	)
	lines = append(lines, header)
	lines = append(lines, "")

	// Active rituals
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf("Active Rituals (%d/%d)",
		len(m.gameState.GetActiveRituals()), m.gameState.PrestigeData.RitualCapacity)))

	for _, ritual := range m.gameState.Rituals {
		status := "Active"
		if !ritual.IsActive {
			status = "Inactive"
		}
		lines = append(lines, fmt.Sprintf("  • %s - %s", ritual.Name, status))
	}

	if len(m.gameState.Rituals) == 0 {
		lines = append(lines, DimStyle.Render("  No rituals created"))
	}
	lines = append(lines, "")

	// Ritual builder
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	lines = append(lines, SubtitleStyle.Render("Create New Ritual (select 3 spells)"))

	// Show selected spells
	lines = append(lines, fmt.Sprintf("Selected: %d/3", len(m.ritualSpells)))
	for _, spellID := range m.ritualSpells {
		spell := m.gameState.GetSpellByID(spellID)
		if spell != nil {
			lines = append(lines, fmt.Sprintf("  ✓ %s", spell.Name))
		}
	}
	lines = append(lines, "")

	// Available spells
	lines = append(lines, DimStyle.Render("Available spells:"))
	for i, spell := range m.gameState.Spells {
		selected := i == m.selectedIndex

		// Check if already in ritual selection
		inSelection := false
		for _, id := range m.ritualSpells {
			if id == spell.ID {
				inSelection = true
				break
			}
		}

		prefix := "  "
		if selected {
			prefix = "> "
		}
		if inSelection {
			prefix = "  ✓ "
		}

		icon := GetElementIcon(string(spell.Element))
		line := fmt.Sprintf("%s%s %s", prefix, icon, spell.Name)

		if selected {
			lines = append(lines, SelectedStyle.Render(line))
		} else if inSelection {
			lines = append(lines, SuccessStyle.Render(line))
		} else {
			lines = append(lines, TextStyle.Render(line))
		}
	}
	lines = append(lines, "")

	// Footer
	lines = append(lines, FooterStyle.Render("[↑↓] Navigate  [Enter] Select  [C] Clear  [B/Esc] Back"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewStats renders the stats view.
func (m Model) viewStats() string {
	if m.gameState == nil {
		return "No game loaded"
	}

	gs := m.gameState
	var lines []string

	header := HeaderStyle.Width(60).Render(
		TitleStyle.Render("📊 STATISTICS"),
	)
	lines = append(lines, header)
	lines = append(lines, "")

	// Tower stats
	lines = append(lines, SubtitleStyle.Render("Tower Progress"))
	lines = append(lines, fmt.Sprintf("  Current Floor: %d", gs.Tower.CurrentFloor))
	lines = append(lines, fmt.Sprintf("  Max Floor Reached: %d", gs.Tower.MaxFloorReached))
	lines = append(lines, fmt.Sprintf("  Current Mana: %s", utils.FormatNumber(gs.Tower.CurrentMana)))
	lines = append(lines, fmt.Sprintf("  Lifetime Mana: %s", utils.FormatNumber(gs.Tower.LifetimeManaEarned)))
	lines = append(lines, "")

	// Mana generation
	if m.engine != nil {
		lines = append(lines, SubtitleStyle.Render("Mana Generation"))
		mps := m.engine.CalculateManaPerSecond(gs)
		lines = append(lines, fmt.Sprintf("  Base Rate: %s/sec", utils.FormatNumber(mps)))
		lines = append(lines, fmt.Sprintf("  Era Multiplier: %s", utils.FormatMultiplier(gs.PrestigeData.EraMultiplier)))
		lines = append(lines, fmt.Sprintf("  Permanent Bonus: %s", utils.FormatMultiplier(gs.PrestigeData.PermanentManaGenMultiplier)))
		ritualBonus := m.engine.GetRitualBonus(gs)
		lines = append(lines, fmt.Sprintf("  Ritual Bonus: %s", utils.FormatMultiplier(ritualBonus)))
		lines = append(lines, "")
	}

	// Prestige stats
	lines = append(lines, SubtitleStyle.Render("Prestige"))
	lines = append(lines, fmt.Sprintf("  Current Era: %d", gs.PrestigeData.CurrentEra))
	lines = append(lines, fmt.Sprintf("  Total Ascensions: %d", gs.PrestigeData.TotalAscensions))
	lines = append(lines, fmt.Sprintf("  Cooldown Reduction: %s", utils.FormatPercent(gs.PrestigeData.SpellCooldownReduction)))
	lines = append(lines, fmt.Sprintf("  Ritual Capacity: %d", gs.PrestigeData.RitualCapacity))
	lines = append(lines, fmt.Sprintf("  Auto-Cast Slots: %d (base 2 + %d bonus)", gs.GetAutoCastSlotCount(), gs.PrestigeData.AutoCastSlotBonus))
	lines = append(lines, "")

	// Spell stats
	lines = append(lines, SubtitleStyle.Render("Spells"))
	lines = append(lines, fmt.Sprintf("  Unlocked: %d", len(gs.Spells)))
	totalCasts := 0
	for _, s := range gs.Spells {
		totalCasts += s.CastCount
	}
	lines = append(lines, fmt.Sprintf("  Total Casts: %d", totalCasts))
	lines = append(lines, "")

	// Footer
	lines = append(lines, FooterStyle.Render("[B/Esc] Back"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewPrestige renders the prestige view.
func (m Model) viewPrestige() string {
	if m.gameState == nil {
		return "No game loaded"
	}

	gs := m.gameState
	var lines []string

	header := HeaderStyle.Width(60).Render(
		TitleStyle.Render("✨ PRESTIGE - ASCENSION"),
	)
	lines = append(lines, header)
	lines = append(lines, "")

	canPrestige := m.engine != nil && m.engine.CanPrestige(gs)

	if canPrestige {
		lines = append(lines, SuccessStyle.Render("You have reached Floor 100!"))
		lines = append(lines, "")
		lines = append(lines, TextStyle.Render("Ascending will:"))
		lines = append(lines, TextStyle.Render("  • Reset your floor to 1"))
		lines = append(lines, TextStyle.Render("  • Reset your current mana"))
		lines = append(lines, TextStyle.Render("  • Remove all rituals"))
		lines = append(lines, "")
		lines = append(lines, HighlightStyle.Render("But you will gain:"))

		newEra := gs.PrestigeData.CurrentEra + 1
		newMultiplier := 1.0 + (0.15 * float64(newEra))
		lines = append(lines, fmt.Sprintf("  • Era %d (%.2fx multiplier)", newEra, newMultiplier))
		lines = append(lines, "  • +5% permanent mana generation")
		lines = append(lines, "  • +5% spell cooldown reduction")
		if gs.PrestigeData.RitualCapacity < 3 {
			lines = append(lines, "  • +1 ritual slot")
		}
		if gs.PrestigeData.AutoCastSlotBonus < 2 {
			lines = append(lines, "  • +1 auto-cast slot")
		}
		lines = append(lines, "")
		lines = append(lines, WarningStyle.Render("Press [Enter] to ascend"))
	} else {
		lines = append(lines, ErrorStyle.Render(fmt.Sprintf("Reach Floor 100 to prestige (currently floor %d)", gs.Tower.CurrentFloor)))
	}
	lines = append(lines, "")

	// Footer
	lines = append(lines, FooterStyle.Render("[B/Esc] Back"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
