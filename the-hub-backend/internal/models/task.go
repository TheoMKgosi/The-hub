package models

import "time"

type Task struct {
	ID          uint       `json:"task_id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Priority    *int       `json:"priority" gorm:"check:priority >= 1 AND priority <= 5"`
	Status      string     `json:"status" gorm:"default:pending"`
	GoalID      *uint      `json:"goal_id" gorm:"foreignKey"`
	UserID      uint       `json:"user_id" gorm:"foreignKey"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   time.Time  `json:"-"`
}
