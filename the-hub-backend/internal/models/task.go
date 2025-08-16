package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	ID          uint           `json:"task_id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	DueDate     *time.Time     `json:"due_date"`
	Priority    *int           `json:"priority" gorm:"check:priority >= 1 AND priority <= 5"`
	Status      string         `json:"status" gorm:"default:pending"`
	OrderIndex  int            `json:"order" gorm:"default:0"`
	GoalID      *uint          `json:"goal_id"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	Goal        Goal           `json:"-" gorm:"foreignKey:GoalID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
