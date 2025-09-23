package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

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

// CleanCompletedTasks removes old completed tasks based on retention policy
func (tc *TaskCleaner) CleanCompletedTasks(retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	if tc.dryRun {
		var count int64
		if err := tc.db.Model(&models.Task{}).Unscoped().Where("status = ? AND updated_at < ? AND deleted_at IS NULL", "complete", cutoffDate).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count completed tasks for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would clean %d completed tasks older than %d days", count, retentionDays)
		return nil
	}

	result := tc.db.Unscoped().Where("status = ? AND updated_at < ? AND deleted_at IS NULL", "complete", cutoffDate).Delete(&models.Task{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean completed tasks: %w", result.Error)
	}

	log.Printf("Cleaned %d completed tasks older than %d days", result.RowsAffected, retentionDays)
	return nil
}

// CleanOrphanedTimeEntries removes time entries for non-existent tasks
func (tc *TaskCleaner) CleanOrphanedTimeEntries() error {
	if tc.dryRun {
		var count int64
		if err := tc.db.Raw(`
			SELECT COUNT(*) FROM time_entries
			WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
		`).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to count orphaned time entries for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would clean %d orphaned time entries", count)
		return nil
	}

	result := tc.db.Exec(`
		DELETE FROM time_entries
		WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
	`)
	if result.Error != nil {
		return fmt.Errorf("failed to clean orphaned time entries: %w", result.Error)
	}

	log.Printf("Cleaned %d orphaned time entries", result.RowsAffected)
	return nil
}

// CleanOrphanedTaskDependencies removes task dependencies for non-existent tasks
func (tc *TaskCleaner) CleanOrphanedTaskDependencies() error {
	if tc.dryRun {
		var count int64
		if err := tc.db.Raw(`
			SELECT COUNT(*) FROM task_dependencies
			WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
			   OR depends_on_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
		`).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to count orphaned task dependencies for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would clean %d orphaned task dependencies", count)
		return nil
	}

	result := tc.db.Exec(`
		DELETE FROM task_dependencies
		WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
		   OR depends_on_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
	`)
	if result.Error != nil {
		return fmt.Errorf("failed to clean orphaned task dependencies: %w", result.Error)
	}

	log.Printf("Cleaned %d orphaned task dependencies", result.RowsAffected)
	return nil
}

// CleanOrphanedScheduledTasks removes scheduled tasks for non-existent tasks
func (tc *TaskCleaner) CleanOrphanedScheduledTasks() error {
	if tc.dryRun {
		var count int64
		if err := tc.db.Raw(`
			SELECT COUNT(*) FROM scheduled_tasks
			WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
		`).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to count orphaned scheduled tasks for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would clean %d orphaned scheduled tasks", count)
		return nil
	}

	result := tc.db.Exec(`
		DELETE FROM scheduled_tasks
		WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
	`)
	if result.Error != nil {
		return fmt.Errorf("failed to clean orphaned scheduled tasks: %w", result.Error)
	}

	log.Printf("Cleaned %d orphaned scheduled tasks", result.RowsAffected)
	return nil
}

// UpdateParentTaskStatuses updates parent task statuses based on subtasks
func (tc *TaskCleaner) UpdateParentTaskStatuses() error {
	if tc.dryRun {
		var count int64
		if err := tc.db.Raw(`
			SELECT COUNT(*) FROM tasks t
			WHERE t.parent_task_id IS NULL
			  AND t.status != 'complete'
			  AND t.deleted_at IS NULL
			  AND NOT EXISTS (
				SELECT 1 FROM tasks st
				WHERE st.parent_task_id = t.id
				  AND st.status != 'complete'
				  AND st.deleted_at IS NULL
			  )
		`).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to count parent tasks to update for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would update %d parent task statuses", count)
		return nil
	}

	// Update parent tasks that should be completed (all subtasks completed)
	result := tc.db.Exec(`
		UPDATE tasks
		SET status = 'complete', updated_at = NOW()
		WHERE id IN (
			SELECT DISTINCT t.id
			FROM tasks t
			WHERE t.parent_task_id IS NULL
			  AND t.status != 'complete'
			  AND t.deleted_at IS NULL
			  AND NOT EXISTS (
				SELECT 1 FROM tasks st
				WHERE st.parent_task_id = t.id
				  AND st.status != 'complete'
				  AND st.deleted_at IS NULL
			  )
		)
	`)
	if result.Error != nil {
		return fmt.Errorf("failed to update parent task statuses: %w", result.Error)
	}

	log.Printf("Updated %d parent task statuses", result.RowsAffected)
	return nil
}

// CleanExpiredSoftDeletes permanently removes soft-deleted records older than retention period
func (tc *TaskCleaner) CleanExpiredSoftDeletes(retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	if tc.dryRun {
		var taskCount, timeEntryCount, dependencyCount int64

		if err := tc.db.Model(&models.Task{}).Unscoped().Where("deleted_at < ?", cutoffDate).Count(&taskCount).Error; err != nil {
			return fmt.Errorf("failed to count expired soft-deleted tasks for dry run: %w", err)
		}

		if err := tc.db.Model(&models.TimeEntry{}).Unscoped().Where("deleted_at < ?", cutoffDate).Count(&timeEntryCount).Error; err != nil {
			return fmt.Errorf("failed to count expired soft-deleted time entries for dry run: %w", err)
		}

		if err := tc.db.Model(&models.TaskDependency{}).Unscoped().Where("deleted_at < ?", cutoffDate).Count(&dependencyCount).Error; err != nil {
			return fmt.Errorf("failed to count expired soft-deleted task dependencies for dry run: %w", err)
		}

		log.Printf("[DRY RUN] Would clean %d expired soft-deleted tasks, %d time entries, %d dependencies",
			taskCount, timeEntryCount, dependencyCount)
		return nil
	}

	// Clean soft-deleted tasks
	taskResult := tc.db.Unscoped().Where("deleted_at < ?", cutoffDate).Delete(&models.Task{})
	if taskResult.Error != nil {
		return fmt.Errorf("failed to clean expired soft-deleted tasks: %w", taskResult.Error)
	}

	// Clean soft-deleted time entries
	timeEntryResult := tc.db.Unscoped().Where("deleted_at < ?", cutoffDate).Delete(&models.TimeEntry{})
	if timeEntryResult.Error != nil {
		return fmt.Errorf("failed to clean expired soft-deleted time entries: %w", timeEntryResult.Error)
	}

	// Clean soft-deleted task dependencies
	dependencyResult := tc.db.Unscoped().Where("deleted_at < ?", cutoffDate).Delete(&models.TaskDependency{})
	if dependencyResult.Error != nil {
		return fmt.Errorf("failed to clean expired soft-deleted task dependencies: %w", dependencyResult.Error)
	}

	log.Printf("Cleaned %d expired soft-deleted tasks, %d time entries, %d dependencies",
		taskResult.RowsAffected, timeEntryResult.RowsAffected, dependencyResult.RowsAffected)
	return nil
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
		completedRetentionDays  = flag.Int("completed-retention", 90, "Days to retain completed tasks")
		softDeleteRetentionDays = flag.Int("soft-delete-retention", 30, "Days to retain soft-deleted records")
		dryRun                  = flag.Bool("dry-run", false, "Show what would be cleaned without actually cleaning")
		optimize                = flag.Bool("optimize", false, "Run database optimization after cleanup")
		cleanCompleted          = flag.Bool("clean-completed", false, "Delete all completed tasks immediately")
	)

	flag.Parse()

	// Load environment variables
	if err := config.InitDBManager("postgres"); err != nil {
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

	// Clean completed tasks
	if err := cleaner.CleanCompletedTasks(*completedRetentionDays); err != nil {
		log.Printf("Error cleaning completed tasks: %v", err)
	}

	// Clean orphaned records
	if err := cleaner.CleanOrphanedTimeEntries(); err != nil {
		log.Printf("Error cleaning orphaned time entries: %v", err)
	}

	if err := cleaner.CleanOrphanedTaskDependencies(); err != nil {
		log.Printf("Error cleaning orphaned task dependencies: %v", err)
	}

	if err := cleaner.CleanOrphanedScheduledTasks(); err != nil {
		log.Printf("Error cleaning orphaned scheduled tasks: %v", err)
	}

	// Update parent task statuses
	if err := cleaner.UpdateParentTaskStatuses(); err != nil {
		log.Printf("Error updating parent task statuses: %v", err)
	}

	// Clean expired soft deletes
	if err := cleaner.CleanExpiredSoftDeletes(*softDeleteRetentionDays); err != nil {
		log.Printf("Error cleaning expired soft deletes: %v", err)
	}

	// Optimize database if requested
	if *optimize {
		if err := cleaner.OptimizeTaskIndexes(); err != nil {
			log.Printf("Error optimizing database: %v", err)
		}
	}

	log.Println("Task cleanup process completed successfully")
}
