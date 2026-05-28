// Main entry point for the DeckForge API server.
// Run with: go run ./cmd/server
package main

import (
	"log"
	"os"

	"github.com/deckforge/backend/internal/config"
	"github.com/deckforge/backend/internal/database"
	"github.com/deckforge/backend/internal/handler"
	"github.com/deckforge/backend/internal/router"
	"github.com/deckforge/backend/internal/service"
)

func main() {
	// 1. Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if cfg.OpenAIKey == "" {
		log.Println("WARNING: OPENAI_API_KEY is not set — /generate will fail")
	}

	// 2. Connect to PostgreSQL
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	// 3. Ensure upload/export directories exist
	for _, dir := range []string{cfg.UploadDir, cfg.ExportDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	// 4. Initialize services (business logic layer)
	authService := service.NewAuthService(db, cfg)
	extractService := service.NewExtractService()
	openaiService := service.NewOpenAIService(cfg)
	presentationService := service.NewPresentationService(db, extractService, openaiService, cfg)
	exportService := service.NewExportService(db, cfg)

	// 5. Initialize handlers (HTTP layer)
	authHandler := handler.NewAuthHandler(authService)
	presentationHandler := handler.NewPresentationHandler(presentationService, exportService)

	// 6. Setup routes and start server
	engine := router.Setup(cfg, authHandler, presentationHandler, authService)

	addr := ":" + cfg.Port
	log.Printf("DeckForge API listening on http://localhost%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
