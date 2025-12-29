package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Player represents a user profile in the game.
type Player struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UUID               string             `bson:"uuid" json:"uuid"`
	Username           string             `bson:"username" json:"username"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	LastPlayed         time.Time          `bson:"last_played" json:"last_played"`
	PlayDurationMs     int64              `bson:"play_duration_ms" json:"play_duration_ms"`
	TotalPrestigeCount int                `bson:"total_prestige_count" json:"total_prestige_count"`
	CurrentSaveSlot    int                `bson:"current_save_slot" json:"current_save_slot"`
	Version            int                `bson:"version" json:"version"`
}

// NewPlayer creates a new Player with default values.
func NewPlayer(uuid, username string) *Player {
	now := time.Now()
	return &Player{
		UUID:               uuid,
		Username:           username,
		CreatedAt:          now,
		LastPlayed:         now,
		PlayDurationMs:     0,
		TotalPrestigeCount: 0,
		CurrentSaveSlot:    0,
		Version:            1,
	}
}

// UpdateLastPlayed updates the last played timestamp.
func (p *Player) UpdateLastPlayed() {
	p.LastPlayed = time.Now()
}

// AddPlayDuration adds play time in milliseconds.
func (p *Player) AddPlayDuration(durationMs int64) {
	p.PlayDurationMs += durationMs
}
