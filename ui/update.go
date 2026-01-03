package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
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
	case ViewSpecialize:
		return m.handleSpecializeKeys(msg)
	case ViewFloorEvent:
		return m.handleFloorEventKeys(msg)
	case ViewRotation:
		return m.handleRotationKeys(msg)
	}

	return m, nil
}

func (m Model) handleFloorEventKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If the event has already expired, exit the view.
	if m.gameState == nil || m.gameState.Session == nil || m.gameState.Session.ActiveFloorEvent == nil {
		m.GoBack()
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case "down", "j":
		if m.selectedIndex < 2 {
			m.selectedIndex++
		}
	case "1":
		m.selectedIndex = 0
		return m.handleFloorEventChoice()
	case "2":
		m.selectedIndex = 1
		return m.handleFloorEventChoice()
	case "3":
		m.selectedIndex = 2
		return m.handleFloorEventChoice()
	case "enter", " ":
		return m.handleFloorEventChoice()
	case "esc", "b":
		// Explicitly ignore the event (no bonus)
		m.gameState.ClearFloorEvent()
		m.GoBack()
		m.ShowNotification("Floor event ignored")
	}

	return m, nil
}

func (m Model) handleFloorEventChoice() (tea.Model, tea.Cmd) {
	if m.gameState == nil || m.gameState.Session == nil || m.gameState.Session.ActiveFloorEvent == nil {
		m.GoBack()
		return m, nil
	}

	var choice models.FloorEventChoice
	switch m.selectedIndex {
	case 0:
		choice = models.FloorEventChoiceManaGen
	case 1:
		choice = models.FloorEventChoiceSigilChargeRate
	case 2:
		choice = models.FloorEventChoiceCooldownReduction
	default:
		choice = models.FloorEventChoiceManaGen
	}

	currentFloor := 0
	if m.gameState.Tower != nil {
		currentFloor = m.gameState.Tower.CurrentFloor
	}

	m.gameState.ApplyFloorEventChoice(choice, currentFloor, game.FloorEventBuffDurationFloors)
	m.GoBack()
	m.ShowNotification("Floor event chosen: " + models.FloorEventChoiceDisplayNames[choice])
	return m, nil
}

// handleConfirmKey handles keys during confirmation.
func (m Model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		action := m.confirmAction
		m.CancelConfirm()
		// Handle confirmed action based on action key
		switch action {
		case "prestige":
			return m.handlePrestigeConfirmed()
		case "reset_rituals":
			if m.engine != nil && m.gameState != nil {
				m.engine.ResetRituals(m.gameState)
				m.ritualSpells = []string{}
				m.ShowNotification("Rituals reset")
			}
			return m, nil
		default:
			return m, nil
		}

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
	case "o":
		// v1.5.0: Navigate to rotation view
		m.Navigate(ViewRotation)
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
	case "enter":
		// Cast selected spell manually
		if m.gameState != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			if err := m.engine.CastSpell(m.gameState, spell, true); err != nil {
				m.ShowNotification(err.Error())
			} else {
				// Check for synergy activation and combine notification
				if m.gameState.HasActiveSynergy() {
					m.ShowNotification(fmt.Sprintf("%s cast! %s SYNERGY!", spell.Name, string(m.gameState.GetActiveSynergy())))
				} else {
					m.ShowNotification(fmt.Sprintf("%s cast!", spell.Name))
				}
			}
		}
	case " ":
		// Toggle auto-cast slot for selected spell
		if m.gameState != nil && m.engine != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			inSlot, err := m.engine.ToggleSpellAutoCast(m.gameState, spell.ID)
			if err != nil {
				m.ShowNotification(err.Error())
			} else if inSlot {
				m.ShowNotification(spell.Name + " added to slot")
			} else {
				m.ShowNotification(spell.Name + " removed from slot")
			}
		}
	case "u", "U":
		// Upgrade selected spell
		if m.gameState != nil && m.engine != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			if err := m.engine.UpgradeSpell(m.gameState, spell); err != nil {
				m.ShowNotification(err.Error())
			} else {
				m.ShowNotification(fmt.Sprintf("%s upgraded to Lv%d!", spell.Name, spell.Level))
			}
		}
	case "<", ",":
		// Move spell up in auto-cast priority
		if m.gameState != nil && m.engine != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			if m.engine.MoveAutoCastSlotUp(m.gameState, spell.ID) {
				m.ShowNotification(fmt.Sprintf("%s priority increased", spell.Name))
			} else if !m.gameState.IsSpellInAutoCast(spell.ID) {
				m.ShowNotification(fmt.Sprintf("%s not in auto-cast slot", spell.Name))
			} else {
				m.ShowNotification(fmt.Sprintf("%s already at top priority", spell.Name))
			}
		}
	case ">", ".":
		// Move spell down in auto-cast priority
		if m.gameState != nil && m.engine != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			if m.engine.MoveAutoCastSlotDown(m.gameState, spell.ID) {
				m.ShowNotification(fmt.Sprintf("%s priority decreased", spell.Name))
			} else if !m.gameState.IsSpellInAutoCast(spell.ID) {
				m.ShowNotification(fmt.Sprintf("%s not in auto-cast slot", spell.Name))
			} else {
				m.ShowNotification(fmt.Sprintf("%s already at lowest priority", spell.Name))
			}
		}
	case "esc", "b":
		m.Navigate(ViewTower)
	case "a":
		// Toggle auto-cast on/off
		if m.engine != nil {
			enabled := m.engine.ToggleAutoCast(m.gameState)
			if enabled {
				m.ShowNotification("Auto-cast enabled")
			} else {
				m.ShowNotification("Auto-cast disabled")
			}
		}
	case "c":
		// Cycle auto-cast condition for selected spell
		if m.gameState != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			if m.gameState.IsSpellInAutoCast(spell.ID) {
				newCond := m.gameState.CycleAutoCastCondition(spell.ID)
				m.ShowNotification(fmt.Sprintf("%s: %s", spell.Name, models.ConditionDisplayNames[newCond]))
			} else {
				m.ShowNotification(spell.Name + " not in auto-cast slot")
			}
		}
	case "x":
		// Open specialization menu if spell needs it
		if m.gameState != nil && m.selectedIndex < len(m.gameState.Spells) {
			spell := m.gameState.Spells[m.selectedIndex]
			tier, needs := spell.NeedsSpecialization()
			if needs {
				m.specSpellID = spell.ID
				m.specTier = tier
				m.specChoiceIdx = 0
				m.Navigate(ViewSpecialize)
			} else if spell.Level >= 5 {
				m.ShowNotification("Already specialized")
			} else {
				m.ShowNotification("Reach level 5 to specialize")
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
	case "x", "X":
		// Reset all rituals (free up ritual slots)
		if m.gameState != nil && len(m.gameState.Rituals) > 0 {
			m.StartConfirmAction("Reset ALL rituals for this save? (y/n)", "reset_rituals")
		} else {
			m.ShowNotification("No rituals to reset")
		}
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
		m.StartConfirmAction("Are you sure you want to prestige? (y/n)", "prestige")
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

// handleSpecializeKeys handles keys in the specialization view.
func (m Model) handleSpecializeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k", "left", "h":
		if m.specChoiceIdx > 0 {
			m.specChoiceIdx--
		}
	case "down", "j", "right", "l":
		if m.specChoiceIdx < 1 {
			m.specChoiceIdx++
		}
	case "enter", " ":
		// Apply specialization
		if m.gameState != nil {
			spell := m.gameState.GetSpellByID(m.specSpellID)
			if spell != nil {
				var spec models.SpellSpecialization
				if m.specTier == 1 {
					if m.specChoiceIdx == 0 {
						spec = models.SpecCritChance
					} else {
						spec = models.SpecManaEfficiency
					}
					spell.Tier1Spec = spec
				} else {
					if m.specChoiceIdx == 0 {
						spec = models.SpecBurstDamage
					} else {
						spec = models.SpecRapidCast
					}
					spell.Tier2Spec = spec
				}
				m.ShowNotification(fmt.Sprintf("%s specialized: %s!", spell.Name, models.SpecializationDisplayNames[spec]))
				m.Navigate(ViewSpells)
			}
		}
	case "esc", "b":
		m.Navigate(ViewSpells)
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

	// Floor events: auto-open when available (but don't interrupt confirmations/specializations)
	if m.gameState != nil && m.gameState.Session != nil {
		if m.gameState.Session.ActiveFloorEvent != nil && m.currentView != ViewFloorEvent && !m.confirming && m.currentView != ViewSpecialize {
			m.Navigate(ViewFloorEvent)
		}
		// If we're showing the event view and it expired/was cleared, return to the previous view.
		if m.currentView == ViewFloorEvent && m.gameState.Session.ActiveFloorEvent == nil {
			m.GoBack()
		}
	}

	// Check for aggregated skip notifications (every 2 seconds)
	if time.Since(m.lastSkipCheckAt) > 2*time.Second {
		skipCount := m.engine.GetAndResetSkipCount(m.gameState)
		if skipCount > 0 {
			m.ShowNotification(fmt.Sprintf("Auto-cast skipped %d times (low mana)", skipCount))
		}
		m.lastSkipCheckAt = time.Now()
	}

	// Clear old notifications (after 3 seconds)
	if m.notification != "" && time.Since(m.notificationTime) > 3*time.Second {
		m.ClearNotification()
	}

	// Auto-save every 30 seconds
	if m.saveStore != nil && time.Since(m.gameState.Session.LastSavedAt) > 30*time.Second {
		return m, tea.Batch(m.tickCmd(), m.saveGameCmd())
	}

	return m, m.tickCmd()
}

// v1.5.0: Rotation View Key Handler
func (m Model) handleRotationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.gameState == nil {
		m.GoBack()
		return m, nil
	}

	m.gameState.EnsureRotation()
	rotation := m.gameState.Session.Rotation

	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case "down", "j":
		if m.selectedIndex < len(rotation.Spells)-1 {
			m.selectedIndex++
		}
	case "o":
		// Toggle rotation on/off
		rotation.Enabled = !rotation.Enabled
		status := "disabled"
		if rotation.Enabled {
			status = "enabled"
		}
		m.ShowNotification(fmt.Sprintf("Rotation system %s", status))
	case "v":
		// Convert auto-cast to rotation
		m.gameState.ConvertAutoCastToRotation()
		m.selectedIndex = 0
		m.ShowNotification("Converted auto-cast slots to rotation system")
	case "w":
		// Toggle cooldown weaving
		rotation.CooldownWeaving = !rotation.CooldownWeaving
		status := "disabled"
		if rotation.CooldownWeaving {
			status = "enabled"
		}
		m.ShowNotification(fmt.Sprintf("Cooldown weaving %s", status))
	case " ":
		// Toggle selected spell enabled/disabled
		if m.selectedIndex < len(rotation.Spells) {
			m.gameState.ToggleRotationSpell(rotation.Spells[m.selectedIndex].SpellID)
			status := "disabled"
			if rotation.Spells[m.selectedIndex].Enabled {
				status = "enabled"
			}
			m.ShowNotification(fmt.Sprintf("Spell %s", status))
		}
	case "esc":
		m.GoBack()
	}

	return m, nil
}

// saveGameCmd returns a command to save the game.
func (m Model) saveGameCmd() tea.Cmd {
	return func() tea.Msg {
		if m.saveStore == nil || m.gameState == nil {
			return SaveCompleteMsg{Error: nil}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		m.gameState.Session.LastSavedAt = time.Now()
		err := m.saveStore.Save(ctx, m.gameState)
		return SaveCompleteMsg{Error: err}
	}
}
