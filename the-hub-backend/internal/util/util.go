package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))                    // should be in env variables in production
var refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET")) // separate secret for refresh tokens

// LoadJWTSecret reloads the JWT secret from environment variables
// This is primarily used for testing
func LoadJWTSecret() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
}

func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// GenerateAccessToken creates a short-lived access token (15 minutes)
func GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // token expires in 15 minutes
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken creates a cryptographically secure refresh token
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashRefreshToken hashes a refresh token for secure storage
func HashRefreshToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckRefreshTokenHash verifies a refresh token against its hash
func CheckRefreshTokenHash(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

// HashRefreshTokenUnsafe is a helper function for database queries (not for storage)
func HashRefreshTokenUnsafe(token string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// ValidateRefreshToken validates a refresh token and returns the token record if valid
func ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	var refreshTokens []models.RefreshToken

	// Get all non-revoked refresh tokens that haven't expired
	if err := config.GetDB().Where("revoked = false AND expires_at > ?", time.Now()).Find(&refreshTokens).Error; err != nil {
		return nil, err
	}

	// Check each token to see if it matches
	for _, refreshToken := range refreshTokens {
		if CheckRefreshTokenHash(tokenString, refreshToken.TokenHash) {
			return &refreshToken, nil
		}
	}

	// If we get here, no matching token was found
	return nil, gorm.ErrRecordNotFound
}

// RevokeRefreshToken marks a refresh token as revoked
func RevokeRefreshToken(tokenID uuid.UUID) error {
	return config.GetDB().Model(&models.RefreshToken{}).Where("id = ?", tokenID).Update("revoked", true).Error
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
func RevokeAllUserRefreshTokens(userID uuid.UUID) error {
	return config.GetDB().Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = false", userID).Update("revoked", true).Error
}

// GetDefaultUserSettings returns the default settings for new users
func GetDefaultUserSettings() map[string]interface{} {
	return map[string]interface{}{
		"theme": map[string]interface{}{
			"mode": "system", // light, dark, system
		},
		"notifications": map[string]interface{}{
			"email": map[string]interface{}{
				"enabled":        true,
				"task_reminders": true,
				"goal_deadlines": true,
				"weekly_reports": true,
			},
			"push": map[string]interface{}{
				"enabled":        false,
				"task_reminders": false,
				"goal_deadlines": false,
			},
		},
		"dashboard": map[string]interface{}{
			"layout": "default",
			"widgets": map[string]interface{}{
				"tasks":    true,
				"goals":    true,
				"schedule": true,
				"stats":    true,
			},
		},
		"task": map[string]interface{}{
			"tri-modal": false,
		},
		"privacy": map[string]interface{}{
			"profile_visibility": "private",
			"activity_sharing":   false,
		},
		"preferences": map[string]interface{}{
			"language":    "en",
			"timezone":    "UTC",
			"date_format": "MM/DD/YYYY",
			"time_format": "12h",
		},
	}
}

func mergeSettings(userSettings, defaults map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Start with defaults
	for key, defaultValue := range defaults {
		result[key] = defaultValue
	}

	// Override with user settings
	for key, userValue := range userSettings {
		if defaultValue, exists := defaults[key]; exists {
			// Both exist - check if both are maps for recursive merge
			if userMap, userIsMap := userValue.(map[string]interface{}); userIsMap {
				if defaultMap, defaultIsMap := defaultValue.(map[string]interface{}); defaultIsMap {
					result[key] = mergeSettings(userMap, defaultMap)
					continue
				}
			}
		}
		// User value takes precedence
		result[key] = userValue
	}

	return result
}

func EnsureCompleteSettings(settingsJSON string) (string, bool) {
	var userSettings map[string]interface{}

	// Parse existing settings safely
	if settingsJSON != "" && settingsJSON != "{}" {
		if err := json.Unmarshal([]byte(settingsJSON), &userSettings); err != nil {
			config.Logger.Warnf("Failed to parse user settings, using defaults: %v", err)
			userSettings = make(map[string]interface{})
		}
	} else {
		userSettings = make(map[string]interface{})
	}

	// Get current defaults
	defaults := GetDefaultUserSettings()

	// Merge settings
	mergedSettings := mergeSettings(userSettings, defaults)

	// Serialize back to JSON
	mergedJSON, err := json.Marshal(mergedSettings)
	if err != nil {
		config.Logger.Errorf("Failed to marshal merged settings: %v", err)
		return settingsJSON, false // Return original on error
	}

	// Check if migration was needed
	migrationNeeded := string(mergedJSON) != settingsJSON

	if migrationNeeded {
		config.Logger.Infof("Settings auto-migration completed")
	}

	return string(mergedJSON), migrationNeeded
}
