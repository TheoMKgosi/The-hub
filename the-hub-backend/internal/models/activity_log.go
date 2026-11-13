package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ActivityEntry struct {
	ID           uuid.UUID      `json:"activity_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Title        string         `json:"title" gorm:"not null"`
	Notes        string         `json:"notes"`
	Category     string         `json:"category"`
	Tags         pq.StringArray `json:"tags" gorm:"type:text[]"`
	ActivityDate time.Time      `json:"activity_date" gorm:"type:date;not null"`
	DurationMin  int            `json:"duration_min" gorm:"default:0"`
	User         User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
