package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	ID            uuid.UUID      `json:"feedback_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	User          User           `json:"-" gorm:"foreignKey:UserID"`
	Type          string         `json:"type" gorm:"not null;check:type IN ('bug', 'feature', 'improvement', 'general')"`
	Subject       string         `json:"subject" gorm:"not null"`
	Description   string         `json:"description" gorm:"not null"`
	Rating        *int           `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
	PageURL       string         `json:"page_url"`
	UserAgent     string         `json:"user_agent"`
	Status        string         `json:"status" gorm:"default:'pending';check:status IN ('pending', 'reviewed', 'implemented', 'declined')"`
	AdminResponse string         `json:"admin_response"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for the Feedback model
func (Feedback) TableName() string {
	return "feedback"
}
