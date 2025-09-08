package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBManager struct {
	DB *gorm.DB
}

var dbManager *DBManager

func getPostgresDSN() string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)
}

func InitDBManager(dbType string) error {
	var dialector gorm.Dialector

	if dbType != "postgres" {
		return fmt.Errorf("only PostgreSQL is supported")
	}

	dialector = postgres.Open(getPostgresDSN())

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	dbManager = &DBManager{DB: db}
	return nil
}

func GetDB() *gorm.DB {
	return dbManager.DB
}

func GetDBManager() *DBManager {
	return dbManager
}

// SetTestDB sets the database manager for testing purposes
func SetTestDB(db *gorm.DB) {
	dbManager = &DBManager{DB: db}
}

func (dm *DBManager) HealthCheck(ctx context.Context) error {
	sqlDB, err := dm.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Legacy functions for backward compatibility
func InitDBPostgreSQL() {
	if err := InitDBManager("postgres"); err != nil {
		log.Fatal("Error initializing PostgreSQL database:", err)
	}
}

func GetDBPostgreSQL() *gorm.DB {
	return GetDB()
}
