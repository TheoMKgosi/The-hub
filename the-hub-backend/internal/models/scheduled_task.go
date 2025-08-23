package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduledTask represents a calendar event for a Task with a default one-hour duration.
type ScheduledTask struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:text"`
	// TaskID      *uuid.UUID `json:"task_id" gorm:"not null;uniqueIndex"`
	// Task        *Task      `json:"-"`
	Title       string    `json:"title" gorm:"not null"`
	Start       time.Time `json:"start" gorm:"not null"`
	End         time.Time `json:"end" gorm:"not null"`
	UserID      uuid.UUID `json:"user_id"`
	User        User      `json:"-" gorm:"foreignKey:UserID"`
	CreatedByAI bool      `json:"created_by_ai" gorm:"default:false"`
	CreateAt    time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// BeforeCreate hook to generate UUID
func (st *ScheduledTask) BeforeCreate(tx *gorm.DB) error {
	if st.ID == uuid.Nil {
		st.ID = uuid.New()
	}
	return nil
}
