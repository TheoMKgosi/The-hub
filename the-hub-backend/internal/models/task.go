package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID      `json:"task_id" gorm:"primaryKey;type:text"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	DueDate     *time.Time     `json:"due_date"`
	Priority    *int           `json:"priority" gorm:"check:priority >= 1 AND priority <= 5"`
	Status      string         `json:"status" gorm:"default:pending"`
	OrderIndex  int            `json:"order" gorm:"default:0"`
	GoalID      *uuid.UUID     `json:"goal_id"`
	UserID      uuid.UUID      `json:"user_id"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	Goal        Goal           `json:"-" gorm:"foreignKey:GoalID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
