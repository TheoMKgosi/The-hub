package config

import (
	"log"
	"os"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbOpen, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host= " + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort,
	}), &gorm.Config{})

	db = dbOpen

	if err != nil {
		log.Fatal("Error opening database")
	}

	db.AutoMigrate(&models.Goal{}, &models.Task{})

}

func GetDB()  *gorm.DB {
	return db
}
