package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushNotificationService struct {
	db *gorm.DB
}

type NotificationEvent struct {
	UserID   uuid.UUID              `json:"user_id"`
	Type     string                 `json:"type"` // task_reminder, goal_deadline, budget_alert, etc.
	Title    string                 `json:"title"`
	Body     string                 `json:"body"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Priority string                 `json:"priority,omitempty"` // high, normal, low
}

func NewPushNotificationService(db *gorm.DB) *PushNotificationService {
	return &PushNotificationService{db: db}
}

// SendNotification sends a push notification to a user if they have enabled notifications
func (s *PushNotificationService) SendNotification(event NotificationEvent) error {
	// Check user preferences
	if !s.shouldSendNotification(event.UserID, event.Type) {
		return nil // User has disabled this type of notification
	}

	// Get user's active subscriptions
	var subscriptions []models.PushSubscription
	if err := s.db.Where("user_id = ? AND is_active = ?", event.UserID, true).Find(&subscriptions).Error; err != nil {
		return fmt.Errorf("failed to fetch subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		return nil // No active subscriptions
	}

	// Send to all subscriptions
	successCount := 0
	for _, sub := range subscriptions {
		if err := s.sendWebPush(sub, event); err != nil {
			config.Logger.Error("Failed to send push notification", "error", err, "endpoint", sub.Endpoint)
		} else {
			successCount++
		}
	}

	config.Logger.Info("Push notification sent", "user_id", event.UserID, "type", event.Type, "success_count", successCount)
	return nil
}

// shouldSendNotification checks if the user has enabled this type of notification
func (s *PushNotificationService) shouldSendNotification(userID uuid.UUID, notificationType string) bool {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		config.Logger.Error("Failed to fetch user for notification preferences", "error", err)
		return false
	}

	// Parse user settings
	var settings map[string]interface{}
	if user.Settings == "" {
		return false
	}
	if err := json.Unmarshal([]byte(user.Settings), &settings); err != nil {
		config.Logger.Error("Failed to parse user settings", "error", err)
		return false
	}

	// Check push notification settings
	pushSettings, ok := settings["notifications"].(map[string]interface{})
	if !ok {
		return false
	}

	push, ok := pushSettings["push"].(map[string]interface{})
	if !ok {
		return false
	}

	// Check if push notifications are enabled globally
	enabled, ok := push["enabled"].(bool)
	if !ok || !enabled {
		return false
	}

	// Check specific notification type
	switch notificationType {
	case "task_reminder":
		return push["task_reminders"].(bool)
	case "goal_deadline":
		return push["goal_deadlines"].(bool)
	case "budget_alert":
		return push["budget_alerts"].(bool)
	case "study_reminder":
		return push["study_reminders"].(bool)
	default:
		return false
	}
}

// sendWebPush sends a web push notification (simplified implementation)
func (s *PushNotificationService) sendWebPush(subscription models.PushSubscription, event NotificationEvent) error {
	// This is a placeholder implementation
	// In a real application, you would use a web push library like:
	// - github.com/SherClockHolmes/webpush-go
	// - github.com/appleboy/go-fcm

	config.Logger.Info("Web push notification would be sent",
		"endpoint", subscription.Endpoint,
		"type", event.Type,
		"title", event.Title,
		"body", event.Body)

	// TODO: Implement actual web push sending
	// Example with webpush-go:
	/*
		sub := &webpush.Subscription{
			Endpoint: subscription.Endpoint,
			Keys: webpush.Keys{
				P256dh: subscription.P256dh,
				Auth:   subscription.Auth,
			},
		}

		payload := map[string]interface{}{
			"title": event.Title,
			"body":  event.Body,
			"data":  event.Data,
		}

		payloadBytes, _ := json.Marshal(payload)

		resp, err := webpush.SendNotification(payloadBytes, sub, &webpush.Options{
			Subscriber:      "mailto:admin@thehub.com", // Your email
			VAPIDPublicKey:  "YOUR_VAPID_PUBLIC_KEY",
			VAPIDPrivateKey: "YOUR_VAPID_PRIVATE_KEY",
			TTL:             30,
		})
	*/

	return nil
}

// SendTaskReminder sends a task reminder notification
func (s *PushNotificationService) SendTaskReminder(taskID uuid.UUID, userID uuid.UUID, taskTitle string, dueDate *time.Time) error {
	var body string
	if dueDate != nil {
		body = fmt.Sprintf("Task '%s' is due on %s", taskTitle, dueDate.Format("Jan 2, 2006"))
	} else {
		body = fmt.Sprintf("Don't forget about task '%s'", taskTitle)
	}

	event := NotificationEvent{
		UserID: userID,
		Type:   "task_reminder",
		Title:  "Task Reminder",
		Body:   body,
		Data: map[string]interface{}{
			"type":    "task_reminder",
			"task_id": taskID,
		},
		Priority: "normal",
	}

	return s.SendNotification(event)
}

// SendGoalDeadlineReminder sends a goal deadline reminder
func (s *PushNotificationService) SendGoalDeadlineReminder(goalID uuid.UUID, userID uuid.UUID, goalTitle string, dueDate time.Time) error {
	body := fmt.Sprintf("Goal '%s' is due on %s", goalTitle, dueDate.Format("Jan 2, 2006"))

	event := NotificationEvent{
		UserID: userID,
		Type:   "goal_deadline",
		Title:  "Goal Deadline",
		Body:   body,
		Data: map[string]interface{}{
			"type":    "goal_deadline",
			"goal_id": goalID,
		},
		Priority: "high",
	}

	return s.SendNotification(event)
}

// SendBudgetAlert sends a budget alert notification
func (s *PushNotificationService) SendBudgetAlert(userID uuid.UUID, categoryName string, alertType string, message string) error {
	event := NotificationEvent{
		UserID: userID,
		Type:   "budget_alert",
		Title:  "Budget Alert",
		Body:   message,
		Data: map[string]interface{}{
			"type":          "budget_alert",
			"category_name": categoryName,
			"alert_type":    alertType,
		},
		Priority: "high",
	}

	return s.SendNotification(event)
}

// SendStudyReminder sends a study reminder notification
func (s *PushNotificationService) SendStudyReminder(userID uuid.UUID, topicTitle string, message string) error {
	event := NotificationEvent{
		UserID: userID,
		Type:   "study_reminder",
		Title:  "Study Reminder",
		Body:   message,
		Data: map[string]interface{}{
			"type":        "study_reminder",
			"topic_title": topicTitle,
		},
		Priority: "normal",
	}

	return s.SendNotification(event)
}

// SendWelcomeNotification sends a welcome notification to new users
func (s *PushNotificationService) SendWelcomeNotification(userID uuid.UUID, userName string) error {
	event := NotificationEvent{
		UserID: userID,
		Type:   "welcome",
		Title:  "Welcome to The Hub!",
		Body:   fmt.Sprintf("Hi %s! Welcome to your personal learning platform. Let's get started!", userName),
		Data: map[string]interface{}{
			"type": "welcome",
		},
		Priority: "normal",
	}

	return s.SendNotification(event)
}

// SendAchievementNotification sends an achievement notification
func (s *PushNotificationService) SendAchievementNotification(userID uuid.UUID, achievementTitle string, description string) error {
	event := NotificationEvent{
		UserID: userID,
		Type:   "achievement",
		Title:  "Achievement Unlocked!",
		Body:   fmt.Sprintf("%s - %s", achievementTitle, description),
		Data: map[string]interface{}{
			"type":              "achievement",
			"achievement_title": achievementTitle,
			"description":       description,
		},
		Priority: "high",
	}

	return s.SendNotification(event)
}
