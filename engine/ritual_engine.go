package engine

import (
	"errors"

	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
)

// Ritual errors
var (
	ErrRitualSlotsFull   = errors.New("ritual slots are full")
	ErrDuplicateSpells   = errors.New("ritual requires 3 unique spells")
	ErrSpellNotOwned     = errors.New("spell not owned")
	ErrRitualNotFound    = errors.New("ritual not found")
	ErrInvalidSpellCount = errors.New("ritual requires exactly 3 spells")
)

// CreateRitual creates a new ritual from 3 spells with v1.2.0 combo effects.
func (e *GameEngine) CreateRitual(gs *models.GameState, spellIDs []string) (*models.Ritual, error) {
	// Validate spell count
	if len(spellIDs) != game.SpellsPerRitual {
		return nil, ErrInvalidSpellCount
	}

	// Check if slots available
	if !gs.CanAddRitual() {
		return nil, ErrRitualSlotsFull
	}

	// Validate uniqueness
	seen := make(map[string]bool)
	for _, id := range spellIDs {
		if seen[id] {
			return nil, ErrDuplicateSpells
		}
		seen[id] = true
	}

	// Validate ownership
	for _, id := range spellIDs {
		if !gs.HasSpell(id) {
			return nil, ErrSpellNotOwned
		}
	}

	// Compute ritual combo effects (v1.2.0)
	comboInfo := game.ComputeRitualCombo(spellIDs)

	// Create the ritual with effects
	ritual := models.NewRitualWithEffects(
		spellIDs,
		comboInfo.Name,
		comboInfo.Composition,
		comboInfo.DominantElement,
		comboInfo.Effects,
		comboInfo.HasSpellEcho,
		comboInfo.SignatureName,
	)
	gs.Rituals = append(gs.Rituals, ritual)
	gs.ActiveRitualCount = len(gs.GetActiveRituals())

	return ritual, nil
}

// RemoveRitual removes a ritual by ID.
func (e *GameEngine) RemoveRitual(gs *models.GameState, ritualID string) error {
	for i, ritual := range gs.Rituals {
		if ritual.ID == ritualID {
			// Remove from slice
			gs.Rituals = append(gs.Rituals[:i], gs.Rituals[i+1:]...)
			gs.ActiveRitualCount = len(gs.GetActiveRituals())
			return nil
		}
	}
	return ErrRitualNotFound
}

// ResetRituals removes all rituals from the save slot.
// This is a safety valve to prevent players from permanently locking their ritual capacity.
func (e *GameEngine) ResetRituals(gs *models.GameState) {
	gs.Rituals = []*models.Ritual{}
	gs.ActiveRitualCount = 0
}

// ToggleRitual activates or deactivates a ritual.
func (e *GameEngine) ToggleRitual(gs *models.GameState, ritualID string) error {
	for _, ritual := range gs.Rituals {
		if ritual.ID == ritualID {
			ritual.IsActive = !ritual.IsActive
			gs.ActiveRitualCount = len(gs.GetActiveRituals())
			return nil
		}
	}
	return ErrRitualNotFound
}

// GetRitualByID returns a ritual by its ID.
func (e *GameEngine) GetRitualByID(gs *models.GameState, ritualID string) *models.Ritual {
	for _, ritual := range gs.Rituals {
		if ritual.ID == ritualID {
			return ritual
		}
	}
	return nil
}

// GetAvailableRitualSlots returns the number of available ritual slots.
func (e *GameEngine) GetAvailableRitualSlots(gs *models.GameState) int {
	return gs.PrestigeData.RitualCapacity - len(gs.GetActiveRituals())
}

// GetMaxRitualSlots returns the maximum ritual slots.
func (e *GameEngine) GetMaxRitualSlots(gs *models.GameState) int {
	return gs.PrestigeData.RitualCapacity
}

// GetRitualBonusPercent returns the total ritual bonus as a percentage.
func (e *GameEngine) GetRitualBonusPercent(gs *models.GameState) float64 {
	activeCount := len(gs.GetActiveRituals())
	return float64(activeCount) * models.RitualBonus * 100
}

// CanCreateRitual checks if a ritual can be created with given spells.
func (e *GameEngine) CanCreateRitual(gs *models.GameState, spellIDs []string) (bool, error) {
	if len(spellIDs) != game.SpellsPerRitual {
		return false, ErrInvalidSpellCount
	}

	if !gs.CanAddRitual() {
		return false, ErrRitualSlotsFull
	}

	seen := make(map[string]bool)
	for _, id := range spellIDs {
		if seen[id] {
			return false, ErrDuplicateSpells
		}
		seen[id] = true

		if !gs.HasSpell(id) {
			return false, ErrSpellNotOwned
		}
	}

	return true, nil
}

// GetSpellsAvailableForRitual returns spells that can be used in a new ritual.
func (e *GameEngine) GetSpellsAvailableForRitual(gs *models.GameState, excludeIDs []string) []*models.Spell {
	available := []*models.Spell{}
	excludeMap := make(map[string]bool)
	for _, id := range excludeIDs {
		excludeMap[id] = true
	}

	for _, spell := range gs.Spells {
		if !excludeMap[spell.ID] {
			available = append(available, spell)
		}
	}

	return available
}

// v1.2.0: Ritual Combo Effect Aggregation

// GetTotalRitualDamageBonus returns combined damage bonus from all active rituals.
func (e *GameEngine) GetTotalRitualDamageBonus(gs *models.GameState) float64 {
	total := 0.0
	for _, ritual := range gs.Rituals {
		if ritual.IsActive {
			total += ritual.GetTotalDamageBonus()
		}
	}
	return total
}

// GetTotalRitualCooldownReduction returns combined cooldown reduction from all active rituals.
func (e *GameEngine) GetTotalRitualCooldownReduction(gs *models.GameState) float64 {
	total := 0.0
	for _, ritual := range gs.Rituals {
		if ritual.IsActive {
			total += ritual.GetTotalCooldownReduction()
		}
	}
	return total
}

// GetTotalRitualManaCostReduction returns combined mana cost reduction from all active rituals.
func (e *GameEngine) GetTotalRitualManaCostReduction(gs *models.GameState) float64 {
	total := 0.0
	for _, ritual := range gs.Rituals {
		if ritual.IsActive {
			total += ritual.GetTotalManaCostReduction()
		}
	}
	return total
}

// GetTotalRitualSigilChargeBonus returns combined sigil charge bonus from all active rituals.
func (e *GameEngine) GetTotalRitualSigilChargeBonus(gs *models.GameState) float64 {
	total := 0.0
	for _, ritual := range gs.Rituals {
		if ritual.IsActive {
			total += ritual.GetTotalSigilChargeBonus()
		}
	}
	return total
}
