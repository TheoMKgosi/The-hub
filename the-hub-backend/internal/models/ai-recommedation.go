package models

import "time"


type AIRecommendation struct {
	ID             uint      `gorm:"primaryKey"`
	TaskID         uint      `gorm:"index"`
	Task           Task
	SuggestedStart time.Time
	SuggestedEnd   time.Time
	Confidence     float32
	Accepted       bool
	CreatedAt      time.Time
}


