package tests

import (
	"log"
	"os"

	"github.com/TheoMKgosi/The-hub/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB holds the test database connection
var TestDB *gorm.DB

// TestUser represents a test user for testing
type TestUser struct {
	ID       uint
	Email    string
	Name     string
	Password string
}

// SetupTestDB initializes the test database
func SetupTestDB() {
	// Use in-memory SQLite for testing
	dsn := ":memory:"

	var err error
	TestDB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
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
