package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deckforge/backend/internal/config"
	"github.com/deckforge/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PresentationService orchestrates upload, generation, and retrieval.
type PresentationService struct {
	db       *gorm.DB
	extract  *ExtractService
	openai   *OpenAIService
	uploadDir string
}

func NewPresentationService(db *gorm.DB, extract *ExtractService, openai *OpenAIService, cfg *config.Config) *PresentationService {
	return &PresentationService{
		db:        db,
		extract:   extract,
		openai:    openai,
		uploadDir: cfg.UploadDir,
	}
}

// SaveUpload stores a file on disk and records metadata in the database.
func (s *PresentationService) SaveUpload(userID uuid.UUID, originalName string, data []byte) (*models.UploadedFile, error) {
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, err
	}

	ext := filepath.Ext(originalName)
	storedName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	storedPath := filepath.Join(s.uploadDir, storedName)

	if err := os.WriteFile(storedPath, data, 0644); err != nil {
		return nil, err
	}

	text, fileType, err := s.extract.ExtractText(storedPath)
	if err != nil {
		_ = os.Remove(storedPath)
		return nil, fmt.Errorf("extract text: %w", err)
	}

	record := models.UploadedFile{
		UserID:        userID,
		OriginalName:  originalName,
		StoredPath:    storedPath,
		FileType:      fileType,
		FileSize:      int64(len(data)),
		ExtractedText: text,
	}
	if err := s.db.Create(&record).Error; err != nil {
		_ = os.Remove(storedPath)
		return nil, err
	}
	return &record, nil
}

// SaveTextUpload creates an uploaded_files record from raw pasted text.
func (s *PresentationService) SaveTextUpload(userID uuid.UUID, text string) (*models.UploadedFile, error) {
	record := models.UploadedFile{
		UserID:        userID,
		OriginalName:  "pasted-text.txt",
		StoredPath:    "",
		FileType:      "txt",
		FileSize:      int64(len(text)),
		ExtractedText: text,
	}
	if err := s.db.Create(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// Generate creates a presentation from an uploaded file using OpenAI.
func (s *PresentationService) Generate(ctx context.Context, userID, fileID uuid.UUID) (*models.Presentation, error) {
	var file models.UploadedFile
	if err := s.db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	if file.ExtractedText == "" {
		return nil, errors.New("no text could be extracted from the file")
	}

	presentation := models.Presentation{
		UserID:         userID,
		UploadedFileID: &file.ID,
		Title:          "Generating...",
		Status:         "processing",
	}
	if err := s.db.Create(&presentation).Error; err != nil {
		return nil, err
	}

	deck, err := s.openai.GenerateSlides(ctx, file.ExtractedText)
	if err != nil {
		s.db.Model(&presentation).Updates(map[string]interface{}{
			"status": "failed",
			"title":  "Generation failed",
		})
		return nil, err
	}

	// Persist slides
	for i, sl := range deck.Slides {
		contentJSON, _ := json.Marshal(sl.Bullets)
		slide := models.Slide{
			PresentationID: presentation.ID,
			SlideOrder:     i + 1,
			SlideType:      sl.SlideType,
			Title:          sl.Title,
			Subtitle:       sl.Subtitle,
			Content:        datatypes.JSON(contentJSON),
		}
		if err := s.db.Create(&slide).Error; err != nil {
			return nil, err
		}
	}

	summary := truncate(file.ExtractedText, 500)
	s.db.Model(&presentation).Updates(map[string]interface{}{
		"title":          deck.Title,
		"status":         "completed",
		"source_summary": summary,
	})

	return s.GetByID(userID, presentation.ID)
}

// GenerateFromText generates directly from pasted text (no file upload).
func (s *PresentationService) GenerateFromText(ctx context.Context, userID uuid.UUID, text string) (*models.Presentation, error) {
	file, err := s.SaveTextUpload(userID, text)
	if err != nil {
		return nil, err
	}
	return s.Generate(ctx, userID, file.ID)
}

// ListByUser returns all presentations for a user, newest first.
func (s *PresentationService) ListByUser(userID uuid.UUID) ([]models.Presentation, error) {
	var list []models.Presentation
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

// GetByID returns one presentation with slides (must belong to user).
func (s *PresentationService) GetByID(userID, presentationID uuid.UUID) (*models.Presentation, error) {
	var p models.Presentation
	err := s.db.Preload("Slides", func(db *gorm.DB) *gorm.DB {
		return db.Order("slide_order ASC")
	}).Where("id = ? AND user_id = ?", presentationID, userID).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("presentation not found")
		}
		return nil, err
	}
	return &p, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
