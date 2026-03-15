package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text"`
	Tags      string         `json:"tags" gorm:"type:jsonb;default:'[]'"` // JSON array of strings
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type NoteResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (n *Note) ToResponse() NoteResponse {
	tags := []string{}
	if n.Tags != "" && n.Tags != "[]" {
		_ = json.Unmarshal([]byte(n.Tags), &tags)
	}
	return NoteResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Content:   n.Content,
		Tags:      tags,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
