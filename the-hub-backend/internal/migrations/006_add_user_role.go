package migrations

import (
	"gorm.io/gorm"
)

func Migrate006AddUserRole(db *gorm.DB) error {
	// Add role column to users table with default value 'user'
	return db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'user'").Error
}
