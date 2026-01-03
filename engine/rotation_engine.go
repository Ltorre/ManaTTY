package engine

import (
	"sort"

	"github.com/Ltorre/ManaTTY/models"
)

// v1.5.0: Advanced Spell Rotation System

// ProcessRotation handles priority-based spell rotation.
// This replaces the simple auto-cast loop with intelligent priority management.
func (e *GameEngine) ProcessRotation(gs *models.GameState) int {
	if gs.Session.Rotation == nil || !gs.Session.Rotation.Enabled {
		// Fall back to legacy auto-cast if rotation not enabled
		return e.ProcessAutoCasts(gs)
	}

	skipped := 0
	rotation := gs.Session.Rotation

	// Calculate mana reserve threshold
	manaReserve := gs.Tower.MaxMana * rotation.ManaThreshold
	availableMana := gs.Tower.CurrentMana - manaReserve

	// Group spells by priority
	highPriority := []models.RotationSpellConfig{}
	mediumPriority := []models.RotationSpellConfig{}
	lowPriority := []models.RotationSpellConfig{}

	for _, config := range rotation.Spells {
		if !config.Enabled {
			continue
		}

		switch config.Priority {
		case models.PriorityHigh:
			highPriority = append(highPriority, config)
		case models.PriorityMedium:
			mediumPriority = append(mediumPriority, config)
		case models.PriorityLow:
			lowPriority = append(lowPriority, config)
		}
	}

	// Process high priority first, then medium, then low (filler)
	priorityGroups := [][]models.RotationSpellConfig{highPriority, mediumPriority, lowPriority}

	for _, group := range priorityGroups {
		if rotation.CooldownWeaving {
			// Sort by cooldown remaining (cast spell with shortest CD first)
			sort.Slice(group, func(i, j int) bool {
				spellI := gs.GetSpellByID(group[i].SpellID)
				spellJ := gs.GetSpellByID(group[j].SpellID)
				if spellI == nil || spellJ == nil {
					return false
				}
				return spellI.CooldownRemainingMs < spellJ.CooldownRemainingMs
			})
		}

		for _, config := range group {
			spell := gs.GetSpellByID(config.SpellID)
			if spell == nil || !spell.IsReady() {
				continue
			}

			// Check rotation condition
			if !e.checkRotationCondition(gs, spell, config.Condition, availableMana) {
				continue
			}

			// Try to cast
			if err := e.CastSpell(gs, spell, false); err == ErrInsufficientMana {
				skipped++
				continue
			} else if err != nil {
				// Skip if other error (cooldown, needs specialization, etc.)
				continue
			}

			// If OptimizeForIdle is enabled, only cast one spell per tick
			// This spreads out casts for sustained DPS rather than burst
			if rotation.OptimizeForIdle {
				gs.Session.AutoCastSkipCount += skipped
				return skipped
			}
		}
	}

	gs.Session.AutoCastSkipCount += skipped
	return skipped
}

// checkRotationCondition evaluates advanced rotation conditions.
func (e *GameEngine) checkRotationCondition(gs *models.GameState, spell *models.Spell, cond models.RotationCondition, availableMana float64) bool {
	switch cond {
	case models.RotationConditionAlways:
		return true

	case models.RotationConditionManaAbove50:
		return gs.Tower.CurrentMana >= gs.Tower.MaxMana*0.5

	case models.RotationConditionManaAbove75:
		return gs.Tower.CurrentMana >= gs.Tower.MaxMana*0.75

	case models.RotationConditionSigilNotFull:
		return !gs.Tower.IsSigilCharged()

	case models.RotationConditionSynergyActive:
		return gs.HasActiveSynergy()

	case models.RotationConditionManaEfficient:
		// Only cast if current mana / spell cost > 2.0 (efficient)
		manaCost := e.calculateSpellManaCost(gs, spell)
		if manaCost <= 0 {
			return false
		}
		return (availableMana / manaCost) > 2.0

	case models.RotationConditionDuringSynergy:
		// Check if any active ritual synergy benefits this spell's element
		synergies := e.GetActiveSynergies(gs)
		for _, synergy := range synergies {
			for _, elem := range synergy.Elements {
				if elem == spell.Element {
					return true
				}
			}
		}
		return false

	case models.RotationConditionSigilAlmostFul:
		return gs.Tower.SigilCharge >= float64(gs.Tower.SigilRequired)*0.8

	case models.RotationConditionHighPriority:
		return true // Always cast when ready

	case models.RotationConditionFillerOnly:
		// Only cast if no higher priority spells are ready
		// This is handled by priority grouping, so always allow here
		return true

	default:
		return true
	}
}

// calculateSpellManaCost computes the effective mana cost for a spell.
// (Helper for mana efficiency calculations)
func (e *GameEngine) calculateSpellManaCost(gs *models.GameState, spell *models.Spell) float64 {
	// This duplicates logic from spell_engine.go CastSpell
	// Consider refactoring to share this calculation
	manaCost := spell.GetEffectiveManaCost(0.08) // Use game.SpellManaCostPerLevel = 0.08

	// Apply specialization
	if spell.HasSpecialization(models.SpecManaEfficiency) {
		manaCost *= 0.8 // -20%
	}

	// Apply elemental resonance
	resCounts := gs.GetAutoCastElementCounts()
	if resCounts[spell.Element] >= 2 && spell.Element == models.ElementThunder {
		manaCost *= 0.95 // -5% (game.ResonanceThunderManaCostReduction)
	}

	// Apply ritual bonuses
	ritualReduction := e.GetTotalRitualManaCostReductionWithSynergies(gs)
	if ritualReduction > 0 {
		manaCost *= (1.0 - ritualReduction)
	}

	// Apply active synergy
	if gs.HasActiveSynergy() && gs.GetActiveSynergy() == spell.Element {
		manaCost *= 0.8 // -20%
	}

	return manaCost
}
