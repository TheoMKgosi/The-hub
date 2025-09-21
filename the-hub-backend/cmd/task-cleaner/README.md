# Task Cleaner Tool

A maintenance tool for cleaning up old and orphaned task-related data in The Hub application.

## Overview

The task cleaner performs several maintenance operations on the task database:

- Removes old completed tasks based on retention policy
- Cleans up orphaned time entries, task dependencies, and scheduled tasks
- Updates parent task statuses based on subtask completion
- Permanently removes expired soft-deleted records
- Optimizes database indexes and analyzes tables

## Usage

### Building

```bash
# From the backend directory
make build-task-cleaner
# or
go build -o task-cleaner ./cmd/task-cleaner
```

### Running

```bash
# Run with default settings (90 days retention for completed tasks, 30 days for soft deletes)
./task-cleaner

# Dry run to see what would be cleaned without making changes
./task-cleaner --dry-run

# Custom retention periods
./task-cleaner --completed-retention=60 --soft-delete-retention=14

# Run with database optimization
./task-cleaner --optimize

# Combine options
./task-cleaner --dry-run --completed-retention=30 --optimize
```

### Command Line Options

- `--completed-retention`: Days to retain completed tasks (default: 90)
- `--soft-delete-retention`: Days to retain soft-deleted records before permanent deletion (default: 30)
- `--dry-run`: Show what would be cleaned without actually cleaning
- `--optimize`: Run database optimization after cleanup

### Makefile Targets

```bash
# Build the task cleaner
make build-task-cleaner

# Run task cleaner with defaults
make clean-tasks

# Run task cleaner in dry-run mode
make clean-tasks-dry-run
```

## Operations Performed

### 1. Clean Completed Tasks
Removes tasks with status "completed" that are older than the retention period.

### 2. Clean Orphaned Records
- **Time Entries**: Removes time tracking entries for tasks that no longer exist
- **Task Dependencies**: Removes dependency relationships for deleted tasks
- **Scheduled Tasks**: Removes calendar events for deleted tasks

### 3. Update Parent Task Statuses
Automatically marks parent tasks as completed when all their subtasks are completed.

### 4. Clean Expired Soft Deletes
Permanently removes soft-deleted records that are older than the retention period.

### 5. Database Optimization (optional)
- Rebuilds indexes concurrently to avoid blocking
- Updates table statistics for query optimization

## Safety Features

- **Dry Run Mode**: Use `--dry-run` to preview changes before applying them
- **Transactional Operations**: Each cleanup operation is isolated
- **Logging**: All operations are logged with counts of affected records
- **Error Handling**: Tool stops on critical errors but logs warnings for non-critical issues

## Database Compatibility

This tool is designed specifically for PostgreSQL and uses PostgreSQL-specific features:

- Native UUID support
- `NOW()` function for timestamps
- `REINDEX INDEX CONCURRENTLY` for online index maintenance
- `ANALYZE` for statistics updates

## Scheduling

For production use, consider scheduling this tool to run periodically:

```bash
# Cron job example (runs daily at 2 AM)
0 2 * * * cd /path/to/the-hub-backend && ./task-cleaner --optimize
```

## Monitoring

The tool logs all operations to stdout/stderr. In production, consider:

- Redirecting output to log files
- Setting up monitoring alerts for cleanup operations
- Tracking cleanup metrics over time

## Troubleshooting

### Common Issues

1. **Database Connection Failed**: Ensure environment variables are set correctly
2. **Permission Denied**: Ensure the database user has sufficient privileges
3. **Index Rebuild Fails**: Some indexes may not exist; warnings are logged but operation continues

### Recovery

If something goes wrong during cleanup:

1. The tool performs operations in a specific order to minimize impact
2. Check the logs for which operations completed successfully
3. For critical issues, restore from database backup
4. Run with `--dry-run` first on subsequent executions