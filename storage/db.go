package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// DatabaseName is the default database name.
	DatabaseName = "mage_tower"

	// Collection names
	CollectionPlayers   = "players"
	CollectionGameSaves = "game_saves"
	CollectionSpellDefs = "spell_definitions"
)

// Database holds the MongoDB client and database reference.
type Database struct {
	Client    *mongo.Client
	DB        *mongo.Database
	Players   *mongo.Collection
	Saves     *mongo.Collection
	SpellDefs *mongo.Collection
}

// NewDatabase creates a new Database instance.
func NewDatabase() *Database {
	return &Database{}
}

// Connect establishes connection to MongoDB.
func (db *Database) Connect(ctx context.Context, uri string) error {
	// Set connection options
	opts := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	// Connect
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	// Store references
	db.Client = client
	db.DB = client.Database(DatabaseName)
	db.Players = db.DB.Collection(CollectionPlayers)
	db.Saves = db.DB.Collection(CollectionGameSaves)
	db.SpellDefs = db.DB.Collection(CollectionSpellDefs)

	return nil
}

// Disconnect closes the database connection.
func (db *Database) Disconnect(ctx context.Context) error {
	if db.Client != nil {
		return db.Client.Disconnect(ctx)
	}
	return nil
}

// Ping checks if the database is reachable.
func (db *Database) Ping(ctx context.Context) error {
	if db.Client == nil {
		return mongo.ErrClientDisconnected
	}
	return db.Client.Ping(ctx, readpref.Primary())
}

// IsConnected returns true if connected to the database.
func (db *Database) IsConnected() bool {
	if db.Client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return db.Client.Ping(ctx, readpref.Primary()) == nil
}
