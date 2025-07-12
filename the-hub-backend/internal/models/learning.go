package models

import (
	"time"

	"gorm.io/gorm"
)

type Deck struct {
	ID        uint           `json:"deck_id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uint           `json:"user_id"`
	Cards     []Card         `json:"cards" gorm:"foreignKey:DeckID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeleteAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type Card struct {
	ID           uint           `json:"card_id" gorm:"primaryKey"`
	DeckID       uint           `json:"deck_id" gorm:"not null"`
	Question     string         `json:"question" gorm:"not null"`
	Answer       string         `json:"answer" gorm:"not null"`
	Easiness     float64        `json:"-" gorm:"default:2.5"` // SM-2 easiness factor
	Interval     int            `json:"-" gorm:"default:1"`   // Days until next review
	Repetitions  int            `json:"-" gorm:"default:0"`   // Successful reviews in a row
	LastReviewed time.Time      `json:"-"`                    // Last time card was reviewed
	NextReview   time.Time      `json:"-"`                    // When the card should next appear
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeleteAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
