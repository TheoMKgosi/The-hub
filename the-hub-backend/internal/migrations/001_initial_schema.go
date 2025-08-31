package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate001InitialSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Goal{},
		&models.Task{},
		&models.ScheduledTask{},
		&models.Deck{},
		&models.Card{},
		&models.Budget{},
		&models.BudgetCategory{},
		&models.Income{},
		&models.Topic{},
		&models.Tag{},
		&models.Task_learning{},
		&models.User{},
		&models.TaskStats{},
	)
}
