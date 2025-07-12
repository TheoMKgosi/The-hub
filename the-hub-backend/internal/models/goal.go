package models

import (
	"time"

	"gorm.io/gorm"
)

type Goal struct {
	ID          uint           `json:"goal_id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Tasks       []Task         `json:"tasks"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
