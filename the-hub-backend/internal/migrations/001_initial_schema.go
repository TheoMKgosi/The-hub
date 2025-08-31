package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate001InitialSchema(db *gorm.DB) error {
	// For PostgreSQL with UUID, we need to ensure proper table creation order
	// and handle potential existing tables with TEXT columns

	// First, create tables without foreign keys
	err := db.AutoMigrate(
		&models.User{},
		&models.Goal{},
		&models.Task{},
		&models.ScheduledTask{},
		&models.Deck{},
		&models.Card{},
		&models.BudgetCategory{},
		&models.Budget{},
		&models.Income{},
		&models.Transaction{},
		&models.Topic{},
		&models.Task_learning{},
		&models.Resource{},
		&models.StudySession{},
		&models.Tag{},
		&models.AIRecommendation{},
		&models.TaskStats{},
		&models.RecurrenceRule{},
		&models.RepeatRule{},
	)
	if err != nil {
		return err
	}

	// Create junction tables after main tables
	return db.AutoMigrate(
		&models.DeckUser{},
	)
}
