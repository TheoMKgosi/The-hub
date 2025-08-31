package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MigrationFunc func(*gorm.DB) error

type Migration struct {
	Version string
	Up      MigrationFunc
	Down    MigrationFunc
}

type MigrationRecord struct {
	Version   string    `gorm:"primaryKey"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

var migrations = []Migration{
	{
		Version: "001",
		Up:      Migrate001InitialSchema,
		Down:    nil, // Add down migration if needed
	},
}

func RunMigrations(db *gorm.DB) error {
	// Auto-migrate the migration record table
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return err
	}

	for _, migration := range migrations {
		// Check if migration is already applied
		var count int64
		db.Model(&MigrationRecord{}).Where("version = ?", migration.Version).Count(&count)

		if count == 0 {
			// Run migration
			if err := migration.Up(db); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.Version, err)
			}

			// Record migration
			if err := db.Create(&MigrationRecord{Version: migration.Version}).Error; err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
			}
		}
	}

	return nil
}
