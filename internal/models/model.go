package models

import (
	"time"
)

type Goal struct {
	ID          uint      `json:"goal_id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tasks       []Task    `json:"tasks"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
}

type Task struct {
	ID          uint      `json:"task_id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GoalID      *uint     `json:"goal_id" gorm:"foreignKey"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
}
