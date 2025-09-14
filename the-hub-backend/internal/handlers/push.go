package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushHandler struct {
	db *gorm.DB
}

func NewPushHandler(db *gorm.DB) *PushHandler {
	return &PushHandler{db: db}
}

type PushSubscriptionRequest struct {
	Subscription models.PushSubscription `json:"subscription"`
	UserAgent    string                  `json:"user_agent"`
}

type PushNotificationRequest struct {
	UserID uuid.UUID              `json:"user_id"`
	Title  string                 `json:"title"`
	Body   string                 `json:"body"`
	Icon   string                 `json:"icon,omitempty"`
	Badge  string                 `json:"badge,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

// Subscribe handles push subscription creation/updates
func (h *PushHandler) Subscribe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req PushSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID from context
	req.Subscription.UserID = userID.(uuid.UUID)
	req.Subscription.UserAgent = req.UserAgent

	// Check if subscription already exists
	var existing models.PushSubscription
	result := h.db.Where("user_id = ? AND endpoint = ?", userID, req.Subscription.Endpoint).First(&existing)

	if result.Error == nil {
		// Update existing subscription
		existing.P256dh = req.Subscription.P256dh
		existing.Auth = req.Subscription.Auth
		existing.UserAgent = req.UserAgent
		existing.IsActive = true

		if err := h.db.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
		return
	}

	// Create new subscription
	req.Subscription.IsActive = true
	if err := h.db.Create(&req.Subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subscription created successfully"})
}

// Unsubscribe handles push subscription removal
func (h *PushHandler) Unsubscribe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	endpoint := c.Query("endpoint")
	if endpoint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Endpoint is required"})
		return
	}

	// Soft delete by marking as inactive
	result := h.db.Model(&models.PushSubscription{}).
		Where("user_id = ? AND endpoint = ?", userID, endpoint).
		Update("is_active", false)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

// GetSubscriptions returns user's active subscriptions
func (h *PushHandler) GetSubscriptions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var subscriptions []models.PushSubscription
	if err := h.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": subscriptions})
}

// SendNotification sends a push notification to a user
func (h *PushHandler) SendNotification(c *gin.Context) {
	var req PushNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's active subscriptions
	var subscriptions []models.PushSubscription
	if err := h.db.Where("user_id = ? AND is_active = ?", req.UserID, true).Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	if len(subscriptions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active subscriptions found for user"})
		return
	}

	// Send notification to all subscriptions
	successCount := 0
	for _, sub := range subscriptions {
		if err := h.sendWebPush(sub, req); err != nil {
			// Log error but continue with other subscriptions
			config.Logger.Error("Failed to send push notification", "error", err, "endpoint", sub.Endpoint)
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Notifications sent",
		"success_count": successCount,
		"total_count":   len(subscriptions),
	})
}

// sendNotificationToUser sends a notification to a specific user
func (h *PushHandler) sendNotificationToUser(userID uuid.UUID, title, body string, data map[string]interface{}) error {
	// Get user's active subscriptions
	var subscriptions []models.PushSubscription
	if err := h.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&subscriptions).Error; err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		return nil // No subscriptions, not an error
	}

	// Send to all subscriptions
	for _, sub := range subscriptions {
		req := PushNotificationRequest{
			UserID: userID,
			Title:  title,
			Body:   body,
			Icon:   "/icon-192x192.png",
			Data:   data,
		}
		if err := h.sendWebPush(sub, req); err != nil {
			config.Logger.Error("Failed to send push notification", "error", err, "endpoint", sub.Endpoint)
		}
	}

	return nil
}

// sendWebPush sends a web push notification using the subscription
func (h *PushHandler) sendWebPush(subscription models.PushSubscription, req PushNotificationRequest) error {
	// This is a simplified implementation
	// In a real application, you'd use a library like web-push-go
	// and configure VAPID keys for authentication

	payload := map[string]interface{}{
		"title": req.Title,
		"body":  req.Body,
		"icon":  req.Icon,
		"badge": req.Badge,
		"data":  req.Data,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Here you would use a web push library to send the notification
	// For now, we'll just log it
	config.Logger.Info("Web push notification would be sent",
		"endpoint", subscription.Endpoint,
		"payload", string(payloadBytes))

	// TODO: Implement actual web push sending using a library like:
	// https://github.com/SherClockHolmes/webpush-go

	return nil
}

// SendTaskReminder sends a reminder notification for a task
func (h *PushHandler) SendTaskReminder(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get task details (simplified - you'd fetch from database)
	var task models.Task
	if err := h.db.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check if task belongs to user
	if task.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Send notification directly
	h.sendNotificationToUser(userID.(uuid.UUID), "Task Reminder", task.Title, map[string]interface{}{
		"type":    "task_reminder",
		"task_id": task.ID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Task reminder sent"})
}

// SendGoalReminder sends a reminder notification for a goal
func (h *PushHandler) SendGoalReminder(c *gin.Context) {
	goalID := c.Param("goal_id")
	if goalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Goal ID is required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get goal details (simplified - you'd fetch from database)
	var goal models.Goal
	if err := h.db.First(&goal, goalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Check if goal belongs to user
	if goal.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Send notification directly
	h.sendNotificationToUser(userID.(uuid.UUID), "Goal Reminder", goal.Title, map[string]interface{}{
		"type":    "goal_reminder",
		"goal_id": goal.ID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Goal reminder sent"})
}

// Handler functions for routes
func SubscribePush(c *gin.Context) {
	handler := NewPushHandler(config.GetDB())
	handler.Subscribe(c)
}

func UnsubscribePush(c *gin.Context) {
	handler := NewPushHandler(config.GetDB())
	handler.Unsubscribe(c)
}

func GetPushSubscriptions(c *gin.Context) {
	handler := NewPushHandler(config.GetDB())
	handler.GetSubscriptions(c)
}

func SendPushNotification(c *gin.Context) {
	handler := NewPushHandler(config.GetDB())
	handler.SendNotification(c)
}

func SendTaskReminder(c *gin.Context) {
	handler := NewPushHandler(config.GetDB())
	handler.SendTaskReminder(c)
}

func SendGoalReminder(c *gin.Context) {
	handler := NewPushHandler(config.GetDB())
	handler.SendGoalReminder(c)
}
