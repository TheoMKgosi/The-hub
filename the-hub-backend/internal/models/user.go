package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"task_id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
