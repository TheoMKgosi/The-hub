package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepeatRule struct {
	ID        uuid.UUID `gorm:"primaryKey;type:text"`
	Frequency string    // "daily", "weekly", "monthly"
	Interval  int       // every X units
	ByDay     string    // JSON-encoded array: '["mon", "wed"]'
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook to generate UUID
func (rr *RepeatRule) BeforeCreate(tx *gorm.DB) error {
	if rr.ID == uuid.Nil {
		rr.ID = uuid.New()
	}
	return nil
}
