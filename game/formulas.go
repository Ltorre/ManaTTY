package game

import (
	"math"
)

// CalculateManaPerSecond computes the mana generation rate.
// Formula: BaseMana * EraMultiplier * RitualBonus
func CalculateManaPerSecond(currentFloor int, currentEra int, activeRituals int) float64 {
	// Base mana from floor
	baseMana := BaseManaPerSecond + (ManaPerFloorBonus * float64(currentFloor))

	// Era multiplier from prestige
	eraMultiplier := EraMultiplierBase + (EraMultiplierPerEra * float64(currentEra))

	// Ritual bonus (capped at max)
	ritualBonus := 1.0 + (RitualBonusPerActive * float64(activeRituals))
	if ritualBonus > (1.0 + MaxRitualBonus) {
		ritualBonus = 1.0 + MaxRitualBonus
	}

	return baseMana * eraMultiplier * ritualBonus
}

// CalculateManaPerSecondWithBonuses includes permanent mana gen multiplier.
func CalculateManaPerSecondWithBonuses(currentFloor int, currentEra int, activeRituals int, permanentMultiplier float64) float64 {
	base := CalculateManaPerSecond(currentFloor, currentEra, activeRituals)
	return base * permanentMultiplier
}

// CalculateFloorCost computes the mana required to climb to the next floor.
// Formula: BaseCost * (Floor ^ Exponent)
func CalculateFloorCost(floor int) float64 {
	return BaseFloorCost * math.Pow(float64(floor), FloorCostExponent)
}

// CalculateTotalManaForFloor returns total mana needed to reach a floor from floor 1.
func CalculateTotalManaForFloor(targetFloor int) float64 {
	total := 0.0
	for f := 1; f < targetFloor; f++ {
		total += CalculateFloorCost(f)
	}
	return total
}

// CalculateFloorsFromMana determines how many floors can be climbed with given mana.
func CalculateFloorsFromMana(startFloor int, availableMana float64) (floorsClimbed int, remainingMana float64) {
	remainingMana = availableMana
	floorsClimbed = 0
	currentFloor := startFloor

	for {
		cost := CalculateFloorCost(currentFloor)
		if remainingMana < cost {
			break
		}
		remainingMana -= cost
		floorsClimbed++
		currentFloor++
	}

	return floorsClimbed, remainingMana
}

// CalculateEraMultiplier returns the multiplier for a given era.
func CalculateEraMultiplier(era int) float64 {
	return EraMultiplierBase + (EraMultiplierPerEra * float64(era))
}

// CalculateRitualBonus returns the total ritual bonus multiplier.
func CalculateRitualBonus(activeRituals int) float64 {
	bonus := RitualBonusPerActive * float64(activeRituals)
	if bonus > MaxRitualBonus {
		bonus = MaxRitualBonus
	}
	return 1.0 + bonus
}

// CalculateOfflineMana computes mana earned during offline time.
func CalculateOfflineMana(manaPerSecond float64, offlineSeconds float64) float64 {
	if offlineSeconds < float64(MinOfflineSeconds) {
		return 0
	}
	return manaPerSecond * offlineSeconds * OfflinePenalty
}

// CalculateSpellCooldown returns cooldown after applying reduction bonuses.
func CalculateSpellCooldown(baseCooldownMs int64, cooldownReduction float64) int64 {
	if cooldownReduction > MaxCooldownReduction {
		cooldownReduction = MaxCooldownReduction
	}
	cooldown := float64(baseCooldownMs) * (1.0 - cooldownReduction)
	if cooldown < float64(MinSpellCooldownMs) {
		cooldown = float64(MinSpellCooldownMs)
	}
	return int64(cooldown)
}

// CalculateManualCastCost returns the mana cost for a manual spell cast.
func CalculateManualCastCost(baseManaCost float64) float64 {
	return baseManaCost * (1.0 + ManualCastPenalty)
}

// CanPrestige returns true if the player can prestige at the current floor.
func CanPrestige(currentFloor int) bool {
	return currentFloor >= PrestigeFloor
}

// CalculatePrestigeBonuses returns the bonuses gained from a prestige.
type PrestigeBonuses struct {
	NewEra              int
	NewEraMultiplier    float64
	AddedManaGen        float64
	AddedCooldownRedux  float64
	AddedManaRetention  float64
	NewRitualCapacity   int
}

// GetPrestigeBonuses calculates what bonuses will be gained from prestiging.
func GetPrestigeBonuses(currentEra int, currentRitualCap int) PrestigeBonuses {
	newEra := currentEra + 1
	newCap := currentRitualCap + 1
	if newCap > MaxActiveRituals {
		newCap = MaxActiveRituals
	}

	return PrestigeBonuses{
		NewEra:              newEra,
		NewEraMultiplier:    CalculateEraMultiplier(newEra),
		AddedManaGen:        PrestigeManaGenBonus,
		AddedCooldownRedux:  PrestigeCooldownBonus,
		AddedManaRetention:  PrestigeManaRetention,
		NewRitualCapacity:   newCap,
	}
}
