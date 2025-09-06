package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate003AddCalendarIntegration(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.CalendarIntegration{},
		&models.CalendarEventSync{},
	)
}
