package models

import (
	"time"
)

// PrestigeData contains all prestige/ascension related data.
type PrestigeData struct {
	TotalAscensions            int         `bson:"total_ascensions" json:"total_ascensions"`
	CurrentEra                 int         `bson:"current_era" json:"current_era"`
	EraMultiplier              float64     `bson:"era_multiplier" json:"era_multiplier"`
	PermanentManaGenMultiplier float64     `bson:"permanent_mana_gen_multiplier" json:"permanent_mana_gen_multiplier"`
	SpellCooldownReduction     float64     `bson:"spell_cooldown_reduction" json:"spell_cooldown_reduction"`
	ManaRetention              float64     `bson:"mana_retention" json:"mana_retention"`
	RitualCapacity             int         `bson:"ritual_capacity" json:"ritual_capacity"`
	AutoCastSlotBonus          int         `bson:"auto_cast_slot_bonus" json:"auto_cast_slot_bonus"` // Extra auto-cast slots from prestige
	UnlockedPrestigeSpells     []string    `bson:"unlocked_prestige_spells" json:"unlocked_prestige_spells"`
	PrestigeEvents             []time.Time `bson:"prestige_events" json:"prestige_events"`
}

// PrestigeMilestone is the floor required to prestige.
const PrestigeMilestone = 100

// EraMultiplierBonus is added per era.
const EraMultiplierBonus = 0.15

// PrestigeManaGenBonus is added per prestige.
const PrestigeManaGenBonus = 0.05

// PrestigeCooldownBonus is added per prestige.
const PrestigeCooldownBonus = 0.05

// PrestigeManaRetentionBonus is added per prestige.
const PrestigeManaRetentionBonus = 0.10

// NewPrestigeData creates initial prestige data.
func NewPrestigeData() *PrestigeData {
	return &PrestigeData{
		TotalAscensions:            0,
		CurrentEra:                 0,
		EraMultiplier:              1.0,
		PermanentManaGenMultiplier: 1.0,
		SpellCooldownReduction:     0.0,
		ManaRetention:              0.0,
		RitualCapacity:             1,
		AutoCastSlotBonus:          0,
		UnlockedPrestigeSpells:     []string{},
		PrestigeEvents:             []time.Time{},
	}
}

// CanPrestige returns true if prestige is available.
func (p *PrestigeData) CanPrestige(currentFloor int) bool {
	return currentFloor >= PrestigeMilestone
}

// ProcessPrestige applies prestige bonuses and increments era.
func (p *PrestigeData) ProcessPrestige() {
	p.TotalAscensions++
	p.CurrentEra++

	// Update era multiplier
	p.EraMultiplier = 1.0 + (EraMultiplierBonus * float64(p.CurrentEra))

	// Add permanent bonuses
	p.PermanentManaGenMultiplier += PrestigeManaGenBonus
	p.SpellCooldownReduction += PrestigeCooldownBonus
	p.ManaRetention += PrestigeManaRetentionBonus

	// Cap cooldown reduction at 50%
	if p.SpellCooldownReduction > 0.50 {
		p.SpellCooldownReduction = 0.50
	}

	// Cap mana retention at 90%
	if p.ManaRetention > 0.90 {
		p.ManaRetention = 0.90
	}

	// Unlock ritual slot (max 3)
	if p.RitualCapacity < MaxActiveRituals {
		p.RitualCapacity++
	}

	// Unlock auto-cast slot (max +2 bonus = 4 total)
	if p.AutoCastSlotBonus < 2 {
		p.AutoCastSlotBonus++
	}

	// Record prestige event
	p.PrestigeEvents = append(p.PrestigeEvents, time.Now())

	// Unlock prestige-exclusive spells based on era
	if p.CurrentEra == 3 && !contains(p.UnlockedPrestigeSpells, "spell_meteor_strike") {
		p.UnlockedPrestigeSpells = append(p.UnlockedPrestigeSpells, "spell_meteor_strike")
	}
}

// GetTotalMultiplier returns the combined era + permanent mana multiplier.
func (p *PrestigeData) GetTotalMultiplier() float64 {
	return p.EraMultiplier * p.PermanentManaGenMultiplier
}

// contains checks if a string is in a slice.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
