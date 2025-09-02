package models

import (
	"time"

	"github.com/google/uuid"
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
