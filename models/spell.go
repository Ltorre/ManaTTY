package models

import (
	"time"
)

// Element represents a spell's magical element type.
type Element string

const (
	ElementFire    Element = "fire"
	ElementIce     Element = "ice"
	ElementThunder Element = "thunder"
	ElementArcane  Element = "arcane"
)

// SpellSpecialization represents a chosen perk for a spell at milestone levels.
type SpellSpecialization string

const (
	// Tier 1 specializations (unlocked at level 5)
	SpecNone           SpellSpecialization = ""
	SpecCritChance     SpellSpecialization = "crit_chance"     // +15% crit chance (2x damage on crit)
	SpecManaEfficiency SpellSpecialization = "mana_efficiency" // +20% reduced mana cost

	// Tier 2 specializations (unlocked at level 10)
	SpecBurstDamage SpellSpecialization = "burst_damage" // +30% damage
	SpecRapidCast   SpellSpecialization = "rapid_cast"   // +25% cooldown reduction
)

// SpellDefinition represents a spell template (stored in DB).
// Note: Spell level scaling uses global constants from game/constants.go:
// - SpellDamagePerLevel (+15% per level)
// - SpellCooldownPerLevel (-5% per level)
// - SpellManaCostPerLevel (-8% per level)
type SpellDefinition struct {
	ID                string  `bson:"_id" json:"id"`
	Name              string  `bson:"name" json:"name"`
	Description       string  `bson:"description" json:"description"`
	FlavorText        string  `bson:"flavor_text" json:"flavor_text"`
	Element           Element `bson:"element" json:"element"`
	BaseDamage        float64 `bson:"base_damage" json:"base_damage"`
	BaseCooldownMs    int64   `bson:"base_cooldown_ms" json:"base_cooldown_ms"`
	BaseManaCost      float64 `bson:"base_mana_cost" json:"base_mana_cost"`
	RequiredFloor     int     `bson:"required_floor" json:"required_floor"`
	UnlockedByDefault bool    `bson:"unlocked_by_default" json:"unlocked_by_default"`
	PrestigeExclusive bool    `bson:"prestige_exclusive" json:"prestige_exclusive"`
	Version           int     `bson:"version" json:"version"`
}

// Spell represents a player's instance of a spell with progress.
type Spell struct {
	ID                  string    `bson:"id" json:"id"`
	Name                string    `bson:"name" json:"name"`
	Element             Element   `bson:"element" json:"element"`
	Level               int       `bson:"level" json:"level"`
	BaseDamage          float64   `bson:"base_damage" json:"base_damage"`
	BaseCooldownMs      int64     `bson:"base_cooldown_ms" json:"base_cooldown_ms"`
	BaseManaRequirement float64   `bson:"base_mana_requirement" json:"base_mana_requirement"`
	RequiredFloor       int       `bson:"required_floor" json:"required_floor"`
	CooldownRemainingMs int64     `bson:"cooldown_remaining_ms" json:"cooldown_remaining_ms"`
	LastCastTime        time.Time `bson:"last_cast_time" json:"last_cast_time"`
	CastCount           int       `bson:"cast_count" json:"cast_count"`

	// Specializations chosen at milestone levels
	Tier1Spec SpellSpecialization `bson:"tier1_spec" json:"tier1_spec"` // Chosen at level 5
	Tier2Spec SpellSpecialization `bson:"tier2_spec" json:"tier2_spec"` // Chosen at level 10
}

// NewSpellFromDefinition creates a Spell instance from a SpellDefinition.
func NewSpellFromDefinition(def *SpellDefinition) *Spell {
	return &Spell{
		ID:                  def.ID,
		Name:                def.Name,
		Element:             def.Element,
		Level:               1,
		BaseDamage:          def.BaseDamage,
		BaseCooldownMs:      def.BaseCooldownMs,
		BaseManaRequirement: def.BaseManaCost,
		RequiredFloor:       def.RequiredFloor,
		CooldownRemainingMs: 0,
		LastCastTime:        time.Time{},
		CastCount:           0,
	}
}

// IsReady returns true if the spell can be cast (cooldown finished).
func (s *Spell) IsReady() bool {
	return s.CooldownRemainingMs <= 0
}

// GetCurrentDamage returns the damage at current level.
func (s *Spell) GetCurrentDamage(damagePerLevel float64) float64 {
	return s.BaseDamage + (damagePerLevel * float64(s.Level-1))
}

// GetCurrentCooldown returns the cooldown at current level (in ms).
func (s *Spell) GetCurrentCooldown(reductionPerLevel int64) int64 {
	cooldown := s.BaseCooldownMs - (reductionPerLevel * int64(s.Level-1))
	if cooldown < 1000 { // Minimum 1 second
		cooldown = 1000
	}
	return cooldown
}

// StartCooldown sets the spell on cooldown.
func (s *Spell) StartCooldown(cooldownReduction float64) {
	cooldown := float64(s.BaseCooldownMs) * (1.0 - cooldownReduction)
	s.CooldownRemainingMs = int64(cooldown)
	s.LastCastTime = time.Now()
	s.CastCount++
}

// UpdateCooldown reduces cooldown by elapsed time.
func (s *Spell) UpdateCooldown(elapsedMs int64) {
	if s.CooldownRemainingMs > 0 {
		s.CooldownRemainingMs -= elapsedMs
		if s.CooldownRemainingMs < 0 {
			s.CooldownRemainingMs = 0
		}
	}
}

// GetEffectiveCooldown returns the cooldown after level bonuses (in ms).
func (s *Spell) GetEffectiveCooldown(cooldownPerLevel float64) int64 {
	// Each level reduces cooldown by cooldownPerLevel %
	reduction := cooldownPerLevel * float64(s.Level-1)
	cooldown := float64(s.BaseCooldownMs) * (1.0 - reduction)
	if cooldown < 1000 {
		cooldown = 1000 // Minimum 1 second
	}
	return int64(cooldown)
}

// GetEffectiveManaCost returns the mana cost after level bonuses.
func (s *Spell) GetEffectiveManaCost(manaCostPerLevel float64) float64 {
	// Each level reduces mana cost by manaCostPerLevel %
	reduction := manaCostPerLevel * float64(s.Level-1)
	cost := s.BaseManaRequirement * (1.0 - reduction)
	if cost < 1 {
		cost = 1 // Minimum 1 mana
	}
	return cost
}

// GetEffectiveDamage returns the damage after level bonuses.
func (s *Spell) GetEffectiveDamage(damagePerLevel float64) float64 {
	// Each level increases damage by damagePerLevel %
	bonus := damagePerLevel * float64(s.Level-1)
	return s.BaseDamage * (1.0 + bonus)
}

// CanLevelUp returns true if the spell can be upgraded.
func (s *Spell) CanLevelUp(maxLevel int) bool {
	return s.Level < maxLevel
}

// LevelUp increases the spell level by 1.
// Returns false if already at or above maxLevel (currently SpellMaxLevel=10).
func (s *Spell) LevelUp(maxLevel int) bool {
	if s.Level >= maxLevel {
		return false
	}
	s.Level++
	return true
}

// NeedsSpecialization returns true if spell is at a milestone level without a chosen spec.
func (s *Spell) NeedsSpecialization() (tier int, needs bool) {
	if s.Level >= 5 && s.Tier1Spec == SpecNone {
		return 1, true
	}
	if s.Level >= 10 && s.Tier2Spec == SpecNone {
		return 2, true
	}
	return 0, false
}

// HasSpecialization checks if spell has a specific specialization.
func (s *Spell) HasSpecialization(spec SpellSpecialization) bool {
	return s.Tier1Spec == spec || s.Tier2Spec == spec
}
