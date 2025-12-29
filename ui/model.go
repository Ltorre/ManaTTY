package ui

import (
	"time"

	"github.com/Ltorre/ManaTTY/engine"
	"github.com/Ltorre/ManaTTY/models"
	"github.com/Ltorre/ManaTTY/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// ViewType represents the current screen being displayed.
type ViewType string

const (
	ViewTower    ViewType = "tower"
	ViewSpells   ViewType = "spells"
	ViewRituals  ViewType = "rituals"
	ViewStats    ViewType = "stats"
	ViewPrestige ViewType = "prestige"
	ViewMenu     ViewType = "menu"
)

// Model is the main Bubble Tea model for the game.
type Model struct {
	// Game state
	gameState *models.GameState
	player    *models.Player
	engine    *engine.GameEngine

	// Storage (optional, can be nil for offline play)
	db         *storage.Database
	playerRepo *storage.PlayerRepository
	saveRepo   *storage.SaveRepository

	// UI state
	currentView  ViewType
	previousView ViewType
	width        int
	height       int
	ready        bool

	// Selection state
	selectedIndex int
	maxIndex      int

	// Confirmation dialogs
	confirming  bool
	confirmText string

	// Notifications
	notification     string
	notificationTime time.Time

	// Timing
	lastUpdate       time.Time
	tickInterval     time.Duration
	lastSkipCheckAt  time.Time // For aggregated mana hints

	// Error handling
	lastError error

	// Ritual builder state
	ritualSpells []string
}

// NewModel creates a new UI model.
func NewModel() *Model {
	now := time.Now()
	return &Model{
		currentView:    ViewTower,
		previousView:   ViewTower,
		lastUpdate:     now,
		lastSkipCheckAt: now,
		tickInterval:   100 * time.Millisecond, // 10 FPS
		ritualSpells:   make([]string, 0, 3),
	}
}

// SetGameState sets the game state for the model.
func (m *Model) SetGameState(gs *models.GameState) {
	m.gameState = gs
}

// SetPlayer sets the player for the model.
func (m *Model) SetPlayer(p *models.Player) {
	m.player = p
}

// SetEngine sets the game engine.
func (m *Model) SetEngine(e *engine.GameEngine) {
	m.engine = e
}

// SetDatabase sets the database connections.
func (m *Model) SetDatabase(db *storage.Database) {
	m.db = db
	if db != nil {
		m.playerRepo = storage.NewPlayerRepository(db)
		m.saveRepo = storage.NewSaveRepository(db)
	}
}

// Init initializes the model and returns initial commands.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.tickCmd(),
		tea.EnterAltScreen,
	)
}

// tickCmd returns a command that sends a tick after the interval.
func (m Model) tickCmd() tea.Cmd {
	return tea.Tick(m.tickInterval, func(t time.Time) tea.Msg {
		return TickMsg{Timestamp: t}
	})
}

// ShowNotification displays a notification message.
func (m *Model) ShowNotification(msg string) {
	m.notification = msg
	m.notificationTime = time.Now()
}

// ClearNotification clears the current notification.
func (m *Model) ClearNotification() {
	m.notification = ""
}

// Navigate changes the current view.
func (m *Model) Navigate(view ViewType) {
	m.previousView = m.currentView
	m.currentView = view
	m.selectedIndex = 0
}

// GoBack returns to the previous view.
func (m *Model) GoBack() {
	m.currentView = m.previousView
	m.selectedIndex = 0
}

// StartConfirm starts a confirmation dialog.
func (m *Model) StartConfirm(text string) {
	m.confirming = true
	m.confirmText = text
}

// CancelConfirm cancels the confirmation dialog.
func (m *Model) CancelConfirm() {
	m.confirming = false
	m.confirmText = ""
}

// Message types for Bubble Tea

// TickMsg is sent on each game tick.
type TickMsg struct {
	Timestamp time.Time
}

// NavigateMsg requests navigation to a view.
type NavigateMsg struct {
	View ViewType
}

// SaveGameMsg requests saving the game.
type SaveGameMsg struct{}

// SaveCompleteMsg indicates save completed.
type SaveCompleteMsg struct {
	Error error
}

// LoadGameMsg requests loading a game.
type LoadGameMsg struct {
	PlayerUUID string
	Slot       int
}

// LoadCompleteMsg indicates load completed.
type LoadCompleteMsg struct {
	GameState *models.GameState
	Error     error
}

// CastSpellMsg requests casting a spell.
type CastSpellMsg struct {
	SpellIndex int
	Manual     bool
}

// CreateRitualMsg requests creating a ritual.
type CreateRitualMsg struct {
	SpellIDs []string
}

// PrestigeMsg requests prestige action.
type PrestigeMsg struct{}

// NotificationMsg displays a notification.
type NotificationMsg struct {
	Text     string
	Duration time.Duration
}

// ErrorMsg indicates an error occurred.
type ErrorMsg struct {
	Error error
}
