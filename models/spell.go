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

// SpellDefinition represents a spell template (stored in DB).
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
	Scaling           struct {
		DamagePerLevel           float64 `bson:"damage_per_level" json:"damage_per_level"`
		CooldownReductionPerLvl  int64   `bson:"cooldown_reduction_per_level" json:"cooldown_reduction_per_level"`
		ManaCostReductionPerLvl  float64 `bson:"mana_cost_reduction_per_level" json:"mana_cost_reduction_per_level"`
	} `bson:"scaling" json:"scaling"`
	Version int `bson:"version" json:"version"`
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
