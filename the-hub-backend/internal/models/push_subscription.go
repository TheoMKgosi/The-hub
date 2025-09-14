package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushSubscription struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Endpoint  string         `json:"endpoint" gorm:"unique;not null"`
	P256dh    string         `json:"p256dh" gorm:"not null"`
	Auth      string         `json:"auth" gorm:"not null"`
	UserAgent string         `json:"user_agent"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
