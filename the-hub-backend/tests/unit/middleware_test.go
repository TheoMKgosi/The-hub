package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-tokens")
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	router := gin.New()
	router.Use(util.JWTAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
	})

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing authorization header",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization header missing or malformed",
		},
		{
			name:           "malformed authorization header",
			authorization:  "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Authorization header missing or malformed",
		},
		{
			name:           "invalid token",
			authorization:  "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid or expired token",
		},
		{
			name:           "valid token",
			authorization:  "Bearer " + createTestToken(1),
			expectedStatus: http.StatusOK,
			expectedBody:   "Access granted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body contains expected message
			if tt.expectedBody != "" && !contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("Expected response body to contain '%s', got '%s'", tt.expectedBody, w.Body.String())
			}

			// For valid token, check if userID is set in context
			if tt.name == "valid token" && tt.expectedStatus == http.StatusOK {
				// The middleware should have set userID in context
				// We can't easily test this without more complex setup, but we can verify the request succeeded
				if w.Code != http.StatusOK {
					t.Error("Valid token should allow access to protected route")
				}
			}
		})
	}
}

func TestJWTAuthMiddlewareTokenValidation(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-tokens")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(util.JWTAuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Test with valid token
	validToken := createTestToken(123)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response to check if userID is correctly extracted
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if userID, exists := response["user_id"]; exists {
		if userID != float64(123) {
			t.Errorf("Expected user_id 123, got %v", userID)
		}
	} else {
		t.Error("Expected user_id in response")
	}
}

// Helper function to create a test JWT token
func createTestToken(userID uint) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     9999999999, // Far future expiry for testing
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString
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
