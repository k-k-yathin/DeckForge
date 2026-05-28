// Package database connects to PostgreSQL using GORM (Go ORM).
package database

import (
	"fmt"
	"log"

	"github.com/deckforge/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a PostgreSQL connection and auto-migrates tables.
func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// AutoMigrate creates/updates tables to match our Go structs
	if err := db.AutoMigrate(
		&models.User{},
		&models.UploadedFile{},
		&models.Presentation{},
		&models.Slide{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}
