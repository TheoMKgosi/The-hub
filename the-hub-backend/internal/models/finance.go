package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Budget struct {
	ID         uuid.UUID      `json:"budget_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CategoryID uuid.UUID      `json:"category_id" gorm:"type:uuid;not null"`
	Category   BudgetCategory `gorm:"foreignKey:CategoryID"`
	Amount     float64        `json:"amount" gorm:"not null"`
	StartDate  time.Time      `json:"start_date" gorm:"not null"`
	EndDate    time.Time      `json:"end_date" gorm:"not null"`
	UserID     uuid.UUID      `json:"-" gorm:"type:uuid"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	IncomeID   *uuid.UUID     `json:"income_id" gorm:"type:uuid"` // optional: link budget to income
	Income     Income         `json:"-" gorm:"foreignKey:IncomeID"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type BudgetCategory struct {
	ID        uuid.UUID      `json:"budget_category_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uuid.UUID      `json:"-" gorm:"type:uuid"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Budgets   []Budget       `json:"-" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Income struct {
	ID         uuid.UUID      `json:"income_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Source     string         `json:"source" gorm:"not null"`
	Amount     float64        `json:"amount" gorm:"not null"`
	UserID     uuid.UUID      `json:"-" gorm:"type:uuid"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	Budgets    []Budget       `json:"budgets" gorm:"foreignKey:IncomeID"`
	ReceivedAt time.Time      `json:"received_at" gorm:"not null"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type Transaction struct {
	ID          uuid.UUID      `json:"transaction_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Description string         `json:"description" gorm:"not null"`
	Amount      float64        `json:"amount" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"` // "income" or "expense"
	Date        time.Time      `json:"date" gorm:"not null"`
	CategoryID  *uuid.UUID     `json:"category_id" gorm:"type:uuid"` // optional: link to budget category
	Category    BudgetCategory `json:"-" gorm:"foreignKey:CategoryID"`
	UserID      uuid.UUID      `json:"-" gorm:"type:uuid"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// BudgetPerformance tracks historical budget performance data
type BudgetPerformance struct {
	ID              uuid.UUID `json:"budget_performance_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BudgetID        uuid.UUID `json:"budget_id" gorm:"type:uuid;not null"`
	Budget          Budget    `json:"-" gorm:"foreignKey:BudgetID"`
	PeriodStart     time.Time `json:"period_start" gorm:"not null"`
	PeriodEnd       time.Time `json:"period_end" gorm:"not null"`
	BudgetAmount    float64   `json:"budget_amount" gorm:"not null"`
	SpentAmount     float64   `json:"spent_amount" gorm:"not null"`
	UtilizationRate float64   `json:"utilization_rate" gorm:"not null"`
	Status          string    `json:"status" gorm:"not null"` // "on_track", "warning", "over_budget"
	UserID          uuid.UUID `json:"-" gorm:"type:uuid"`
	User            User      `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

// BudgetAlertLog tracks budget alerts that have been sent
type BudgetAlertLog struct {
	ID        uuid.UUID `json:"alert_log_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BudgetID  uuid.UUID `json:"budget_id" gorm:"type:uuid;not null"`
	Budget    Budget    `json:"-" gorm:"foreignKey:BudgetID"`
	AlertType string    `json:"alert_type" gorm:"not null"` // "warning", "danger", "over_budget"
	Message   string    `json:"message" gorm:"not null"`
	Threshold float64   `json:"threshold" gorm:"not null"` // percentage threshold that triggered alert
	UserID    uuid.UUID `json:"-" gorm:"type:uuid"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	SentAt    time.Time `json:"sent_at" gorm:"not null"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
