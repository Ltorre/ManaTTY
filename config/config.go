package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	// MongoDB settings
	MongoDBURI string

	// Game settings
	GameTickRate     int // Ticks per second
	AutoSaveInterval int // Seconds between auto-saves

	// Logging
	LogLevel string

	// Development
	Debug bool
}

// DefaultConfig returns configuration with default values.
func DefaultConfig() *Config {
	return &Config{
		MongoDBURI:       "mongodb://localhost:27017/mage_tower",
		GameTickRate:     10,
		AutoSaveInterval: 30,
		LogLevel:         "info",
		Debug:            false,
	}
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := DefaultConfig()

	// MongoDB
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		cfg.MongoDBURI = uri
	}

	// Game settings
	if rate := os.Getenv("GAME_TICK_RATE"); rate != "" {
		if r, err := strconv.Atoi(rate); err == nil && r > 0 {
			cfg.GameTickRate = r
		}
	}

	if interval := os.Getenv("AUTO_SAVE_INTERVAL"); interval != "" {
		if i, err := strconv.Atoi(interval); err == nil && i > 0 {
			cfg.AutoSaveInterval = i
		}
	}

	// Logging
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}

	// Debug
	if debug := os.Getenv("DEBUG"); debug == "true" || debug == "1" {
		cfg.Debug = true
	}

	return cfg, nil
}

// GetEnv retrieves an environment variable with a default value.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retrieves an integer environment variable with a default.
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// GetEnvBool retrieves a boolean environment variable with a default.
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
