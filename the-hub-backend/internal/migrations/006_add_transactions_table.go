package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate006AddTransactionsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Transaction{})
}

func Migrate006AddTransactionsTableDown(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Transaction{})
}
