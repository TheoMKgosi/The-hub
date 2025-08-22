package models

import (
	"time"
)

type Topic struct {
	ID          uint   `json:"topic_id" gorm:"primaryKey"`
	UserID      uint   `json:"-"`
	User        User   `json:"-"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:not_started"` // or use enum logic
	// EstimatedHours int        `json:"estimated_hours"`
	Deadline  *time.Time `json:"deadline"`
	CreatedAt time.Time  `json:"-"`
	// Tasks          []Task_learning `json:"tasks"`
	Tags []Tag `json:"tags" gorm:"many2many:topic_tags"`
	// Reflections []Reflection
}

type Task_learning struct {
	ID         uint   `json:"task_learning_id" gorm:"primaryKey"`
	TopicID    uint   `json:"-"`
	Topic      Topic  `json:"-"`
	Title      string `json:"title" gorm:"not null"`
	Notes      string `json:"notes"`
	Status     string `json:"status" gorm:"default:pending"` // in_progress, done
	OrderIndex int    `json:"-"`
	// EstimatedTime int // in minutes
	CreatedAt time.Time `json:"-"`
	// Resources []Resource `json:"resources"`
	// Reflections []Reflection
}

type Resource struct {
	ID      uint `gorm:"primaryKey"`
	TopicID *uint
	TaskID  *uint

	Title string
	Link  string
	Type  string // video, article, doc
	Notes string
}

type StudySession struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	User        User
	TopicID     *uint
	TaskID      *uint
	DurationMin int
	StartedAt   time.Time
	EndedAt     time.Time
}

/*
type Reflection struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User
	TopicID   *uint
	TaskID    *uint
	Content   string
	CreatedAt time.Time
}
*/

type Tag struct {
	ID     uint   `json:"tag_id" gorm:"primaryKey"`
	UserID uint   `json:"user_id"`
	Name   string `json:"name" gorm:"unique;not null"`
	Color  string `json:"color"`
}
