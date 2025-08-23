package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Budget struct {
	ID         uuid.UUID      `json:"budget_id" gorm:"primaryKey;type:text"`
	CategoryID uuid.UUID      `json:"category_id" gorm:"not null"`
	Category   BudgetCategory `gorm:"foreignKey:CategoryID"`
	Amount     float64        `json:"amount" gorm:"not null"`
	StartDate  time.Time      `json:"start_date" gorm:"not null"`
	EndDate    time.Time      `json:"end_date" gorm:"not null"`
	UserID     uuid.UUID      `json:"-"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	IncomeID   *uuid.UUID     `json:"income_id"` // optional: link budget to income
	Income     Income         `json:"-" gorm:"foreignKey:IncomeID"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BudgetCategory struct {
	ID        uuid.UUID      `json:"budget_category_id" gorm:"primaryKey;type:text"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uuid.UUID      `json:"-"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Budgets   []Budget       `json:"-" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (bc *BudgetCategory) BeforeCreate(tx *gorm.DB) error {
	if bc.ID == uuid.Nil {
		bc.ID = uuid.New()
	}
	return nil
}

type Income struct {
	ID         uuid.UUID      `json:"income_id" gorm:"primaryKey;type:text"`
	Source     string         `json:"source" gorm:"not null"`
	Amount     float64        `json:"amount" gorm:"not null"`
	UserID     uuid.UUID      `json:"-"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	Budgets    []Budget       `json:"budgets" gorm:"foreignKey:IncomeID"`
	ReceivedAt time.Time      `json:"received_at" gorm:"not null"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (i *Income) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type Transaction struct {
}
