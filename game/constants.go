package game

// Game Balance Constants
const (
	// Mana Generation
	BaseManaPerSecond    = 10.0 // Base mana generated per second
	ManaPerFloorBonus    = 2.0  // Additional mana per floor
	RitualBonusPerActive = 0.15 // +15% mana gen per active ritual
	MaxRitualBonus       = 0.45 // Maximum ritual bonus (3 rituals)

	// Floor Climbing
	BaseFloorCost     = 100.0 // Base mana cost for floor 1
	FloorCostExponent = 1.8   // Exponential scaling for floor costs
	PrestigeFloor     = 100   // Floor required for prestige

	// Prestige System
	EraMultiplierBase     = 1.0  // Starting multiplier
	EraMultiplierPerEra   = 0.15 // +15% per era
	PrestigeManaGenBonus  = 0.05 // +5% permanent mana gen per prestige
	PrestigeCooldownBonus = 0.05 // +5% cooldown reduction per prestige
	PrestigeManaRetention = 0.10 // +10% mana retained per prestige
	MaxCooldownReduction  = 0.50 // Cap at 50% cooldown reduction
	MaxManaRetention      = 0.90 // Cap at 90% mana retention

	// Rituals
	RitualCooldownMs = 60000 // 60 seconds
	MaxActiveRituals = 3     // Maximum active rituals
	SpellsPerRitual  = 3     // Spells needed to form a ritual

	// Spells
	MinSpellCooldownMs = 1000 // Minimum 1 second cooldown
	ManualCastPenalty  = 0.10 // 10% extra mana cost for manual cast

	// Offline Progress
	OfflinePenalty    = 0.50 // 50% mana generation while offline
	MinOfflineSeconds = 1    // Minimum offline time to process

	// Game Loop
	DefaultTickRateHz   = 10 // 10 ticks per second
	AutoSaveIntervalSec = 30 // Auto-save every 30 seconds

	// UI
	ProgressBarWidth = 20 // Characters in progress bar
)
