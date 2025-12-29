package storage

import (
	"context"
	"errors"
	"time"

	"github.com/Ltorre/ManaTTY/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrSaveNotFound is returned when a game save is not found.
var ErrSaveNotFound = errors.New("game save not found")

// SaveRepository handles game save database operations.
type SaveRepository struct {
	collection *mongo.Collection
}

// NewSaveRepository creates a new SaveRepository.
func NewSaveRepository(db *Database) *SaveRepository {
	return &SaveRepository{
		collection: db.Saves,
	}
}

// Save upserts a game save (insert or update).
func (r *SaveRepository) Save(ctx context.Context, save *models.GameState) error {
	filter := bson.M{
		"player_uuid": save.PlayerUUID,
		"slot":        save.Slot,
	}

	save.SavedAt = time.Now()
	save.Version++

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": save}, opts)
	return err
}

// Load retrieves a game save by player UUID and slot.
func (r *SaveRepository) Load(ctx context.Context, playerUUID string, slot int) (*models.GameState, error) {
	var save models.GameState
	filter := bson.M{
		"player_uuid": playerUUID,
		"slot":        slot,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&save)
	if err == mongo.ErrNoDocuments {
		return nil, ErrSaveNotFound
	}
	if err != nil {
		return nil, err
	}

	return &save, nil
}

// LoadLatest loads the most recently saved game for a player.
func (r *SaveRepository) LoadLatest(ctx context.Context, playerUUID string) (*models.GameState, error) {
	var save models.GameState

	opts := options.FindOne().SetSort(bson.D{{Key: "saved_at", Value: -1}})
	filter := bson.M{"player_uuid": playerUUID}

	err := r.collection.FindOne(ctx, filter, opts).Decode(&save)
	if err == mongo.ErrNoDocuments {
		return nil, ErrSaveNotFound
	}
	if err != nil {
		return nil, err
	}

	return &save, nil
}

// ListSaves returns all saves for a player.
func (r *SaveRepository) ListSaves(ctx context.Context, playerUUID string) ([]*models.GameState, error) {
	opts := options.Find().SetSort(bson.D{{Key: "slot", Value: 1}})
	filter := bson.M{"player_uuid": playerUUID}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var saves []*models.GameState
	if err = cursor.All(ctx, &saves); err != nil {
		return nil, err
	}

	return saves, nil
}

// Delete removes a specific game save.
func (r *SaveRepository) Delete(ctx context.Context, playerUUID string, slot int) error {
	filter := bson.M{
		"player_uuid": playerUUID,
		"slot":        slot,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// DeleteAllForPlayer removes all saves for a player.
func (r *SaveRepository) DeleteAllForPlayer(ctx context.Context, playerUUID string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"player_uuid": playerUUID})
	return err
}

// Exists checks if a save exists for a player and slot.
func (r *SaveRepository) Exists(ctx context.Context, playerUUID string, slot int) (bool, error) {
	filter := bson.M{
		"player_uuid": playerUUID,
		"slot":        slot,
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountSaves returns the number of saves for a player.
func (r *SaveRepository) CountSaves(ctx context.Context, playerUUID string) (int, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"player_uuid": playerUUID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetLastSavedTime returns the last save time for a player.
func (r *SaveRepository) GetLastSavedTime(ctx context.Context, playerUUID string, slot int) (time.Time, error) {
	save, err := r.Load(ctx, playerUUID, slot)
	if err != nil {
		return time.Time{}, err
	}
	return save.SavedAt, nil
}
