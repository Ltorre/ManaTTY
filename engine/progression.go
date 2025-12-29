package engine

import (
	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
)

// CalculateFloorsFromMana determines how many floors can be climbed.
func (e *GameEngine) CalculateFloorsFromMana(startFloor int, mana float64) (int, float64) {
	return game.CalculateFloorsFromMana(startFloor, mana)
}

// GetManaCostForFloor returns the mana needed for a specific floor.
func (e *GameEngine) GetManaCostForFloor(floor int) float64 {
	return game.CalculateFloorCost(floor)
}

// GetTotalManaToFloor returns total mana to reach a floor from floor 1.
func (e *GameEngine) GetTotalManaToFloor(targetFloor int) float64 {
	return game.CalculateTotalManaForFloor(targetFloor)
}

// CanPrestige checks if the player can prestige.
func (e *GameEngine) CanPrestige(gs *models.GameState) bool {
	return game.CanPrestige(gs.Tower.CurrentFloor)
}

// GetPrestigePreview returns what bonuses will be gained from prestige.
func (e *GameEngine) GetPrestigePreview(gs *models.GameState) game.PrestigeBonuses {
	return game.GetPrestigeBonuses(
		gs.PrestigeData.CurrentEra,
		gs.PrestigeData.RitualCapacity,
	)
}

// ProcessPrestige performs the prestige/ascension action.
func (e *GameEngine) ProcessPrestige(gs *models.GameState) bool {
	if !e.CanPrestige(gs) {
		return false
	}

	// Get base spells for reset
	baseSpells := game.GetBaseSpells()

	// Process prestige (resets tower, applies bonuses)
	gs.ResetForPrestige(baseSpells)

	if e.OnPrestige != nil {
		e.OnPrestige(gs.PrestigeData.CurrentEra)
	}

	return true
}

// GetEraMultiplier returns the current era multiplier.
func (e *GameEngine) GetEraMultiplier(gs *models.GameState) float64 {
	return gs.PrestigeData.EraMultiplier
}

// GetTotalMultiplier returns the combined prestige multiplier.
func (e *GameEngine) GetTotalMultiplier(gs *models.GameState) float64 {
	return gs.PrestigeData.GetTotalMultiplier()
}

// GetRitualBonus returns the current ritual bonus multiplier.
func (e *GameEngine) GetRitualBonus(gs *models.GameState) float64 {
	return game.CalculateRitualBonus(len(gs.GetActiveRituals()))
}

// GetProgressStats returns comprehensive progression statistics.
func (e *GameEngine) GetProgressStats(gs *models.GameState) ProgressStats {
	return ProgressStats{
		CurrentFloor:    gs.Tower.CurrentFloor,
		MaxFloorReached: gs.Tower.MaxFloorReached,
		CurrentMana:     gs.Tower.CurrentMana,
		ManaToNextFloor: e.GetManaToNextFloor(gs),
		FloorProgress:   e.GetFloorProgress(gs),
		ManaPerSecond:   e.CalculateManaPerSecond(gs),
		LifetimeMana:    gs.Tower.LifetimeManaEarned,
		TotalSpells:     len(gs.Spells),
		ActiveRituals:   len(gs.GetActiveRituals()),
		CurrentEra:      gs.PrestigeData.CurrentEra,
		TotalAscensions: gs.PrestigeData.TotalAscensions,
		EraMultiplier:   e.GetEraMultiplier(gs),
		TotalMultiplier: e.GetTotalMultiplier(gs),
		RitualBonus:     e.GetRitualBonus(gs),
		CanPrestige:     e.CanPrestige(gs),
		TimeToNextFloor: e.GetTimeToNextFloor(gs),
	}
}

// ProgressStats holds comprehensive progression data.
type ProgressStats struct {
	CurrentFloor    int
	MaxFloorReached int
	CurrentMana     float64
	ManaToNextFloor float64
	FloorProgress   float64
	ManaPerSecond   float64
	LifetimeMana    float64
	TotalSpells     int
	ActiveRituals   int
	CurrentEra      int
	TotalAscensions int
	EraMultiplier   float64
	TotalMultiplier float64
	RitualBonus     float64
	CanPrestige     bool
	TimeToNextFloor interface{}
}
