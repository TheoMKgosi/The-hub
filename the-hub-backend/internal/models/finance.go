package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID         uint           `json:"budget_id" gorm:"primaryKey"`
	CategoryID uint           `json:"category_id" gorm:"not null"`
	Category   BudgetCategory `gorm:"foreignKey:CategoryID"`
	Amount     float64        `json:"amount" gorm:"not null"`
	StartDate  time.Time      `json:"start_date" gorm:"not null"`
	EndDate    time.Time      `json:"end_date" gorm:"not null"`
	UserID     uint           `json:"-"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	IncomeID   *uint          `json:"income_id"` // optional: link budget to income
	Income     Income         `json:"-" gorm:"foreignKey:IncomeID"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type BudgetCategory struct {
	ID        uint           `json:"budget_category_id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uint           `json:"-"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Budgets   []Budget       `json:"-" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Income struct {
	ID         uint           `json:"income_id" gorm:"primaryKey"`
	Source     string         `json:"source" gorm:"not null"`
	Amount     float64        `json:"amount" gorm:"not null"`
	UserID     uint           `json:"-"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	Budgets    []Budget       `json:"budgets" gorm:"foreignKey:IncomeID"`
	ReceivedAt time.Time      `json:"received_at" gorm:"not null"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
