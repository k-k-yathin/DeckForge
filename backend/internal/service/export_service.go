package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deckforge/backend/internal/config"
	"github.com/deckforge/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

// ExportService builds PPTX and PDF files from stored presentations.
type ExportService struct {
	db        *gorm.DB
	exportDir string
}

func NewExportService(db *gorm.DB, cfg *config.Config) *ExportService {
	return &ExportService{db: db, exportDir: cfg.ExportDir}
}

// ExportPPTX creates a PowerPoint file and returns its path on disk.
func (s *ExportService) ExportPPTX(userID, presentationID uuid.UUID) (string, error) {
	p, err := s.loadPresentation(userID, presentationID)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(s.exportDir, 0755); err != nil {
		return "", err
	}

	outPath := filepath.Join(s.exportDir, fmt.Sprintf("%s.pptx", presentationID))
	if err := writePPTX(p, outPath); err != nil {
		return "", err
	}
	return outPath, nil
}

// ExportPDF creates a PDF with one page per slide.
func (s *ExportService) ExportPDF(userID, presentationID uuid.UUID) (string, error) {
	p, err := s.loadPresentation(userID, presentationID)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(s.exportDir, 0755); err != nil {
		return "", err
	}

	outPath := filepath.Join(s.exportDir, fmt.Sprintf("%s.pdf", presentationID))
	if err := writePDF(p, outPath); err != nil {
		return "", err
	}
	return outPath, nil
}

func (s *ExportService) loadPresentation(userID, id uuid.UUID) (*models.Presentation, error) {
	var p models.Presentation
	err := s.db.Preload("Slides", func(db *gorm.DB) *gorm.DB {
		return db.Order("slide_order ASC")
	}).Where("id = ? AND user_id = ?", id, userID).First(&p).Error
	if err != nil {
		return nil, fmt.Errorf("presentation not found")
	}
	return &p, nil
}

func writePDF(p *models.Presentation, outPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)

	for _, slide := range p.Slides {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 22)
		pdf.CellFormat(0, 12, slide.Title, "", 1, "C", false, 0, "")

		if slide.Subtitle != "" {
			pdf.SetFont("Arial", "", 14)
			pdf.CellFormat(0, 10, slide.Subtitle, "", 1, "C", false, 0, "")
		}

		pdf.Ln(8)
		pdf.SetFont("Arial", "", 12)

		var bullets []string
		_ = json.Unmarshal(slide.Content, &bullets)
		for _, b := range bullets {
			pdf.MultiCell(0, 7, "• "+b, "", "L", false)
		}
	}

	return pdf.OutputFileAndClose(outPath)
}
