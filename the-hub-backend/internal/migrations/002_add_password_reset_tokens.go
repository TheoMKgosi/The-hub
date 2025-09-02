package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate002AddPasswordResetTokens(db *gorm.DB) error {
	return db.AutoMigrate(&models.PasswordResetToken{})
}
