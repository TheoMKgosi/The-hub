package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate005UpdateCalendarZoneDefaults(db *gorm.DB) error {
	// Create calendar zone tables and insert default categories
	err := db.AutoMigrate(
		&models.CalendarZone{},
		&models.ZoneCategory{},
	)
	if err != nil {
		return err
	}

	// Insert default zone categories
	defaultCategories := models.GetDefaultZoneCategories()
	for _, category := range defaultCategories {
		if err := db.Create(&category).Error; err != nil {
			// Ignore duplicate key errors
			continue
		}
	}

	// Update the default value for allow_scheduling column to false
	// This will affect new records created after this migration
	if err := db.Exec("ALTER TABLE calendar_zones ALTER COLUMN allow_scheduling SET DEFAULT false").Error; err != nil {
		return err
	}

	return nil
}
