package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSettings is a custom type that handles JSONB scanning for user settings
type UserSettings map[string]interface{}

// Scan implements the Scanner interface for database deserialization
func (us *UserSettings) Scan(value interface{}) error {
	if value == nil {
		*us = make(UserSettings)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		*us = make(UserSettings)
		return nil
	}

	// Try to unmarshal as JSON object first
	var m map[string]interface{}
	if err := json.Unmarshal(bytes, &m); err != nil {
		// If that fails, try to unmarshal as string containing JSON
		var s string
		if err2 := json.Unmarshal(bytes, &s); err2 != nil {
			// If both fail, initialize as empty map
			*us = make(UserSettings)
			return nil
		}
		// Unmarshal the string content as JSON
		if err3 := json.Unmarshal([]byte(s), &m); err3 != nil {
			*us = make(UserSettings)
			return nil
		}
	}

	*us = UserSettings(m)
	return nil
}

// Value implements the Valuer interface for database serialization
func (us UserSettings) Value() (driver.Value, error) {
	if us == nil {
		return "{}", nil
	}
	return json.Marshal(us)
}

type User struct {
	ID        uuid.UUID      `json:"user_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"`
	Settings  UserSettings   `json:"settings" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
