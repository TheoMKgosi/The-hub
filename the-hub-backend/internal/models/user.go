package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint                   `json:"user_id" gorm:"primaryKey"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email" gorm:"unique"`
	Password  string                 `json:"-"`
	Settings  map[string]interface{} `json:"settings" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"-"`
	UpdatedAt time.Time              `json:"-"`
	DeletedAt gorm.DeletedAt         `json:"-" gorm:"index"`
}
