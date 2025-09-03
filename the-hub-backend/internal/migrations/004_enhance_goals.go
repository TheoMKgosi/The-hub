package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate004EnhanceGoals(db *gorm.DB) error {
	// AutoMigrate will add new columns to existing goal table
	return db.AutoMigrate(&models.Goal{})
}
