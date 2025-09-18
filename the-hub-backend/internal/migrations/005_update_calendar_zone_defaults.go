package migrations

import (
	"gorm.io/gorm"
)

func Migrate005UpdateCalendarZoneDefaults(db *gorm.DB) error {
	// Update the default value for allow_scheduling column to false
	// This will affect new records created after this migration
	if err := db.Exec("ALTER TABLE calendar_zones ALTER COLUMN allow_scheduling SET DEFAULT false").Error; err != nil {
		return err
	}

	return nil
}
