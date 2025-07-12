package config

import (
	"log"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbLite *gorm.DB

func InitDBSQLite() {

	dbOpen, err := gorm.Open(sqlite.Open("the-hub.db"), &gorm.Config{})

	dbLite = dbOpen

	if err != nil {
		log.Fatal("Error opening database")
	}

	dbLite.AutoMigrate(&models.Goal{}, &models.Task{}, &models.ScheduledTask{}, &models.User{}, &models.Deck{}, &models.Card{})

}

func GetDB() *gorm.DB {
	return dbLite
}
