package game

// GameEvent represents different events that can occur in the game.
type GameEvent int

const (
	EventNone GameEvent = iota
	EventFloorClimbed
	EventSpellUnlocked
	EventSpellCast
	EventRitualCreated
	EventRitualActivated
	EventRitualExpired
	EventPrestige
	EventGameSaved
	EventGameLoaded
	EventOfflineProgress
	EventLevelUp
	EventAchievement
)

// EventData contains information about a game event.
type EventData struct {
	Type    GameEvent
	Message string
	Data    map[string]interface{}
}

// NewEvent creates a new event.
func NewEvent(eventType GameEvent, message string) EventData {
	return EventData{
		Type:    eventType,
		Message: message,
		Data:    make(map[string]interface{}),
	}
}

// WithData adds data to an event (builder pattern).
func (e EventData) WithData(key string, value interface{}) EventData {
	e.Data[key] = value
	return e
}

// Event constructors for common events
func EventFloorUp(newFloor int) EventData {
	return NewEvent(EventFloorClimbed, "Climbed to floor").
		WithData("floor", newFloor)
}

func EventNewSpell(spellID string, spellName string) EventData {
	return NewEvent(EventSpellUnlocked, "Unlocked new spell: "+spellName).
		WithData("spell_id", spellID).
		WithData("spell_name", spellName)
}

func EventSpellCasted(spellID string, manaCost float64) EventData {
	return NewEvent(EventSpellCast, "Spell cast").
		WithData("spell_id", spellID).
		WithData("mana_cost", manaCost)
}

func EventNewRitual(ritualID string, spellIDs []string) EventData {
	return NewEvent(EventRitualCreated, "Created new ritual").
		WithData("ritual_id", ritualID).
		WithData("spell_ids", spellIDs)
}

func EventPrestiged(newEra int, multiplier float64) EventData {
	return NewEvent(EventPrestige, "Ascended to a new era!").
		WithData("era", newEra).
		WithData("multiplier", multiplier)
}

func EventOffline(manaEarned float64, floorsClimbed int, timeOffline float64) EventData {
	return NewEvent(EventOfflineProgress, "Offline progress calculated").
		WithData("mana_earned", manaEarned).
		WithData("floors_climbed", floorsClimbed).
		WithData("time_offline_seconds", timeOffline)
}
