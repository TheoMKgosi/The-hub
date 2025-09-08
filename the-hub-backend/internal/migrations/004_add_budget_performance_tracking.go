package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate004AddBudgetPerformanceTracking(db *gorm.DB) error {
	// Create budget_performance table
	err := db.AutoMigrate(&models.BudgetPerformance{})
	if err != nil {
		return err
	}

	// Create budget_alert_logs table
	err = db.AutoMigrate(&models.BudgetAlertLog{})
	if err != nil {
		return err
	}

	return nil
}
