package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecurrenceRule defines how an event repeats
type RecurrenceRule struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:text"`
	Frequency  string     `json:"frequency" gorm:"not null"` // daily, weekly, monthly
	Interval   int        `json:"interval" gorm:"default:1"` // every N days/weeks/months
	EndDate    *time.Time `json:"end_date"`
	Count      *int       `json:"count"`        // number of occurrences
	ByDay      string     `json:"by_day"`       // e.g., "MO,TU,WE" for weekly
	ByMonthDay *int       `json:"by_month_day"` // day of month
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
}

// BeforeCreate hook to generate UUID
func (rr *RecurrenceRule) BeforeCreate(tx *gorm.DB) error {
	if rr.ID == uuid.Nil {
		rr.ID = uuid.New()
	}
	return nil
}
