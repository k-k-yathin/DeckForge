// Package models defines database entities (GORM models).
// Each struct maps to a table in PostgreSQL.
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// User represents a registered DeckForge account.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	FullName     string    `gorm:"column:full_name;not null" json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// UploadedFile stores metadata about a user's uploaded document.
type UploadedFile struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	OriginalName  string    `gorm:"column:original_name;not null" json:"original_name"`
	StoredPath    string    `gorm:"column:stored_path;not null" json:"stored_path"`
	FileType      string    `gorm:"column:file_type;not null" json:"file_type"`
	FileSize      int64     `gorm:"column:file_size;not null" json:"file_size"`
	ExtractedText string    `gorm:"column:extracted_text;type:text" json:"extracted_text,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (UploadedFile) TableName() string { return "uploaded_files" }

// Presentation is a generated pitch deck.
type Presentation struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	UploadedFileID *uuid.UUID `gorm:"type:uuid" json:"uploaded_file_id,omitempty"`
	Title          string     `gorm:"not null" json:"title"`
	Status         string     `gorm:"not null;default:pending" json:"status"`
	SourceSummary  string     `gorm:"column:source_summary;type:text" json:"source_summary,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Slides         []Slide    `gorm:"foreignKey:PresentationID" json:"slides,omitempty"`
}

func (Presentation) TableName() string { return "presentations" }

// Slide is one slide inside a presentation.
type Slide struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PresentationID uuid.UUID      `gorm:"type:uuid;not null;index" json:"presentation_id"`
	SlideOrder     int            `gorm:"column:slide_order;not null" json:"slide_order"`
	SlideType      string         `gorm:"column:slide_type;not null" json:"slide_type"`
	Title          string         `gorm:"not null" json:"title"`
	Subtitle       string         `json:"subtitle,omitempty"`
	Content        datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (Slide) TableName() string { return "slides" }
