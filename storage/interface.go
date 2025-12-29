package storage

import (
	"context"
	"errors"
	"time"

	"github.com/Ltorre/ManaTTY/models"
)

// ErrSaveNotFound is returned when a save file doesn't exist.
var ErrSaveNotFound = errors.New("save not found")

// SaveStore defines the interface for game save storage.
// Both MongoDB and JSON file storage implement this interface.
type SaveStore interface {
	// Save upserts a game save (insert or update).
	Save(ctx context.Context, save *models.GameState) error

	// Load retrieves a game save by player UUID and slot.
	Load(ctx context.Context, playerUUID string, slot int) (*models.GameState, error)

	// LoadLatest loads the most recently saved game for a player.
	LoadLatest(ctx context.Context, playerUUID string) (*models.GameState, error)

	// ListSaves returns all saves for a player.
	ListSaves(ctx context.Context, playerUUID string) ([]*models.GameState, error)

	// Delete removes a specific game save.
	Delete(ctx context.Context, playerUUID string, slot int) error

	// DeleteAllForPlayer removes all saves for a player.
	DeleteAllForPlayer(ctx context.Context, playerUUID string) error

	// Exists checks if a save exists for a player and slot.
	Exists(ctx context.Context, playerUUID string, slot int) (bool, error)

	// CountSaves returns the number of saves for a player.
	CountSaves(ctx context.Context, playerUUID string) (int, error)

	// GetLastSavedTime returns the last save time for a player.
	GetLastSavedTime(ctx context.Context, playerUUID string, slot int) (time.Time, error)
}

// PlayerStore defines the interface for player storage.
type PlayerStore interface {
	// Create creates a new player.
	Create(ctx context.Context, player *models.Player) error

	// GetByUUID retrieves a player by UUID.
	GetByUUID(ctx context.Context, uuid string) (*models.Player, error)

	// GetByUsername retrieves a player by username.
	GetByUsername(ctx context.Context, username string) (*models.Player, error)

	// Update updates an existing player.
	Update(ctx context.Context, player *models.Player) error

	// Delete removes a player.
	Delete(ctx context.Context, uuid string) error
}
