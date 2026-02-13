package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID           uuid.UUID  `json:"task_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title        string     `json:"title" gorm:"not null"`
	Description  string     `json:"description"`
	DueDate      *time.Time `json:"due_date"`
	Priority     *int       `json:"priority" gorm:"check:priority >= 1 AND priority <= 5"`
	Status       string     `json:"status" gorm:"default:pending"`
	OrderIndex   int        `json:"order" gorm:"default:0"`
	GoalID       *uuid.UUID `json:"goal_id" gorm:"type:uuid"`
	ParentTaskID *uuid.UUID `json:"parent_task_id" gorm:"type:uuid"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	// Time tracking fields
	StartTime        *time.Time `json:"start_time"`
	TimeEstimate     *int       `json:"time_estimate_minutes"`               // Estimated time in minutes
	TimeSpent        int        `json:"time_spent_minutes" gorm:"default:0"` // Total time spent in minutes
	IsRecurring      bool       `json:"is_recurring" gorm:"default:false"`
	RecurrenceRuleID *uuid.UUID `json:"recurrence_rule_id" gorm:"type:uuid"`

	// Task classification fields
	Category       string          `json:"category"`                // work, study, personal, creative, etc.
	TaskType       string          `json:"task_type"`               // meeting, development, learning, exercise, etc.
	Tags           []string        `json:"tags" gorm:"type:text[]"` // Flexible tagging system
	User           User            `json:"-" gorm:"foreignKey:UserID"`
	Goal           Goal            `json:"-" gorm:"foreignKey:GoalID"`
	ParentTask     *Task           `json:"-" gorm:"foreignKey:ParentTaskID"`
	Subtasks       []Task          `json:"subtasks" gorm:"foreignKey:ParentTaskID"`
	TimeEntries    []TimeEntry     `json:"time_entries" gorm:"foreignKey:TaskID"`
	RecurrenceRule *RecurrenceRule `json:"-" gorm:"foreignKey:RecurrenceRuleID"`
	CreatedAt      time.Time       `json:"-"`
	UpdatedAt      time.Time       `json:"-"`
	DeletedAt      gorm.DeletedAt  `json:"-" gorm:"index"`
}

// IsSubtask returns true if this task is a subtask
func (t *Task) IsSubtask() bool {
	return t.ParentTaskID != nil
}

// GetSubtasks returns all subtasks for this task
func (t *Task) GetSubtasks(db *gorm.DB) ([]Task, error) {
	var subtasks []Task
	err := db.Where("parent_task_id = ? AND user_id = ?", t.ID, t.UserID).Find(&subtasks).Error
	return subtasks, err
}

// GetDependencies returns all tasks that this task depends on
// TODO: Fix this duplication and find out if it is necessary
func (t *Task) GetDependencies(db *gorm.DB) ([]Task, error) {
	var dependencies []Task
	err := db.Joins("JOIN task_dependencies ON tasks.id = task_dependencies.depends_on_id").
		Where("task_dependencies.task_id = ? AND task_dependencies.user_id = ?", t.ID, t.UserID).
		Find(&dependencies).Error
	return dependencies, err
}

// GetDependents returns all tasks that depend on this task
func (t *Task) GetDependents(db *gorm.DB) ([]Task, error) {
	var dependents []Task
	err := db.Joins("JOIN task_dependencies ON tasks.id = task_dependencies.task_id").
		Where("task_dependencies.depends_on_id = ? AND task_dependencies.user_id = ?", t.ID, t.UserID).
		Find(&dependents).Error
	return dependents, err
}

// CanBeCompleted checks if this task can be marked as completed
// A task can be completed if all its dependencies are completed
func (t *Task) CanBeCompleted(db *gorm.DB) (bool, error) {
	if t.Status == "completed" {
		return true, nil
	}

	dependencies, err := t.GetDependencies(db)
	if err != nil {
		return false, err
	}

	// Check if all dependencies are completed
	for _, dep := range dependencies {
		if dep.Status != "completed" {
			return false, nil
		}
	}

	return true, nil
}

// UpdateParentStatus updates the parent task status based on subtasks
func (t *Task) UpdateParentStatus(db *gorm.DB) error {
	if t.ParentTaskID == nil {
		return nil
	}

	var parentTask Task
	if err := db.Where("id = ?", *t.ParentTaskID).First(&parentTask).Error; err != nil {
		return err
	}

	subtasks, err := parentTask.GetSubtasks(db)
	if err != nil {
		return err
	}

	// If all subtasks are completed, mark parent as completed
	allCompleted := true
	for _, subtask := range subtasks {
		if subtask.Status != "completed" {
			allCompleted = false
			break
		}
	}

	if allCompleted && parentTask.Status != "completed" {
		return db.Model(&parentTask).Update("status", "completed").Error
	} else if !allCompleted && parentTask.Status == "completed" {
		return db.Model(&parentTask).Update("status", "pending").Error
	}

	return nil
}

// StartTimeTracking starts a new time entry for the task
func (t *Task) StartTimeTracking(db *gorm.DB, userID uuid.UUID, description string) (*TimeEntry, error) {
	// Stop any currently running time entries for this user
	if err := db.Model(&TimeEntry{}).Where("user_id = ? AND is_running = ?", userID, true).Update("is_running", false).Error; err != nil {
		return nil, err
	}

	timeEntry := TimeEntry{
		TaskID:      t.ID,
		UserID:      userID,
		Description: description,
		StartTime:   time.Now(),
		IsRunning:   true,
		Duration:    0,
	}

	if err := db.Create(&timeEntry).Error; err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// StopTimeTracking stops the currently running time entry for the task
func (t *Task) StopTimeTracking(db *gorm.DB, userID uuid.UUID) (*TimeEntry, error) {
	var timeEntry TimeEntry
	if err := db.Where("task_id = ? AND user_id = ? AND is_running = ?", t.ID, userID, true).First(&timeEntry).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	duration := int(now.Sub(timeEntry.StartTime).Minutes())

	timeEntry.EndTime = &now
	timeEntry.Duration = duration
	timeEntry.IsRunning = false

	if err := db.Save(&timeEntry).Error; err != nil {
		return nil, err
	}

	// Update task's total time spent
	if err := db.Model(t).Update("time_spent_minutes", gorm.Expr("time_spent_minutes + ?", duration)).Error; err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// GetTotalTimeSpent calculates the total time spent on the task including all time entries
func (t *Task) GetTotalTimeSpent(db *gorm.DB) (int, error) {
	var total int
	if err := db.Model(&TimeEntry{}).Where("task_id = ?", t.ID).Select("COALESCE(SUM(duration), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetRunningTimeEntry returns the currently running time entry for this task and user
func (t *Task) GetRunningTimeEntry(db *gorm.DB, userID uuid.UUID) (*TimeEntry, error) {
	var timeEntry TimeEntry
	if err := db.Where("task_id = ? AND user_id = ? AND is_running = ?", t.ID, userID, true).First(&timeEntry).Error; err != nil {
		return nil, err
	}
	return &timeEntry, nil
}

// CreateFromTemplate creates a new task from a template
func (tt *TaskTemplate) CreateFromTemplate(userID uuid.UUID) *Task {
	return &Task{
		UserID:       userID,
		Title:        tt.TitleTemplate,
		Description:  tt.DescriptionTemplate,
		Priority:     tt.Priority,
		TimeEstimate: tt.TimeEstimate,
		Category:     tt.Category,
	}
}

// TaskStats represents aggregated statistics for task analytics
type TaskStats struct {
	ID     uuid.UUID `json:"stats_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
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

// TaskDependency represents a dependency relationship between tasks
type TaskDependency struct {
	ID          uuid.UUID      `json:"dependency_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TaskID      uuid.UUID      `json:"task_id" gorm:"type:uuid;not null"`       // The task that depends on another
	DependsOnID uuid.UUID      `json:"depends_on_id" gorm:"type:uuid;not null"` // The task that must be completed first
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Task        Task           `json:"-" gorm:"foreignKey:TaskID"`
	DependsOn   Task           `json:"-" gorm:"foreignKey:DependsOnID"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TimeEntry represents a time tracking entry for a task
type TimeEntry struct {
	ID          uuid.UUID      `json:"time_entry_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TaskID      uuid.UUID      `json:"task_id" gorm:"type:uuid;not null"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Description string         `json:"description"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     *time.Time     `json:"end_time"`
	Duration    int            `json:"duration_minutes" gorm:"not null"` // Duration in minutes
	IsRunning   bool           `json:"is_running" gorm:"default:false"`
	Task        Task           `json:"-" gorm:"foreignKey:TaskID"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TaskTemplate represents a reusable task template
type TaskTemplate struct {
	ID          uuid.UUID `json:"template_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	// Template fields
	TitleTemplate       string         `json:"title_template"`
	DescriptionTemplate string         `json:"description_template"`
	Priority            *int           `json:"priority"`
	TimeEstimate        *int           `json:"time_estimate_minutes"`
	Tags                string         `json:"tags"` // JSON string of tags
	IsPublic            bool           `json:"is_public" gorm:"default:false"`
	UsageCount          int            `json:"usage_count" gorm:"default:0"`
	User                User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt           time.Time      `json:"-"`
	UpdatedAt           time.Time      `json:"-"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

// TaskAnalytics represents productivity analytics for a user
type TaskAnalytics struct {
	UserID uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	Date   time.Time `json:"date" gorm:"primaryKey"`
	// Task metrics
	TotalTasks     int `json:"total_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	PendingTasks   int `json:"pending_tasks"`
	OverdueTasks   int `json:"overdue_tasks"`
	// Time metrics
	TotalTimeEstimated int     `json:"total_time_estimated_minutes"`
	TotalTimeSpent     int     `json:"total_time_spent_minutes"`
	AverageTaskTime    float64 `json:"average_task_time_minutes"`
	// Productivity scores
	CompletionRate     float64 `json:"completion_rate"`
	ProductivityScore  float64 `json:"productivity_score"`
	OnTimeDeliveryRate float64 `json:"on_time_delivery_rate"`
	// Goal metrics
	ActiveGoals        int     `json:"active_goals"`
	CompletedGoals     int     `json:"completed_goals"`
	GoalCompletionRate float64 `json:"goal_completion_rate"`
	// Template usage
	TemplatesUsed int `json:"templates_used"`
	// Recurring tasks
	RecurringTasks int       `json:"recurring_tasks"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// CalculateProductivityScore calculates a productivity score based on various metrics
func (ta *TaskAnalytics) CalculateProductivityScore() {
	if ta.TotalTasks == 0 {
		ta.ProductivityScore = 0
		return
	}

	// Base completion rate (40% weight)
	completionScore := ta.CompletionRate * 0.4

	// On-time delivery (30% weight)
	onTimeScore := ta.OnTimeDeliveryRate * 0.3

	// Goal completion (20% weight)
	goalScore := ta.GoalCompletionRate * 0.2

	// Template usage bonus (10% weight) - encourages using templates
	templateBonus := 0.0
	if ta.TemplatesUsed > 0 {
		templateBonus = float64(ta.TemplatesUsed) / float64(ta.TotalTasks) * 0.1
		if templateBonus > 0.1 {
			templateBonus = 0.1
		}
	}

	ta.ProductivityScore = completionScore + onTimeScore + goalScore + templateBonus
}

// UserProductivityInsights represents aggregated insights for a user
type UserProductivityInsights struct {
	UserID    uuid.UUID `json:"user_id"`
	Period    string    `json:"period"` // daily, weekly, monthly
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	// Overall metrics
	TotalTasks               int     `json:"total_tasks"`
	TotalCompleted           int     `json:"total_completed"`
	AverageCompletionRate    float64 `json:"average_completion_rate"`
	AverageProductivityScore float64 `json:"average_productivity_score"`
	// Time metrics
	TotalTimeSpent     int     `json:"total_time_spent_minutes"`
	AverageTimePerTask float64 `json:"average_time_per_task_minutes"`
	// Trends
	CompletionRateTrend float64 `json:"completion_rate_trend"`
	ProductivityTrend   float64 `json:"productivity_trend"`
	// Peak performance
	BestDay      *time.Time `json:"best_day"`
	BestDayScore float64    `json:"best_day_score"`
	// Patterns
	MostProductiveHour      *int `json:"most_productive_hour"`
	MostProductiveDayOfWeek *int `json:"most_productive_day_of_week"`
	// Goals
	GoalsCompleted  int     `json:"goals_completed"`
	GoalSuccessRate float64 `json:"goal_success_rate"`
	// Recommendations
	Recommendations []string `json:"recommendations"`
}

// TaskShare represents sharing a task with another user
type TaskShare struct {
	ID           uuid.UUID      `json:"share_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TaskID       uuid.UUID      `json:"task_id" gorm:"type:uuid;not null"`
	OwnerID      uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	SharedWithID uuid.UUID      `json:"shared_with_id" gorm:"type:uuid;not null"`
	Permission   string         `json:"permission" gorm:"not null"` // view, edit, admin
	Task         Task           `json:"-" gorm:"foreignKey:TaskID"`
	Owner        User           `json:"-" gorm:"foreignKey:OwnerID"`
	SharedWith   User           `json:"-" gorm:"foreignKey:SharedWithID"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// GoalShare represents sharing a goal with another user
type GoalShare struct {
	ID           uuid.UUID      `json:"share_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GoalID       uuid.UUID      `json:"goal_id" gorm:"type:uuid;not null"`
	OwnerID      uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	SharedWithID uuid.UUID      `json:"shared_with_id" gorm:"type:uuid;not null"`
	Permission   string         `json:"permission" gorm:"not null"` // view, edit, admin
	Goal         Goal           `json:"-" gorm:"foreignKey:GoalID"`
	Owner        User           `json:"-" gorm:"foreignKey:OwnerID"`
	SharedWith   User           `json:"-" gorm:"foreignKey:SharedWithID"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TaskComment represents a comment on a task
type TaskComment struct {
	ID        uuid.UUID      `json:"comment_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TaskID    uuid.UUID      `json:"task_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Content   string         `json:"content" gorm:"not null"`
	Task      Task           `json:"-" gorm:"foreignKey:TaskID"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// GoalComment represents a comment on a goal
type GoalComment struct {
	ID        uuid.UUID      `json:"comment_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GoalID    uuid.UUID      `json:"goal_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Content   string         `json:"content" gorm:"not null"`
	Goal      Goal           `json:"-" gorm:"foreignKey:GoalID"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
