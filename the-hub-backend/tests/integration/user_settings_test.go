package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/TheoMKgosi/The-hub/internal/routes"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUserSettingsTestDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_USER", "postgres"),
		getEnvOrDefault("DB_PASSWORD", "postgres"),
		getEnvOrDefault("DB_NAME", "the_hub_test"),
		getEnvOrDefault("DB_PORT", "5432"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	// Initialize the database manager
	config.SetTestDB(db)

	return db, nil
}

func TestUserSettingsAPI(t *testing.T) {
	// Setup test database
	db, err := setupUserSettingsTestDB()
	assert.NoError(t, err)

	// Create test user with JSON string for settings
	settingsJSON := `{"theme": "light", "language": "en"}`
	testUser := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	// Set settings using JSON string
	result := db.Model(&testUser).Create(&testUser)
	assert.NoError(t, result.Error)

	// Update settings separately
	result = db.Model(&testUser).Update("settings", settingsJSON)
	assert.NoError(t, result.Error)

	// Generate JWT token for the test user
	token, err := util.GenerateJWT(testUser.ID)
	assert.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.RegisterRoutes(router)

	t.Run("GetUserSettings", func(t *testing.T) {
		// Create request
		req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d/settings", testUser.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		settings, exists := response["settings"]
		assert.True(t, exists, "Settings should be present in response")

		settingsMap, ok := settings.(map[string]interface{})
		assert.True(t, ok, "Settings should be a map")

		assert.Equal(t, "light", settingsMap["theme"])
		assert.Equal(t, "en", settingsMap["language"])
	})

	t.Run("UpdateUserSettings", func(t *testing.T) {
		// Prepare update data
		updateData := map[string]interface{}{
			"theme":         "dark",
			"notifications": true,
		}
		jsonData, _ := json.Marshal(updateData)

		// Create request
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d/settings", testUser.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		settings, exists := response["settings"]
		assert.True(t, exists, "Settings should be present in response")

		settingsMap, ok := settings.(map[string]interface{})
		assert.True(t, ok, "Settings should be a map")

		assert.Equal(t, "dark", settingsMap["theme"])
		assert.Equal(t, true, settingsMap["notifications"])
	})

	t.Run("PatchUserSettings", func(t *testing.T) {
		// Prepare patch data (only update theme)
		patchData := map[string]interface{}{
			"theme": "auto",
		}
		jsonData, _ := json.Marshal(patchData)

		// Create request
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/users/%d/settings", testUser.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		settings, exists := response["settings"]
		assert.True(t, exists, "Settings should be present in response")

		settingsMap, ok := settings.(map[string]interface{})
		assert.True(t, ok, "Settings should be a map")

		assert.Equal(t, "auto", settingsMap["theme"])
		// notifications should still be there from previous update
		assert.Equal(t, true, settingsMap["notifications"])
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		// Create another user
		otherUser := models.User{
			Name:     "Other User",
			Email:    "other@example.com",
			Password: "hashedpassword",
		}
		db.Create(&otherUser)

		// Try to access other user's settings with testUser's token
		req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d/settings", otherUser.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should be forbidden
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
