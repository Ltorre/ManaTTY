package models

import "time"

// FloorEventChoice is the chosen bonus type for a floor event.
type FloorEventChoice string

const (
	FloorEventChoiceManaGen           FloorEventChoice = "mana_gen"
	FloorEventChoiceSigilChargeRate   FloorEventChoice = "sigil_charge_rate"
	FloorEventChoiceCooldownReduction FloorEventChoice = "cooldown_reduction"
)

// FloorEventChoiceDisplayNames provides consistent UI names.
var FloorEventChoiceDisplayNames = map[FloorEventChoice]string{
	FloorEventChoiceManaGen:           "Mana Surge",
	FloorEventChoiceSigilChargeRate:   "Sigil Attunement",
	FloorEventChoiceCooldownReduction: "Timewarp",
}

// FloorEventState represents an active, time-limited choice awaiting input.
type FloorEventState struct {
	Floor        int   `bson:"floor" json:"floor"`
	AppearedAtMs int64 `bson:"appeared_at_ms" json:"appeared_at_ms"`
	ExpiresAtMs  int64 `bson:"expires_at_ms" json:"expires_at_ms"`
}

// FloorEventBuff is a temporary bonus earned from a floor event.
type FloorEventBuff struct {
	Choice         FloorEventChoice `bson:"choice" json:"choice"`
	StartedAtFloor int              `bson:"started_at_floor" json:"started_at_floor"`
	ExpiresAtFloor int              `bson:"expires_at_floor" json:"expires_at_floor"`
}

func (gs *GameState) HasActiveFloorEvent(nowMs int64) bool {
	if gs == nil || gs.Session == nil || gs.Session.ActiveFloorEvent == nil {
		return false
	}
	return nowMs < gs.Session.ActiveFloorEvent.ExpiresAtMs
}

// ClearFloorEvent dismisses the current event without granting a bonus.
func (gs *GameState) ClearFloorEvent() {
	if gs == nil || gs.Session == nil {
		return
	}
	gs.Session.ActiveFloorEvent = nil
}

// EnsureFloorEventExpiry clears an expired event (no bonus) and returns true if it expired.
func (gs *GameState) EnsureFloorEventExpiry(nowMs int64) bool {
	if gs == nil || gs.Session == nil || gs.Session.ActiveFloorEvent == nil {
		return false
	}
	if nowMs >= gs.Session.ActiveFloorEvent.ExpiresAtMs {
		gs.Session.ActiveFloorEvent = nil
		return true
	}
	return false
}

// StartFloorEvent creates a new time-limited event.
func (gs *GameState) StartFloorEvent(floor int, now time.Time, timeoutMs int64) {
	if gs == nil || gs.Session == nil {
		return
	}
	nowMs := now.UnixMilli()
	gs.Session.ActiveFloorEvent = &FloorEventState{
		Floor:        floor,
		AppearedAtMs: nowMs,
		ExpiresAtMs:  nowMs + timeoutMs,
	}
}

func (gs *GameState) HasActiveFloorEventBuff(currentFloor int) bool {
	if gs == nil || gs.Session == nil || gs.Session.ActiveFloorBuff == nil {
		return false
	}
	return currentFloor < gs.Session.ActiveFloorBuff.ExpiresAtFloor
}

// GetActiveFloorBuffChoice returns the active floor-buff choice, or "" if none is active.
// This mirrors the pattern used by GetActiveSynergy().
func (gs *GameState) GetActiveFloorBuffChoice(currentFloor int) FloorEventChoice {
	if !gs.HasActiveFloorEventBuff(currentFloor) {
		return ""
	}
	// HasActiveFloorEventBuff ensures Session and ActiveFloorBuff are non-nil.
	return gs.Session.ActiveFloorBuff.Choice
}

func (gs *GameState) MaybeExpireFloorEventBuff(currentFloor int) bool {
	if gs == nil || gs.Session == nil || gs.Session.ActiveFloorBuff == nil {
		return false
	}
	if currentFloor >= gs.Session.ActiveFloorBuff.ExpiresAtFloor {
		gs.Session.ActiveFloorBuff = nil
		return true
	}
	return false
}

// ApplyFloorEventChoice grants a temporary bonus for a fixed number of floors.
// This also clears the pending floor event.
func (gs *GameState) ApplyFloorEventChoice(choice FloorEventChoice, currentFloor int, durationFloors int) {
	if gs == nil || gs.Session == nil {
		return
	}
	gs.Session.ActiveFloorBuff = &FloorEventBuff{
		Choice:         choice,
		StartedAtFloor: currentFloor,
		ExpiresAtFloor: currentFloor + durationFloors,
	}
	gs.Session.ActiveFloorEvent = nil
}
