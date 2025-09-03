package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate006TimeTrackingTemplates(db *gorm.DB) error {
	// AutoMigrate will add new columns to existing task table and create new tables
	return db.AutoMigrate(&models.Task{}, &models.TimeEntry{}, &models.TaskTemplate{}, &models.RecurrenceRule{})
}
