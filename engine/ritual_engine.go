package engine

import (
	"errors"

	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
)

// Ritual errors
var (
	ErrRitualSlotsFull     = errors.New("ritual slots are full")
	ErrDuplicateSpells     = errors.New("ritual requires 3 unique spells")
	ErrSpellNotOwned       = errors.New("spell not owned")
	ErrRitualNotFound      = errors.New("ritual not found")
	ErrInvalidSpellCount   = errors.New("ritual requires exactly 3 spells")
)

// CreateRitual creates a new ritual from 3 spells.
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

	// Create the ritual
	ritual := models.NewRitual(spellIDs)
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
