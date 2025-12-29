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

	// Spell Leveling
	SpellUpgradeBaseCost     = 500.0 // Base mana cost to upgrade a spell
	SpellUpgradeCostExponent = 1.5   // Cost scaling per level
	SpellMaxLevel            = 10    // Maximum spell level
	SpellCooldownPerLevel    = 0.05  // -5% cooldown per level
	SpellManaCostPerLevel    = 0.08  // -8% mana cost per level
	SpellDamagePerLevel      = 0.15  // +15% damage per level

	// Spell Specializations (chosen at levels 5 and 10)
	SpecTier1Level          = 5    // Level to unlock tier 1 specialization
	SpecTier2Level          = 10   // Level to unlock tier 2 specialization
	SpecCritChanceBonus     = 0.15 // +15% crit chance (Tier 1)
	SpecCritDamageMulti     = 2.0  // 2x damage on crit
	SpecManaEfficiencyBonus = 0.20 // -20% mana cost (Tier 1)
	SpecBurstDamageBonus    = 0.30 // +30% damage (Tier 2)
	SpecRapidCastBonus      = 0.25 // -25% cooldown (Tier 2)

	// Element Synergies
	ElementStreakRequired  = 3    // Casts of same element to trigger synergy
	ElementSynergyDuration = 10.0 // Seconds the synergy buff lasts
	ElementSynergyBonus    = 0.20 // +20% bonus during synergy

	// Ascension Sigil - damage requirement to climb floors
	SigilBaseDamage    = 200.0 // Base damage used in sigil requirement formula (floor 1 baseline)
	SigilScaleExponent = 1.6   // Scaling per floor (tuned to stay relevant vs mana costs)
	SigilFloorFactor   = 1.0   // Multiplier to tune overall gate strength

	// Offline Progress
	OfflinePenalty    = 0.50 // 50% mana generation while offline
	MinOfflineSeconds = 1    // Minimum offline time to process

	// Game Loop
	DefaultTickRateHz   = 10 // 10 ticks per second
	AutoSaveIntervalSec = 30 // Auto-save every 30 seconds

	// UI
	ProgressBarWidth = 20 // Characters in progress bar
)
