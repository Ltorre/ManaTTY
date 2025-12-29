package models

// Ritual represents a 3-spell combination that boosts mana generation.
type Ritual struct {
	ID                string   `bson:"id" json:"id"`
	Name              string   `bson:"name" json:"name"`
	SpellIDs          []string `bson:"spell_ids" json:"spell_ids"`
	IsActive          bool     `bson:"is_active" json:"is_active"`
	CooldownMs        int64    `bson:"cooldown_ms" json:"cooldown_ms"`
	CooldownRemaining int64    `bson:"cooldown_remaining" json:"cooldown_remaining"`
	BoostMultiplier   float64  `bson:"boost_multiplier" json:"boost_multiplier"`
	CastCount         int      `bson:"cast_count" json:"cast_count"`
}

// RitualBonus is the mana generation multiplier per active ritual.
const RitualBonus = 0.15

// MaxActiveRituals is the maximum number of rituals (at full prestige).
const MaxActiveRituals = 3

// NewRitual creates a new ritual from 3 spells.
func NewRitual(spellIDs []string) *Ritual {
	if len(spellIDs) != 3 {
		return nil
	}

	// Generate a simple ID from spell IDs
	id := "ritual_" + spellIDs[0] + "_" + spellIDs[1] + "_" + spellIDs[2]

	return &Ritual{
		ID:                id,
		Name:              generateRitualName(spellIDs),
		SpellIDs:          spellIDs,
		IsActive:          true,
		CooldownMs:        60000, // 60 seconds
		CooldownRemaining: 0,
		BoostMultiplier:   1.0 + RitualBonus,
		CastCount:         0,
	}
}

// generateRitualName creates a display name from spell IDs.
func generateRitualName(spellIDs []string) string {
	if len(spellIDs) < 3 {
		return "Unknown Ritual"
	}
	return spellIDs[0] + " + " + spellIDs[1] + " + " + spellIDs[2]
}

// IsReady returns true if the ritual can be activated.
func (r *Ritual) IsReady() bool {
	return r.CooldownRemaining <= 0
}

// StartCooldown puts the ritual on cooldown.
func (r *Ritual) StartCooldown() {
	r.CooldownRemaining = r.CooldownMs
	r.CastCount++
}

// UpdateCooldown reduces cooldown by elapsed time.
func (r *Ritual) UpdateCooldown(elapsedMs int64) {
	if r.CooldownRemaining > 0 {
		r.CooldownRemaining -= elapsedMs
		if r.CooldownRemaining < 0 {
			r.CooldownRemaining = 0
		}
	}
}

// CalculateTotalRitualBonus returns the combined bonus from active rituals.
func CalculateTotalRitualBonus(rituals []*Ritual) float64 {
	activeCount := 0
	for _, r := range rituals {
		if r.IsActive {
			activeCount++
		}
	}
	return 1.0 + (RitualBonus * float64(activeCount))
}
