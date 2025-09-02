package models

import (
	"time"

	"github.com/google/uuid"
)

type RepeatRule struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Frequency string    // "daily", "weekly", "monthly"
	Interval  int       // every X units
	ByDay     string    // JSON-encoded array: '["mon", "wed"]'
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
