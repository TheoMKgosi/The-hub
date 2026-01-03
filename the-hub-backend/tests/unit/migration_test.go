package unit

import (
	// "testing"

	//"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var migrationTestDB *gorm.DB

// setupMigrationTestDB initializes the test database for migration tests
func setupMigrationTestDB() {
	if migrationTestDB != nil {
		return
	}

	dsn := "host=" + getEnvOrDefault("DB_HOST", "localhost") +
		" user=" + getEnvOrDefault("DB_USER", "postgres") +
		" password=" + getEnvOrDefault("DB_PASSWORD", "postgres") +
		" dbname=" + getEnvOrDefault("DB_NAME", "the_hub_test") +
		" port=" + getEnvOrDefault("DB_PORT", "5432") +
		" sslmode=disable"

	var err error
	migrationTestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to connect to test database")
	}
}

// func TestGetMigrationVersion(t *testing.T) {
// 	setupMigrationTestDB()
//
// 	// Run migrations first
// 	err := migrations.RunLegacyMigrations(migrationTestDB)
// 	assert.NoError(t, err)
//
// 	// Get migration version (this will test the golang-migrate integration)
// 	version, dirty, err := migrations.GetMigrationVersion(migrationTestDB)
// 	// Since we're using legacy migrations, this might return an error
// 	// but the function should not panic
// 	if err == nil {
// 		assert.False(t, dirty, "Database should not be in dirty state")
// 		t.Logf("Migration version: %d", version)
// 	}
// }
