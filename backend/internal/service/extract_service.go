package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz"
)

// ExtractService pulls plain text from uploaded documents.
type ExtractService struct{}

func NewExtractService() *ExtractService {
	return &ExtractService{}
}

// ExtractText detects file type by extension and extracts readable text.
func (s *ExtractService) ExtractText(filePath string) (string, string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".pdf":
		text, err := s.extractPDF(filePath)
		return text, "pdf", err
	case ".docx":
		text, err := s.extractDOCX(filePath)
		return text, "docx", err
	case ".txt":
		text, err := s.extractTXT(filePath)
		return text, "txt", err
	default:
		return "", "", fmt.Errorf("unsupported file type: %s (use PDF, DOCX, or TXT)", ext)
	}
}

func (s *ExtractService) extractPDF(path string) (string, error) {
	doc, err := fitz.New(path)
	if err != nil {
		return "", fmt.Errorf("open pdf: %w", err)
	}
	defer doc.Close()

	var b strings.Builder
	for i := 0; i < doc.NumPage(); i++ {
		text, err := doc.Text(i)
		if err != nil {
			continue
		}
		b.WriteString(text)
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String()), nil
}

func (s *ExtractService) extractDOCX(path string) (string, error) {
	// DOCX is a ZIP archive; we read document.xml via a lightweight parser
	content, err := readDocxText(path)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *ExtractService) extractTXT(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
