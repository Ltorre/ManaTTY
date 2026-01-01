package models

// RitualEffectType defines the type of passive bonus a ritual provides.
type RitualEffectType string

const (
	RitualEffectDamage      RitualEffectType = "damage"      // Fire signature: +X% spell damage
	RitualEffectCooldown    RitualEffectType = "cooldown"    // Ice signature: -X% spell cooldown
	RitualEffectManaCost    RitualEffectType = "mana_cost"   // Thunder signature: -X% mana cost
	RitualEffectSigilRate   RitualEffectType = "sigil_rate"  // Arcane signature: +X% sigil charge rate
	RitualEffectManaGenRate RitualEffectType = "mana_gen"    // Hybrid combos: +X% mana generation
)

// RitualComposition indicates the element distribution in a ritual.
type RitualComposition string

const (
	CompositionPure   RitualComposition = "pure"   // 3 same element
	CompositionHybrid RitualComposition = "hybrid" // 2+1 elements
	CompositionTriad  RitualComposition = "triad"  // 1/1/1 elements
)

// RitualEffect represents a single effect bonus from a ritual.
type RitualEffect struct {
	Type      RitualEffectType `bson:"type" json:"type"`
	Magnitude float64          `bson:"magnitude" json:"magnitude"` // e.g., 0.18 for +18%
}

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

	// v1.2.0: Ritual combo effects
	Composition     RitualComposition `bson:"composition" json:"composition"`
	Effects         []RitualEffect    `bson:"effects" json:"effects"`
	HasSpellEcho    bool              `bson:"has_spell_echo" json:"has_spell_echo"`     // +5% kicker
	SignatureName   string            `bson:"signature_name" json:"signature_name"`     // Special flavor name (if any)
	DominantElement Element           `bson:"dominant_element" json:"dominant_element"` // For hybrid/pure
}

// RitualBonus is the mana generation multiplier per active ritual.
const RitualBonus = 0.15

// MaxActiveRituals is the maximum number of rituals (at full prestige).
const MaxActiveRituals = 3

// NewRitual creates a new ritual from 3 spells.
// Note: This creates a basic ritual. Use NewRitualWithEffects for v1.2.0 combo effects.
func NewRitual(spellIDs []string) *Ritual {
	if len(spellIDs) != 3 {
		return nil
	}

	// Generate a simple ID from spell IDs
	id := "ritual_" + spellIDs[0] + "_" + spellIDs[1] + "_" + spellIDs[2]

	return &Ritual{
		ID:                id,
		Name:              legacyRitualName(spellIDs),
		SpellIDs:          spellIDs,
		IsActive:          true,
		CooldownMs:        60000, // 60 seconds
		CooldownRemaining: 0,
		BoostMultiplier:   1.0 + RitualBonus,
		CastCount:         0,
		// v1.2.0 fields will be populated by game.ComputeRitualCombo
		Composition:     CompositionTriad, // Default
		Effects:         []RitualEffect{},
		HasSpellEcho:    false,
		SignatureName:   "",
		DominantElement: "",
	}
}

// NewRitualWithEffects creates a ritual with computed v1.2.0 combo effects.
// The info parameter should come from game.ComputeRitualCombo().
func NewRitualWithEffects(spellIDs []string, name string, composition RitualComposition, dominant Element, effects []RitualEffect, hasEcho bool, signatureName string) *Ritual {
	if len(spellIDs) != 3 {
		return nil
	}

	id := "ritual_" + spellIDs[0] + "_" + spellIDs[1] + "_" + spellIDs[2]

	return &Ritual{
		ID:                id,
		Name:              name,
		SpellIDs:          spellIDs,
		IsActive:          true,
		CooldownMs:        60000,
		CooldownRemaining: 0,
		BoostMultiplier:   1.0 + RitualBonus,
		CastCount:         0,
		Composition:       composition,
		Effects:           effects,
		HasSpellEcho:      hasEcho,
		SignatureName:     signatureName,
		DominantElement:   dominant,
	}
}

// legacyRitualName creates a simple display name (fallback).
func legacyRitualName(spellIDs []string) string {
	if len(spellIDs) < 3 {
		return "Unknown Ritual"
	}
	return spellIDs[0] + " + " + spellIDs[1] + " + " + spellIDs[2]
}

// GetEffectByType returns the effect of a specific type, if present.
func (r *Ritual) GetEffectByType(effectType RitualEffectType) (RitualEffect, bool) {
	for _, effect := range r.Effects {
		if effect.Type == effectType {
			return effect, true
		}
	}
	return RitualEffect{}, false
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
