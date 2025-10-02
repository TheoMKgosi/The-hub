package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate006AddSubscriptionSystem(db *gorm.DB) error {
	// Create subscription-related tables
	err := db.AutoMigrate(
		&models.SubscriptionPlan{},
		&models.Subscription{},
		&models.Payment{},
		&models.PayPalWebhookEvent{},
	)
	if err != nil {
		return err
	}

	return nil
}
