package engine

import (
	"errors"

	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
)

// Spell casting errors
var (
	ErrSpellOnCooldown  = errors.New("spell is on cooldown")
	ErrInsufficientMana = errors.New("insufficient mana")
	ErrSpellNotFound    = errors.New("spell not found")
	ErrSpellNotUnlocked = errors.New("spell not unlocked")
)

// CastSpell attempts to cast a spell.
// Both manual and auto-cast now require mana. Manual costs +10% more.
func (e *GameEngine) CastSpell(gs *models.GameState, spell *models.Spell, manual bool) error {
	// Check cooldown
	if !spell.IsReady() {
		return ErrSpellOnCooldown
	}

	// Calculate mana cost (manual casts cost more due to ManualCastPenalty)
	manaCost := spell.BaseManaRequirement
	if manual {
		manaCost = game.CalculateManualCastCost(spell.BaseManaRequirement)
	}

	// Check mana - both auto and manual require mana now
	if gs.Tower.CurrentMana < manaCost {
		return ErrInsufficientMana
	}

	// Deduct mana
	gs.Tower.SpendMana(manaCost)

	// Apply cooldown with prestige reduction
	cooldownReduction := gs.PrestigeData.SpellCooldownReduction
	spell.StartCooldown(cooldownReduction)

	return nil
}

// CastSpellByID attempts to cast a spell by its ID.
func (e *GameEngine) CastSpellByID(gs *models.GameState, spellID string, manual bool) error {
	spell := gs.GetSpellByID(spellID)
	if spell == nil {
		return ErrSpellNotFound
	}

	return e.CastSpell(gs, spell, manual)
}

// GetSpellCooldownRemaining returns the remaining cooldown for a spell.
func (e *GameEngine) GetSpellCooldownRemaining(spell *models.Spell) int64 {
	return spell.CooldownRemainingMs
}

// GetSpellEffectiveCooldown returns the cooldown after prestige reductions.
func (e *GameEngine) GetSpellEffectiveCooldown(spell *models.Spell, cooldownReduction float64) int64 {
	return game.CalculateSpellCooldown(spell.BaseCooldownMs, cooldownReduction)
}

// GetReadySpells returns all spells that can be cast.
func (e *GameEngine) GetReadySpells(gs *models.GameState) []*models.Spell {
	ready := []*models.Spell{}
	for _, spell := range gs.Spells {
		if spell.IsReady() {
			ready = append(ready, spell)
		}
	}
	return ready
}

// GetSpellsOnCooldown returns all spells currently on cooldown.
func (e *GameEngine) GetSpellsOnCooldown(gs *models.GameState) []*models.Spell {
	onCooldown := []*models.Spell{}
	for _, spell := range gs.Spells {
		if !spell.IsReady() {
			onCooldown = append(onCooldown, spell)
		}
	}
	return onCooldown
}

// CountReadySpells returns the number of spells ready to cast.
func (e *GameEngine) CountReadySpells(gs *models.GameState) int {
	count := 0
	for _, spell := range gs.Spells {
		if spell.IsReady() {
			count++
		}
	}
	return count
}

// GetTotalCastCount returns total spell casts across all spells.
func (e *GameEngine) GetTotalCastCount(gs *models.GameState) int {
	total := 0
	for _, spell := range gs.Spells {
		total += spell.CastCount
	}
	return total
}

// ToggleAutoCast toggles auto-casting on or off.
func (e *GameEngine) ToggleAutoCast(gs *models.GameState) bool {
	gs.Session.AutoCastEnabled = !gs.Session.AutoCastEnabled
	return gs.Session.AutoCastEnabled
}

// SetAutoCast sets auto-casting to a specific state.
func (e *GameEngine) SetAutoCast(gs *models.GameState, enabled bool) {
	gs.Session.AutoCastEnabled = enabled
}

// ErrNoAutoCastSlots is returned when all auto-cast slots are full.
var ErrNoAutoCastSlots = errors.New("no auto-cast slots available")

// ToggleSpellAutoCast adds or removes a spell from auto-cast slots.
// Returns (isNowInSlot, error). Error is non-nil if slots are full when trying to add.
func (e *GameEngine) ToggleSpellAutoCast(gs *models.GameState, spellID string) (bool, error) {
	// Verify spell exists
	if gs.GetSpellByID(spellID) == nil {
		return false, ErrSpellNotFound
	}
	
	result := gs.ToggleSpellAutoCast(spellID)
	switch result {
	case models.AutoCastAdded:
		return true, nil
	case models.AutoCastRemoved:
		return false, nil
	case models.AutoCastSlotsFull:
		return false, ErrNoAutoCastSlots
	default:
		return false, nil
	}
}

// AddSpellToAutoCast adds a spell to an auto-cast slot.
func (e *GameEngine) AddSpellToAutoCast(gs *models.GameState, spellID string) error {
	if gs.GetSpellByID(spellID) == nil {
		return ErrSpellNotFound
	}
	if gs.IsSpellInAutoCast(spellID) {
		return nil // Already in slot
	}
	if !gs.AddSpellToAutoCast(spellID) {
		return errors.New("no auto-cast slots available")
	}
	return nil
}

// RemoveSpellFromAutoCast removes a spell from auto-cast slots.
func (e *GameEngine) RemoveSpellFromAutoCast(gs *models.GameState, spellID string) {
	gs.RemoveSpellFromAutoCast(spellID)
}

// GetAutoCastSlots returns the current auto-cast slot configuration.
func (e *GameEngine) GetAutoCastSlots(gs *models.GameState) []string {
	return gs.Session.AutoCastSlots
}

// GetMaxAutoCastSlots returns total available slots.
func (e *GameEngine) GetMaxAutoCastSlots(gs *models.GameState) int {
	return gs.GetAutoCastSlotCount()
}

// GetUsedAutoCastSlots returns number of slots in use.
func (e *GameEngine) GetUsedAutoCastSlots(gs *models.GameState) int {
	return len(gs.Session.AutoCastSlots)
}
