# Backend Testing

This directory contains comprehensive tests for the Go backend application.

## Test Structure

```
tests/
├── unit/                    # Unit tests
│   ├── util_test.go        # Utility function tests (password hashing, JWT)
│   ├── user_model_test.go  # User model and database tests
│   └── middleware_test.go  # JWT middleware tests
├── integration/            # Integration tests
│   └── user_handler_test.go # HTTP handler integration tests
├── fixtures/               # Test data fixtures
└── test_setup.go           # Common test setup and utilities
```

## Running Tests

### Using Make (Recommended)

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests with verbose output
make test-verbose

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Clean test artifacts
make clean
```

### Using Go directly

```bash
# Run all tests
go test ./tests/... -v

# Run with coverage
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run specific test packages
go test ./tests/unit/... -v
go test ./tests/integration/... -v
```

## Test Categories

### Unit Tests

**Utility Functions (`util_test.go`)**
- Password hashing and verification
- JWT token generation and validation
- Password hash consistency

**User Model (`user_model_test.go`)**
- User creation with valid/invalid data
- Email uniqueness constraints
- CRUD operations (Create, Read, Update, Delete)
- Soft delete functionality

**Middleware (`middleware_test.go`)**
- JWT authentication middleware
- Token validation
- Authorization header handling
- User context extraction

### Integration Tests

**User Handlers (`user_handler_test.go`)**
- User registration endpoint
- User login endpoint
- HTTP request/response validation
- Error handling for invalid inputs
- Database integration

## Test Database

Tests use an in-memory SQLite database for isolation and speed:
- Each test gets a fresh database instance
- No external dependencies required
- Fast execution and reliable results

## Environment Setup

Tests automatically configure:
- JWT secret key for token testing
- Gin test mode
- In-memory database
- Test-specific configurations

## Writing New Tests

### Unit Tests

```go
func TestMyFunction(t *testing.T) {
    // Test setup
    // ...

    // Test execution
    result := MyFunction(input)

    // Assertions
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Integration Tests

```go
func TestMyHandler(t *testing.T) {
    // Setup router
    router := gin.New()
    router.POST("/endpoint", MyHandler)

    // Create request
    requestBody := map[string]interface{}{"key": "value"}
    jsonBody, _ := json.Marshal(requestBody)
    req, _ := http.NewRequest("POST", "/endpoint", bytes.NewBuffer(jsonBody))

    // Execute request
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert response
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

## Test Coverage

Run `make test-coverage` to generate an HTML coverage report. The report will be saved as `coverage.html` in the project root.

## Best Practices

1. **Test Isolation**: Each test should be independent
2. **Descriptive Names**: Use clear, descriptive test names
3. **Table-Driven Tests**: Use table-driven tests for similar test cases
4. **Cleanup**: Clean up test data after each test
5. **Error Messages**: Provide clear error messages in assertions
6. **Edge Cases**: Test edge cases and error conditions
7. **Database Tests**: Use transactions or in-memory databases for database tests

## Continuous Integration

These tests are designed to run in CI environments:
- No external dependencies
- Fast execution
- Deterministic results
- Good coverage of critical paths