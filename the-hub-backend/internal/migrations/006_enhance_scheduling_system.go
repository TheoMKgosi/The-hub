package migrations

import (
	"fmt"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate006EnhanceSchedulingSystem(db *gorm.DB) error {
	// Add new columns to tasks table
	if err := db.Exec(`
		ALTER TABLE tasks
		ADD COLUMN IF NOT EXISTS category TEXT,
		ADD COLUMN IF NOT EXISTS task_type TEXT,
		ADD COLUMN IF NOT EXISTS tags TEXT[]
	`).Error; err != nil {
		return err
	}

	// Verify columns were added by checking if they exist
	var columnCount int64
	if err := db.Raw(`
		SELECT COUNT(*) FROM information_schema.columns
		WHERE table_name = 'tasks' AND column_name IN ('category', 'task_type', 'tags')
	`).Scan(&columnCount).Error; err != nil {
		return err
	}

	if columnCount < 3 {
		return fmt.Errorf("failed to add required columns to tasks table")
	}

	// Add new columns to calendar_zones table
	if err := db.Exec(`
		ALTER TABLE calendar_zones
		ADD COLUMN IF NOT EXISTS scheduling_mode TEXT DEFAULT 'none',
		ADD COLUMN IF NOT EXISTS allowed_task_categories TEXT[],
		ADD COLUMN IF NOT EXISTS allowed_task_types TEXT[],
		ADD COLUMN IF NOT EXISTS blocked_task_categories TEXT[],
		ADD COLUMN IF NOT EXISTS blocked_task_types TEXT[],
		ADD COLUMN IF NOT EXISTS allow_non_zone_scheduling BOOLEAN DEFAULT true,
		ADD COLUMN IF NOT EXISTS non_zone_start_time TIMESTAMP,
		ADD COLUMN IF NOT EXISTS non_zone_end_time TIMESTAMP,
		ADD COLUMN IF NOT EXISTS non_zone_days_of_week TEXT[]
	`).Error; err != nil {
		return err
	}

	// Set default non-zone times (9 AM - 6 PM) for existing zones
	defaultStartTime := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
	defaultEndTime := time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)
	defaultDays := []string{"monday", "tuesday", "wednesday", "thursday", "friday"}

	if err := db.Model(&models.CalendarZone{}).Where("non_zone_start_time IS NULL").Update("non_zone_start_time", defaultStartTime).Error; err != nil {
		return err
	}
	if err := db.Model(&models.CalendarZone{}).Where("non_zone_end_time IS NULL").Update("non_zone_end_time", defaultEndTime).Error; err != nil {
		return err
	}
	if err := db.Model(&models.CalendarZone{}).Where("non_zone_days_of_week IS NULL OR array_length(non_zone_days_of_week, 1) IS NULL").Update("non_zone_days_of_week", defaultDays).Error; err != nil {
		return err
	}

	// Auto-migrate to ensure all changes are applied
	return db.AutoMigrate(
		&models.Task{},
		&models.CalendarZone{},
	)
}
