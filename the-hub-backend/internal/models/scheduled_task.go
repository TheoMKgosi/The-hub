package models

import "time"

// ScheduledTask represents a calendar event for a Task with a default one-hour duration.
type ScheduledTask struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	// TaskID      *uint      `json:"task_id" gorm:"not null;uniqueIndex"`
	// Task        *Task      `json:"-"`
	Title       string    `json:"title" gorm:"not null"`
	Start       time.Time `json:"start" gorm:"not null"`
	End         time.Time `json:"end" gorm:"not null"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"-" gorm:"foreignKey:UserID"`
	CreatedByAI bool      `json:"created_by_ai" gorm:"default:false"`
	CreateAt    time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
