package models

import "time"

// Account represents a financial account owned by a user
type Account struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"type:uuid;not null;index"`
	Name         string  `gorm:"size:100;not null"`
	Balance      float64 `gorm:"not null;default:0"`
	User         User    `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Transactions []Transaction
}

// Transaction represents a financial transaction on an account
type Transaction struct {
	ID          uint    `gorm:"primaryKey"`
	AccountID   uint    `gorm:"not null;index"`
	Amount      float64 `gorm:"not null"`
	Type        string  `gorm:"size:50;not null"` // e.g., "credit" or "debit"
	Description string  `gorm:"size:255"`
	Account     Account `gorm:"foreignKey:AccountID"`
	CreatedAt   time.Time
}

// Payment represents a payment record
type Payment struct {
	ID          uint      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PaymentCode string    `gorm:"size:100;uniqueIndex;not null"`
	UserID      uint      `gorm:"type:uuid;not null;index"`
	Amount      float64   `gorm:"not null"`
	Status      string    `gorm:"size:50;not null"` // e.g., "pending", "completed"
	User        User      `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
