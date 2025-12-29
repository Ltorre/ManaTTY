package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Ltorre/ManaTTY/models"
)

// JSONSaveStore implements SaveStore using local JSON files.
// Saves are stored in ~/.manatty/saves/<player_uuid>/<slot>.json
type JSONSaveStore struct {
	baseDir string
	mu      sync.RWMutex
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
func (s *JSONSaveStore) playerDir(playerUUID string) string {
	return filepath.Join(s.baseDir, playerUUID)
}

// savePath returns the file path for a specific save slot.
func (s *JSONSaveStore) savePath(playerUUID string, slot int) string {
	return filepath.Join(s.playerDir(playerUUID), slotFilename(slot))
}

func slotFilename(slot int) string {
	return "slot_" + itoa(slot) + ".json"
}

// Simple int to string without importing strconv
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	negative := i < 0
	if negative {
		i = -i
	}
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}

// Save upserts a game save to a JSON file.
func (s *JSONSaveStore) Save(ctx context.Context, save *models.GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure player directory exists
	playerDir := s.playerDir(save.PlayerUUID)
	if err := os.MkdirAll(playerDir, 0755); err != nil {
		return err
	}

	save.SavedAt = time.Now()
	save.Version++

	data, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return err
	}

	savePath := s.savePath(save.PlayerUUID, save.Slot)
	return os.WriteFile(savePath, data, 0644)
}

// Load retrieves a game save from a JSON file.
func (s *JSONSaveStore) Load(ctx context.Context, playerUUID string, slot int) (*models.GameState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	savePath := s.savePath(playerUUID, slot)
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	playerDir := s.playerDir(playerUUID)
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
	s.mu.Lock()
	defer s.mu.Unlock()

	savePath := s.savePath(playerUUID, slot)
	err := os.Remove(savePath)
	if os.IsNotExist(err) {
		return nil // Already deleted
	}
	return err
}

// DeleteAllForPlayer removes all saves for a player.
func (s *JSONSaveStore) DeleteAllForPlayer(ctx context.Context, playerUUID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	playerDir := s.playerDir(playerUUID)
	err := os.RemoveAll(playerDir)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Exists checks if a save exists for a player and slot.
func (s *JSONSaveStore) Exists(ctx context.Context, playerUUID string, slot int) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	savePath := s.savePath(playerUUID, slot)
	_, err := os.Stat(savePath)
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

func (s *JSONPlayerStore) playerPath(uuid string) string {
	return filepath.Join(s.baseDir, uuid+".json")
}

// Create creates a new player.
func (s *JSONPlayerStore) Create(ctx context.Context, player *models.Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.playerPath(player.UUID), data, 0644)
}

// GetByUUID retrieves a player by UUID.
func (s *JSONPlayerStore) GetByUUID(ctx context.Context, uuid string) (*models.Player, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.playerPath(uuid))
	if os.IsNotExist(err) {
		return nil, errors.New("player not found")
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

	return nil, errors.New("player not found")
}

// Update updates an existing player.
func (s *JSONPlayerStore) Update(ctx context.Context, player *models.Player) error {
	return s.Create(ctx, player) // Same as create (overwrite)
}

// Delete removes a player.
func (s *JSONPlayerStore) Delete(ctx context.Context, uuid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(s.playerPath(uuid))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
