package unit

import (
	"testing"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskScore(t *testing.T) {
	now := time.Now()
	future := now.Add(48 * time.Hour) // Use 48 hours to ensure it's not treated as urgent

	tests := []struct {
		name     string
		task     models.Task
		expected int
	}{
		{
			name: "high priority task",
			task: models.Task{
				Title:    "High Priority",
				Priority: intPtr(5),
				DueDate:  &future,
			},
			expected: 75, // 5 * 10 + 25 (for due date in 24-72 hours)
		},
		{
			name: "overdue task",
			task: models.Task{
				Title:   "Overdue",
				DueDate: &now,
			},
			expected: 100, // overdue bonus
		},
		{
			name: "recurring task",
			task: models.Task{
				Title:       "Recurring",
				IsRecurring: true,
			},
			expected: 3, // recurring bonus
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := taskScore(tt.task)
			assert.Equal(t, tt.expected, score)
		})
	}
}

func intPtr(i int) *int {
	return &i
}

// taskScore calculates a priority score for a task (copied from ai package for testing)
func taskScore(task models.Task) int {
	score := 0

	// Priority score (higher priority = higher score)
	if task.Priority != nil {
		score += *task.Priority * 10
	}

	// Deadline proximity score
	if task.DueDate != nil {
		hoursUntilDue := time.Until(*task.DueDate).Hours()
		if hoursUntilDue < 0 {
			// Overdue tasks get highest priority
			score += 100
		} else if hoursUntilDue < 24 {
			score += 50
		} else if hoursUntilDue < 72 {
			score += 25
		}
	}

	// Time estimate consideration (shorter tasks get slight preference for scheduling)
	if task.TimeEstimate != nil && *task.TimeEstimate <= 30 {
		score += 5
	}

	// Recurring tasks get slight boost
	if task.IsRecurring {
		score += 3
	}

	return score
}
