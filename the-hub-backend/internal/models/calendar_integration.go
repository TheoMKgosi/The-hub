package models

import (
	"time"

	"github.com/google/uuid"
)

// CalendarProvider represents different calendar providers
type CalendarProvider string

const (
	ProviderGoogle  CalendarProvider = "google"
	ProviderOutlook CalendarProvider = "outlook"
	ProviderApple   CalendarProvider = "apple"
)

// CalendarIntegration represents a user's calendar integration
type CalendarIntegration struct {
	ID             uuid.UUID        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"`
	User           User             `json:"-" gorm:"foreignKey:UserID"`
	Provider       CalendarProvider `json:"provider" gorm:"not null"`
	ProviderUserID string           `json:"provider_user_id" gorm:"not null"` // User's ID in the provider system
	AccessToken    string           `json:"-" gorm:"not null"`                // Encrypted access token
	RefreshToken   string           `json:"-" gorm:"not null"`                // Encrypted refresh token
	TokenExpiry    time.Time        `json:"-" gorm:"not null"`
	CalendarID     string           `json:"calendar_id" gorm:"not null"` // Primary calendar ID to sync with
	IsActive       bool             `json:"is_active" gorm:"default:true"`
	LastSyncAt     *time.Time       `json:"last_sync_at"`
	SyncEnabled    bool             `json:"sync_enabled" gorm:"default:true"`
	CreatedAt      time.Time        `json:"-"`
	UpdatedAt      time.Time        `json:"-"`
}

// CalendarEventSync represents the mapping between local scheduled tasks and external calendar events
type CalendarEventSync struct {
	ID                    uuid.UUID           `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ScheduledTaskID       uuid.UUID           `json:"scheduled_task_id" gorm:"type:uuid;not null"`
	ScheduledTask         ScheduledTask       `json:"-" gorm:"foreignKey:ScheduledTaskID"`
	CalendarIntegrationID uuid.UUID           `json:"calendar_integration_id" gorm:"type:uuid;not null"`
	CalendarIntegration   CalendarIntegration `json:"-" gorm:"foreignKey:CalendarIntegrationID"`
	ExternalEventID       string              `json:"external_event_id" gorm:"not null"` // Event ID in external calendar
	LastSyncedAt          time.Time           `json:"last_synced_at"`
	CreatedAt             time.Time           `json:"-"`
	UpdatedAt             time.Time           `json:"-"`
}
