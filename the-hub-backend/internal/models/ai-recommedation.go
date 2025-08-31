package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIRecommendation struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TaskID         uuid.UUID `gorm:"type:uuid;index"`
	Task           Task
	SuggestedStart time.Time
	SuggestedEnd   time.Time
	Confidence     float32
	Accepted       bool
	CreatedAt      time.Time
}

// BeforeCreate hook to generate UUID
func (ar *AIRecommendation) BeforeCreate(tx *gorm.DB) error {
	if ar.ID == uuid.Nil {
		ar.ID = uuid.New()
	}
	return nil
}
