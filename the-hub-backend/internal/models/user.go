package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID              `json:"user_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email" gorm:"unique"`
	Password  string                 `json:"-"`
	Settings  map[string]interface{} `json:"settings" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"-"`
	UpdatedAt time.Time              `json:"-"`
	DeletedAt gorm.DeletedAt         `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
