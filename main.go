package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
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

	// Prompt for nickname
	nickname := promptNickname()

	fmt.Println("Loading...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		utils.Error("Failed to load config: %v", err)
	}

	// Set log level
	utils.SetLogLevel(utils.ParseLogLevel(cfg.LogLevel))

	// Initialize storage based on config
	var saveStore storage.SaveStore
	var playerStore storage.PlayerStore
	var db *storage.Database

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if cfg.StorageMode == "mongodb" && cfg.MongoDBURI != "" {
		// Use MongoDB storage
		db = storage.NewDatabase()
		if err := db.Connect(ctx, cfg.MongoDBURI); err != nil {
			utils.Warn("Database connection failed: %v", err)
			utils.Info("Falling back to local storage")
			cfg.StorageMode = "local"
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

			saveStore = storage.NewSaveRepository(db)
			playerStore = storage.NewPlayerRepository(db)
		}
	}

	// Fall back to local storage if MongoDB not available
	if cfg.StorageMode == "local" || saveStore == nil {
		utils.Info("Using local storage (~/.manatty/)")
		var err error
		saveStore, err = storage.NewJSONSaveStore()
		if err != nil {
			utils.Error("Failed to create local save store: %v", err)
			os.Exit(1)
		}
		playerStore, err = storage.NewJSONPlayerStore()
		if err != nil {
			utils.Error("Failed to create local player store: %v", err)
			os.Exit(1)
		}
	}

	// Create or load game state
	gameState, player := initializeGame(ctx, saveStore, playerStore, nickname)

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
	model.SetSaveStore(saveStore)
	model.SetDatabase(db) // Keep for backward compatibility (may be nil)

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

// promptNickname asks the user for their nickname to load or create a save.
func promptNickname() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter your nickname (or press Enter for 'Wizard'): ")
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		nickname = "Wizard"
	}
	fmt.Printf("Welcome, %s!\n\n", nickname)
	return nickname
}

// initializeGame creates or loads a game state for the given nickname.
func initializeGame(ctx context.Context, saveStore storage.SaveStore, playerStore storage.PlayerStore, nickname string) (*models.GameState, *models.Player) {
	// Try to find existing player by username
	player, err := playerStore.GetByUsername(ctx, nickname)
	if err == nil && player != nil {
		// Load their latest save
		gameState, err := saveStore.LoadLatest(ctx, player.UUID)
		if err == nil {
			utils.Info("Loaded save for %s (Floor %d)", player.Username, gameState.Tower.CurrentFloor)
			return gameState, player
		}
	}

	// Create new game for this nickname
	utils.Info("Creating new game for %s...", nickname)

	playerUUID := uuid.New().String()
	player = models.NewPlayer(playerUUID, nickname)

	gameState := models.NewGameState(playerUUID, 0)

	// Add starting spells
	for _, spell := range game.GetBaseSpells() {
		gameState.AddSpell(spell)
	}

	// Save new player and game
	if err := playerStore.Create(ctx, player); err != nil {
		utils.Warn("Failed to create player: %v", err)
	}

	if err := saveStore.Save(ctx, gameState); err != nil {
		utils.Warn("Failed to save game: %v", err)
	}

	return gameState, player
}
