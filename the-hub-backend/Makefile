# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build directory
BUILD_DIR=build

# Test coverage
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all build clean test test-coverage test-verbose deps help

# Default target
all: test

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/the-hub .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -rf $(BUILD_DIR)

# Run all tests
#test:
#	$(GOTEST) ./tests/... -v

# Run tests with coverage
#test-coverage:
	#$(GOTEST) ./tests/... -v -coverprofile=$(COVERAGE_FILE)
	#$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	#@echo "Coverage report generated: $(COVERAGE_HTML)"

# Run tests with verbose output
#test-verbose:
	#$(GOTEST) ./tests/... -v -race

# Run unit tests only
#test-unit:
#	$(GOTEST) ./tests/unit/... -v

# Run integration tests only
#test-integration:
	#$(GOTEST) ./tests/integration/... -v

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests and generate coverage badge
#test-with-badge: test-coverage
	#@echo "Generating coverage badge..."
	#@grep -o 'coverage: [0-9.]*%' $(COVERAGE_FILE) | head -1 | sed 's/coverage: //' | xargs -I {} echo "Coverage: {}%"

# Build task cleaner tool
build-task-cleaner:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/task-cleaner ./cmd/task-cleaner

# Run task cleaner with default settings
clean-tasks:
	$(BUILD_DIR)/task-cleaner

# Run task cleaner in dry-run mode
clean-tasks-dry-run:
	$(BUILD_DIR)/task-cleaner --dry-run

# Build admin CLI tool
build-admin-cli:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/admin-cli ./cmd/admin-cli

# Run admin CLI tool interactively
create-admin:
	$(BUILD_DIR)/admin-cli

# Create admin with command line arguments
create-admin-args:
	@if [ -z "$(EMAIL)" ] || [ -z "$(NAME)" ] || [ -z "$(PASSWORD)" ]; then \
		echo "Usage: make create-admin-args EMAIL=user@example.com NAME='Admin User' PASSWORD=password"; \
		exit 1; \
	fi
	$(BUILD_DIR)/admin-cli -email $(EMAIL) -name "$(NAME)" -password $(PASSWORD)

# Build seeder tool
build-seeder:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/seeder ./cmd/seeder

# Seed database with sample data
seed:
	@echo "Seeding database with sample data..."
	$(GOCMD) run ./cmd/seeder -command all

# Seed users only
seed-users:
	@echo "Seeding users..."
	$(GOCMD) run ./cmd/seeder -command users

# Seed goals and tasks only
seed-goals:
	@echo "Seeding goals and tasks..."
	$(GOCMD) run ./cmd/seeder -command goals

# Clean seeded data
seed-clean:
	@echo "Cleaning seeded data..."
	$(GOCMD) run ./cmd/seeder -command clean -force

# Migration commands
migrate-up:
	@echo "Running database migrations..."
	$(GOCMD) run -tags migrate ./cmd/migrate up

migrate-down:
	@echo "Rolling back last migration..."
	$(GOCMD) run -tags migrate ./cmd/migrate down 1

migrate-version:
	@echo "Current migration version:"
	$(GOCMD) run -tags migrate ./cmd/migrate version

migrate-status:
	@echo "Migration status:"
	$(GOCMD) run -tags migrate ./cmd/migrate status

# Help target
help:
	@echo "Available targets:"
	@echo "  build              - Build application"
	@echo "  build-task-cleaner - Build task cleaner tool"
	@echo "  build-admin-cli    - Build admin CLI tool"
	@echo "  build-seeder       - Build seeder tool"
	@echo "  clean              - Clean build artifacts"
	@echo "  clean-tasks        - Run task cleaner with default settings"
	@echo "  clean-tasks-dry-run - Run task cleaner in dry-run mode"
	@echo "  create-admin       - Run admin CLI tool interactively"
	@echo "  create-admin-args  - Create admin with command line arguments"
	@echo "  migrate-up         - Run all pending migrations"
	@echo "  migrate-down       - Rollback last migration"
	@echo "  migrate-version    - Show current migration version"
	@echo "  migrate-status     - Show migration status"
	@echo "  seed               - Seed database with all sample data"
	@echo "  seed-users         - Seed database with sample users only"
	@echo "  seed-goals         - Seed database with sample goals and tasks"
	@echo "  seed-clean         - Clean all seeded data"
	@echo "  test               - Run all tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  test-verbose       - Run tests with verbose output and race detection"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run integration tests only"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  help               - Show this help message"
