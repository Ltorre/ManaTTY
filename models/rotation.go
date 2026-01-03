package models

// RotationCondition defines advanced conditions for spell rotation.
type RotationCondition string

const (
	// Basic conditions (v1.2)
	RotationConditionAlways        RotationCondition = "always"
	RotationConditionManaAbove50   RotationCondition = "mana_above_50"
	RotationConditionManaAbove75   RotationCondition = "mana_above_75"
	RotationConditionSigilNotFull  RotationCondition = "sigil_not_full"
	RotationConditionSynergyActive RotationCondition = "synergy_active"

	// v1.5.0: Advanced conditions
	RotationConditionManaEfficient  RotationCondition = "mana_efficient"    // Only when mana/cost ratio > 2.0
	RotationConditionDuringSynergy  RotationCondition = "during_synergy"    // Only when ritual synergy active for this element
	RotationConditionSigilAlmostFul RotationCondition = "sigil_almost_full" // Only when sigil > 80%
	RotationConditionHighPriority   RotationCondition = "high_priority"     // Always cast ASAP (ignore mana efficiency)
	RotationConditionFillerOnly     RotationCondition = "filler_only"       // Only when nothing else can cast
)

// RotationPriority defines spell priority tiers.
type RotationPriority int

const (
	PriorityHigh   RotationPriority = 1 // Cast ASAP when ready
	PriorityMedium RotationPriority = 2 // Normal rotation priority
	PriorityLow    RotationPriority = 3 // Filler spell
)

// RotationSpellConfig defines advanced rotation settings for a spell.
type RotationSpellConfig struct {
	SpellID   string            `bson:"spell_id" json:"spell_id"`
	Priority  RotationPriority  `bson:"priority" json:"priority"`
	Condition RotationCondition `bson:"condition" json:"condition"`
	Enabled   bool              `bson:"enabled" json:"enabled"` // Can disable without removing from rotation
}

// SpellRotation holds the complete rotation configuration.
type SpellRotation struct {
	Enabled         bool                  `bson:"enabled" json:"enabled"`                     // Master toggle
	Spells          []RotationSpellConfig `bson:"spells" json:"spells"`                       // Spell configurations
	CooldownWeaving bool                  `bson:"cooldown_weaving" json:"cooldown_weaving"`   // Smart cooldown management
	ManaThreshold   float64               `bson:"mana_threshold" json:"mana_threshold"`       // Reserve this % of mana
	OptimizeForIdle bool                  `bson:"optimize_for_idle" json:"optimize_for_idle"` // Prioritize sustained DPS over burst
}

// GetConditionDescription returns a human-readable description for a rotation condition.
func GetConditionDescription(cond RotationCondition) string {
	switch cond {
	case RotationConditionAlways:
		return "Always cast when ready"
	case RotationConditionManaAbove50:
		return "Only when mana > 50%"
	case RotationConditionManaAbove75:
		return "Only when mana > 75%"
	case RotationConditionSigilNotFull:
		return "Only when sigil not full"
	case RotationConditionSynergyActive:
		return "Only during element synergy"
	case RotationConditionManaEfficient:
		return "Only when mana-efficient (ratio > 2.0)"
	case RotationConditionDuringSynergy:
		return "Only during ritual synergy for element"
	case RotationConditionSigilAlmostFul:
		return "Only when sigil > 80%"
	case RotationConditionHighPriority:
		return "High priority - cast ASAP"
	case RotationConditionFillerOnly:
		return "Filler - only when nothing else ready"
	default:
		return "Unknown condition"
	}
}

// GetPriorityLabel returns a label for priority tier.
func GetPriorityLabel(priority RotationPriority) string {
	switch priority {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "Unknown"
	}
}

// DefaultRotation creates a default rotation configuration from auto-cast slots.
func DefaultRotation() *SpellRotation {
	return &SpellRotation{
		Enabled:         false, // Disabled by default - user must enable
		Spells:          []RotationSpellConfig{},
		CooldownWeaving: true,
		ManaThreshold:   0.1, // Reserve 10% mana by default
		OptimizeForIdle: true,
	}
}
