package models

// TowerState represents the player's progress in the tower.
type TowerState struct {
	CurrentFloor       int     `bson:"current_floor" json:"current_floor"`
	MaxFloorReached    int     `bson:"max_floor_reached" json:"max_floor_reached"`
	CurrentMana        float64 `bson:"current_mana" json:"current_mana"`
	MaxMana            float64 `bson:"max_mana" json:"max_mana"`
	LifetimeManaEarned float64 `bson:"lifetime_mana_earned" json:"lifetime_mana_earned"`
}

// NewTowerState creates a new tower state at floor 1.
func NewTowerState() *TowerState {
	return &TowerState{
		CurrentFloor:       1,
		MaxFloorReached:    1,
		CurrentMana:        0,
		MaxMana:            100,
		LifetimeManaEarned: 0,
	}
}

// AddMana adds mana to the current total and lifetime count.
func (t *TowerState) AddMana(amount float64) {
	t.CurrentMana += amount
	t.LifetimeManaEarned += amount
}

// SpendMana deducts mana if available.
func (t *TowerState) SpendMana(amount float64) bool {
	if t.CurrentMana >= amount {
		t.CurrentMana -= amount
		return true
	}
	return false
}

// ClimbFloor increments the floor and updates max reached.
func (t *TowerState) ClimbFloor() {
	t.CurrentFloor++
	if t.CurrentFloor > t.MaxFloorReached {
		t.MaxFloorReached = t.CurrentFloor
	}
}

// Reset resets the tower state for prestige.
func (t *TowerState) Reset() {
	t.CurrentFloor = 1
	t.CurrentMana = 0
}

// GetFloorProgress returns progress towards next floor (0.0 to 1.0).
func (t *TowerState) GetFloorProgress() float64 {
	if t.MaxMana <= 0 {
		return 0
	}
	progress := t.CurrentMana / t.MaxMana
	if progress > 1.0 {
		progress = 1.0
	}
	return progress
}
