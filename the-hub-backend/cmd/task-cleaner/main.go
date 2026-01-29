package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type TaskCleaner struct {
	db     *gorm.DB
	dryRun bool
}

func NewTaskCleaner(db *gorm.DB, dryRun bool) *TaskCleaner {
	return &TaskCleaner{db: db, dryRun: dryRun}
}

// OptimizeTaskIndexes rebuilds indexes for better performance
func (tc *TaskCleaner) OptimizeTaskIndexes() error {
	if tc.dryRun {
		log.Println("[DRY RUN] Would optimize database indexes and analyze tables")
		return nil
	}

	// PostgreSQL-specific index optimization
	queries := []string{
		"REINDEX INDEX CONCURRENTLY idx_tasks_user_id",
		"REINDEX INDEX CONCURRENTLY idx_tasks_status",
		"REINDEX INDEX CONCURRENTLY idx_tasks_deleted_at",
		"REINDEX INDEX CONCURRENTLY idx_tasks_parent_task_id",
		"REINDEX INDEX CONCURRENTLY idx_time_entries_task_id",
		"REINDEX INDEX CONCURRENTLY idx_task_dependencies_task_id",
		"REINDEX INDEX CONCURRENTLY idx_task_dependencies_depends_on_id",
		"REINDEX INDEX CONCURRENTLY idx_scheduled_tasks_task_id",
		"ANALYZE tasks, time_entries, task_dependencies, scheduled_tasks",
	}

	for _, query := range queries {
		if err := tc.db.Exec(query).Error; err != nil {
			log.Printf("Warning: Failed to execute optimization query '%s': %v", query, err)
			// Continue with other optimizations even if one fails
		}
	}

	log.Println("Completed database optimization")
	return nil
}

// CleanAllCompletedTasks removes all completed tasks immediately
func (tc *TaskCleaner) CleanAllCompletedTasks() error {
	if tc.dryRun {
		var count int64
		if err := tc.db.Model(&models.Task{}).Where("status = ? AND deleted_at IS NULL", "complete").Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count completed tasks for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would delete %d completed tasks", count)
		return nil
	}

	result := tc.db.Unscoped().Where("status = ? AND deleted_at IS NULL", "complete").Delete(&models.Task{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete completed tasks: %w", result.Error)
	}

	log.Printf("Deleted %d completed tasks", result.RowsAffected)
	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		dryRun                  = flag.Bool("dry-run", false, "Show what would be cleaned without actually cleaning")
		optimize                = flag.Bool("optimize", false, "Run database optimization after cleanup")
		cleanCompleted          = flag.Bool("clean-completed", false, "Delete all completed tasks immediately")
	)

	flag.Parse()

	// Load environment variables
	if err := config.InitDBManager(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	db := config.GetDB()
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	// Health check
	if err := config.GetDBManager().HealthCheck(context.Background()); err != nil {
		log.Fatal("Database health check failed:", err)
	}

	cleaner := NewTaskCleaner(db, *dryRun)

	if *dryRun {
		log.Println("DRY RUN MODE - No changes will be made")
	}

	log.Println("Starting task cleanup process...")

	if *cleanCompleted {
		if err := cleaner.CleanAllCompletedTasks(); err != nil {
			log.Printf("Error cleaning completed tasks: %v", err)
		}
		return // Exit after cleaning
	}

	// Optimize database if requested
	if *optimize {
		if err := cleaner.OptimizeTaskIndexes(); err != nil {
			log.Printf("Error optimizing database: %v", err)
		}
	}

	log.Println("Task cleanup process completed successfully")
}
