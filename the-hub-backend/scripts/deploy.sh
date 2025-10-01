#!/bin/bash

set -e  # Exit on any error

echo "🚀 Starting deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required environment variables are set
check_env_vars() {
    local required_vars=("DB_HOST" "DB_USER" "DB_PASSWORD" "DB_NAME")
    local missing_vars=()

    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            missing_vars+=("$var")
        fi
    done

    if [[ ${#missing_vars[@]} -ne 0 ]]; then
        print_error "Missing required environment variables: ${missing_vars[*]}"
        exit 1
    fi
}

# Wait for database to be ready
wait_for_db() {
    print_status "Waiting for database to be ready..."
    local max_attempts=30
    local attempt=1

    while [[ $attempt -le $max_attempts ]]; do
        if pg_isready -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" 2>/dev/null; then
            print_status "Database is ready!"
            return 0
        fi

        print_warning "Database not ready yet (attempt $attempt/$max_attempts). Waiting..."
        sleep 2
        ((attempt++))
    done

    print_error "Database failed to become ready after $max_attempts attempts"
    exit 1
}

# Run database migrations
run_migrations() {
    print_status "Running database migrations..."

    # Create database backup before migration (optional but recommended)
    if [[ "${CREATE_BACKUP:-true}" == "true" ]]; then
        print_status "Creating database backup..."
        local backup_file="/tmp/backup_$(date +%Y%m%d_%H%M%S).sql"
        if pg_dump -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" > "$backup_file" 2>/dev/null; then
            print_status "Backup created: $backup_file"
            # In production, you might want to upload this to cloud storage
        else
            print_warning "Failed to create backup, but continuing with migration..."
        fi
    fi

    # Run migrations
    if ./migrate up; then
        print_status "Migrations completed successfully!"
    else
        print_error "Migration failed!"
        # Attempt rollback if migration failed
        if [[ "${AUTO_ROLLBACK:-false}" == "true" ]]; then
            print_warning "Attempting to rollback last migration..."
            ./migrate down 1 || print_error "Rollback also failed!"
        fi
        exit 1
    fi
}

# Check migration status
check_migration_status() {
    print_status "Checking migration status..."
    if ./migrate version; then
        print_status "Migration status check passed"
    else
        print_error "Failed to check migration status"
        exit 1
    fi
}

# Health check
health_check() {
    print_status "Running health checks..."

    # Wait for application to start
    local max_attempts=30
    local attempt=1

    while [[ $attempt -le $max_attempts ]]; do
        if curl -f -s "http://localhost:8080/health" > /dev/null 2>&1; then
            print_status "Application health check passed!"
            return 0
        fi

        print_warning "Health check failed (attempt $attempt/$max_attempts). Waiting..."
        sleep 3
        ((attempt++))
    done

    print_error "Application failed health check after $max_attempts attempts"
    exit 1
}

# Main deployment function
main() {
    print_status "Starting deployment process..."

    # Check environment
    check_env_vars

    # Wait for database
    wait_for_db

    # Run migrations
    run_migrations

    # Check migration status
    check_migration_status

    # Start the application (this would be handled by your service manager)
    print_status "Starting application..."
    # systemctl start the-hub-backend  # or your service manager command

    # Health check
    health_check

    print_status "🎉 Deployment completed successfully!"
}

# Run main function
main "$@"