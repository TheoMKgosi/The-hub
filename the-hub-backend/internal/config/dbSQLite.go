package config

import (
	"log"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbLite *gorm.DB

func InitDBSQLite() {

	dsn := "file:the-hub.db?_foreign_keys=on"
	dbOpen, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	dbLite = dbOpen

	if err != nil {
		log.Fatal("Error opening database")
	}

	dbLite.AutoMigrate(&models.Goal{}, &models.Task{}, &models.ScheduledTask{}, &models.User{}, &models.Deck{}, &models.Card{}, &models.Budget{}, &models.BudgetCategory{})

}

func GetDB() *gorm.DB {
	return dbLite
}
