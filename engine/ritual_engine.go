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

// getRitualEffects returns ritual effects, computing dynamically for legacy rituals.
func getRitualEffects(ritual *models.Ritual) []models.RitualEffect {
	// If effects are already computed, use them
	if len(ritual.Effects) > 0 {
		return ritual.Effects
	}
	// For legacy rituals, compute dynamically
	comboInfo := game.ComputeRitualCombo(ritual.SpellIDs)
	return comboInfo.Effects
}

// getTotalRitualEffectBonus aggregates a specific effect type across all active rituals.
func (e *GameEngine) getTotalRitualEffectBonus(gs *models.GameState, effectType models.RitualEffectType) float64 {
	total := 0.0
	for _, ritual := range gs.Rituals {
		if ritual.IsActive {
			for _, effect := range getRitualEffects(ritual) {
				if effect.Type == effectType {
					total += effect.Magnitude
				}
			}
		}
	}
	return total
}

// GetTotalRitualDamageBonus returns combined damage bonus from all active rituals.
func (e *GameEngine) GetTotalRitualDamageBonus(gs *models.GameState) float64 {
	return e.getTotalRitualEffectBonus(gs, models.RitualEffectDamage)
}

// GetTotalRitualCooldownReduction returns combined cooldown reduction from all active rituals.
func (e *GameEngine) GetTotalRitualCooldownReduction(gs *models.GameState) float64 {
	return e.getTotalRitualEffectBonus(gs, models.RitualEffectCooldown)
}

// GetTotalRitualManaCostReduction returns combined mana cost reduction from all active rituals.
func (e *GameEngine) GetTotalRitualManaCostReduction(gs *models.GameState) float64 {
	return e.getTotalRitualEffectBonus(gs, models.RitualEffectManaCost)
}

// GetTotalRitualSigilChargeBonus returns combined sigil charge bonus from all active rituals.
func (e *GameEngine) GetTotalRitualSigilChargeBonus(gs *models.GameState) float64 {
	return e.getTotalRitualEffectBonus(gs, models.RitualEffectSigilRate)
}

// GetTotalRitualManaGenBonus returns combined mana generation bonus from all active rituals.
func (e *GameEngine) GetTotalRitualManaGenBonus(gs *models.GameState) float64 {
	return e.getTotalRitualEffectBonus(gs, models.RitualEffectManaGenRate)
}

// v1.4.0: Ritual Synergy System

// GetActiveSynergies detects and returns all active synergies based on ritual elements.
func (e *GameEngine) GetActiveSynergies(gs *models.GameState) []models.RitualSynergy {
	// Count elements from all active rituals
	elementCounts := make(map[models.Element]int)
	for _, ritual := range gs.Rituals {
		if ritual.IsActive {
			effects := getRitualEffects(ritual)
			for _, effect := range effects {
				// Map effect type back to element
				switch effect.Type {
				case models.RitualEffectDamage:
					elementCounts[models.ElementFire]++
				case models.RitualEffectCooldown:
					elementCounts[models.ElementIce]++
				case models.RitualEffectManaCost:
					elementCounts[models.ElementThunder]++
				case models.RitualEffectManaGenRate:
					elementCounts[models.ElementArcane]++
				}
			}
		}
	}

	// Check each synergy definition
	activeSynergies := []models.RitualSynergy{}
	for _, synergy := range models.SynergyDefinitions {
		// Check if all required elements are present
		hasAllElements := true
		for _, element := range synergy.Elements {
			if elementCounts[element] == 0 {
				hasAllElements = false
				break
			}
		}
		if hasAllElements {
			activeSynergies = append(activeSynergies, synergy)
		}
	}

	return activeSynergies
}

// GetSynergyBonusForElement returns the total synergy bonus multiplier for a specific element.
// This stacks multiplicatively with ritual effects.
func (e *GameEngine) GetSynergyBonusForElement(gs *models.GameState, element models.Element) float64 {
	activeSynergies := e.GetActiveSynergies(gs)
	totalBonus := 0.0

	for _, synergy := range activeSynergies {
		// Check if this synergy applies to the element
		for _, synergyElement := range synergy.Elements {
			if synergyElement == element {
				totalBonus += synergy.Magnitude
				break
			}
		}
	}

	return totalBonus
}

// GetTotalRitualDamageBonusWithSynergies returns damage bonus with synergy multipliers applied.
func (e *GameEngine) GetTotalRitualDamageBonusWithSynergies(gs *models.GameState) float64 {
	baseBonus := e.GetTotalRitualDamageBonus(gs)
	synergyBonus := e.GetSynergyBonusForElement(gs, models.ElementFire)
	return baseBonus * (1.0 + synergyBonus)
}

// GetTotalRitualCooldownReductionWithSynergies returns cooldown reduction with synergy multipliers applied.
func (e *GameEngine) GetTotalRitualCooldownReductionWithSynergies(gs *models.GameState) float64 {
	baseBonus := e.GetTotalRitualCooldownReduction(gs)
	synergyBonus := e.GetSynergyBonusForElement(gs, models.ElementIce)
	return baseBonus * (1.0 + synergyBonus)
}

// GetTotalRitualManaCostReductionWithSynergies returns mana cost reduction with synergy multipliers applied.
func (e *GameEngine) GetTotalRitualManaCostReductionWithSynergies(gs *models.GameState) float64 {
	baseBonus := e.GetTotalRitualManaCostReduction(gs)
	synergyBonus := e.GetSynergyBonusForElement(gs, models.ElementThunder)
	return baseBonus * (1.0 + synergyBonus)
}

// GetTotalRitualManaGenBonusWithSynergies returns mana generation bonus with synergy multipliers applied.
func (e *GameEngine) GetTotalRitualManaGenBonusWithSynergies(gs *models.GameState) float64 {
	baseBonus := e.GetTotalRitualManaGenBonus(gs)
	synergyBonus := e.GetSynergyBonusForElement(gs, models.ElementArcane)
	return baseBonus * (1.0 + synergyBonus)
}
