# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Test coverage
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all build clean test test-coverage test-verbose deps help

# Default target
all: test

# Build the application
build:
	$(GOBUILD) -o the-hub .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -f the-hub

# Run all tests
test:
	$(GOTEST) ./tests/... -v

# Run tests with coverage
test-coverage:
	$(GOTEST) ./tests/... -v -coverprofile=$(COVERAGE_FILE)
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Run tests with verbose output
test-verbose:
	$(GOTEST) ./tests/... -v -race

# Run unit tests only
test-unit:
	$(GOTEST) ./tests/unit/... -v

# Run integration tests only
test-integration:
	$(GOTEST) ./tests/integration/... -v

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests and generate coverage badge
test-with-badge: test-coverage
	@echo "Generating coverage badge..."
	@grep -o 'coverage: [0-9.]*%' $(COVERAGE_FILE) | head -1 | sed 's/coverage: //' | xargs -I {} echo "Coverage: {}%"

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-verbose   - Run tests with verbose output and race detection"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  help           - Show this help message"