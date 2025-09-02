package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Goal struct {
	ID          uuid.UUID      `json:"goal_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Tasks       []Task         `json:"tasks"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
