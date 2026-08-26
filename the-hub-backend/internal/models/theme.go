package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Theme represents a group of tasks shown together on the days assigned to it.
// With the one-theme-per-day rule, a weekday may only belong to a single theme.
type Theme struct {
	ID          uuid.UUID      `json:"theme_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Color       string         `json:"color" gorm:"default:'#3B82F6'"`
	Icon        string         `json:"icon"`
	Days        pq.StringArray `json:"days" gorm:"type:text[]"` // lowercase weekday names: monday..sunday
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateThemeRequest represents the request body for creating a theme.
type CreateThemeRequest struct {
	Name        string   `json:"name" binding:"required" example:"Deep Work"`
	Description string   `json:"description" example:"Focused development and writing work"`
	Color       string   `json:"color" example:"#3B82F6"`
	Icon        string   `json:"icon" example:"i-lucide-briefcase"`
	Days        []string `json:"days" example:"monday,tuesday"`
}

// UpdateThemeRequest represents the request body for updating a theme.
type UpdateThemeRequest struct {
	Name        *string  `json:"name" example:"Deep Work"`
	Description *string  `json:"description"`
	Color       *string  `json:"color"`
	Icon        *string  `json:"icon"`
	Days        []string `json:"days"`
}
