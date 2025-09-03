package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Goal struct {
	ID             uuid.UUID      `json:"goal_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid"`
	Title          string         `json:"title" gorm:"not null"`
	Description    string         `json:"description"`
	DueDate        *time.Time     `json:"due_date"`
	Priority       *int           `json:"priority" gorm:"check:priority >= 1 AND priority <= 5"`
	Status         string         `json:"status" gorm:"default:active"`
	Category       string         `json:"category"`
	Color          string         `json:"color" gorm:"default:#3B82F6"`
	Progress       float64        `json:"progress" gorm:"default:0"` // Calculated field: 0-100
	TotalTasks     int            `json:"total_tasks" gorm:"default:0"`
	CompletedTasks int            `json:"completed_tasks" gorm:"default:0"`
	Tasks          []Task         `json:"tasks" gorm:"-"`
	User           User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// CalculateProgress calculates and updates the goal's progress based on task completion
func (g *Goal) CalculateProgress(db *gorm.DB) error {
	var totalTasks int64
	var completedTasks int64

	// Count total tasks for this goal
	if err := db.Model(&Task{}).Where("goal_id = ?", g.ID).Count(&totalTasks).Error; err != nil {
		return err
	}

	// Count completed tasks for this goal
	if err := db.Model(&Task{}).Where("goal_id = ? AND status = ?", g.ID, "completed").Count(&completedTasks).Error; err != nil {
		return err
	}

	g.TotalTasks = int(totalTasks)
	g.CompletedTasks = int(completedTasks)

	// Calculate progress percentage
	if totalTasks > 0 {
		g.Progress = float64(completedTasks) / float64(totalTasks) * 100
	} else {
		g.Progress = 0
	}

	// Update goal status based on progress
	if g.Progress == 100 && g.Status == "active" {
		g.Status = "completed"
	} else if g.Progress < 100 && g.Status == "completed" {
		g.Status = "active"
	}

	return nil
}

// IsOverdue checks if the goal is overdue
func (g *Goal) IsOverdue() bool {
	if g.DueDate == nil {
		return false
	}
	return g.DueDate.Before(time.Now()) && g.Status != "completed"
}

// GetDaysRemaining returns the number of days remaining until the goal is due
func (g *Goal) GetDaysRemaining() *int {
	if g.DueDate == nil {
		return nil
	}

	duration := time.Until(*g.DueDate)
	days := int(duration.Hours() / 24)
	return &days
}
