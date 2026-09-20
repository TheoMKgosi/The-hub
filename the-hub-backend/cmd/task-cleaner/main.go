package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/getsentry/sentry-go"
	// "github.com/joho/godotenv"
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
		if err := tc.db.Model(&models.Task{}).Unscoped().Where("status IN ? AND updated_at < ? AND deleted_at IS NULL", []string{"completed", "complete"}, cutoffDate).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count completed tasks for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would clean %d completed tasks older than %d days", count, retentionDays)
		return nil
	}

	result := tc.db.Unscoped().Where("status IN ? AND updated_at < ? AND deleted_at IS NULL", []string{"completed", "complete"}, cutoffDate).Delete(&models.Task{})
	if result.Error != nil {
		return fmt.Errorf("failed to clean completed tasks: %w", result.Error)
	}

	log.Printf("Cleaned %d completed tasks older than %d days", result.RowsAffected, retentionDays)
	return nil
}

// CleanOrphanedTaskDependencies removes task dependencies for non-existent tasks
func (tc *TaskCleaner) CleanOrphanedTaskDependencies() error {
	if tc.dryRun {
		var count int64
		if err := tc.db.Raw(`
			SELECT COUNT(*) FROM task_dependencies
			WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
			   OR dependency_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
		`).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to count orphaned task dependencies for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would clean %d orphaned task dependencies", count)
		return nil
	}

	result := tc.db.Exec(`
		DELETE FROM task_dependencies
		WHERE task_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
		   OR dependency_id NOT IN (SELECT id FROM tasks WHERE deleted_at IS NULL)
	`)
	if result.Error != nil {
		return fmt.Errorf("failed to clean orphaned task dependencies: %w", result.Error)
	}

	log.Printf("Cleaned %d orphaned task dependencies", result.RowsAffected)
	return nil
}

// UpdateParentTaskStatuses updates parent task statuses based on subtasks
func (tc *TaskCleaner) UpdateParentTaskStatuses() error {
	if tc.dryRun {
		var count int64
		if err := tc.db.Raw(`
			SELECT COUNT(*) FROM tasks t
			WHERE t.parent_task_id IS NULL
			  AND t.status != 'completed'
			  AND t.deleted_at IS NULL
			  AND NOT EXISTS (
				SELECT 1 FROM tasks st
				WHERE st.parent_task_id = t.id
				  AND st.status != 'completed'
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
		SET status = 'completed', updated_at = NOW()
		WHERE id IN (
			SELECT DISTINCT t.id
			FROM tasks t
			WHERE t.parent_task_id IS NULL
			  AND t.status != 'completed'
			  AND t.deleted_at IS NULL
			  AND NOT EXISTS (
				SELECT 1 FROM tasks st
				WHERE st.parent_task_id = t.id
				  AND st.status != 'completed'
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
		var taskCount int64

		if err := tc.db.Model(&models.Task{}).Unscoped().Where("deleted_at < ?", cutoffDate).Count(&taskCount).Error; err != nil {
			return fmt.Errorf("failed to count expired soft-deleted tasks for dry run: %w", err)
		}

		log.Printf("[DRY RUN] Would clean %d expired soft-deleted tasks",
			taskCount)
		return nil
	}

	// Clean soft-deleted tasks
	taskResult := tc.db.Unscoped().Where("deleted_at < ?", cutoffDate).Delete(&models.Task{})
	if taskResult.Error != nil {
		return fmt.Errorf("failed to clean expired soft-deleted tasks: %w", taskResult.Error)
	}

	log.Printf("Cleaned %d expired soft-deleted tasks",
		taskResult.RowsAffected)
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
		if err := tc.db.Model(&models.Task{}).Where("status IN ? AND deleted_at IS NULL", []string{"completed", "complete"}).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count completed tasks for dry run: %w", err)
		}
		log.Printf("[DRY RUN] Would delete %d completed tasks", count)
		return nil
	}

	result := tc.db.Unscoped().Where("status IN ? AND deleted_at IS NULL", []string{"completed", "complete"}).Delete(&models.Task{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete completed tasks: %w", result.Error)
	}

	log.Printf("Deleted %d completed tasks", result.RowsAffected)
	return nil
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	var (
		completedRetentionDays  = flag.Int("completed-retention", 90, "Days to retain completed tasks")
		softDeleteRetentionDays = flag.Int("soft-delete-retention", 30, "Days to retain soft-deleted records")
		dryRun                  = flag.Bool("dry-run", false, "Show what would be cleaned without actually cleaning")
		optimize                = flag.Bool("optimize", false, "Run database optimization after cleanup")
		cleanCompleted          = flag.Bool("clean-completed", false, "Delete all completed tasks immediately")
	)

	flag.Parse()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:             "https://90c441f99b2cf5c0b23b00665dc6313b@o4509804910936064.ingest.de.sentry.io/4512038764871760",
		EnableTracing:   true,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Printf("Sentry initialization failed: %v\n", err)
		return
	}
	defer sentry.Flush(2 * time.Second)

	transaction := sentry.StartTransaction(context.Background(), "task-cleaner", sentry.WithOpName("task"))
	defer transaction.Finish()
	ctx := transaction.Context()
	log.SetOutput(io.MultiWriter(os.Stderr, sentry.NewLogger(ctx)))

	const monitorSlug = "task-cleaner"
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("monitor", sentry.Context{"slug": monitorSlug})
	})
	monitorConfig := &sentry.MonitorConfig{
		Schedule: sentry.CrontabSchedule("0 2 * * *"),
	}
	checkInID := sentry.CaptureCheckIn(
		&sentry.CheckIn{
			MonitorSlug: monitorSlug,
			Status:      sentry.CheckInStatusInProgress,
		},
		monitorConfig,
	)
	checkInStatus := sentry.CheckInStatusError
	defer func() {
		checkIn := &sentry.CheckIn{
			MonitorSlug: monitorSlug,
			Status:      checkInStatus,
		}
		if checkInID != nil {
			checkIn.ID = *checkInID
		}
		sentry.CaptureCheckIn(checkIn, monitorConfig)
	}()

	hadErrors := false
	captureError := func(message string, err error) {
		hadErrors = true
		log.Printf("%s: %v", message, err)
		sentry.CaptureException(err)
	}

	// Load environment variables
	if err := config.InitDBManager(); err != nil {
		captureError("Failed to initialize database", err)
		return
	}

	db := config.GetDB()
	if db == nil {
		captureError("Failed to initialize database", fmt.Errorf("database connection is nil"))
		return
	}

	// Health check
	if err := config.GetDBManager().HealthCheck(ctx); err != nil {
		captureError("Database health check failed", err)
		return
	}

	cleaner := NewTaskCleaner(db, *dryRun)

	if *dryRun {
		log.Println("DRY RUN MODE - No changes will be made")
	}

	log.Println("Starting task cleanup process...")

	if *cleanCompleted {
		if err := cleaner.CleanAllCompletedTasks(); err != nil {
			captureError("Error cleaning completed tasks", err)
		} else {
			checkInStatus = sentry.CheckInStatusOK
			log.Println("Task cleanup process completed successfully")
		}
		return // Exit after cleaning
	}

	// Clean completed tasks
	if err := cleaner.CleanCompletedTasks(*completedRetentionDays); err != nil {
		captureError("Error cleaning completed tasks", err)
	}

	if err := cleaner.CleanOrphanedTaskDependencies(); err != nil {
		captureError("Error cleaning orphaned task dependencies", err)
	}

	// Update parent task statuses
	if err := cleaner.UpdateParentTaskStatuses(); err != nil {
		captureError("Error updating parent task statuses", err)
	}

	// Clean expired soft deletes
	if err := cleaner.CleanExpiredSoftDeletes(*softDeleteRetentionDays); err != nil {
		captureError("Error cleaning expired soft deletes", err)
	}

	// Optimize database if requested
	if *optimize {
		if err := cleaner.OptimizeTaskIndexes(); err != nil {
			captureError("Error optimizing database", err)
		}
	}

	if hadErrors {
		log.Println("Task cleanup process completed with errors")
		return
	}

	checkInStatus = sentry.CheckInStatusOK
	log.Println("Task cleanup process completed successfully")
}
