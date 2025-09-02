package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

func Migrate003AddRefreshTokens(db *gorm.DB) error {
	// Create refresh tokens table
	err := db.AutoMigrate(
		&models.RefreshToken{},
		&models.PasswordResetToken{},
	)
	if err != nil {
		return err
	}

	// Add foreign key constraint for refresh tokens
	return db.Exec(`
		ALTER TABLE refresh_tokens
		ADD CONSTRAINT fk_refresh_tokens_user
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	`).Error
}
