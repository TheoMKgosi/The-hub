package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate008Collaboration(db *gorm.DB) error {
	// AutoMigrate will create the collaboration tables
	return db.AutoMigrate(&models.TaskShare{}, &models.GoalShare{}, &models.TaskComment{}, &models.GoalComment{})
}
