package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetToken struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"unique;not null;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	Used      bool      `json:"used" gorm:"default:false"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
}

// IsExpired checks if the token has expired
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// IsValid checks if the token is valid (not used and not expired)
func (p *PasswordResetToken) IsValid() bool {
	return !p.Used && !p.IsExpired()
}
