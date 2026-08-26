package models

import (
	"time"

	"github.com/google/uuid"
)

type FocusSession struct {
	ID          uuid.UUID  `json:"focus_session_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	TaskID      *uuid.UUID `json:"task_id" gorm:"type:uuid"`
	GoalID      *uuid.UUID `json:"goal_id" gorm:"type:uuid"`
	DurationMin int        `json:"duration_min" gorm:"default:0"`
	StartedAt   time.Time  `json:"started_at" gorm:"not null"`
	EndedAt     *time.Time `json:"ended_at"`
	Status      string     `json:"status" gorm:"default:'active'"`       // active, completed, interrupted
	SessionType string     `json:"session_type" gorm:"default:'focus'"`  // focus, break
	CycleGroup  *uuid.UUID `json:"cycle_group" gorm:"type:uuid"`         // groups work+break sessions from one pomo cycle
	CycleOrder  int        `json:"cycle_order" gorm:"default:0"`         // position within a cycle group
	Notes       string     `json:"notes"`
	User        User       `json:"-" gorm:"foreignKey:UserID"`
	Task        Task       `json:"-" gorm:"foreignKey:TaskID"`
	Goal        Goal       `json:"-" gorm:"foreignKey:GoalID"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
}
