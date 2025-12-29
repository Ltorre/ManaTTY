package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Game tick
	case TickMsg:
		return m.handleTick(msg)

	// Save complete
	case SaveCompleteMsg:
		if msg.Error != nil {
			m.ShowNotification("Save failed!")
		} else {
			m.ShowNotification("Game saved!")
		}
		return m, nil

	// Load complete
	case LoadCompleteMsg:
		if msg.Error != nil {
			m.ShowNotification("Load failed!")
		} else if msg.GameState != nil {
			m.gameState = msg.GameState
			m.ShowNotification("Game loaded!")
		}
		return m, nil

	// Notification
	case NotificationMsg:
		m.ShowNotification(msg.Text)
		return m, nil

	// Error
	case ErrorMsg:
		m.lastError = msg.Error
		m.ShowNotification("Error: " + msg.Error.Error())
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle confirmation mode
	if m.confirming {
		return m.handleConfirmKey(msg)
	}

	// Global keys
	switch msg.String() {
	case "ctrl+c", "q":
		// Save and quit
		return m, tea.Sequence(m.saveGameCmd(), tea.Quit)

	case "ctrl+s":
		// Manual save
		return m, m.saveGameCmd()
	}

	// View-specific keys
	switch m.currentView {
	case ViewTower:
		return m.handleTowerKeys(msg)
	case ViewSpells:
		return m.handleSpellsKeys(msg)
	case ViewRituals:
		return m.handleRitualsKeys(msg)
	case ViewStats:
		return m.handleStatsKeys(msg)
	case ViewPrestige:
		return m.handlePrestigeKeys(msg)
	case ViewMenu:
		return m.handleMenuKeys(msg)
	}

	return m, nil
}

// handleConfirmKey handles keys during confirmation.
func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		m.CancelConfirm()
		// Handle confirmed action based on view
		if m.currentView == ViewPrestige {
			return m.handlePrestigeConfirmed()
		}
		return m, nil

	case "n", "N", "esc":
		m.CancelConfirm()
		return m, nil
	}
	return m, nil
}

// handleTowerKeys handles keys in the tower view.
func (m Model) handleTowerKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		m.Navigate(ViewSpells)
	case "r":
		m.Navigate(ViewRituals)
	case "t":
		m.Navigate(ViewStats)
	case "p":
		if m.engine != nil && m.engine.CanPrestige(m.gameState) {
			m.Navigate(ViewPrestige)
		} else {
			m.ShowNotification("Reach floor 100 to prestige!")
		}
	case "m":
		m.Navigate(ViewMenu)
	case "a":
		// Toggle auto-cast
		if m.engine != nil {
			enabled := m.engine.ToggleAutoCast(m.gameState)
			if enabled {
				m.ShowNotification("Auto-cast enabled")
			} else {
				m.ShowNotification("Auto-cast disabled")
			}
		}
	}
	return m, nil
}

// handleSpellsKeys handles keys in the spells view.
func (m Model) handleSpellsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case "down", "j":
		if m.gameState != nil && m.selectedIndex < len(m.gameState.Spells)-1 {
			m.selectedIndex++
		}
	case "enter", " ":
		// Cast selected spell
		if m.gameState != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			if err := m.engine.CastSpell(m.gameState, spell, true); err != nil {
				m.ShowNotification(err.Error())
			} else {
				m.ShowNotification(spell.Name + " cast!")
			}
		}
	case "esc", "b":
		m.Navigate(ViewTower)
	case "a":
		// Toggle auto-cast
		if m.engine != nil {
			enabled := m.engine.ToggleAutoCast(m.gameState)
			if enabled {
				m.ShowNotification("Auto-cast enabled")
			} else {
				m.ShowNotification("Auto-cast disabled")
			}
		}
	}
	return m, nil
}

// handleRitualsKeys handles keys in the rituals view.
func (m Model) handleRitualsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case "down", "j":
		m.selectedIndex++
	case "enter", " ":
		// Add spell to ritual builder or select ritual
		if m.gameState != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			if len(m.ritualSpells) < 3 {
				// Check if already selected
				alreadySelected := false
				for _, id := range m.ritualSpells {
					if id == spell.ID {
						alreadySelected = true
						break
					}
				}
				if !alreadySelected {
					m.ritualSpells = append(m.ritualSpells, spell.ID)
					if len(m.ritualSpells) == 3 {
						// Create ritual
						_, err := m.engine.CreateRitual(m.gameState, m.ritualSpells)
						if err != nil {
							m.ShowNotification(err.Error())
						} else {
							m.ShowNotification("Ritual created!")
						}
						m.ritualSpells = []string{}
					}
				}
			}
		}
	case "c":
		// Clear ritual builder
		m.ritualSpells = []string{}
		m.ShowNotification("Selection cleared")
	case "esc", "b":
		m.ritualSpells = []string{}
		m.Navigate(ViewTower)
	}
	return m, nil
}

// handleStatsKeys handles keys in the stats view.
func (m Model) handleStatsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.Navigate(ViewTower)
	}
	return m, nil
}

// handlePrestigeKeys handles keys in the prestige view.
func (m Model) handlePrestigeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.StartConfirm("Are you sure you want to prestige? (y/n)")
	case "esc", "b":
		m.Navigate(ViewTower)
	}
	return m, nil
}

// handlePrestigeConfirmed handles confirmed prestige action.
func (m Model) handlePrestigeConfirmed() (tea.Model, tea.Cmd) {
	if m.engine != nil && m.engine.ProcessPrestige(m.gameState) {
		m.ShowNotification("Ascended to Era " + string(rune('0'+m.gameState.PrestigeData.CurrentEra)) + "!")
		m.Navigate(ViewTower)
	}
	return m, nil
}

// handleMenuKeys handles keys in the menu view.
func (m Model) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b":
		m.Navigate(ViewTower)
	case "s":
		return m, m.saveGameCmd()
	case "q":
		return m, tea.Sequence(m.saveGameCmd(), tea.Quit)
	}
	return m, nil
}

// handleTick processes a game tick.
func (m Model) handleTick(msg TickMsg) (tea.Model, tea.Cmd) {
	if m.gameState == nil || m.engine == nil {
		return m, m.tickCmd()
	}

	// Calculate elapsed time
	elapsed := msg.Timestamp.Sub(m.lastUpdate)
	m.lastUpdate = msg.Timestamp

	// Update game state
	m.engine.Tick(m.gameState, elapsed)

	// Clear old notifications (after 3 seconds)
	if m.notification != "" && time.Since(m.notificationTime) > 3*time.Second {
		m.ClearNotification()
	}

	// Auto-save every 30 seconds
	if m.saveRepo != nil && time.Since(m.gameState.Session.LastSavedAt) > 30*time.Second {
		return m, tea.Batch(m.tickCmd(), m.saveGameCmd())
	}

	return m, m.tickCmd()
}

// saveGameCmd returns a command to save the game.
func (m Model) saveGameCmd() tea.Cmd {
	return func() tea.Msg {
		if m.saveRepo == nil || m.gameState == nil {
			return SaveCompleteMsg{Error: nil}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		m.gameState.Session.LastSavedAt = time.Now()
		err := m.saveRepo.Save(ctx, m.gameState)
		return SaveCompleteMsg{Error: err}
	}
}
