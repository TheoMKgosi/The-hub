package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate007Analytics(db *gorm.DB) error {
	// AutoMigrate will create the analytics tables
	return db.AutoMigrate(&models.TaskAnalytics{})
}
