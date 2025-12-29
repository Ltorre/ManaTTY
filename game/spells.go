package game

import "github.com/Ltorre/ManaTTY/models"

// DefaultSpells returns the list of all spell definitions.
func DefaultSpells() []*models.SpellDefinition {
	return []*models.SpellDefinition{
		// FIRE Element
		{
			ID:                "spell_fireball",
			Name:              "Fireball",
			Description:       "Hurl flames at enemies, dealing fire damage.",
			FlavorText:        "A classic spell of arcane mastery.",
			Element:           models.ElementFire,
			BaseDamage:        100,
			BaseCooldownMs:    3000,
			BaseManaCost:      50,
			RequiredFloor:     1,
			UnlockedByDefault: true,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_inferno",
			Name:              "Inferno",
			Description:       "Unleash a devastating wave of fire.",
			FlavorText:        "The flames of destruction.",
			Element:           models.ElementFire,
			BaseDamage:        350,
			BaseCooldownMs:    5000,
			BaseManaCost:      150,
			RequiredFloor:     25,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_meteor_strike",
			Name:              "Meteor Strike",
			Description:       "Call down a meteor from the heavens.",
			FlavorText:        "The ultimate fire spell, reserved for the worthy.",
			Element:           models.ElementFire,
			BaseDamage:        1000,
			BaseCooldownMs:    15000,
			BaseManaCost:      500,
			RequiredFloor:     75,
			UnlockedByDefault: false,
			PrestigeExclusive: true,
			Version:           1,
		},

		// ICE Element
		{
			ID:                "spell_frostbolt",
			Name:              "Frostbolt",
			Description:       "Launch a bolt of freezing ice.",
			FlavorText:        "Cold as the northern winds.",
			Element:           models.ElementIce,
			BaseDamage:        80,
			BaseCooldownMs:    4000,
			BaseManaCost:      60,
			RequiredFloor:     3,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_blizzard",
			Name:              "Blizzard",
			Description:       "Summon a raging blizzard around your enemies.",
			FlavorText:        "Winter's wrath unleashed.",
			Element:           models.ElementIce,
			BaseDamage:        450,
			BaseCooldownMs:    6000,
			BaseManaCost:      200,
			RequiredFloor:     35,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_frost_nova",
			Name:              "Frost Nova",
			Description:       "Release an expanding ring of ice.",
			FlavorText:        "Freeze them in their tracks.",
			Element:           models.ElementIce,
			BaseDamage:        600,
			BaseCooldownMs:    8000,
			BaseManaCost:      300,
			RequiredFloor:     60,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},

		// THUNDER Element
		{
			ID:                "spell_lightning",
			Name:              "Lightning",
			Description:       "Strike with a bolt of lightning.",
			FlavorText:        "Swift as the storm.",
			Element:           models.ElementThunder,
			BaseDamage:        120,
			BaseCooldownMs:    5000,
			BaseManaCost:      75,
			RequiredFloor:     5,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_chain_lightning",
			Name:              "Chain Lightning",
			Description:       "Lightning that jumps between enemies.",
			FlavorText:        "The storm spreads.",
			Element:           models.ElementThunder,
			BaseDamage:        400,
			BaseCooldownMs:    7000,
			BaseManaCost:      250,
			RequiredFloor:     45,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_thunderstorm",
			Name:              "Thunderstorm",
			Description:       "Call down a devastating thunderstorm.",
			FlavorText:        "Nature's fury.",
			Element:           models.ElementThunder,
			BaseDamage:        700,
			BaseCooldownMs:    10000,
			BaseManaCost:      400,
			RequiredFloor:     70,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},

		// ARCANE Element
		{
			ID:                "spell_vortex",
			Name:              "Arcane Vortex",
			Description:       "Create a swirling vortex of arcane energy.",
			FlavorText:        "Pure magical chaos.",
			Element:           models.ElementArcane,
			BaseDamage:        150,
			BaseCooldownMs:    4500,
			BaseManaCost:      100,
			RequiredFloor:     10,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_echo",
			Name:              "Spell Echo",
			Description:       "Echo your previous spell for double effect.",
			FlavorText:        "Magic repeats itself.",
			Element:           models.ElementArcane,
			BaseDamage:        500,
			BaseCooldownMs:    6000,
			BaseManaCost:      225,
			RequiredFloor:     55,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
		{
			ID:                "spell_arcane_blast",
			Name:              "Arcane Blast",
			Description:       "A concentrated blast of pure arcane power.",
			FlavorText:        "The pinnacle of arcane mastery.",
			Element:           models.ElementArcane,
			BaseDamage:        800,
			BaseCooldownMs:    12000,
			BaseManaCost:      450,
			RequiredFloor:     80,
			UnlockedByDefault: false,
			PrestigeExclusive: false,
			Version:           1,
		},
	}
}

// GetSpellDefinition returns a spell definition by ID.
func GetSpellDefinition(id string) *models.SpellDefinition {
	for _, spell := range DefaultSpells() {
		if spell.ID == id {
			return spell
		}
	}
	return nil
}

// GetBaseSpells returns spells that are unlocked by default.
func GetBaseSpells() []*models.Spell {
	spells := []*models.Spell{}
	for _, def := range DefaultSpells() {
		if def.UnlockedByDefault {
			spells = append(spells, models.NewSpellFromDefinition(def))
		}
	}
	return spells
}

// GetSpellsForFloor returns all spells unlocked at or before a floor.
func GetSpellsForFloor(floor int, includePrestige bool) []*models.SpellDefinition {
	spells := []*models.SpellDefinition{}
	for _, def := range DefaultSpells() {
		if def.RequiredFloor <= floor {
			if !def.PrestigeExclusive || includePrestige {
				spells = append(spells, def)
			}
		}
	}
	return spells
}

// GetNewSpellsAtFloor returns spells unlocked exactly at a floor.
func GetNewSpellsAtFloor(floor int) []*models.SpellDefinition {
	spells := []*models.SpellDefinition{}
	for _, def := range DefaultSpells() {
		if def.RequiredFloor == floor && !def.PrestigeExclusive {
			spells = append(spells, def)
		}
	}
	return spells
}
