package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionPlan represents different subscription tiers
type SubscriptionPlan struct {
	ID           uuid.UUID      `json:"plan_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string         `json:"name" gorm:"not null;unique"`
	Description  string         `json:"description"`
	Price        float64        `json:"price" gorm:"not null"`           // Monthly price in USD
	Currency     string         `json:"currency" gorm:"default:'USD'"`   // Currency code
	Interval     string         `json:"interval" gorm:"default:'month'"` // month, year
	PayPalPlanID string         `json:"paypal_plan_id" gorm:"unique"`    // PayPal plan ID
	Features     []string       `json:"features" gorm:"type:text[]"`     // JSON array of features
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Subscription represents a user's active subscription
type Subscription struct {
	ID                   uuid.UUID        `json:"subscription_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID               uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"`
	User                 User             `json:"-" gorm:"foreignKey:UserID"`
	PlanID               uuid.UUID        `json:"plan_id" gorm:"type:uuid;not null"`
	Plan                 SubscriptionPlan `json:"plan" gorm:"foreignKey:PlanID"`
	PayPalSubscriptionID string           `json:"paypal_subscription_id" gorm:"unique"` // PayPal subscription ID
	Status               string           `json:"status" gorm:"not null"`               // active, cancelled, expired, suspended
	StartDate            time.Time        `json:"start_date" gorm:"not null"`
	EndDate              *time.Time       `json:"end_date"` // Null for ongoing subscriptions
	CancelledAt          *time.Time       `json:"cancelled_at"`
	AutoRenew            bool             `json:"auto_renew" gorm:"default:true"`
	CreatedAt            time.Time        `json:"-"`
	UpdatedAt            time.Time        `json:"-"`
	DeletedAt            gorm.DeletedAt   `json:"-" gorm:"index"`
}

// Payment represents individual payment transactions
type Payment struct {
	ID              uuid.UUID      `json:"payment_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SubscriptionID  uuid.UUID      `json:"subscription_id" gorm:"type:uuid;not null"`
	Subscription    Subscription   `json:"-" gorm:"foreignKey:SubscriptionID"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	User            User           `json:"-" gorm:"foreignKey:UserID"`
	PayPalPaymentID string         `json:"paypal_payment_id" gorm:"unique"` // PayPal payment/transaction ID
	Amount          float64        `json:"amount" gorm:"not null"`
	Currency        string         `json:"currency" gorm:"default:'USD'"`
	Status          string         `json:"status" gorm:"not null"` // completed, pending, failed, refunded
	PaymentDate     time.Time      `json:"payment_date" gorm:"not null"`
	Description     string         `json:"description"`
	CreatedAt       time.Time      `json:"-"`
	UpdatedAt       time.Time      `json:"-"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// PayPalWebhookEvent represents webhook events from PayPal
type PayPalWebhookEvent struct {
	ID          uuid.UUID      `json:"webhook_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EventType   string         `json:"event_type" gorm:"not null"`  // PayPal webhook event type
	ResourceID  string         `json:"resource_id" gorm:"not null"` // PayPal resource ID
	EventData   string         `json:"event_data" gorm:"type:text"` // Raw JSON event data
	Processed   bool           `json:"processed" gorm:"default:false"`
	ProcessedAt *time.Time     `json:"processed_at"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
