package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations executes all pending database migrations
func RunMigrations(db *gorm.DB) error {
	// For now, use legacy migrations until we fully convert to SQL migrations
	return RunLegacyMigrations(db)
}

// RunLegacyMigrations executes the existing Go-based migrations
// This is a temporary function to support the transition
func RunLegacyMigrations(db *gorm.DB) error {
	migrations := []func(*gorm.DB) error{
		Migrate001InitialSchema,
		Migrate002AddDefaultSettings,
		Migrate003AddCalendarIntegration,
		Migrate004AddBudgetPerformanceTracking,
		Migrate005UpdateCalendarZoneDefaults,
		Migrate006EnhanceSchedulingSystem,
		Migrate007AddSubscriptionSystem,
		Migrate008AddUserRole,
	}

	for i, migration := range migrations {
		if err := migration(db); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(db *gorm.DB) (uint, bool, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return 0, false, err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	return version, dirty, err
}
