package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Ltorre/ManaTTY/models"
	"github.com/google/uuid"
)

// JSONSaveStore implements SaveStore using local JSON files.
// Saves are stored in ~/.manatty/saves/<player_uuid>/<slot>.json
// Note: Context cancellation is not supported for file-based storage operations.
type JSONSaveStore struct {
	baseDir string
	mu      sync.RWMutex
}

// ErrInvalidUUID is returned when a UUID fails validation.
var ErrInvalidUUID = fmt.Errorf("invalid UUID format")

// validateUUID checks if the given string is a valid UUID to prevent path traversal.
func validateUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidUUID
	}
	return nil
}

// NewJSONSaveStore creates a new JSON-based save store.
func NewJSONSaveStore() (*JSONSaveStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(homeDir, ".manatty", "saves")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	return &JSONSaveStore{baseDir: baseDir}, nil
}

// playerDir returns the directory for a player's saves.
// Returns an error if the UUID is invalid (prevents path traversal).
func (s *JSONSaveStore) playerDir(playerUUID string) (string, error) {
	if err := validateUUID(playerUUID); err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(s.baseDir, playerUUID)), nil
}

// savePath returns the file path for a specific save slot.
func (s *JSONSaveStore) savePath(playerUUID string, slot int) (string, error) {
	dir, err := s.playerDir(playerUUID)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, slotFilename(slot)), nil
}

func slotFilename(slot int) string {
	return "slot_" + strconv.Itoa(slot) + ".json"
}

// Save upserts a game save to a JSON file.
func (s *JSONSaveStore) Save(ctx context.Context, save *models.GameState) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure player directory exists
	playerDir, err := s.playerDir(save.PlayerUUID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(playerDir, 0755); err != nil {
		return err
	}

	save.SavedAt = time.Now()
	save.Version++

	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}

	savePath, err := s.savePath(save.PlayerUUID, save.Slot)
	if err != nil {
		return err
	}
	return os.WriteFile(savePath, data, 0644)
}

// Load retrieves a game save from a JSON file.
func (s *JSONSaveStore) Load(ctx context.Context, playerUUID string, slot int) (*models.GameState, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	savePath, err := s.savePath(playerUUID, slot)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(savePath)
	if os.IsNotExist(err) {
		return nil, ErrSaveNotFound
	}
	if err != nil {
		return nil, err
	}

	var save models.GameState
	if err := json.Unmarshal(data, &save); err != nil {
		return nil, err
	}

	return &save, nil
}

// LoadLatest loads the most recently saved game for a player.
func (s *JSONSaveStore) LoadLatest(ctx context.Context, playerUUID string) (*models.GameState, error) {
	saves, err := s.ListSaves(ctx, playerUUID)
	if err != nil {
		return nil, err
	}
	if len(saves) == 0 {
		return nil, ErrSaveNotFound
	}

	// Sort by saved_at descending
	sort.Slice(saves, func(i, j int) bool {
		return saves[i].SavedAt.After(saves[j].SavedAt)
	})

	return saves[0], nil
}

// ListSaves returns all saves for a player.
func (s *JSONSaveStore) ListSaves(ctx context.Context, playerUUID string) ([]*models.GameState, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	playerDir, err := s.playerDir(playerUUID)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(playerDir)
	if os.IsNotExist(err) {
		return []*models.GameState{}, nil
	}
	if err != nil {
		return nil, err
	}

	var saves []*models.GameState
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(playerDir, entry.Name()))
		if err != nil {
			continue
		}

		var save models.GameState
		if err := json.Unmarshal(data, &save); err != nil {
			continue
		}

		saves = append(saves, &save)
	}

	// Sort by slot
	sort.Slice(saves, func(i, j int) bool {
		return saves[i].Slot < saves[j].Slot
	})

	return saves, nil
}

// Delete removes a specific game save.
func (s *JSONSaveStore) Delete(ctx context.Context, playerUUID string, slot int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	savePath, err := s.savePath(playerUUID, slot)
	if err != nil {
		return err
	}
	err = os.Remove(savePath)
	if os.IsNotExist(err) {
		return nil // Already deleted
	}
	return err
}

// DeleteAllForPlayer removes all saves for a player.
func (s *JSONSaveStore) DeleteAllForPlayer(ctx context.Context, playerUUID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	playerDir, err := s.playerDir(playerUUID)
	if err != nil {
		return err
	}
	err = os.RemoveAll(playerDir)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Exists checks if a save exists for a player and slot.
func (s *JSONSaveStore) Exists(ctx context.Context, playerUUID string, slot int) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	savePath, err := s.savePath(playerUUID, slot)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(savePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CountSaves returns the number of saves for a player.
func (s *JSONSaveStore) CountSaves(ctx context.Context, playerUUID string) (int, error) {
	saves, err := s.ListSaves(ctx, playerUUID)
	if err != nil {
		return 0, err
	}
	return len(saves), nil
}

// GetLastSavedTime returns the last save time for a player.
func (s *JSONSaveStore) GetLastSavedTime(ctx context.Context, playerUUID string, slot int) (time.Time, error) {
	save, err := s.Load(ctx, playerUUID, slot)
	if err != nil {
		return time.Time{}, err
	}
	return save.SavedAt, nil
}

// JSONPlayerStore implements PlayerStore using local JSON files.
// Note: Context cancellation is not supported for file-based storage operations.
type JSONPlayerStore struct {
	baseDir string
	mu      sync.RWMutex
}

// NewJSONPlayerStore creates a new JSON-based player store.
func NewJSONPlayerStore() (*JSONPlayerStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(homeDir, ".manatty", "players")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	return &JSONPlayerStore{baseDir: baseDir}, nil
}

func (s *JSONPlayerStore) playerPath(id string) (string, error) {
	if err := validateUUID(id); err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(s.baseDir, id+".json")), nil
}

// Create creates a new player.
func (s *JSONPlayerStore) Create(ctx context.Context, player *models.Player) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return err
	}

	path, err := s.playerPath(player.UUID)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GetByUUID retrieves a player by UUID.
func (s *JSONPlayerStore) GetByUUID(ctx context.Context, id string) (*models.Player, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	path, err := s.playerPath(id)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, ErrPlayerNotFound
	}
	if err != nil {
		return nil, err
	}

	var player models.Player
	if err := json.Unmarshal(data, &player); err != nil {
		return nil, err
	}

	return &player, nil
}

// GetByUsername retrieves a player by username.
func (s *JSONPlayerStore) GetByUsername(ctx context.Context, username string) (*models.Player, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.baseDir, entry.Name()))
		if err != nil {
			continue
		}

		var player models.Player
		if err := json.Unmarshal(data, &player); err != nil {
			continue
		}

		if player.Username == username {
			return &player, nil
		}
	}

	return nil, ErrPlayerNotFound
}

// Update updates an existing player.
func (s *JSONPlayerStore) Update(ctx context.Context, player *models.Player) error {
	return s.Create(ctx, player) // Same as create (overwrite)
}

// Delete removes a player.
func (s *JSONPlayerStore) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	path, err := s.playerPath(id)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
