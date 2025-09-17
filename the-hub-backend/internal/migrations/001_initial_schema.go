package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate001InitialSchema(db *gorm.DB) error {
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
		&models.LearningPath{},
		&models.LearningPathTopic{},
		&models.AIRecommendation{},
		&models.TaskStats{},
		&models.RecurrenceRule{},
		&models.RepeatRule{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.PasswordResetToken{},
		&models.TimeEntry{},
		&models.TaskAnalytics{},
		&models.TaskShare{},
		&models.GoalShare{},
		&models.TaskComment{},
		&models.GoalComment{},
		&models.CalendarZone{},
		&models.ZoneCategory{},
	)
	if err != nil {
		return err
	}

	// Create junction tables after main tables
	if err := db.AutoMigrate(
		&models.DeckUser{},
		&models.TaskDependency{},
		&models.TaskTemplate{},
	); err != nil {
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

	return nil
}
