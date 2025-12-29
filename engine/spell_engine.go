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
