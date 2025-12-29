package storage

import (
	"context"

	"github.com/Ltorre/ManaTTY/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates all required database indexes.
func (db *Database) EnsureIndexes(ctx context.Context) error {
	// Players collection indexes
	playerIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "last_played", Value: -1}},
		},
	}

	if _, err := db.Players.Indexes().CreateMany(ctx, playerIndexes); err != nil {
		return err
	}

	// Game saves collection indexes
	saveIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "player_uuid", Value: 1},
				{Key: "slot", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "saved_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "player_uuid", Value: 1}},
		},
	}

	if _, err := db.Saves.Indexes().CreateMany(ctx, saveIndexes); err != nil {
		return err
	}

	// Spell definitions collection indexes
	spellIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "required_floor", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "element", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "unlocked_by_default", Value: 1}},
		},
	}

	if _, err := db.SpellDefs.Indexes().CreateMany(ctx, spellIndexes); err != nil {
		return err
	}

	return nil
}

// SeedSpellDefinitions inserts default spell definitions if not present.
func (db *Database) SeedSpellDefinitions(ctx context.Context, spells []*models.SpellDefinition) error {
	for _, spell := range spells {
		filter := bson.M{"_id": spell.ID}
		opts := options.Update().SetUpsert(true)
		update := bson.M{"$setOnInsert": spell}

		_, err := db.SpellDefs.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSpellDefinitions retrieves all spell definitions from the database.
func (db *Database) GetSpellDefinitions(ctx context.Context) ([]*models.SpellDefinition, error) {
	cursor, err := db.SpellDefs.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var spells []*models.SpellDefinition
	if err = cursor.All(ctx, &spells); err != nil {
		return nil, err
	}

	return spells, nil
}

// GetSpellDefinition retrieves a single spell definition by ID.
func (db *Database) GetSpellDefinition(ctx context.Context, id string) (*models.SpellDefinition, error) {
	var spell models.SpellDefinition
	err := db.SpellDefs.FindOne(ctx, bson.M{"_id": id}).Decode(&spell)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &spell, nil
}

// DropAllData removes all data from all collections (for testing).
func (db *Database) DropAllData(ctx context.Context) error {
	if err := db.Players.Drop(ctx); err != nil {
		return err
	}
	if err := db.Saves.Drop(ctx); err != nil {
		return err
	}
	if err := db.SpellDefs.Drop(ctx); err != nil {
		return err
	}
	return nil
}

// GetStats returns database statistics.
func (db *Database) GetStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	playerCount, err := db.Players.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats["players"] = playerCount

	saveCount, err := db.Saves.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats["saves"] = saveCount

	spellCount, err := db.SpellDefs.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	stats["spell_definitions"] = spellCount

	return stats, nil
}
