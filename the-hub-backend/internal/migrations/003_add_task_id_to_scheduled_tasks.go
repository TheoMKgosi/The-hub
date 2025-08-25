package migrations

import (
	"gorm.io/gorm"
)

func Migrate003AddTaskIDToScheduledTasks(db *gorm.DB) error {
	// Add task_id column to scheduled_tasks table
	// Note: This uses raw SQL since GORM's AutoMigrate might not handle column additions
	// as explicitly as needed for this migration

	// For SQLite
	if err := db.Exec("ALTER TABLE scheduled_tasks ADD COLUMN task_id TEXT").Error; err != nil {
		return err
	}

	// For PostgreSQL, you might need:
	// ALTER TABLE scheduled_tasks ADD COLUMN task_id UUID;

	// Optional: Add foreign key constraint if needed
	// db.Exec("ALTER TABLE scheduled_tasks ADD CONSTRAINT fk_scheduled_tasks_task_id FOREIGN KEY (task_id) REFERENCES tasks(id)")

	// If you need to populate existing data, you can do it here
	// For example, if there are existing scheduled tasks that should be linked to tasks

	return nil
}

func Migrate003AddTaskIDToScheduledTasksDown(db *gorm.DB) error {
	// Remove the task_id column
	if err := db.Exec("ALTER TABLE scheduled_tasks DROP COLUMN task_id").Error; err != nil {
		return err
	}

	return nil
}
