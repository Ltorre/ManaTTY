package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GameState represents the complete state of a game save.
type GameState struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PlayerUUID        string             `bson:"player_uuid" json:"player_uuid"`
	Slot              int                `bson:"slot" json:"slot"`
	Tower             *TowerState        `bson:"tower" json:"tower"`
	Spells            []*Spell           `bson:"spells" json:"spells"`
	UnlockedSpellIDs  []string           `bson:"unlocked_spell_ids" json:"unlocked_spell_ids"`
	Rituals           []*Ritual          `bson:"rituals" json:"rituals"`
	ActiveRitualCount int                `bson:"active_ritual_count" json:"active_ritual_count"`
	PassiveBonuses    *PassiveBonuses    `bson:"passive_bonuses" json:"passive_bonuses"`
	PrestigeData      *PrestigeData      `bson:"prestige" json:"prestige"`
	Session           *SessionData       `bson:"session" json:"session"`
	SavedAt           time.Time          `bson:"saved_at" json:"saved_at"`
	Version           int                `bson:"version" json:"version"`
}

// PassiveBonuses contains modifiers that affect gameplay.
type PassiveBonuses struct {
	ManaGenMultiplier      float64 `bson:"mana_gen_multiplier" json:"mana_gen_multiplier"`
	FloorClimbSpeed        float64 `bson:"floor_climb_speed" json:"floor_climb_speed"`
	SpellCooldownReduction float64 `bson:"spell_cooldown_reduction" json:"spell_cooldown_reduction"`
	RitualCapacity         int     `bson:"ritual_capacity" json:"ritual_capacity"`
}

// AutoCastCondition defines when an auto-cast slot should trigger.
type AutoCastCondition string

const (
	ConditionAlways        AutoCastCondition = "always"         // Always cast when ready
	ConditionManaAbove50   AutoCastCondition = "mana_above_50"  // Only if mana > 50%
	ConditionManaAbove75   AutoCastCondition = "mana_above_75"  // Only if mana > 75%
	ConditionSigilNotFull  AutoCastCondition = "sigil_not_full" // Only if sigil not charged
	ConditionSynergyActive AutoCastCondition = "synergy_active" // Only during element synergy
)

// AutoCastSlotConfig holds the spell ID and condition for an auto-cast slot.
type AutoCastSlotConfig struct {
	SpellID   string            `bson:"spell_id" json:"spell_id"`
	Condition AutoCastCondition `bson:"condition" json:"condition"`
}

// SessionData contains current play session information.
type SessionData struct {
	SessionStartMs  int64     `bson:"session_start_ms" json:"session_start_ms"`
	SessionDuration int64     `bson:"session_duration_ms" json:"session_duration_ms"`
	LastTickMs      int64     `bson:"last_tick_ms" json:"last_tick_ms"`
	LastSavedAt     time.Time `bson:"last_saved_at" json:"last_saved_at"`
	AutoCastEnabled bool      `bson:"auto_cast_enabled" json:"auto_cast_enabled"`
	AutoCastSlots   []string  `bson:"auto_cast_slots" json:"auto_cast_slots"` // Spell IDs in auto-cast slots (legacy, kept for compat)

	// Auto-cast slot configurations with conditions
	AutoCastConfigs []AutoCastSlotConfig `bson:"auto_cast_configs" json:"auto_cast_configs"`

	// Element synergy tracking
	LastCastElements   []Element `bson:"last_cast_elements" json:"last_cast_elements"`       // Recent cast elements (up to 3)
	ActiveSynergy      Element   `bson:"active_synergy" json:"active_synergy"`               // Currently active synergy element
	SynergyExpiresAtMs int64     `bson:"synergy_expires_at_ms" json:"synergy_expires_at_ms"` // When synergy expires

	// Aggregated notifications
	AutoCastSkipCount int `bson:"-" json:"-"` // Transient: skipped auto-casts this second
}

// NewGameState creates a new game state with defaults.
func NewGameState(playerUUID string, slot int) *GameState {
	now := time.Now()
	return &GameState{
		PlayerUUID:        playerUUID,
		Slot:              slot,
		Tower:             NewTowerState(),
		Spells:            []*Spell{},
		UnlockedSpellIDs:  []string{},
		Rituals:           []*Ritual{},
		ActiveRitualCount: 0,
		PassiveBonuses:    NewPassiveBonuses(),
		PrestigeData:      NewPrestigeData(),
		Session:           NewSessionData(),
		SavedAt:           now,
		Version:           1,
	}
}

// NewPassiveBonuses creates default passive bonuses.
func NewPassiveBonuses() *PassiveBonuses {
	return &PassiveBonuses{
		ManaGenMultiplier:      1.0,
		FloorClimbSpeed:        1.0,
		SpellCooldownReduction: 0.0,
		RitualCapacity:         1,
	}
}

// NewSessionData creates a new session.
func NewSessionData() *SessionData {
	now := time.Now()
	return &SessionData{
		SessionStartMs:     now.UnixMilli(),
		SessionDuration:    0,
		LastTickMs:         now.UnixMilli(),
		LastSavedAt:        now,
		AutoCastEnabled:    true,
		AutoCastSlots:      []string{},             // Legacy, kept for compatibility
		AutoCastConfigs:    []AutoCastSlotConfig{}, // New: slots with conditions
		LastCastElements:   []Element{},
		ActiveSynergy:      "",
		SynergyExpiresAtMs: 0,
		AutoCastSkipCount:  0,
	}
}

// GetSpellByID returns a spell from the player's list by ID.
func (gs *GameState) GetSpellByID(spellID string) *Spell {
	for _, spell := range gs.Spells {
		if spell.ID == spellID {
			return spell
		}
	}
	return nil
}

// HasSpell returns true if the player has unlocked a spell.
func (gs *GameState) HasSpell(spellID string) bool {
	for _, id := range gs.UnlockedSpellIDs {
		if id == spellID {
			return true
		}
	}
	return false
}

// AddSpell adds a new spell to the player's collection.
func (gs *GameState) AddSpell(spell *Spell) {
	if !gs.HasSpell(spell.ID) {
		gs.Spells = append(gs.Spells, spell)
		gs.UnlockedSpellIDs = append(gs.UnlockedSpellIDs, spell.ID)
	}
}

// GetActiveRituals returns only active rituals.
func (gs *GameState) GetActiveRituals() []*Ritual {
	active := []*Ritual{}
	for _, r := range gs.Rituals {
		if r.IsActive {
			active = append(active, r)
		}
	}
	return active
}

// CanAddRitual returns true if there's capacity for another ritual.
func (gs *GameState) CanAddRitual() bool {
	activeCount := len(gs.GetActiveRituals())
	return activeCount < gs.PrestigeData.RitualCapacity
}

// UpdateSession updates session timing data.
func (gs *GameState) UpdateSession() {
	now := time.Now()
	gs.Session.LastTickMs = now.UnixMilli()
	gs.Session.SessionDuration = now.UnixMilli() - gs.Session.SessionStartMs
}

// IsSpellInAutoCast returns true if a spell is in an auto-cast slot.
func (gs *GameState) IsSpellInAutoCast(spellID string) bool {
	// Check new config system first
	for _, cfg := range gs.Session.AutoCastConfigs {
		if cfg.SpellID == spellID {
			return true
		}
	}
	// Fallback to legacy slots for backward compatibility
	for _, id := range gs.Session.AutoCastSlots {
		if id == spellID {
			return true
		}
	}
	return false
}

// GetAutoCastSlotCount returns max auto-cast slots (base 2 + prestige bonuses).
func (gs *GameState) GetAutoCastSlotCount() int {
	base := 2
	prestigeBonus := gs.PrestigeData.AutoCastSlotBonus
	return base + prestigeBonus
}

// GetAvailableAutoCastSlots returns remaining slot capacity.
func (gs *GameState) GetAvailableAutoCastSlots() int {
	usedSlots := len(gs.Session.AutoCastConfigs)
	if usedSlots == 0 {
		usedSlots = len(gs.Session.AutoCastSlots) // Legacy fallback
	}
	return gs.GetAutoCastSlotCount() - usedSlots
}

// AddSpellToAutoCast adds a spell to auto-cast slots with default condition.
func (gs *GameState) AddSpellToAutoCast(spellID string) bool {
	return gs.AddSpellToAutoCastWithCondition(spellID, ConditionAlways)
}

// AddSpellToAutoCastWithCondition adds a spell to auto-cast slots with a specific condition.
func (gs *GameState) AddSpellToAutoCastWithCondition(spellID string, condition AutoCastCondition) bool {
	if gs.IsSpellInAutoCast(spellID) {
		return false // Already in slot
	}
	if gs.GetAvailableAutoCastSlots() <= 0 {
		return false // No slots available
	}
	gs.Session.AutoCastConfigs = append(gs.Session.AutoCastConfigs, AutoCastSlotConfig{
		SpellID:   spellID,
		Condition: condition,
	})
	// Also update legacy slots for backward compat
	gs.Session.AutoCastSlots = append(gs.Session.AutoCastSlots, spellID)
	return true
}

// RemoveSpellFromAutoCast removes a spell from auto-cast slots.
func (gs *GameState) RemoveSpellFromAutoCast(spellID string) bool {
	// Remove from new config system
	for i, cfg := range gs.Session.AutoCastConfigs {
		if cfg.SpellID == spellID {
			gs.Session.AutoCastConfigs = append(gs.Session.AutoCastConfigs[:i], gs.Session.AutoCastConfigs[i+1:]...)
			break
		}
	}
	// Remove from legacy slots
	for i, id := range gs.Session.AutoCastSlots {
		if id == spellID {
			gs.Session.AutoCastSlots = append(gs.Session.AutoCastSlots[:i], gs.Session.AutoCastSlots[i+1:]...)
			return true
		}
	}
	return false
}

// GetAutoCastCondition returns the condition for a spell's auto-cast slot.
func (gs *GameState) GetAutoCastCondition(spellID string) AutoCastCondition {
	for _, cfg := range gs.Session.AutoCastConfigs {
		if cfg.SpellID == spellID {
			return cfg.Condition
		}
	}
	return ConditionAlways // Default for legacy slots
}

// SetAutoCastCondition updates the condition for a spell's auto-cast slot.
func (gs *GameState) SetAutoCastCondition(spellID string, condition AutoCastCondition) bool {
	for i, cfg := range gs.Session.AutoCastConfigs {
		if cfg.SpellID == spellID {
			gs.Session.AutoCastConfigs[i].Condition = condition
			return true
		}
	}
	return false
}

// CycleAutoCastCondition cycles through available conditions for a slot.
func (gs *GameState) CycleAutoCastCondition(spellID string) AutoCastCondition {
	conditions := []AutoCastCondition{
		ConditionAlways,
		ConditionManaAbove50,
		ConditionManaAbove75,
		ConditionSigilNotFull,
		ConditionSynergyActive,
	}
	current := gs.GetAutoCastCondition(spellID)
	for i, cond := range conditions {
		if cond == current {
			nextIdx := (i + 1) % len(conditions)
			gs.SetAutoCastCondition(spellID, conditions[nextIdx])
			return conditions[nextIdx]
		}
	}
	return ConditionAlways
}

// AutoCastToggleResult represents the result of toggling auto-cast.
type AutoCastToggleResult int

const (
	AutoCastRemoved   AutoCastToggleResult = iota // Spell was removed from slot
	AutoCastAdded                                 // Spell was added to slot
	AutoCastSlotsFull                             // Failed: no slots available
)

// ToggleSpellAutoCast adds or removes a spell from auto-cast.
// Returns the result of the toggle operation.
func (gs *GameState) ToggleSpellAutoCast(spellID string) AutoCastToggleResult {
	if gs.IsSpellInAutoCast(spellID) {
		gs.RemoveSpellFromAutoCast(spellID)
		return AutoCastRemoved
	}
	if gs.AddSpellToAutoCast(spellID) {
		return AutoCastAdded
	}
	return AutoCastSlotsFull
}

// MoveAutoCastSlot moves a spell in the auto-cast slots (for priority ordering).
// direction: -1 = move up (higher priority), +1 = move down (lower priority)
func (gs *GameState) MoveAutoCastSlot(spellID string, direction int) bool {
	// Move in new config system
	configs := gs.Session.AutoCastConfigs
	for i, cfg := range configs {
		if cfg.SpellID == spellID {
			newIndex := i + direction
			if newIndex < 0 || newIndex >= len(configs) {
				return false // Can't move further
			}
			// Swap
			configs[i], configs[newIndex] = configs[newIndex], configs[i]
			break
		}
	}

	// Also move in legacy slots for backward compat
	slots := gs.Session.AutoCastSlots
	for i, id := range slots {
		if id == spellID {
			newIndex := i + direction
			if newIndex < 0 || newIndex >= len(slots) {
				return false // Can't move further
			}
			// Swap
			slots[i], slots[newIndex] = slots[newIndex], slots[i]
			return true
		}
	}
	return false
}

// RecordSpellCast records a spell cast for element synergy tracking.
// Resets the streak if a different element is cast (consecutive same-element tracking).
func (gs *GameState) RecordSpellCast(element Element) {
	// If casting a different element, reset the streak
	if len(gs.Session.LastCastElements) > 0 {
		lastElement := gs.Session.LastCastElements[len(gs.Session.LastCastElements)-1]
		if lastElement != element {
			gs.Session.LastCastElements = []Element{}
		}
	}

	gs.Session.LastCastElements = append(gs.Session.LastCastElements, element)
	// Keep only last 3 (enough to trigger synergy)
	if len(gs.Session.LastCastElements) > 3 {
		gs.Session.LastCastElements = gs.Session.LastCastElements[1:]
	}
}

// CheckElementSynergy checks if a synergy should activate.
// Returns the element if synergy triggered, empty string otherwise.
func (gs *GameState) CheckElementSynergy() Element {
	elements := gs.Session.LastCastElements
	if len(elements) < 3 {
		return ""
	}
	// Check if last 3 are same element
	last := elements[len(elements)-1]
	for _, e := range elements[len(elements)-3:] {
		if e != last {
			return ""
		}
	}
	return last
}

// ActivateSynergy activates an element synergy buff.
func (gs *GameState) ActivateSynergy(element Element, durationMs int64) {
	gs.Session.ActiveSynergy = element
	gs.Session.SynergyExpiresAtMs = time.Now().UnixMilli() + durationMs
	// Clear streak so it must be rebuilt
	gs.Session.LastCastElements = []Element{}
}

// HasActiveSynergy returns true if a synergy buff is currently active.
func (gs *GameState) HasActiveSynergy() bool {
	if gs.Session.ActiveSynergy == "" {
		return false
	}
	return time.Now().UnixMilli() < gs.Session.SynergyExpiresAtMs
}

// GetActiveSynergy returns the active synergy element, or empty if none.
func (gs *GameState) GetActiveSynergy() Element {
	if gs.HasActiveSynergy() {
		return gs.Session.ActiveSynergy
	}
	return ""
}

// GetSynergyTimeRemaining returns milliseconds remaining on synergy buff.
func (gs *GameState) GetSynergyTimeRemaining() int64 {
	if !gs.HasActiveSynergy() {
		return 0
	}
	remaining := gs.Session.SynergyExpiresAtMs - time.Now().UnixMilli()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ResetForPrestige resets appropriate data for prestige.
func (gs *GameState) ResetForPrestige(baseSpells []*Spell) {
	// Process prestige bonuses first
	gs.PrestigeData.ProcessPrestige()

	// Reset tower
	gs.Tower.Reset()

	// Reset to base spells only
	gs.Spells = baseSpells
	gs.UnlockedSpellIDs = make([]string, len(baseSpells))
	for i, s := range baseSpells {
		gs.UnlockedSpellIDs[i] = s.ID
	}

	// Clear rituals
	gs.Rituals = []*Ritual{}
	gs.ActiveRitualCount = 0

	// Update passive bonuses from prestige
	gs.PassiveBonuses.ManaGenMultiplier = gs.PrestigeData.PermanentManaGenMultiplier
	gs.PassiveBonuses.SpellCooldownReduction = gs.PrestigeData.SpellCooldownReduction
	gs.PassiveBonuses.RitualCapacity = gs.PrestigeData.RitualCapacity
}
