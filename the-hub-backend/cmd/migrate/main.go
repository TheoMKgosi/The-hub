package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  up      - Run all pending migrations")
		fmt.Println("  down N  - Rollback N migrations")
		fmt.Println("  version - Show current migration version")
		fmt.Println("  status  - Show migration status")
		os.Exit(1)
	}

	command := os.Args[1]

	// Initialize database connection
	if err := config.InitDBManager("postgres"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := config.GetDB()

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	// Create postgres driver
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create postgres driver: %v", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch command {
	case "up":
		fmt.Println("Running migrations...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		if len(os.Args) < 3 {
			log.Fatal("Please specify number of migrations to rollback")
		}
		steps := 1 // default
		if len(os.Args) >= 3 {
			fmt.Sscanf(os.Args[2], "%d", &steps)
		}
		fmt.Printf("Rolling back %d migration(s)...\n", steps)
		if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Current version: %d\n", version)
		if dirty {
			fmt.Println("Warning: Database is in dirty state")
		}

	case "status":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Version: %d", version)
		if dirty {
			fmt.Printf(" (dirty)")
		}
		fmt.Println()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
