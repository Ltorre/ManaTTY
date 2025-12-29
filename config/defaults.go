package config

// DefaultGameConstants contains default game balance values.
var DefaultGameConstants = GameConstants{
	// Mana Generation
	BaseManaPerSecond:    10.0,
	ManaPerFloorBonus:    2.0,
	RitualBonusPerActive: 0.15,
	MaxRitualBonus:       0.45,

	// Floor Climbing
	BaseFloorCost:     100.0,
	FloorCostExponent: 1.8,
	PrestigeFloor:     100,

	// Prestige System
	EraMultiplierBase:     1.0,
	EraMultiplierPerEra:   0.15,
	PrestigeManaGenBonus:  0.05,
	PrestigeCooldownBonus: 0.05,
	PrestigeManaRetention: 0.10,
	MaxCooldownReduction:  0.50,
	MaxManaRetention:      0.90,

	// Rituals
	RitualCooldownMs:  60000,
	MaxActiveRituals:  3,
	SpellsPerRitual:   3,

	// Spells
	MinSpellCooldownMs: 1000,
	ManualCastPenalty:  0.10,

	// Offline Progress
	OfflinePenalty:    0.50,
	MinOfflineSeconds: 1,
}

// GameConstants holds all game balance constants.
type GameConstants struct {
	// Mana Generation
	BaseManaPerSecond    float64
	ManaPerFloorBonus    float64
	RitualBonusPerActive float64
	MaxRitualBonus       float64

	// Floor Climbing
	BaseFloorCost     float64
	FloorCostExponent float64
	PrestigeFloor     int

	// Prestige System
	EraMultiplierBase     float64
	EraMultiplierPerEra   float64
	PrestigeManaGenBonus  float64
	PrestigeCooldownBonus float64
	PrestigeManaRetention float64
	MaxCooldownReduction  float64
	MaxManaRetention      float64

	// Rituals
	RitualCooldownMs int64
	MaxActiveRituals int
	SpellsPerRitual  int

	// Spells
	MinSpellCooldownMs int64
	ManualCastPenalty  float64

	// Offline Progress
	OfflinePenalty    float64
	MinOfflineSeconds int
}
