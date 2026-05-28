// Package config loads environment variables into a typed Config struct.
// Centralizing config makes it easy to see what the app needs to run.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application settings from environment variables.
type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	OpenAIKey    string
	OpenAIModel  string
	UploadDir    string
	ExportDir    string
	CORSOrigin   string
	JWTExpiryHrs int
}

// Load reads .env file (if present) and populates Config from environment.
func Load() (*Config, error) {
	// Ignore error if .env doesn't exist (e.g. in Docker)
	_ = godotenv.Load()

	expiry := 72
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			expiry = n
		}
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://deckforge:deckforge_secret@localhost:5432/deckforge?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "dev-secret-change-me"),
		OpenAIKey:    os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:  getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
		ExportDir:    getEnv("EXPORT_DIR", "./exports"),
		CORSOrigin:   getEnv("CORS_ORIGIN", "http://localhost:5173"),
		JWTExpiryHrs: expiry,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
