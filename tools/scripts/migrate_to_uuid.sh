#!/bin/bash

# UUID Migration Script for Production SQLite Database
# This script provides an easy way to migrate from uint to UUID primary keys

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to backup database
backup_database() {
    local db_path="$1"
    local backup_path="${db_path}.$(date +%Y%m%d_%H%M%S).backup"

    print_status "Creating backup of database..."
    cp "$db_path" "$backup_path"

    if [ $? -eq 0 ]; then
        print_success "Backup created: $backup_path"
        echo "$backup_path"
    else
        print_error "Failed to create backup"
        exit 1
    fi
}

# Function to check database integrity
check_database() {
    local db_path="$1"

    print_status "Checking database integrity..."

    if command_exists sqlite3; then
        sqlite3 "$db_path" "PRAGMA integrity_check;" > /dev/null 2>&1
        if [ $? -eq 0 ]; then
            print_success "Database integrity check passed"
            return 0
        else
            print_error "Database integrity check failed"
            return 1
        fi
    else
        print_warning "sqlite3 command not found, skipping integrity check"
        return 0
    fi
}

# Function to count records before migration
count_records_before() {
    local db_path="$1"
    local output_file="$2"

    print_status "Counting records before migration..."

    if command_exists sqlite3; then
        sqlite3 "$db_path" << 'EOF' > "$output_file"
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'goals', COUNT(*) FROM goals
UNION ALL
SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL
SELECT 'scheduled_tasks', COUNT(*) FROM scheduled_tasks
UNION ALL
SELECT 'decks', COUNT(*) FROM decks
UNION ALL
SELECT 'deck_users', COUNT(*) FROM deck_users
UNION ALL
SELECT 'cards', COUNT(*) FROM cards
UNION ALL
SELECT 'budget_categories', COUNT(*) FROM budget_categories
UNION ALL
SELECT 'budgets', COUNT(*) FROM budgets
UNION ALL
SELECT 'incomes', COUNT(*) FROM incomes
UNION ALL
SELECT 'topics', COUNT(*) FROM topics
UNION ALL
SELECT 'task_learning', COUNT(*) FROM task_learning
UNION ALL
SELECT 'tags', COUNT(*) FROM tags
UNION ALL
SELECT 'resources', COUNT(*) FROM resources
UNION ALL
SELECT 'study_sessions', COUNT(*) FROM study_sessions
UNION ALL
SELECT 'repeat_rules', COUNT(*) FROM repeat_rules
UNION ALL
SELECT 'ai_recommendations', COUNT(*) FROM ai_recommendations;
EOF
        print_success "Record counts saved to: $output_file"
    else
        print_error "sqlite3 command not found"
        exit 1
    fi
}

# Function to count records after migration
count_records_after() {
    local db_path="$1"
    local before_file="$2"
    local after_file="$3"

    print_status "Counting records after migration..."

    if command_exists sqlite3; then
        sqlite3 "$db_path" << 'EOF' > "$after_file"
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'goals', COUNT(*) FROM goals
UNION ALL
SELECT 'tasks', COUNT(*) FROM tasks
UNION ALL
SELECT 'scheduled_tasks', COUNT(*) FROM scheduled_tasks
UNION ALL
SELECT 'decks', COUNT(*) FROM decks
UNION ALL
SELECT 'deck_users', COUNT(*) FROM deck_users
UNION ALL
SELECT 'cards', COUNT(*) FROM cards
UNION ALL
SELECT 'budget_categories', COUNT(*) FROM budget_categories
UNION ALL
SELECT 'budgets', COUNT(*) FROM budgets
UNION ALL
SELECT 'incomes', COUNT(*) FROM incomes
UNION ALL
SELECT 'topics', COUNT(*) FROM topics
UNION ALL
SELECT 'task_learning', COUNT(*) FROM task_learning
UNION ALL
SELECT 'tags', COUNT(*) FROM tags
UNION ALL
SELECT 'resources', COUNT(*) FROM resources
UNION ALL
SELECT 'study_sessions', COUNT(*) FROM study_sessions
UNION ALL
SELECT 'repeat_rules', COUNT(*) FROM repeat_rules
UNION ALL
SELECT 'ai_recommendations', COUNT(*) FROM ai_recommendations;
EOF

        # Compare record counts
        if cmp -s "$before_file" "$after_file"; then
            print_success "Record counts match - no data loss detected"
        else
            print_warning "Record counts differ - please verify data integrity"
            print_status "Before migration:"
            cat "$before_file"
            print_status "After migration:"
            cat "$after_file"
        fi
    else
        print_error "sqlite3 command not found"
        exit 1
    fi
}

# Function to generate UUIDs using Python
generate_uuids_python() {
    local db_path="$1"

    print_status "Generating UUIDs for existing records using Python..."

    if command_exists python3; then
        python3 - << 'EOF' "$db_path"
import sqlite3
import uuid
import sys

db_path = sys.argv[1]
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

# List of tables to migrate
tables = [
    'users', 'goals', 'tasks', 'scheduled_tasks',
    'decks', 'deck_users', 'cards', 'budget_categories',
    'budgets', 'incomes', 'topics', 'task_learning',
    'tags', 'resources', 'study_sessions', 'repeat_rules',
    'ai_recommendations'
]

for table in tables:
    try:
        # Add UUID column if it doesn't exist
        cursor.execute(f"ALTER TABLE {table} ADD COLUMN id_uuid TEXT")

        # Generate UUIDs for existing records
        cursor.execute(f"SELECT id FROM {table}")
        rows = cursor.fetchall()

        for (old_id,) in rows:
            new_uuid = str(uuid.uuid4())
            cursor.execute(f"UPDATE {table} SET id_uuid = ? WHERE id = ?", (new_uuid, old_id))

        print(f"Generated UUIDs for {len(rows)} records in {table}")
        conn.commit()

    except sqlite3.Error as e:
        print(f"Error processing table {table}: {e}")
        conn.rollback()
        sys.exit(1)

conn.close()
print("UUID generation completed successfully")
EOF

        if [ $? -eq 0 ]; then
            print_success "UUID generation completed"
        else
            print_error "UUID generation failed"
            exit 1
        fi
    else
        print_error "Python 3 not found. Please install Python 3 or generate UUIDs manually."
        exit 1
    fi
}

# Function to run SQL migration
run_sql_migration() {
    local db_path="$1"
    local sql_file="$2"

    print_status "Running SQL migration..."

    if command_exists sqlite3; then
        sqlite3 "$db_path" < "$sql_file"

        if [ $? -eq 0 ]; then
            print_success "SQL migration completed"
        else
            print_error "SQL migration failed"
            exit 1
        fi
    else
        print_error "sqlite3 command not found"
        exit 1
    fi
}

# Main migration function
migrate_database() {
    local db_path="$1"
    local method="$2"

    # Check if database exists
    if [ ! -f "$db_path" ]; then
        print_error "Database file not found: $db_path"
        exit 1
    fi

    # Check database integrity
    if ! check_database "$db_path"; then
        print_error "Database integrity check failed. Please fix database issues before migration."
        exit 1
    fi

    # Create backup
    backup_path=$(backup_database "$db_path")

    # Count records before migration
    before_file="/tmp/before_migration_counts.txt"
    count_records_before "$db_path" "$before_file"

    case "$method" in
        "python")
            print_status "Using Python-based migration..."

            # Generate UUIDs
            generate_uuids_python "$db_path"

            # Run SQL migration
            sql_file="$(dirname "$0")/uuid_migration.sql"
            if [ -f "$sql_file" ]; then
                run_sql_migration "$db_path" "$sql_file"
            else
                print_error "SQL migration file not found: $sql_file"
                exit 1
            fi
            ;;

        "sql")
            print_status "Using SQL-based migration..."
            print_warning "Please ensure you have manually generated UUIDs for all records before proceeding"

            sql_file="$(dirname "$0")/uuid_migration.sql"
            if [ -f "$sql_file" ]; then
                run_sql_migration "$db_path" "$sql_file"
            else
                print_error "SQL migration file not found: $sql_file"
                exit 1
            fi
            ;;

        "recovery")
            print_status "Using recovery migration for failed previous attempts..."

            if command_exists python3; then
                recovery_script="$(dirname "$0")/recovery_migration.py"
                if [ -f "$recovery_script" ]; then
                    python3 "$recovery_script" "$db_path"
                    if [ $? -eq 0 ]; then
                        print_success "Recovery migration completed"
                    else
                        print_error "Recovery migration failed"
                        exit 1
                    fi
                else
                    print_error "Recovery migration script not found: $recovery_script"
                    exit 1
                fi
            else
                print_error "Python 3 not found. Recovery migration requires Python 3."
                exit 1
            fi
            ;;

        *)
            print_error "Invalid migration method. Use 'python', 'sql', or 'recovery'"
            exit 1
            ;;
    esac

    # Count records after migration
    after_file="/tmp/after_migration_counts.txt"
    count_records_after "$db_path" "$before_file" "$after_file"

    # Final integrity check
    if check_database "$db_path"; then
        print_success "Migration completed successfully!"
        print_status "Backup saved at: $backup_path"
        print_status "You can now start your application with UUID support"
    else
        print_error "Migration completed but database integrity check failed"
        print_warning "Please restore from backup and investigate the issue"
        print_status "Backup location: $backup_path"
        exit 1
    fi
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 <database_path> [method]

Arguments:
    database_path    Path to your SQLite database file
    method           Migration method: 'python' (recommended), 'sql', or 'recovery'

Examples:
    $0 myapp.db python    # Use Python script for automated migration
    $0 myapp.db sql       # Use SQL script (requires manual UUID generation)
    $0 myapp.db recovery  # Use recovery script for failed previous attempts

Description:
    This script migrates your SQLite database from uint primary keys to UUID primary keys.
    It will:
    1. Create a backup of your database
    2. Check database integrity
    3. Count records before migration
    4. Generate UUIDs for existing records (Python method only)
    5. Run the migration
    6. Verify record counts after migration
    7. Check database integrity

Methods:
    python   - Fully automated migration using Python script
    sql      - Manual migration using SQL script (requires UUID generation)
    recovery - Recovery migration for databases with failed previous attempts

Prerequisites:
    - sqlite3 command-line tool (for SQL method)
    - Python 3 with sqlite3 and uuid modules (for Python and recovery methods)

Safety:
    - Always creates a backup before migration
    - Uses transactions for atomicity
    - Provides rollback instructions if needed

EOF
}

# Main script logic
main() {
    if [ $# -lt 1 ] || [ $# -gt 2 ]; then
        show_usage
        exit 1
    fi

    local db_path="$1"
    local method="${2:-python}"

    print_warning "This script will modify your database. Make sure you have a backup!"
    print_status "Database: $db_path"
    print_status "Method: $method"

    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_status "Migration cancelled"
        exit 0
    fi

    migrate_database "$db_path" "$method"
}

# Run main function with all arguments
main "$@"