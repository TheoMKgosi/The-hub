package models

import "time"

type RepeatRule struct {
	ID        uint      `gorm:"primaryKey"`
	Frequency string    // "daily", "weekly", "monthly"
	Interval  int       // every X units
	ByDay     string    // JSON-encoded array: '["mon", "wed"]'
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

