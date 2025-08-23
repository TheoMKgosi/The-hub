package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate002UUIDTransition(db *gorm.DB) error {
	// This migration transitions all primary keys from uint to UUID
	// Since this is a complex migration involving primary keys, we'll use
	// GORM's AutoMigrate which will handle the schema changes appropriately
	// for the new UUID-based models

	// Note: For SQLite, we need to disable foreign key constraints during migration
	// and then re-enable them afterward to avoid issues with primary key changes

	// Disable foreign key constraints for SQLite
	db.Exec("PRAGMA foreign_keys = OFF")

	err := db.AutoMigrate(
		&models.User{},             // Already uses UUID, but include for consistency
		&models.Goal{},             // Updated to use UUID
		&models.Task{},             // Updated to use UUID
		&models.ScheduledTask{},    // Updated to use UUID
		&models.Deck{},             // Updated to use UUID
		&models.DeckUser{},         // Updated to use UUID
		&models.Card{},             // Updated to use UUID
		&models.Budget{},           // Updated to use UUID
		&models.BudgetCategory{},   // Updated to use UUID
		&models.Income{},           // Updated to use UUID
		&models.Topic{},            // Updated to use UUID
		&models.Task_learning{},    // Updated to use UUID
		&models.Resource{},         // Updated to use UUID
		&models.StudySession{},     // Updated to use UUID
		&models.Tag{},              // Updated to use UUID
		&models.RepeatRule{},       // Updated to use UUID
		&models.AIRecommendation{}, // Updated to use UUID
	)

	// Re-enable foreign key constraints for SQLite
	db.Exec("PRAGMA foreign_keys = ON")

	return err
}

func Migrate002UUIDTransitionDown(db *gorm.DB) error {
	// Note: Down migration for UUID transition is complex and may not be fully reversible
	// In production, you would need to carefully plan the rollback strategy
	// For now, we'll leave this as nil since UUID transition is typically one-way
	return nil
}
