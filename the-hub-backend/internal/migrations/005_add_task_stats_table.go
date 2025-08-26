package migrations

import (
	"gorm.io/gorm"
)

func Migrate005AddTaskStatsTable(db *gorm.DB) error {
	// Create task_stats table for storing aggregated task statistics
	return db.Exec(`
		CREATE TABLE task_stats (
			stats_id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			date DATE NOT NULL,
			total_tasks INTEGER DEFAULT 0,
			completed_tasks INTEGER DEFAULT 0,
			pending_tasks INTEGER DEFAULT 0,
			overdue_tasks INTEGER DEFAULT 0,
			priority_1_tasks INTEGER DEFAULT 0,
			priority_2_tasks INTEGER DEFAULT 0,
			priority_3_tasks INTEGER DEFAULT 0,
			priority_4_tasks INTEGER DEFAULT 0,
			priority_5_tasks INTEGER DEFAULT 0,
			priority_1_completed INTEGER DEFAULT 0,
			priority_2_completed INTEGER DEFAULT 0,
			priority_3_completed INTEGER DEFAULT 0,
			priority_4_completed INTEGER DEFAULT 0,
			priority_5_completed INTEGER DEFAULT 0,
			tasks_with_goals INTEGER DEFAULT 0,
			tasks_without_goals INTEGER DEFAULT 0,
			avg_completion_time_hours REAL,
			tasks_due_today INTEGER DEFAULT 0,
			tasks_due_tomorrow INTEGER DEFAULT 0,
			tasks_due_this_week INTEGER DEFAULT 0,
			completion_rate REAL DEFAULT 0,
			productivity_score REAL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, date)
		)
	`).Error
}

func Migrate005AddTaskStatsTableDown(db *gorm.DB) error {
	// Drop the task_stats table
	return db.Exec("DROP TABLE IF EXISTS task_stats").Error
}
