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
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type BudgetCategory struct {
	ID        uint           `json:"budget_category_id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uint           `json:"-"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Budgets   []Budget       `json:"budgets" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
