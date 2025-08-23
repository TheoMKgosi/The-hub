package migrations

import (
	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

// Example: Adding a new table and migrating data from existing tables
func Migrate002AddUserSettings(db *gorm.DB) error {
	// 1. Create the new table
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS user_settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		theme VARCHAR(50) DEFAULT 'light',
		notifications_enabled BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`).Error; err != nil {
		return err
	}

	// 2. Migrate existing data from users.settings JSON field
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		// Extract theme from existing settings (if it exists)
		theme := "light" // default
		if user.Settings != nil {
			if t, exists := user.Settings["theme"]; exists {
				if themeStr, ok := t.(string); ok {
					theme = themeStr
				}
			}
		}

		// Insert into new table
		if err := db.Exec("INSERT INTO user_settings (user_id, theme) VALUES (?, ?)",
			user.ID, theme).Error; err != nil {
			return err
		}
	}

	return nil
}

// Example: Renaming a table and preserving data
func Migrate003RenameTable(db *gorm.DB) error {
	// For SQLite
	if err := db.Exec("ALTER TABLE old_table RENAME TO new_table").Error; err != nil {
		return err
	}

	// For PostgreSQL, you might need:
	// ALTER TABLE old_table RENAME TO new_table;

	return nil
}

// Example: Adding a column with default value
func Migrate004AddColumnWithData(db *gorm.DB) error {
	// Add column with default
	if err := db.Exec("ALTER TABLE users ADD COLUMN status VARCHAR(50) DEFAULT 'active'").Error; err != nil {
		return err
	}

	// Update existing data based on conditions
	if err := db.Model(&models.User{}).Where("created_at < ?", "2024-01-01").Update("status", "legacy").Error; err != nil {
		return err
	}

	return nil
}
