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

// TaskStats represents aggregated statistics for task analytics
type TaskStats struct {
	ID     uuid.UUID `json:"stats_id" gorm:"primaryKey;type:text"`
	UserID uuid.UUID `json:"user_id" gorm:"not null"`
	Date   time.Time `json:"date" gorm:"not null"` // Date for which stats are calculated

	// Completion metrics
	TotalTasks     int `json:"total_tasks" gorm:"default:0"`
	CompletedTasks int `json:"completed_tasks" gorm:"default:0"`
	PendingTasks   int `json:"pending_tasks" gorm:"default:0"`
	OverdueTasks   int `json:"overdue_tasks" gorm:"default:0"`

	// Priority distribution
	Priority1Tasks int `json:"priority_1_tasks" gorm:"default:0"`
	Priority2Tasks int `json:"priority_2_tasks" gorm:"default:0"`
	Priority3Tasks int `json:"priority_3_tasks" gorm:"default:0"`
	Priority4Tasks int `json:"priority_4_tasks" gorm:"default:0"`
	Priority5Tasks int `json:"priority_5_tasks" gorm:"default:0"`

	// Completion rates by priority
	Priority1Completed int `json:"priority_1_completed" gorm:"default:0"`
	Priority2Completed int `json:"priority_2_completed" gorm:"default:0"`
	Priority3Completed int `json:"priority_3_completed" gorm:"default:0"`
	Priority4Completed int `json:"priority_4_completed" gorm:"default:0"`
	Priority5Completed int `json:"priority_5_completed" gorm:"default:0"`

	// Goal-related metrics
	TasksWithGoals    int `json:"tasks_with_goals" gorm:"default:0"`
	TasksWithoutGoals int `json:"tasks_without_goals" gorm:"default:0"`

	// Time-based metrics
	AvgCompletionTime *float64 `json:"avg_completion_time_hours"` // Average hours to complete tasks
	TasksDueToday     int      `json:"tasks_due_today" gorm:"default:0"`
	TasksDueTomorrow  int      `json:"tasks_due_tomorrow" gorm:"default:0"`
	TasksDueThisWeek  int      `json:"tasks_due_this_week" gorm:"default:0"`

	// Productivity scores
	CompletionRate    float64 `json:"completion_rate" gorm:"default:0"`    // Percentage of tasks completed
	ProductivityScore float64 `json:"productivity_score" gorm:"default:0"` // Calculated productivity metric

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// BeforeCreate hook to generate UUID
func (ts *TaskStats) BeforeCreate(tx *gorm.DB) error {
	if ts.ID == uuid.Nil {
		ts.ID = uuid.New()
	}
	return nil
}
