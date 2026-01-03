package models

// RitualSynergyType defines types of ritual synergies.
type RitualSynergyType string

const (
	// Element pair synergies
	SynergyThermalShock    RitualSynergyType = "thermal_shock"    // Fire + Ice
	SynergyManaConduit     RitualSynergyType = "mana_conduit"     // Thunder + Arcane
	SynergyVolcanicFury    RitualSynergyType = "volcanic_fury"    // Fire + Thunder
	SynergyFrozenLightning RitualSynergyType = "frozen_lightning" // Ice + Thunder
	SynergyArcaneInferno   RitualSynergyType = "arcane_inferno"   // Fire + Arcane
	SynergyGlacialMystic   RitualSynergyType = "glacial_mystic"   // Ice + Arcane
)

// RitualSynergy represents an active synergy bonus between rituals.
type RitualSynergy struct {
	Type        RitualSynergyType `bson:"type" json:"type"`
	Name        string            `bson:"name" json:"name"`
	Description string            `bson:"description" json:"description"`
	Elements    []Element         `bson:"elements" json:"elements"`   // Required elements
	Magnitude   float64           `bson:"magnitude" json:"magnitude"` // Bonus multiplier (e.g., 0.15 for +15%)
}

// SynergyDefinitions maps synergy types to their configurations.
var SynergyDefinitions = map[RitualSynergyType]RitualSynergy{
	SynergyThermalShock: {
		Type:        SynergyThermalShock,
		Name:        "Thermal Shock",
		Description: "Fire and Ice combine for explosive power",
		Elements:    []Element{ElementFire, ElementIce},
		Magnitude:   0.15, // +15% to both Fire and Ice ritual effects
	},
	SynergyManaConduit: {
		Type:        SynergyManaConduit,
		Name:        "Mana Conduit",
		Description: "Thunder and Arcane create perfect energy flow",
		Elements:    []Element{ElementThunder, ElementArcane},
		Magnitude:   0.20, // +20% to both Thunder and Arcane ritual effects
	},
	SynergyVolcanicFury: {
		Type:        SynergyVolcanicFury,
		Name:        "Volcanic Fury",
		Description: "Fire and Thunder unleash devastating force",
		Elements:    []Element{ElementFire, ElementThunder},
		Magnitude:   0.12, // +12% to both Fire and Thunder ritual effects
	},
	SynergyFrozenLightning: {
		Type:        SynergyFrozenLightning,
		Name:        "Frozen Lightning",
		Description: "Ice and Thunder create crystalline energy",
		Elements:    []Element{ElementIce, ElementThunder},
		Magnitude:   0.12, // +12% to both Ice and Thunder ritual effects
	},
	SynergyArcaneInferno: {
		Type:        SynergyArcaneInferno,
		Name:        "Arcane Inferno",
		Description: "Fire and Arcane merge raw power with mystic energy",
		Elements:    []Element{ElementFire, ElementArcane},
		Magnitude:   0.18, // +18% to both Fire and Arcane ritual effects
	},
	SynergyGlacialMystic: {
		Type:        SynergyGlacialMystic,
		Name:        "Glacial Mystic",
		Description: "Ice and Arcane weave time and space",
		Elements:    []Element{ElementIce, ElementArcane},
		Magnitude:   0.18, // +18% to both Ice and Arcane ritual effects
	},
}
