package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TheoMKgosi/The-hub/internal/handlers"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/TheoMKgosi/The-hub/tests"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDB *gorm.DB
var testRouter *gin.Engine

func TestMain(m *testing.M) {
	// Set up test environment
	tests.SetupTestEnvironment()

	// Initialize test database
	dsn := ":memory:"
	var err error
	testDB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Auto-migrate the schema
	err = testDB.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate test database")
	}

	// Set the test database in the config by directly modifying the package variable
	// This is a workaround since we can't easily modify the config package

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test router
	testRouter = gin.New()
	testRouter.POST("/auth/register", handlers.Register)
	testRouter.POST("/auth/login", handlers.Login)
	testRouter.GET("/users", handlers.GetUsers)
	testRouter.PUT("/users/:ID", handlers.UpdateUser)
	testRouter.DELETE("/users/:ID", handlers.DeleteUser)

	code := m.Run()

	// Cleanup
	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}

func TestUserRegistration(t *testing.T) {
	tests.CleanTestData()

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful registration",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"name":     "Test User",
				"password": "testpassword123",
				"settings": map[string]interface{}{"theme": "dark"},
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "missing email",
			requestBody: map[string]interface{}{
				"name":     "Test User",
				"password": "testpassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"email": "test2@example.com",
				"name":  "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
		{
			name: "duplicate email",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"name":     "Another User",
				"password": "anotherpassword",
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "Email already registered",
		},
		{
			name: "short password",
			requestBody: map[string]interface{}{
				"email":    "test3@example.com",
				"name":     "Test User",
				"password": "123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert request body to JSON
			requestBody, _ := json.Marshal(tt.requestBody)

			// Create request
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			testRouter.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check error message if expected
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				if errorMsg, exists := response["error"]; exists {
					if errorStr, ok := errorMsg.(string); ok {
						if !contains(errorStr, tt.expectedError) {
							t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, errorStr)
						}
					}
				} else {
					t.Errorf("Expected error message, but none found in response")
				}
			}

			// For successful registration, check if user was created
			if tt.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				if _, exists := response["token"]; !exists {
					t.Error("Expected token in successful registration response")
				}

				if _, exists := response["user_id"]; !exists {
					t.Error("Expected user_id in successful registration response")
				}
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	tests.CleanTestData()

	// Create a test user first
	testUser, err := tests.CreateTestUser("login@example.com", "Login User", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Hash the password for the test user
	hashedPassword, err := util.HashPassword("testpassword123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	testUser.Password = hashedPassword
	testDB.Save(testUser)

	loginTests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful login",
			requestBody: map[string]interface{}{
				"email":    "login@example.com",
				"password": "testpassword123",
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "wrong email",
			requestBody: map[string]interface{}{
				"email":    "wrong@example.com",
				"password": "testpassword123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid email or password",
		},
		{
			name: "wrong password",
			requestBody: map[string]interface{}{
				"email":    "login@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid email or password",
		},
		{
			name: "missing email",
			requestBody: map[string]interface{}{
				"password": "testpassword123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"email": "login@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
	}

	for _, tt := range loginTests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert request body to JSON
			requestBody, _ := json.Marshal(tt.requestBody)

			// Create request
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			testRouter.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check error message if expected
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				if errorMsg, exists := response["error"]; exists {
					if errorStr, ok := errorMsg.(string); ok {
						if !contains(errorStr, tt.expectedError) {
							t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, errorStr)
						}
					}
				}
			}

			// For successful login, check if token is returned
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				if _, exists := response["token"]; !exists {
					t.Error("Expected token in successful login response")
				}

				if _, exists := response["user"]; !exists {
					t.Error("Expected user data in successful login response")
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
