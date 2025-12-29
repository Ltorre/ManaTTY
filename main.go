package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"github.com/Ltorre/ManaTTY/config"
	"github.com/Ltorre/ManaTTY/engine"
	"github.com/Ltorre/ManaTTY/game"
	"github.com/Ltorre/ManaTTY/models"
	"github.com/Ltorre/ManaTTY/storage"
	"github.com/Ltorre/ManaTTY/ui"
	"github.com/Ltorre/ManaTTY/utils"
)

func main() {
	fmt.Println("🏰 Mage Tower Ascension")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Loading...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		utils.Error("Failed to load config: %v", err)
	}

	// Set log level
	utils.SetLogLevel(utils.ParseLogLevel(cfg.LogLevel))

	// Initialize database (optional)
	var db *storage.Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db = storage.NewDatabase()
	if err := db.Connect(ctx, cfg.MongoDBURI); err != nil {
		utils.Warn("Database connection failed: %v", err)
		utils.Info("Running in offline mode (no saves)")
		db = nil
	} else {
		utils.Info("Connected to MongoDB")

		// Ensure indexes
		if err := db.EnsureIndexes(ctx); err != nil {
			utils.Warn("Failed to create indexes: %v", err)
		}

		// Seed spell definitions
		spellDefs := game.DefaultSpells()
		if err := db.SeedSpellDefinitions(ctx, spellDefs); err != nil {
			utils.Warn("Failed to seed spells: %v", err)
		}
	}

	// Create or load game state
	gameState, player := initializeGame(ctx, db)

	// Apply offline progress if we loaded a save
	gameEngine := engine.NewGameEngine()
	if gameState.SavedAt.After(time.Time{}) {
		offlineProgress := gameEngine.ApplyOfflineProgress(gameState)
		if offlineProgress.TimeOffline > time.Minute {
			fmt.Printf("\n📊 Offline Progress: %s\n", engine.FormatOfflineProgress(offlineProgress))
			if offlineProgress.ManaGenerated > 0 {
				fmt.Printf("   Mana earned: %s\n", utils.FormatNumber(offlineProgress.ManaGenerated))
			}
			if offlineProgress.FloorsClimbed > 0 {
				fmt.Printf("   Floors climbed: %d\n", offlineProgress.FloorsClimbed)
			}
			fmt.Println()
			time.Sleep(2 * time.Second)
		}
	}

	// Create UI model
	model := ui.NewModel()
	model.SetGameState(gameState)
	model.SetPlayer(player)
	model.SetEngine(gameEngine)
	model.SetDatabase(db)

	// Run the TUI
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		utils.Error("Error running program: %v", err)
		os.Exit(1)
	}

	// Cleanup
	if db != nil {
		_ = db.Disconnect(context.Background())
	}

	fmt.Println("\n👋 Thanks for playing Mage Tower Ascension!")
}

// initializeGame creates or loads a game state.
func initializeGame(ctx context.Context, db *storage.Database) (*models.GameState, *models.Player) {
	// Try to load existing save
	if db != nil {
		saveRepo := storage.NewSaveRepository(db)
		playerRepo := storage.NewPlayerRepository(db)

		// Try to find existing player (for simplicity, use first found)
		players, err := playerRepo.List(ctx, 1, 0)
		if err == nil && len(players) > 0 {
			player := players[0]

			// Load their latest save
			gameState, err := saveRepo.LoadLatest(ctx, player.UUID)
			if err == nil {
				utils.Info("Loaded save for %s (Floor %d)", player.Username, gameState.Tower.CurrentFloor)
				return gameState, player
			}
		}
	}

	// Create new game
	utils.Info("Creating new game...")

	playerUUID := uuid.New().String()
	player := models.NewPlayer(playerUUID, "Wizard")

	gameState := models.NewGameState(playerUUID, 0)

	// Add starting spells
	for _, spell := range game.GetBaseSpells() {
		gameState.AddSpell(spell)
	}

	// Save new player and game
	if db != nil {
		playerRepo := storage.NewPlayerRepository(db)
		saveRepo := storage.NewSaveRepository(db)

		if err := playerRepo.Create(ctx, player); err != nil {
			utils.Warn("Failed to create player: %v", err)
		}

		if err := saveRepo.Save(ctx, gameState); err != nil {
			utils.Warn("Failed to save game: %v", err)
		}
	}

	return gameState, player
}
