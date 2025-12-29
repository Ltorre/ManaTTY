package engine

import (
	"time"

	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
)

// OfflineProgress holds the results of offline progress calculation.
type OfflineProgress struct {
	TimeOffline    time.Duration
	ManaGenerated  float64
	FloorsClimbed  int
	FinalFloor     int
	FinalMana      float64
	SpellsUnlocked []string
}

// CalculateOfflineProgress processes offline time and returns the results.
func (e *GameEngine) CalculateOfflineProgress(gs *models.GameState) *OfflineProgress {
	// Calculate time offline
	lastSaved := gs.Session.LastSavedAt
	if lastSaved.IsZero() {
		lastSaved = gs.SavedAt
	}

	timeOffline := time.Since(lastSaved)
	offlineSeconds := timeOffline.Seconds()

	// Skip if minimal offline time
	if offlineSeconds < float64(game.MinOfflineSeconds) {
		return &OfflineProgress{
			TimeOffline: timeOffline,
			FinalFloor:  gs.Tower.CurrentFloor,
			FinalMana:   gs.Tower.CurrentMana,
		}
	}

	// Calculate mana generation rate (at time of disconnect)
	manaPerSecond := e.CalculateManaPerSecond(gs)

	// Apply offline penalty
	offlineMana := game.CalculateOfflineMana(manaPerSecond, offlineSeconds)

	// Add the offline mana
	gs.Tower.AddMana(offlineMana)

	// Update floor requirement
	gs.Tower.MaxMana = game.CalculateFloorCost(gs.Tower.CurrentFloor)

	// Try to climb floors with accumulated mana
	floorsClimbed := 0
	spellsUnlocked := []string{}

	for e.TryClimbFloor(gs) {
		floorsClimbed++

		// Track any spells unlocked
		newSpells := game.GetNewSpellsAtFloor(gs.Tower.CurrentFloor)
		for _, spell := range newSpells {
			if !gs.HasSpell(spell.ID) && !spell.PrestigeExclusive {
				spellsUnlocked = append(spellsUnlocked, spell.ID)
			}
		}
	}

	return &OfflineProgress{
		TimeOffline:    timeOffline,
		ManaGenerated:  offlineMana,
		FloorsClimbed:  floorsClimbed,
		FinalFloor:     gs.Tower.CurrentFloor,
		FinalMana:      gs.Tower.CurrentMana,
		SpellsUnlocked: spellsUnlocked,
	}
}

// ApplyOfflineProgress processes and applies offline progress to game state.
func (e *GameEngine) ApplyOfflineProgress(gs *models.GameState) *OfflineProgress {
	progress := e.CalculateOfflineProgress(gs)

	// Unlock any spells that should have been unlocked
	for _, spellID := range progress.SpellsUnlocked {
		spellDef := game.GetSpellDefinition(spellID)
		if spellDef != nil && !gs.HasSpell(spellID) {
			spell := models.NewSpellFromDefinition(spellDef)
			gs.AddSpell(spell)
		}
	}

	// Reset session for new play session
	gs.Session = models.NewSessionData()

	return progress
}

// FormatOfflineProgress returns a human-readable summary of offline progress.
func FormatOfflineProgress(progress *OfflineProgress) string {
	if progress.TimeOffline < time.Minute {
		return "Welcome back!"
	}

	// Format time offline
	var timeStr string
	hours := int(progress.TimeOffline.Hours())
	minutes := int(progress.TimeOffline.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		timeStr = formatPlural(days, "day") + " " + formatPlural(hours, "hour")
	} else if hours > 0 {
		timeStr = formatPlural(hours, "hour") + " " + formatPlural(minutes, "minute")
	} else {
		timeStr = formatPlural(minutes, "minute")
	}

	return timeStr + " offline"
}

// formatPlural formats a number with singular/plural suffix.
func formatPlural(n int, singular string) string {
	if n == 1 {
		return "1 " + singular
	}
	return formatInt(n) + " " + singular + "s"
}

// formatInt formats an integer with comma separators.
func formatInt(n int) string {
	if n < 1000 {
		return string(rune('0'+n%10) + rune('0'+n/10%10)*10 + rune('0'+n/100)*100)
	}
	// Simple implementation for small numbers
	return string(rune(n))
}

// EstimateOfflineProgress estimates what progress would be made offline.
func (e *GameEngine) EstimateOfflineProgress(gs *models.GameState, duration time.Duration) *OfflineProgress {
	// Calculate mana that would be generated
	manaPerSecond := e.CalculateManaPerSecond(gs)
	offlineMana := game.CalculateOfflineMana(manaPerSecond, duration.Seconds())

	// Calculate floors that could be climbed
	floorsClimbed, remainingMana := game.CalculateFloorsFromMana(gs.Tower.CurrentFloor, gs.Tower.CurrentMana+offlineMana)

	return &OfflineProgress{
		TimeOffline:   duration,
		ManaGenerated: offlineMana,
		FloorsClimbed: floorsClimbed,
		FinalFloor:    gs.Tower.CurrentFloor + floorsClimbed,
		FinalMana:     remainingMana,
	}
}
