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

// ErrPlayerNotFound is returned when a player is not found.
var ErrPlayerNotFound = errors.New("player not found")

// ErrPlayerExists is returned when trying to create a duplicate player.
var ErrPlayerExists = errors.New("player already exists")

// PlayerRepository handles player database operations.
// Implements the PlayerStore interface.
type PlayerRepository struct {
	collection *mongo.Collection
}

// NewPlayerRepository creates a new PlayerRepository.
func NewPlayerRepository(db *Database) *PlayerRepository {
	return &PlayerRepository{
		collection: db.Players,
	}
}

// Create inserts a new player into the database.
func (r *PlayerRepository) Create(ctx context.Context, player *models.Player) error {
	_, err := r.collection.InsertOne(ctx, player)
	if mongo.IsDuplicateKeyError(err) {
		return ErrPlayerExists
	}
	return err
}

// GetByUUID retrieves a player by their UUID.
func (r *PlayerRepository) GetByUUID(ctx context.Context, uuid string) (*models.Player, error) {
	var player models.Player
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&player)
	if err == mongo.ErrNoDocuments {
		return nil, ErrPlayerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// GetByUsername retrieves a player by their username.
func (r *PlayerRepository) GetByUsername(ctx context.Context, username string) (*models.Player, error) {
	var player models.Player
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&player)
	if err == mongo.ErrNoDocuments {
		return nil, ErrPlayerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// Update updates an existing player.
func (r *PlayerRepository) Update(ctx context.Context, player *models.Player) error {
	filter := bson.M{"uuid": player.UUID}
	update := bson.M{"$set": player}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateLastPlayed updates the player's last played timestamp.
func (r *PlayerRepository) UpdateLastPlayed(ctx context.Context, uuid string) error {
	filter := bson.M{"uuid": uuid}
	update := bson.M{
		"$set": bson.M{"last_played": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdatePlayDuration adds to the player's total play duration.
func (r *PlayerRepository) UpdatePlayDuration(ctx context.Context, uuid string, durationMs int64) error {
	filter := bson.M{"uuid": uuid}
	update := bson.M{
		"$inc": bson.M{"play_duration_ms": durationMs},
		"$set": bson.M{"last_played": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// IncrementPrestigeCount increments the player's total prestige count.
func (r *PlayerRepository) IncrementPrestigeCount(ctx context.Context, uuid string) error {
	filter := bson.M{"uuid": uuid}
	update := bson.M{
		"$inc": bson.M{"total_prestige_count": 1},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete removes a player from the database.
func (r *PlayerRepository) Delete(ctx context.Context, uuid string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"uuid": uuid})
	return err
}

// List returns all players (with pagination).
func (r *PlayerRepository) List(ctx context.Context, limit, offset int64) ([]*models.Player, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.D{{Key: "last_played", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var players []*models.Player
	if err = cursor.All(ctx, &players); err != nil {
		return nil, err
	}

	return players, nil
}

// Exists checks if a player exists by UUID.
func (r *PlayerRepository) Exists(ctx context.Context, uuid string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"uuid": uuid})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByUsername checks if a username is taken.
func (r *PlayerRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
