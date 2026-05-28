package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deckforge/backend/internal/middleware"
	"github.com/deckforge/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PresentationHandler handles uploads, generation, and exports.
type PresentationHandler struct {
	presentations *service.PresentationService
	exports       *service.ExportService
}

func NewPresentationHandler(p *service.PresentationService, e *service.ExportService) *PresentationHandler {
	return &PresentationHandler{presentations: p, exports: e}
}

// Upload godoc
// POST /api/v1/upload (multipart file)
func (h *PresentationHandler) Upload(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".pdf": true, ".docx": true, ".txt": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only PDF, DOCX, and TXT files are allowed"})
		return
	}

	if file.Size > 10<<20 { // 10 MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 10MB)"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read file"})
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read file"})
		return
	}

	record, err := h.presentations.SaveUpload(userID, file.Filename, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"file": gin.H{
			"id":             record.ID,
			"original_name":  record.OriginalName,
			"file_type":      record.FileType,
			"file_size":      record.FileSize,
			"extracted_text": truncatePreview(record.ExtractedText, 300),
		},
	})
}

type generateRequest struct {
	FileID string `json:"file_id"`
	Text   string `json:"text"`
}

// Generate godoc
// POST /api/v1/generate
func (h *PresentationHandler) Generate(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req generateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result interface{}
	var err error

	if req.Text != "" {
		result, err = h.presentations.GenerateFromText(c.Request.Context(), userID, req.Text)
	} else if req.FileID != "" {
		fileID, parseErr := uuid.Parse(req.FileID)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file_id"})
			return
		}
		result, err = h.presentations.Generate(c.Request.Context(), userID, fileID)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide file_id or text"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"presentation": result})
}

// ListPresentations godoc
// GET /api/v1/presentations
func (h *PresentationHandler) ListPresentations(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	list, err := h.presentations.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"presentations": list})
}

// GetPresentation godoc
// GET /api/v1/presentation/:id
func (h *PresentationHandler) GetPresentation(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	p, err := h.presentations.GetByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"presentation": p})
}

// ExportPPTX godoc
// GET /api/v1/presentation/:id/export/pptx
func (h *PresentationHandler) ExportPPTX(c *gin.Context) {
	h.serveExport(c, "pptx")
}

// ExportPDF godoc
// GET /api/v1/presentation/:id/export/pdf
func (h *PresentationHandler) ExportPDF(c *gin.Context) {
	h.serveExport(c, "pdf")
}

func (h *PresentationHandler) serveExport(c *gin.Context, format string) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var path string
	switch format {
	case "pptx":
		path, err = h.exports.ExportPPTX(userID, id)
	case "pdf":
		path, err = h.exports.ExportPDF(userID, id)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown format"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := id.String() + "." + format
	c.FileAttachment(path, filename)
}

func truncatePreview(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
