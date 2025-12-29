package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
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
	case ViewSpecialize:
		content = m.viewSpecialize()
	case ViewFloorEvent:
		content = m.viewFloorEvent()
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

func (m Model) viewFloorEvent() string {
	if m.gameState == nil || m.gameState.Session == nil || m.gameState.Session.ActiveFloorEvent == nil {
		return "No active floor event"
	}

	evt := m.gameState.Session.ActiveFloorEvent
	nowMs := time.Now().UnixMilli()
	remainingMs := evt.ExpiresAtMs - nowMs
	if remainingMs < 0 {
		remainingMs = 0
	}
	remaining := time.Duration(remainingMs) * time.Millisecond
	mins := int(remaining.Minutes())
	secs := int(remaining.Seconds()) % 60

	sym := GetSymbols()
	var lines []string
	header := HeaderStyle.Width(60).Render(
		TitleStyle.Render(sym.Event + " FLOOR EVENT"),
	)
	lines = append(lines, header)
	lines = append(lines, "")
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf("A choice appears on Floor %d", evt.Floor)))
	lines = append(lines, DimStyle.Render(fmt.Sprintf("Expires in %dm%02ds (auto-dismisses with no bonus)", mins, secs)))
	lines = append(lines, "")

	opts := []struct {
		Choice models.FloorEventChoice
		Text   string
	}{
		{models.FloorEventChoiceManaGen, fmt.Sprintf("+%.0f%% mana/sec for %d floors", game.FloorEventManaGenBonus*100, game.FloorEventBuffDurationFloors)},
		{models.FloorEventChoiceSigilChargeRate, fmt.Sprintf("+%.0f%% sigil charge for %d floors", game.FloorEventSigilChargeRateBonus*100, game.FloorEventBuffDurationFloors)},
		{models.FloorEventChoiceCooldownReduction, fmt.Sprintf("-%.0f%% spell cooldown for %d floors", game.FloorEventCooldownReduction*100, game.FloorEventBuffDurationFloors)},
	}

	lines = append(lines, SubtitleStyle.Render("Choose one:"))
	for i, opt := range opts {
		prefix := "  "
		style := TextStyle
		if i == m.selectedIndex {
			prefix = "> "
			style = HighlightStyle
		}
		name := models.FloorEventChoiceDisplayNames[opt.Choice]
		lines = append(lines, style.Render(fmt.Sprintf("%s[%d] %s — %s", prefix, i+1, name, opt.Text)))
	}

	lines = append(lines, "")
	lines = append(lines, DimStyle.Render("[1/2/3 or Enter] Select  |  [Esc] Ignore"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewTower renders the main tower view.
func (m Model) viewTower() string {
	if m.gameState == nil {
		return "No game loaded"
	}

	gs := m.gameState
	var lines []string

	// Header
	sym := GetSymbols()
	era := ""
	if gs.PrestigeData.CurrentEra > 0 {
		era = fmt.Sprintf(" - ERA %d", gs.PrestigeData.CurrentEra)
	}
	header := HeaderStyle.Width(60).Render(
		TitleStyle.Render(sym.Tower + " MAGE TOWER ASCENSION" + era),
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
	manaStr := fmt.Sprintf(sym.Mana+" Mana   [%s] %d%% (%s / %s)",
		bar,
		percentage,
		utils.FormatNumber(gs.Tower.CurrentMana),
		utils.FormatNumber(gs.Tower.MaxMana),
	)
	lines = append(lines, manaStr)

	// Ascension Sigil progress bar
	sigilProgress := gs.Tower.GetSigilProgress()
	sigilFilled := int(sigilProgress * float64(barWidth))
	sigilBar := ProgressBarFilled.Render(strings.Repeat("█", sigilFilled)) +
		ProgressBarEmpty.Render(strings.Repeat("░", barWidth-sigilFilled))
	sigilPercent := int(sigilProgress * 100)

	sigilStatus := ""
	if gs.Tower.IsSigilCharged() {
		sigilStatus = SuccessStyle.Render(" ✓ READY")
	}
	sigilStr := fmt.Sprintf("⚔️ Sigil  [%s] %d%% (%.0f / %.0f)%s",
		sigilBar,
		sigilPercent,
		gs.Tower.SigilCharge,
		gs.Tower.SigilRequired,
		sigilStatus,
	)
	lines = append(lines, sigilStr)
	lines = append(lines, "")

	// Stats section
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	lines = append(lines, SubtitleStyle.Render(sym.Stats+" STATS"))

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

	// Element synergy status
	if gs.HasActiveSynergy() {
		element := gs.GetActiveSynergy()
		remaining := gs.GetSynergyTimeRemaining() / 1000 // convert to seconds
		synergyStr := fmt.Sprintf("  %s Synergy: %ds remaining (20%% bonus)",
			string(element), remaining)
		lines = append(lines, HighlightStyle.Render(synergyStr))
	} else if len(gs.Session.LastCastElements) > 0 {
		// Show element streak progress
		streakLen := len(gs.Session.LastCastElements)
		lastElement := gs.Session.LastCastElements[streakLen-1]
		lines = append(lines, DimStyle.Render(fmt.Sprintf("  Element streak: %s ×%d/3",
			string(lastElement), streakLen)))
	}
	lines = append(lines, "")

	// Active rituals
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	activeRituals := gs.GetActiveRituals()
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf(sym.Ritual+" ACTIVE RITUALS (%d/%d)",
		len(activeRituals), gs.PrestigeData.RitualCapacity)))

	if len(activeRituals) == 0 {
		lines = append(lines, DimStyle.Render("  No active rituals"))
	} else {
		for i, ritual := range activeRituals {
			cooldownStr := SuccessStyle.Render("Ready")
			if ritual.CooldownRemaining > 0 {
				cooldownStr = WarningStyle.Render(utils.FormatCooldown(ritual.CooldownRemaining))
			}

			// Compute effects dynamically for legacy rituals or use stored ones
			comboInfo := game.ComputeRitualCombo(ritual.SpellIDs)
			ritualName := comboInfo.Name
			if comboInfo.SignatureName != "" {
				ritualName += " \"" + comboInfo.SignatureName + "\""
			}

			lines = append(lines, fmt.Sprintf("  [%d] %s (%s)", i+1, ritualName, cooldownStr))

			// Show effect summary inline
			if len(comboInfo.Effects) > 0 {
				effectStrs := []string{}
				for _, effect := range comboInfo.Effects {
					icon := game.GetRitualEffectIcon(effect.Type)
					effectStr := game.GetEffectDisplayString(effect)
					effectStrs = append(effectStrs, icon+" "+effectStr)
				}
				lines = append(lines, HighlightStyle.Render("      "+strings.Join(effectStrs, "  |  ")))
			}
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
	sym := GetSymbols()
	var lines []string

	header := HeaderStyle.Width(50).Render(
		TitleStyle.Render(sym.Bullet + " MENU"),
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

	sym := GetSymbols()
	var lines []string

	header := HeaderStyle.Width(70).Render(
		TitleStyle.Render(sym.Bullet + " SPELLS"),
	)
	lines = append(lines, header)

	// Current mana display
	manaStr := fmt.Sprintf(sym.Mana+" Mana: %s / %s",
		utils.FormatNumber(m.gameState.Tower.CurrentMana),
		utils.FormatNumber(m.gameState.Tower.MaxMana))
	lines = append(lines, HighlightStyle.Render(manaStr))
	lines = append(lines, "")

	// Element Synergy Status
	if m.gameState.HasActiveSynergy() {
		synergy := m.gameState.GetActiveSynergy()
		remaining := m.gameState.GetSynergyTimeRemaining() / 1000
		icon := GetElementIcon(string(synergy))
		lines = append(lines, SuccessStyle.Render(fmt.Sprintf("%s %s SYNERGY ACTIVE! +20%% bonus (%ds remaining)", sym.Synergy, icon, remaining)))
		lines = append(lines, "")
	}

	// Auto-Cast Loadout Panel
	slotSpellIDs := m.gameState.GetAutoCastSpellIDs()
	usedSlots := len(slotSpellIDs)
	maxSlots := m.gameState.GetAutoCastSlotCount()
	autoCastStatus := "OFF"
	if m.gameState.Session.AutoCastEnabled {
		autoCastStatus = "ON"
	}
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf(sym.AutoCast+" Auto-Cast Loadout [%s] (%d/%d slots)", autoCastStatus, usedSlots, maxSlots)))

	idxBySpellID := map[string]int{}
	for i, id := range slotSpellIDs {
		idxBySpellID[id] = i
	}

	if len(slotSpellIDs) == 0 {
		lines = append(lines, DimStyle.Render("  (empty - press Space on a spell to add)"))
	} else {
		for i, spellID := range slotSpellIDs {
			spell := m.gameState.GetSpellByID(spellID)
			if spell != nil {
				icon := GetElementIcon(string(spell.Element))
				cond := m.gameState.GetAutoCastCondition(spellID)
				condStr := models.ConditionShortNames[cond]
				lines = append(lines, TextStyle.Render(fmt.Sprintf("  %d. %s %s (Lv%d) [%s]", i+1, icon, spell.Name, spell.Level, condStr)))
			}
		}
	}

	// Elemental Resonance: passive bonuses from themed loadouts (2+ spells of same element)
	counts := m.gameState.GetAutoCastElementCounts()
	resLines := []string{}
	if counts[models.ElementFire] >= game.ElementalResonanceMinSpells {
		resLines = append(resLines, fmt.Sprintf("%s Fire x%d (+%.0f%% dmg)", GetElementIcon(string(models.ElementFire)), counts[models.ElementFire], game.ResonanceFireDamageBonus*100))
	}
	if counts[models.ElementIce] >= game.ElementalResonanceMinSpells {
		resLines = append(resLines, fmt.Sprintf("%s Ice x%d (-%.0f%% CD)", GetElementIcon(string(models.ElementIce)), counts[models.ElementIce], game.ResonanceIceCooldownReduction*100))
	}
	if counts[models.ElementThunder] >= game.ElementalResonanceMinSpells {
		resLines = append(resLines, fmt.Sprintf("%s Thunder x%d (-%.0f%% cost)", GetElementIcon(string(models.ElementThunder)), counts[models.ElementThunder], game.ResonanceThunderManaCostReduction*100))
	}
	if counts[models.ElementArcane] >= game.ElementalResonanceMinSpells {
		resLines = append(resLines, fmt.Sprintf("%s Arcane x%d (+%.0f%% sigil)", GetElementIcon(string(models.ElementArcane)), counts[models.ElementArcane], game.ResonanceArcaneSigilChargeBonus*100))
	}
	if len(resLines) > 0 {
		lines = append(lines, DimStyle.Render("Resonance: "+strings.Join(resLines, "  |  ")))
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
		if idx, ok := idxBySpellID[spell.ID]; ok {
			autoIndicator = fmt.Sprintf("[%d]", idx+1)
		}

		// Status
		status := SuccessStyle.Render("Ready")
		if !spell.IsReady() {
			status = WarningStyle.Render(utils.FormatCooldown(spell.CooldownRemainingMs))
		}

		// Level indicator with max check
		levelStr := fmt.Sprintf("Lv%d", spell.Level)
		if spell.Level >= game.SpellMaxLevel {
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
			upgradeStr := fmt.Sprintf("Upgrade: %s mana", utils.FormatNumber(stats.UpgradeCost))
			if !stats.CanUpgrade {
				upgradeStr = "MAX LEVEL"
			}

			// Specialization info
			specStr := ""
			tier, needs := spell.NeedsSpecialization()
			if needs {
				specStr = fmt.Sprintf(" | ★ Tier %d SPEC READY!", tier)
			} else {
				specs := []string{}
				if spell.Tier1Spec != models.SpecNone {
					specs = append(specs, models.SpecializationShortNames[spell.Tier1Spec])
				}
				if spell.Tier2Spec != models.SpecNone {
					specs = append(specs, models.SpecializationShortNames[spell.Tier2Spec])
				}
				if len(specs) > 0 {
					specStr = " | ★" + strings.Join(specs, ",")
				}
			}

			details := DimStyle.Render(fmt.Sprintf(
				"      %s | DMG: %.0f | CD: %s | Cost: %.0f | %s%s",
				spell.Element,
				stats.Damage,
				utils.FormatMilliseconds(stats.CooldownMs),
				stats.ManaCost,
				upgradeStr,
				specStr,
			))
			lines = append(lines, details)
		}
	}

	lines = append(lines, "")

	// Footer with contextual help
	lines = append(lines, FooterStyle.Render("[↑↓] Select  [Enter] Cast  [U] Upgrade  [Space] Slot  [<>] Order  [C] Cond  [X] Spec  [A] Auto  [B] Back"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// viewRituals renders the rituals view.
func (m Model) viewRituals() string {
	if m.gameState == nil {
		return "No game loaded"
	}

	var lines []string

	header := HeaderStyle.Width(70).Render(
		TitleStyle.Render("⚡ RITUALS"),
	)
	lines = append(lines, header)
	lines = append(lines, "")

	// Active rituals with v1.2.0 effects
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf("Active Rituals (%d/%d)",
		len(m.gameState.GetActiveRituals()), m.gameState.PrestigeData.RitualCapacity)))

	for _, ritual := range m.gameState.Rituals {
		status := SuccessStyle.Render("Active")
		if !ritual.IsActive {
			status = DimStyle.Render("Inactive")
		}

		// Build ritual display with name and signature
		ritualName := ritual.Name
		if ritual.SignatureName != "" {
			ritualName = fmt.Sprintf("%s \"%s\"", ritual.Name, ritual.SignatureName)
		}
		lines = append(lines, fmt.Sprintf("  • %s - %s", ritualName, status))

		// Show ritual effects (v1.2.0)
		if len(ritual.Effects) > 0 {
			effectStrs := []string{}
			for _, effect := range ritual.Effects {
				icon := game.GetRitualEffectIcon(effect.Type)
				effectStr := game.GetEffectDisplayString(effect)
				effectStrs = append(effectStrs, icon+" "+effectStr)
			}
			effectLine := "    " + strings.Join(effectStrs, "  |  ")
			if ritual.HasSpellEcho {
				effectLine += " (Echo ✨)"
			}
			lines = append(lines, HighlightStyle.Render(effectLine))
		}
	}

	if len(m.gameState.Rituals) == 0 {
		lines = append(lines, DimStyle.Render("  No rituals created"))
	}
	lines = append(lines, "")

	// Total ritual bonuses summary (v1.2.0)
	if len(m.gameState.GetActiveRituals()) > 0 && m.engine != nil {
		lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		lines = append(lines, SubtitleStyle.Render("Combined Bonuses"))

		dmgBonus := m.engine.GetTotalRitualDamageBonus(m.gameState)
		cdBonus := m.engine.GetTotalRitualCooldownReduction(m.gameState)
		manaBonus := m.engine.GetTotalRitualManaCostReduction(m.gameState)
		sigilBonus := m.engine.GetTotalRitualSigilChargeBonus(m.gameState)

		bonusLines := []string{}
		if dmgBonus > 0 {
			bonusLines = append(bonusLines, fmt.Sprintf("🔥 +%.0f%% damage", dmgBonus*100))
		}
		if cdBonus > 0 {
			bonusLines = append(bonusLines, fmt.Sprintf("❄️ -%.0f%% cooldown", cdBonus*100))
		}
		if manaBonus > 0 {
			bonusLines = append(bonusLines, fmt.Sprintf("⚡ -%.0f%% mana cost", manaBonus*100))
		}
		if sigilBonus > 0 {
			bonusLines = append(bonusLines, fmt.Sprintf("✨ +%.0f%% sigil charge", sigilBonus*100))
		}

		if len(bonusLines) > 0 {
			lines = append(lines, "  "+strings.Join(bonusLines, "  |  "))
		}
		lines = append(lines, "")
	}

	// Ritual builder
	lines = append(lines, DimStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	lines = append(lines, SubtitleStyle.Render("Create New Ritual (select 3 spells)"))

	// Show selected spells with preview
	lines = append(lines, fmt.Sprintf("Selected: %d/3", len(m.ritualSpells)))
	for _, spellID := range m.ritualSpells {
		spell := m.gameState.GetSpellByID(spellID)
		if spell != nil {
			icon := GetElementIcon(string(spell.Element))
			lines = append(lines, fmt.Sprintf("  ✓ %s %s", icon, spell.Name))
		}
	}

	// Preview ritual combo (v1.2.0)
	if len(m.ritualSpells) == 3 {
		comboInfo := game.ComputeRitualCombo(m.ritualSpells)
		lines = append(lines, "")
		previewName := comboInfo.Name
		if comboInfo.SignatureName != "" {
			previewName = fmt.Sprintf("%s \"%s\"", comboInfo.Name, comboInfo.SignatureName)
		}
		lines = append(lines, HighlightStyle.Render("  Preview: "+previewName))

		effectStrs := []string{}
		for _, effect := range comboInfo.Effects {
			icon := game.GetRitualEffectIcon(effect.Type)
			effectStr := game.GetEffectDisplayString(effect)
			effectStrs = append(effectStrs, icon+" "+effectStr)
		}
		if len(effectStrs) > 0 {
			lines = append(lines, SuccessStyle.Render("  Effects: "+strings.Join(effectStrs, "  |  ")))
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
	lines = append(lines, FooterStyle.Render("[↑↓] Navigate  [Enter] Select  [C] Clear  [X] Reset Rituals  [B/Esc] Back"))

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

// viewSpecialize renders the spell specialization selection.
func (m Model) viewSpecialize() string {
	if m.gameState == nil {
		return "Loading..."
	}

	spell := m.gameState.GetSpellByID(m.specSpellID)
	if spell == nil {
		return "Spell not found"
	}

	sym := GetSymbols()
	lines := []string{}
	lines = append(lines, TitleStyle.Render(sym.Star+" Spell Specialization"))
	lines = append(lines, "")
	lines = append(lines, SubtitleStyle.Render(fmt.Sprintf("Choose a Tier %d specialization for %s (Lv%d)", m.specTier, spell.Name, spell.Level)))
	lines = append(lines, "")

	type specOption struct {
		spec models.SpellSpecialization
		name string
		desc string
	}

	var options []specOption
	if m.specTier == 1 {
		options = []specOption{
			{models.SpecCritChance, "⚔️  Critical Strike", fmt.Sprintf("+%.0f%% chance for %.1fx damage", game.SpecCritChanceBonus*100, game.SpecCritDamageMulti)},
			{models.SpecManaEfficiency, "💧 Mana Efficiency", fmt.Sprintf("-%.0f%% mana cost", game.SpecManaEfficiencyBonus*100)},
		}
	} else {
		options = []specOption{
			{models.SpecBurstDamage, "💥 Burst Damage", fmt.Sprintf("+%.0f%% spell damage", game.SpecBurstDamageBonus*100)},
			{models.SpecRapidCast, "⚡ Rapid Cast", fmt.Sprintf("-%.0f%% cooldown", game.SpecRapidCastBonus*100)},
		}
	}

	for i, opt := range options {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.specChoiceIdx {
			style = style.Bold(true).Foreground(ColorAccent).
				Border(lipgloss.RoundedBorder()).BorderForeground(ColorAccent)
		}
		optText := fmt.Sprintf("%s\n%s", opt.name, lipgloss.NewStyle().Foreground(ColorTextDim).Render(opt.desc))
		lines = append(lines, style.Render(optText))
		lines = append(lines, "")
	}

	lines = append(lines, "")
	lines = append(lines, FooterStyle.Render("[↑/↓] Select  [Enter] Confirm  [Esc] Cancel"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
