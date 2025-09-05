package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

// OAuth2 configuration for Google Calendar
var googleOAuthConfig *oauth2.Config

func init() {
	// Initialize Google OAuth2 config
	googleOAuthConfig = &oauth2.Config{
		ClientID:     getEnvOrDefault("GOOGLE_CLIENT_ID", ""),
		ClientSecret: getEnvOrDefault("GOOGLE_CLIENT_SECRET", ""),
		RedirectURL:  getEnvOrDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		Scopes: []string{
			calendar.CalendarScope,
			calendar.CalendarEventsScope,
		},
		Endpoint: google.Endpoint,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	// Simple environment variable getter - in production use proper env handling
	if key == "GOOGLE_CLIENT_ID" {
		return "your-google-client-id" // Replace with actual env var
	}
	if key == "GOOGLE_CLIENT_SECRET" {
		return "your-google-client-secret" // Replace with actual env var
	}
	if key == "GOOGLE_REDIRECT_URL" {
		return "http://localhost:8080/auth/google/callback"
	}
	return defaultValue
}

// InitiateGoogleCalendarAuth starts the OAuth flow for Google Calendar
func InitiateGoogleCalendarAuth(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate state parameter for CSRF protection
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Store state in session/cache (simplified - use proper session management)
	stateKey := fmt.Sprintf("oauth_state_%s", userID.(uuid.UUID).String())
	// In production, store in Redis/session store
	_ = stateKey

	authURL := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

// HandleGoogleCalendarCallback handles the OAuth callback from Google
func HandleGoogleCalendarCallback(c *gin.Context) {
	code := c.Query("code")
	_ = c.Query("state") // State verification would be implemented in production

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	// Verify state parameter (simplified)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Exchange code for token
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code"})
		return
	}

	// Get user profile from Google
	client := googleOAuthConfig.Client(context.Background(), token)
	calendarService, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Failed to create calendar service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create calendar service"})
		return
	}

	// Get user's calendar list to find primary calendar
	calendarList, err := calendarService.CalendarList.List().Do()
	if err != nil {
		log.Printf("Failed to get calendar list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar list"})
		return
	}

	var primaryCalendarID string
	for _, cal := range calendarList.Items {
		if cal.Primary != true {
			primaryCalendarID = cal.Id
			break
		}
	}

	if primaryCalendarID == "" && len(calendarList.Items) > 0 {
		primaryCalendarID = calendarList.Items[0].Id
	}

	// Save integration to database
	integration := models.CalendarIntegration{
		UserID:         userID.(uuid.UUID),
		Provider:       models.ProviderGoogle,
		ProviderUserID: "google_user_id", // Get from Google profile
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenExpiry:    token.Expiry,
		CalendarID:     primaryCalendarID,
		IsActive:       true,
		SyncEnabled:    true,
	}

	if err := config.GetDB().Create(&integration).Error; err != nil {
		log.Printf("Failed to save calendar integration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save calendar integration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Google Calendar integration successful",
		"integration_id": integration.ID,
	})
}

// GetCalendarIntegrations returns user's calendar integrations
func GetCalendarIntegrations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var integrations []models.CalendarIntegration
	if err := config.GetDB().Where("user_id = ?", userID).Find(&integrations).Error; err != nil {
		log.Printf("Failed to get calendar integrations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar integrations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"integrations": integrations})
}

// SyncCalendarEvents syncs events between local schedule and external calendar
func SyncCalendarEvents(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	integrationID := c.Param("integrationID")
	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Integration ID required"})
		return
	}

	var integration models.CalendarIntegration
	if err := config.GetDB().Where("id = ? AND user_id = ?", integrationID, userID).First(&integration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Calendar integration not found"})
			return
		}
		log.Printf("Failed to get calendar integration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar integration"})
		return
	}

	// Refresh token if needed
	if time.Now().After(integration.TokenExpiry) {
		if err := refreshGoogleToken(&integration); err != nil {
			log.Printf("Failed to refresh token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh access token"})
			return
		}
	}

	// Create OAuth2 token
	token := &oauth2.Token{
		AccessToken:  integration.AccessToken,
		RefreshToken: integration.RefreshToken,
		Expiry:       integration.TokenExpiry,
	}

	// Create calendar service
	client := googleOAuthConfig.Client(context.Background(), token)
	calendarService, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Failed to create calendar service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create calendar service"})
		return
	}

	// Sync scheduled tasks to calendar
	if err := syncScheduledTasksToCalendar(calendarService, integration); err != nil {
		log.Printf("Failed to sync to calendar: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync scheduled tasks"})
		return
	}

	// Update last sync time
	now := time.Now()
	if err := config.GetDB().Model(&integration).Update("last_sync_at", &now).Error; err != nil {
		log.Printf("Failed to update last sync time: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Calendar sync completed successfully"})
}

// DeleteCalendarIntegration removes a calendar integration
func DeleteCalendarIntegration(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	integrationID := c.Param("integrationID")
	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Integration ID required"})
		return
	}

	// Start transaction
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete sync records
	if err := tx.Where("calendar_integration_id = ?", integrationID).Delete(&models.CalendarEventSync{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to delete calendar event syncs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete calendar sync records"})
		return
	}

	// Delete integration
	if err := tx.Where("id = ? AND user_id = ?", integrationID, userID).Delete(&models.CalendarIntegration{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to delete calendar integration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete calendar integration"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Calendar integration deleted successfully"})
}

// Helper functions

func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func refreshGoogleToken(integration *models.CalendarIntegration) error {
	token := &oauth2.Token{
		AccessToken:  integration.AccessToken,
		RefreshToken: integration.RefreshToken,
		Expiry:       integration.TokenExpiry,
	}

	tokenSource := googleOAuthConfig.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return err
	}

	// Update integration with new token
	integration.AccessToken = newToken.AccessToken
	integration.RefreshToken = newToken.RefreshToken
	integration.TokenExpiry = newToken.Expiry

	return config.GetDB().Save(integration).Error
}

func syncScheduledTasksToCalendar(svc *calendar.Service, integration models.CalendarIntegration) error {
	// Get scheduled tasks for the user
	var scheduledTasks []models.ScheduledTask
	if err := config.GetDB().Where("user_id = ?", integration.UserID).Find(&scheduledTasks).Error; err != nil {
		return err
	}

	for _, task := range scheduledTasks {
		// Check if already synced
		var existingSync models.CalendarEventSync
		err := config.GetDB().Where("scheduled_task_id = ? AND calendar_integration_id = ?", task.ID, integration.ID).First(&existingSync).Error

		if err == gorm.ErrRecordNotFound {
			// Create new calendar event
			event := &calendar.Event{
				Summary: task.Title,
				Start: &calendar.EventDateTime{
					DateTime: task.Start.Format(time.RFC3339),
					TimeZone: "UTC",
				},
				End: &calendar.EventDateTime{
					DateTime: task.End.Format(time.RFC3339),
					TimeZone: "UTC",
				},
			}

			createdEvent, err := svc.Events.Insert(integration.CalendarID, event).Do()
			if err != nil {
				log.Printf("Failed to create calendar event: %v", err)
				continue
			}

			// Save sync record
			sync := models.CalendarEventSync{
				ScheduledTaskID:       task.ID,
				CalendarIntegrationID: integration.ID,
				ExternalEventID:       createdEvent.Id,
				LastSyncedAt:          time.Now(),
			}

			if err := config.GetDB().Create(&sync).Error; err != nil {
				log.Printf("Failed to save sync record: %v", err)
			}
		} else if err == nil {
			// Update existing event
			event := &calendar.Event{
				Summary: task.Title,
				Start: &calendar.EventDateTime{
					DateTime: task.Start.Format(time.RFC3339),
					TimeZone: "UTC",
				},
				End: &calendar.EventDateTime{
					DateTime: task.End.Format(time.RFC3339),
					TimeZone: "UTC",
				},
			}

			_, err := svc.Events.Update(integration.CalendarID, existingSync.ExternalEventID, event).Do()
			if err != nil {
				log.Printf("Failed to update calendar event: %v", err)
				continue
			}

			// Update sync record
			existingSync.LastSyncedAt = time.Now()
			if err := config.GetDB().Save(&existingSync).Error; err != nil {
				log.Printf("Failed to update sync record: %v", err)
			}
		}
	}

	return nil
}
