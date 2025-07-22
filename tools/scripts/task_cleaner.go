package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Define your Task model (adjust field names/types to match yours)
type Task struct {
	ID          uint           `json:"task_id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	DueDate     *time.Time     `json:"due_date"`
	Priority    *int           `json:"priority" gorm:"check:priority >= 1 AND priority <= 5"`
	Status      string         `json:"status" gorm:"default:pending"`
	GoalID      *uint          `json:"goal_id"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	Goal        Goal           `json:"-" gorm:"foreignKey:GoalID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Goal struct {
	ID          uint           `json:"goal_id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Tasks       []Task         `json:"tasks"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type User struct {
	ID        uint           `json:"task_id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func clearCompletedTasks(db *gorm.DB) {
	result := db.Where("status = ?", "complete").Delete(&Task{})
	if result.Error != nil {
		log.Printf("Failed to delete completed tasks: %v", result.Error)
	} else {
		log.Printf("Deleted %d completed tasks", result.RowsAffected)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Connect to DB
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	clearCompletedTasks(db)
}
