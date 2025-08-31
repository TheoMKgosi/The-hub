package tests

import (
	"fmt"
	"log"
	"os"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB holds the test database connection
var TestDB *gorm.DB

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestUser represents a test user for testing
type TestUser struct {
	ID       uuid.UUID
	Email    string
	Name     string
	Password string
}

// SetupTestDB initializes the test database
func SetupTestDB() {
	// Use PostgreSQL for testing
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_USER", "postgres"),
		getEnvOrDefault("DB_PASSWORD", "postgres"),
		getEnvOrDefault("DB_NAME", "the_hub_test"),
		getEnvOrDefault("DB_PORT", "5432"),
	)

	var err error
	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Suppress logs during tests
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the schema
	err = TestDB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}
}

// TeardownTestDB cleans up the test database
func TeardownTestDB() {
	if TestDB != nil {
		sqlDB, _ := TestDB.DB()
		sqlDB.Close()
	}
}

// CreateTestUser creates a test user in the database
func CreateTestUser(email, name, password string) (*models.User, error) {
	user := &models.User{
		Email:    email,
		Name:     name,
		Password: password,
		Settings: map[string]interface{}{"theme": "light"},
	}

	result := TestDB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// CleanTestData removes all test data
func CleanTestData() {
	TestDB.Unscoped().Where("1=1").Delete(&models.User{})
}

// SetupTestEnvironment sets up environment variables for testing
func SetupTestEnvironment() {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-tokens")
	os.Setenv("GIN_MODE", "test")
}
