package migrations

import (
	"gorm.io/gorm"
)

func Migrate004AddRecurrenceRuleIDToScheduledTasks(db *gorm.DB) error {
	// Add recurrence_rule_id column to scheduled_tasks table
	// Note: This uses raw SQL since GORM's AutoMigrate might not handle column additions
	// as explicitly as needed for this migration

	// For SQLite
	if err := db.Exec("ALTER TABLE scheduled_tasks ADD COLUMN recurrence_rule_id TEXT").Error; err != nil {
		return err
	}

	// For PostgreSQL, you might need:
	// ALTER TABLE scheduled_tasks ADD COLUMN recurrence_rule_id UUID;

	// Optional: Add foreign key constraint if needed
	// db.Exec("ALTER TABLE scheduled_tasks ADD CONSTRAINT fk_scheduled_tasks_recurrence_rule_id FOREIGN KEY (recurrence_rule_id) REFERENCES recurrence_rules(id)")

	// If you need to populate existing data, you can do it here
	// For example, if there are existing scheduled tasks that should be linked to recurrence rules

	return nil
}

func Migrate004AddRecurrenceRuleIDToScheduledTasksDown(db *gorm.DB) error {
	// Remove the recurrence_rule_id column
	if err := db.Exec("ALTER TABLE scheduled_tasks DROP COLUMN recurrence_rule_id").Error; err != nil {
		return err
	}

	return nil
}
