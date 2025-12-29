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
	ErrSpellMaxLevel    = errors.New("spell already at max level")
)

// CastSpell attempts to cast a spell.
// Both manual and auto-cast now require mana. Manual costs +10% more.
func (e *GameEngine) CastSpell(gs *models.GameState, spell *models.Spell, manual bool) error {
	// Check cooldown
	if !spell.IsReady() {
		return ErrSpellOnCooldown
	}

	// Calculate effective mana cost (with level bonuses)
	manaCost := game.CalculateSpellEffectiveManaCost(spell.BaseManaRequirement, spell.Level)
	if manual {
		manaCost = game.CalculateManualCastCost(manaCost)
	}

	// Apply synergy bonus if active and matching element
	if gs.HasActiveSynergy() && gs.GetActiveSynergy() == spell.Element {
		manaCost *= (1.0 - game.ElementSynergyBonus) // 20% cheaper
	}

	// Check mana - both auto and manual require mana now
	if gs.Tower.CurrentMana < manaCost {
		return ErrInsufficientMana
	}

	// Deduct mana
	gs.Tower.SpendMana(manaCost)

	// Calculate effective cooldown (with level bonuses)
	baseCooldown := game.CalculateSpellEffectiveCooldown(spell.BaseCooldownMs, spell.Level)
	cooldownReduction := gs.PrestigeData.SpellCooldownReduction

	// Apply synergy bonus to cooldown if active
	if gs.HasActiveSynergy() && gs.GetActiveSynergy() == spell.Element {
		cooldownReduction += game.ElementSynergyBonus // Additional 20% reduction
	}

	// Cap cooldown reduction to max allowed
	if cooldownReduction > game.MaxCooldownReduction {
		cooldownReduction = game.MaxCooldownReduction
	}

	spell.CooldownRemainingMs = game.CalculateSpellCooldown(baseCooldown, cooldownReduction)
	spell.CastCount++

	// Record for element synergy tracking
	gs.RecordSpellCast(spell.Element)

	// Check if synergy should trigger
	if synergy := gs.CheckElementSynergy(); synergy != "" {
		gs.ActivateSynergy(synergy, int64(game.ElementSynergyDuration*1000))
		if e.OnSynergyActivated != nil {
			e.OnSynergyActivated(synergy)
		}
	}

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

// MoveAutoCastSlotUp moves a spell higher in auto-cast priority.
func (e *GameEngine) MoveAutoCastSlotUp(gs *models.GameState, spellID string) bool {
	return gs.MoveAutoCastSlot(spellID, -1)
}

// MoveAutoCastSlotDown moves a spell lower in auto-cast priority.
func (e *GameEngine) MoveAutoCastSlotDown(gs *models.GameState, spellID string) bool {
	return gs.MoveAutoCastSlot(spellID, 1)
}

// GetSpellUpgradeCost returns the mana cost to upgrade a spell.
func (e *GameEngine) GetSpellUpgradeCost(spell *models.Spell) float64 {
	return game.CalculateSpellUpgradeCost(spell.Level, spell.BaseManaRequirement)
}

// CanUpgradeSpell returns true if the spell can be upgraded and player has enough mana.
func (e *GameEngine) CanUpgradeSpell(gs *models.GameState, spell *models.Spell) bool {
	if spell.Level >= game.SpellMaxLevel {
		return false
	}
	cost := e.GetSpellUpgradeCost(spell)
	return gs.Tower.CurrentMana >= cost
}

// UpgradeSpell upgrades a spell if possible, spending mana.
func (e *GameEngine) UpgradeSpell(gs *models.GameState, spell *models.Spell) error {
	if spell.Level >= game.SpellMaxLevel {
		return ErrSpellMaxLevel
	}

	cost := e.GetSpellUpgradeCost(spell)
	if gs.Tower.CurrentMana < cost {
		return ErrInsufficientMana
	}

	gs.Tower.SpendMana(cost)
	spell.LevelUp(game.SpellMaxLevel)

	if e.OnSpellUpgraded != nil {
		e.OnSpellUpgraded(spell)
	}

	return nil
}

// GetSpellEffectiveStats returns the effective stats for a spell at its current level.
type SpellEffectiveStats struct {
	ManaCost    float64
	CooldownMs  int64
	Damage      float64
	UpgradeCost float64
	CanUpgrade  bool
}

func (e *GameEngine) GetSpellEffectiveStats(gs *models.GameState, spell *models.Spell) SpellEffectiveStats {
	return SpellEffectiveStats{
		ManaCost:    game.CalculateSpellEffectiveManaCost(spell.BaseManaRequirement, spell.Level),
		CooldownMs:  game.CalculateSpellEffectiveCooldown(spell.BaseCooldownMs, spell.Level),
		Damage:      spell.GetEffectiveDamage(game.SpellDamagePerLevel),
		UpgradeCost: e.GetSpellUpgradeCost(spell),
		CanUpgrade:  spell.Level < game.SpellMaxLevel,
	}
}
