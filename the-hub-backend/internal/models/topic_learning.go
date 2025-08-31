package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	ID          uuid.UUID `json:"topic_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"-" gorm:"type:uuid"`
	User        User      `json:"-"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:not_started"` // or use enum logic
	// EstimatedHours int        `json:"estimated_hours"`
	Deadline  *time.Time `json:"deadline"`
	CreatedAt time.Time  `json:"-"`
	// Tasks          []Task_learning `json:"tasks"`
	Tags []Tag `json:"tags" gorm:"many2many:topic_tags"`
	// Reflections []Reflection
}

// BeforeCreate hook to generate UUID
func (t *Topic) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type Task_learning struct {
	ID         uuid.UUID `json:"task_learning_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TopicID    uuid.UUID `json:"-" gorm:"type:uuid"`
	Topic      Topic     `json:"-"`
	Title      string    `json:"title" gorm:"not null"`
	Notes      string    `json:"notes"`
	Status     string    `json:"status" gorm:"default:pending"` // in_progress, done
	OrderIndex int       `json:"-"`
	// EstimatedTime int // in minutes
	CreatedAt time.Time `json:"-"`
	// Resources []Resource `json:"resources"`
	// Reflections []Reflection
}

// BeforeCreate hook to generate UUID
func (tl *Task_learning) BeforeCreate(tx *gorm.DB) error {
	if tl.ID == uuid.Nil {
		tl.ID = uuid.New()
	}
	return nil
}

type Resource struct {
	ID      uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TopicID *uuid.UUID `gorm:"type:uuid"`
	TaskID  *uuid.UUID `gorm:"type:uuid"`

	Title string
	Link  string
	Type  string // video, article, doc
	Notes string
}

// BeforeCreate hook to generate UUID
func (r *Resource) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type StudySession struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	User        User
	TopicID     *uuid.UUID `gorm:"type:uuid"`
	TaskID      *uuid.UUID `gorm:"type:uuid"`
	DurationMin int
	StartedAt   time.Time
	EndedAt     time.Time
}

// BeforeCreate hook to generate UUID
func (ss *StudySession) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	return nil
}

type Tag struct {
	ID     uuid.UUID `json:"tag_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Name   string    `json:"name" gorm:"unique;not null"`
	Color  string    `json:"color"`
}

// BeforeCreate hook to generate UUID
func (tg *Tag) BeforeCreate(tx *gorm.DB) error {
	if tg.ID == uuid.Nil {
		tg.ID = uuid.New()
	}
	return nil
}
