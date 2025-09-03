package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate005EnhanceTasks(db *gorm.DB) error {
	// AutoMigrate will add new columns to existing task table and create task_dependencies table
	return db.AutoMigrate(&models.Task{}, &models.TaskDependency{})
}
