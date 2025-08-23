package config

import (
	"log"

	"gorm.io/gorm"
)

// Legacy function for backward compatibility
func InitDBSQLite() {
	if err := InitDBManager("sqlite"); err != nil {
		log.Fatal("Error initializing SQLite database:", err)
	}
}

// SetTestDB sets the database connection for testing
func SetTestDB(testDB *gorm.DB) {
	if dbManager != nil {
		dbManager.DB = testDB
	}
}
