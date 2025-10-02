package migrations

import (
	"encoding/json"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/gorm"
)

// getDefaultUserSettings returns the default settings for new users
func getDefaultUserSettings() map[string]interface{} {
	return map[string]interface{}{
		"theme": map[string]interface{}{
			"mode": "system", // light, dark, system
		},
		"notifications": map[string]interface{}{
			"email": map[string]interface{}{
				"enabled":        true,
				"task_reminders": true,
				"goal_deadlines": true,
				"weekly_reports": true,
			},
			"push": map[string]interface{}{
				"enabled":        false,
				"task_reminders": false,
				"goal_deadlines": false,
			},
		},
		"dashboard": map[string]interface{}{
			"layout": "default",
			"widgets": map[string]interface{}{
				"tasks":    true,
				"goals":    true,
				"schedule": true,
				"stats":    true,
			},
		},
		"privacy": map[string]interface{}{
			"profile_visibility": "private",
			"activity_sharing":   false,
		},
		"preferences": map[string]interface{}{
			"language":    "en",
			"timezone":    "UTC",
			"date_format": "MM/DD/YYYY",
			"time_format": "12h",
		},
	}
}

func Migrate002AddDefaultSettings(db *gorm.DB) error {
	// Get all users who don't have settings
	var users []models.User
	if err := db.Where("settings IS NULL OR settings = '{}'").Find(&users).Error; err != nil {
		return err
	}

	// Apply default settings to each user
	defaultSettings := getDefaultUserSettings()
	for _, user := range users {
		settingsJSON, err := json.Marshal(defaultSettings)
		if err != nil {
			return err
		}
		if err := db.Model(&user).Update("settings", string(settingsJSON)).Error; err != nil {
			return err
		}
	}

	return nil
}
