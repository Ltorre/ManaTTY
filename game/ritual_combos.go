package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Ltorre/ManaTTY/models"
)

// Ritual effect magnitudes
const (
	RitualPureMagnitude   = 0.18 // +18% for pure (3 same element)
	RitualHybridMagnitude = 0.12 // +12% for hybrid (2+1)
	RitualTriadMagnitude  = 0.08 // +8% each for triad (1/1/1)
	RitualEchoKicker      = 0.05 // +5% bonus when Spell Echo is included
)

// Element adjectives by count (1, 2, 3 of same element)
var elementAdjectives = map[models.Element][]string{
	models.ElementFire:    {"Blazing", "Infernal", "Volcanic"},
	models.ElementIce:     {"Frozen", "Glacial", "Permafrost"},
	models.ElementThunder: {"Storm", "Tempest", "Cataclysm"},
	models.ElementArcane:  {"Runic", "Ethereal", "Astral"},
}

// Spell power nouns (from highest-damage spell)
var spellNouns = map[string]string{
	"spell_fireball":        "Ember",
	"spell_inferno":         "Pyre",
	"spell_meteor_strike":   "Cataclysm",
	"spell_frostbolt":       "Shard",
	"spell_blizzard":        "Gale",
	"spell_frost_nova":      "Nova",
	"spell_lightning":       "Bolt",
	"spell_chain_lightning": "Arc",
	"spell_thunderstorm":    "Tempest",
	"spell_vortex":          "Vortex",
	"spell_echo":            "Echo",
	"spell_arcane_blast":    "Blast",
}

// Spell base damage for determining "highest-damage" spell
var spellBaseDamage = map[string]float64{
	"spell_fireball":        100,
	"spell_inferno":         350,
	"spell_meteor_strike":   1000,
	"spell_frostbolt":       80,
	"spell_blizzard":        450,
	"spell_frost_nova":      600,
	"spell_lightning":       120,
	"spell_chain_lightning": 400,
	"spell_thunderstorm":    700,
	"spell_vortex":          150,
	"spell_echo":            500,
	"spell_arcane_blast":    800,
}

// Spell elements
var spellElements = map[string]models.Element{
	"spell_fireball":        models.ElementFire,
	"spell_inferno":         models.ElementFire,
	"spell_meteor_strike":   models.ElementFire,
	"spell_frostbolt":       models.ElementIce,
	"spell_blizzard":        models.ElementIce,
	"spell_frost_nova":      models.ElementIce,
	"spell_lightning":       models.ElementThunder,
	"spell_chain_lightning": models.ElementThunder,
	"spell_thunderstorm":    models.ElementThunder,
	"spell_vortex":          models.ElementArcane,
	"spell_echo":            models.ElementArcane,
	"spell_arcane_blast":    models.ElementArcane,
}

// Signature combos with special flavor names
var signatureCombos = map[string]string{
	// Elemental Trinity: starter spells
	"spell_fireball+spell_frostbolt+spell_lightning": "Elemental Trinity",
	// Apocalypse: three highest-damage spells
	"spell_meteor_strike+spell_thunderstorm+spell_arcane_blast": "Apocalypse",
	// Convergence: mid-tier multi-element
	"spell_inferno+spell_blizzard+spell_chain_lightning": "Convergence",
}

// RitualComboInfo holds computed ritual information.
type RitualComboInfo struct {
	Name            string
	SignatureName   string
	Composition     models.RitualComposition
	DominantElement models.Element
	Effects         []models.RitualEffect
	HasSpellEcho    bool
}

// ComputeRitualCombo analyzes spell IDs and returns naming/effect info.
func ComputeRitualCombo(spellIDs []string) RitualComboInfo {
	if len(spellIDs) != 3 {
		return RitualComboInfo{Name: "Unknown Ritual"}
	}

	info := RitualComboInfo{}

	// Count elements
	elementCounts := make(map[models.Element]int)
	for _, id := range spellIDs {
		if elem, ok := spellElements[id]; ok {
			elementCounts[elem]++
		}
	}

	// Check for Spell Echo
	for _, id := range spellIDs {
		if id == "spell_echo" {
			info.HasSpellEcho = true
			break
		}
	}

	// Determine composition
	info.Composition, info.DominantElement = determineComposition(elementCounts)

	// Generate effects based on composition
	info.Effects = generateEffects(info.Composition, info.DominantElement, elementCounts, info.HasSpellEcho)

	// Generate name
	info.Name = generateRitualName(spellIDs, elementCounts, info.Composition)

	// Check for signature name
	info.SignatureName = checkSignatureName(spellIDs)

	return info
}

// determineComposition analyzes element counts to classify the ritual.
func determineComposition(counts map[models.Element]int) (models.RitualComposition, models.Element) {
	var maxElement models.Element
	maxCount := 0

	for elem, count := range counts {
		if count > maxCount {
			maxCount = count
			maxElement = elem
		}
	}

	switch maxCount {
	case 3:
		return models.CompositionPure, maxElement
	case 2:
		return models.CompositionHybrid, maxElement
	default:
		return models.CompositionTriad, ""
	}
}

// generateEffects creates the effect list based on composition.
func generateEffects(comp models.RitualComposition, dominant models.Element, counts map[models.Element]int, hasEcho bool) []models.RitualEffect {
	effects := []models.RitualEffect{}
	echoBonus := 0.0
	if hasEcho {
		echoBonus = RitualEchoKicker
	}

	switch comp {
	case models.CompositionPure:
		// Single element's signature bonus at +18%
		effect := models.RitualEffect{
			Type:      getElementEffectType(dominant),
			Magnitude: RitualPureMagnitude + echoBonus,
		}
		effects = append(effects, effect)

	case models.CompositionHybrid:
		// Dominant element's bonus at +12%
		effect := models.RitualEffect{
			Type:      getElementEffectType(dominant),
			Magnitude: RitualHybridMagnitude + echoBonus,
		}
		effects = append(effects, effect)

	case models.CompositionTriad:
		// All three elements get +8% each
		for elem := range counts {
			effect := models.RitualEffect{
				Type:      getElementEffectType(elem),
				Magnitude: RitualTriadMagnitude + echoBonus,
			}
			effects = append(effects, effect)
		}
	}

	return effects
}

// getElementEffectType returns the signature effect for an element.
func getElementEffectType(elem models.Element) models.RitualEffectType {
	switch elem {
	case models.ElementFire:
		return models.RitualEffectDamage
	case models.ElementIce:
		return models.RitualEffectCooldown
	case models.ElementThunder:
		return models.RitualEffectManaCost
	case models.ElementArcane:
		return models.RitualEffectSigilRate
	default:
		return models.RitualEffectDamage
	}
}

// generateRitualName creates the ritual name from element adjectives and spell noun.
func generateRitualName(spellIDs []string, elementCounts map[models.Element]int, comp models.RitualComposition) string {
	// Find highest-damage spell for the noun
	highestDamage := 0.0
	highestSpell := spellIDs[0]
	for _, id := range spellIDs {
		if dmg, ok := spellBaseDamage[id]; ok && dmg > highestDamage {
			highestDamage = dmg
			highestSpell = id
		}
	}

	noun := spellNouns[highestSpell]
	if noun == "" {
		noun = "Power"
	}

	// Build adjective(s)
	adjectives := []string{}

	// Get sorted elements for consistent ordering
	elements := make([]models.Element, 0, len(elementCounts))
	for elem := range elementCounts {
		elements = append(elements, elem)
	}
	sort.Slice(elements, func(i, j int) bool {
		return string(elements[i]) < string(elements[j])
	})

	for _, elem := range elements {
		count := elementCounts[elem]
		adjList := elementAdjectives[elem]
		if adjList == nil || count == 0 {
			continue
		}
		// Use tier based on count (0-indexed: 1->0, 2->1, 3->2)
		tier := count - 1
		if tier >= len(adjList) {
			tier = len(adjList) - 1
		}
		adjectives = append(adjectives, adjList[tier])
	}

	// Join adjectives
	adjStr := strings.Join(adjectives, " ")

	return "Ritual of " + adjStr + " " + noun
}

// checkSignatureName looks up special combo names.
func checkSignatureName(spellIDs []string) string {
	// Sort spell IDs for consistent key
	sorted := make([]string, len(spellIDs))
	copy(sorted, spellIDs)
	sort.Strings(sorted)
	key := strings.Join(sorted, "+")

	if name, ok := signatureCombos[key]; ok {
		return name
	}

	// Check for "Resonant [Element]" pattern: Echo + 2 same-element spells
	hasEcho := false
	elementCounts := make(map[models.Element]int)
	for _, id := range spellIDs {
		if id == "spell_echo" {
			hasEcho = true
		}
		if elem, ok := spellElements[id]; ok {
			elementCounts[elem]++
		}
	}

	if hasEcho {
		for elem, count := range elementCounts {
			if elem != models.ElementArcane && count >= 2 {
				return "Resonant " + elementDisplayName(elem)
			}
		}
	}

	return ""
}

// elementDisplayName returns a capitalized element name.
func elementDisplayName(elem models.Element) string {
	switch elem {
	case models.ElementFire:
		return "Fire"
	case models.ElementIce:
		return "Ice"
	case models.ElementThunder:
		return "Thunder"
	case models.ElementArcane:
		return "Arcane"
	default:
		return string(elem)
	}
}

// GetEffectDisplayString returns a human-readable effect description.
func GetEffectDisplayString(effect models.RitualEffect) string {
	sign := "+"
	suffix := ""
	switch effect.Type {
	case models.RitualEffectDamage:
		suffix = " dmg"
	case models.RitualEffectCooldown:
		sign = "-"
		suffix = " CD"
	case models.RitualEffectManaCost:
		sign = "-"
		suffix = " cost"
	case models.RitualEffectSigilRate:
		suffix = " sigil"
	}

	percent := int(effect.Magnitude * 100)
	return fmt.Sprintf("%s%d%%%s", sign, percent, suffix)
}

// GetRitualEffectIcon returns an emoji icon for the effect type.
func GetRitualEffectIcon(effectType models.RitualEffectType) string {
	switch effectType {
	case models.RitualEffectDamage:
		return "🔥"
	case models.RitualEffectCooldown:
		return "❄️"
	case models.RitualEffectManaCost:
		return "⚡"
	case models.RitualEffectSigilRate:
		return "✨"
	default:
		return "⭐"
	}
}
