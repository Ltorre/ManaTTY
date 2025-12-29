package engine

import (
	"time"

	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
)

// GameEngine handles all game logic and state updates.
type GameEngine struct {
	// Event handlers
	OnFloorClimbed     func(floor int)
	OnSpellUnlocked    func(spell *models.Spell)
	OnManaGenerated    func(amount float64)
	OnPrestige         func(era int)
	OnSynergyActivated func(element models.Element)
	OnSpellUpgraded    func(spell *models.Spell)
}

// NewGameEngine creates a new game engine instance.
func NewGameEngine() *GameEngine {
	return &GameEngine{}
}

// Tick processes a single game tick, updating all game state.
func (e *GameEngine) Tick(gs *models.GameState, elapsed time.Duration) {
	elapsedMs := elapsed.Milliseconds()
	elapsedSec := elapsed.Seconds()

	// Generate mana
	manaPerSec := e.CalculateManaPerSecond(gs)
	manaGenerated := manaPerSec * elapsedSec
	gs.Tower.AddMana(manaGenerated)

	if e.OnManaGenerated != nil && manaGenerated > 0 {
		e.OnManaGenerated(manaGenerated)
	}

	// Update floor mana requirement
	gs.Tower.MaxMana = game.CalculateFloorCost(gs.Tower.CurrentFloor)

	// Try to climb floors
	for e.TryClimbFloor(gs) {
		// Keep climbing
	}

	// Update spell cooldowns
	e.UpdateSpellCooldowns(gs, elapsedMs)

	// Update ritual cooldowns
	e.UpdateRitualCooldowns(gs, elapsedMs)

	// Auto-cast spells if enabled
	if gs.Session.AutoCastEnabled {
		e.ProcessAutoCasts(gs)
	}

	// Update session data
	gs.UpdateSession()
}

// CalculateManaPerSecond returns the current mana generation rate.
func (e *GameEngine) CalculateManaPerSecond(gs *models.GameState) float64 {
	currentFloor := gs.Tower.CurrentFloor
	currentEra := gs.PrestigeData.CurrentEra
	activeRituals := len(gs.GetActiveRituals())
	permanentMultiplier := gs.PrestigeData.PermanentManaGenMultiplier

	return game.CalculateManaPerSecondWithBonuses(
		currentFloor,
		currentEra,
		activeRituals,
		permanentMultiplier,
	)
}

// TryClimbFloor attempts to climb to the next floor.
func (e *GameEngine) TryClimbFloor(gs *models.GameState) bool {
	requiredMana := game.CalculateFloorCost(gs.Tower.CurrentFloor)

	if gs.Tower.CurrentMana >= requiredMana {
		// Spend mana and climb
		gs.Tower.SpendMana(requiredMana)
		gs.Tower.ClimbFloor()

		// Update the new mana requirement
		gs.Tower.MaxMana = game.CalculateFloorCost(gs.Tower.CurrentFloor)

		// Check for spell unlocks
		e.CheckSpellUnlocks(gs)

		if e.OnFloorClimbed != nil {
			e.OnFloorClimbed(gs.Tower.CurrentFloor)
		}

		return true
	}

	return false
}

// CheckSpellUnlocks checks if new spells should be unlocked at the current floor.
func (e *GameEngine) CheckSpellUnlocks(gs *models.GameState) {
	newSpells := game.GetNewSpellsAtFloor(gs.Tower.CurrentFloor)

	for _, spellDef := range newSpells {
		// Skip if already unlocked
		if gs.HasSpell(spellDef.ID) {
			continue
		}

		// Skip prestige-exclusive spells if not unlocked via prestige
		if spellDef.PrestigeExclusive {
			hasPrestigeUnlock := false
			for _, unlocked := range gs.PrestigeData.UnlockedPrestigeSpells {
				if unlocked == spellDef.ID {
					hasPrestigeUnlock = true
					break
				}
			}
			if !hasPrestigeUnlock {
				continue
			}
		}

		// Create and add the spell
		spell := models.NewSpellFromDefinition(spellDef)
		gs.AddSpell(spell)

		if e.OnSpellUnlocked != nil {
			e.OnSpellUnlocked(spell)
		}
	}
}

// UpdateSpellCooldowns reduces cooldowns for all spells.
func (e *GameEngine) UpdateSpellCooldowns(gs *models.GameState, elapsedMs int64) {
	for _, spell := range gs.Spells {
		spell.UpdateCooldown(elapsedMs)
	}
}

// UpdateRitualCooldowns reduces cooldowns for all rituals.
func (e *GameEngine) UpdateRitualCooldowns(gs *models.GameState, elapsedMs int64) {
	for _, ritual := range gs.Rituals {
		ritual.UpdateCooldown(elapsedMs)
	}
}

// ProcessAutoCasts automatically casts spells in auto-cast slots if mana is available.
// Only spells assigned to auto-cast slots will be cast automatically.
// Returns the number of spells skipped due to insufficient mana.
func (e *GameEngine) ProcessAutoCasts(gs *models.GameState) int {
	skipped := 0
	for _, spellID := range gs.Session.AutoCastSlots {
		spell := gs.GetSpellByID(spellID)
		if spell != nil && spell.IsReady() {
			// CastSpell will check mana and skip if insufficient
			if err := e.CastSpell(gs, spell, false); err == ErrInsufficientMana {
				skipped++
				continue
			}
		}
	}
	gs.Session.AutoCastSkipCount += skipped
	return skipped
}

// GetAndResetSkipCount returns the accumulated skip count and resets it.
func (e *GameEngine) GetAndResetSkipCount(gs *models.GameState) int {
	count := gs.Session.AutoCastSkipCount
	gs.Session.AutoCastSkipCount = 0
	return count
}

// GetFloorProgress returns progress towards the next floor (0.0 to 1.0).
func (e *GameEngine) GetFloorProgress(gs *models.GameState) float64 {
	return gs.Tower.GetFloorProgress()
}

// GetManaToNextFloor returns mana needed for the next floor.
func (e *GameEngine) GetManaToNextFloor(gs *models.GameState) float64 {
	return gs.Tower.MaxMana - gs.Tower.CurrentMana
}

// GetTimeToNextFloor estimates time to reach the next floor.
func (e *GameEngine) GetTimeToNextFloor(gs *models.GameState) time.Duration {
	manaNeeded := e.GetManaToNextFloor(gs)
	if manaNeeded <= 0 {
		return 0
	}

	manaPerSec := e.CalculateManaPerSecond(gs)
	if manaPerSec <= 0 {
		return time.Hour * 999 // Effectively infinite
	}

	secondsNeeded := manaNeeded / manaPerSec
	return time.Duration(secondsNeeded * float64(time.Second))
}
