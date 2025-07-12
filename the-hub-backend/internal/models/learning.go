package models

import (
	"time"

	"gorm.io/gorm"
)

type Deck struct {
	ID        uint           `json:"deck_id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Cards     []Card         `json:"-"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type DeckUser struct {
	ID     uint `gorm:"primaryKey"`
	DeckID uint
	UserID uint
	Role   string // "owner", "editor", "viewer"
	Deck   Deck   `gorm:"foreignKey:DeckID"`
	User   User   `gorm:"foreignKey:UserID"`
}

type Card struct {
	ID           uint           `json:"card_id" gorm:"primaryKey"`
	DeckID       uint           `json:"deck_id" gorm:"not null"`
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
