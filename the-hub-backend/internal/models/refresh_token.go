package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null;index"`
	TokenHash string    `json:"-" gorm:"unique;not null"` // Hashed refresh token (not exposed in JSON)
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	Revoked   bool      `json:"revoked" gorm:"default:false;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
}


// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.Revoked && !rt.IsExpired()
}
