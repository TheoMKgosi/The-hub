package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deck struct {
	ID        uuid.UUID      `json:"deck_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Cards     []Card         `json:"-"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type DeckUser struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DeckID uuid.UUID `gorm:"type:uuid"`
	UserID uuid.UUID `gorm:"type:uuid"`
	Role   string    // "owner", "editor", "viewer"
	Deck   Deck      `gorm:"foreignKey:DeckID"`
	User   User      `gorm:"foreignKey:UserID"`
}

type Card struct {
	ID           uuid.UUID      `json:"card_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DeckID       uuid.UUID      `json:"deck_id" gorm:"type:uuid;not null"`
	Question     string         `json:"question" gorm:"not null"`
	Answer       string         `json:"answer" gorm:"not null"`
	Easiness     float64        `json:"-" gorm:"default:2.5"`     // SM-2 easiness factor
	Interval     int            `json:"-" gorm:"default:1"`       // Days until next review
	Repetitions  int            `json:"-" gorm:"default:0"`       // Successful reviews in a row
	LastReviewed time.Time      `json:"last_review"`              // Last time card was reviewed
	NextReview   time.Time      `json:"next_review" gorm:"index"` // When the card should next appear
	Deck         Deck           `json:"-" gorm:"foreignKey:DeckID"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
