package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plutov/paypal/v4"
)

// PayPalWebhookEvent represents a webhook event from PayPal
type PayPalWebhookEvent struct {
	ID           string `json:"id"`
	EventType    string `json:"event_type"`
	ResourceType string `json:"resource_type"`
	Resource     struct {
		ID string `json:"id"`
	} `json:"resource"`
}

// GetSubscriptionPlans returns all available subscription plans
func GetSubscriptionPlans(c *gin.Context) {
	var plans []models.SubscriptionPlan

	if err := config.GetDB().Where("is_active = ?", true).Find(&plans).Error; err != nil {
		log.Println("Error fetching subscription plans:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch subscription plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

// GetUserSubscription returns the current user's subscription
func GetUserSubscription(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var subscription models.Subscription
	if err := config.GetDB().Preload("Plan").Where("user_id = ? AND status IN (?)", userIDUUID, []string{"active", "pending"}).First(&subscription).Error; err != nil {
		// No active subscription found
		c.JSON(http.StatusOK, gin.H{"subscription": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription": subscription})
}

// CreateSubscription creates a new PayPal subscription
func CreateSubscription(c *gin.Context) {
	var input struct {
		PlanID string `json:"plan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Check if user already has an active subscription
	var existingSubscription models.Subscription
	if err := config.GetDB().Where("user_id = ? AND status IN (?)", userIDUUID, []string{"active", "pending"}).First(&existingSubscription).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has an active subscription"})
		return
	}

	// Get the plan
	planID, err := uuid.Parse(input.PlanID)
	if err != nil {
		log.Println("Error parsing plan ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var plan models.SubscriptionPlan
	if err := config.GetDB().Where("id = ? AND is_active = ?", planID, true).First(&plan).Error; err != nil {
		log.Println("Error fetching plan:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription plan not found"})
		return
	}

	// Check if PayPal client is configured
	paypalClient := config.GetPayPalClient()
	if paypalClient == nil {
		log.Println("PayPal client not configured")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment system not configured"})
		return
	}

	// Create PayPal subscription
	subscriptionReq := paypal.SubscriptionBase{
		PlanID: plan.PayPalPlanID,
		Subscriber: &paypal.Subscriber{
			EmailAddress: "", // Will be set by PayPal during approval
		},
		ApplicationContext: &paypal.ApplicationContext{
			ReturnURL: "https://yourapp.com/subscription/success",
			CancelURL: "https://yourapp.com/subscription/cancel",
		},
	}

	paypalSub, err := paypalClient.CreateSubscription(context.Background(), subscriptionReq)
	if err != nil {
		log.Println("Error creating PayPal subscription:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	// Create local subscription record
	subscription := models.Subscription{
		UserID:               userIDUUID,
		PlanID:               planID,
		PayPalSubscriptionID: paypalSub.ID,
		Status:               "pending",
		StartDate:            time.Now(),
		AutoRenew:            true,
	}

	if err := config.GetDB().Create(&subscription).Error; err != nil {
		log.Println("Error creating subscription record:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"subscription":        subscription,
		"paypal_approval_url": paypalSub.Links[0].Href, // Usually the first link is the approval URL
	})
}

// CancelSubscription cancels a user's subscription
func CancelSubscription(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Find active subscription
	var subscription models.Subscription
	if err := config.GetDB().Where("user_id = ? AND status = ?", userIDUUID, "active").First(&subscription).Error; err != nil {
		log.Println("Error finding active subscription:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "No active subscription found"})
		return
	}

	// Check if PayPal client is configured
	paypalClient := config.GetPayPalClient()
	if paypalClient != nil {
		// Cancel PayPal subscription
		err := paypalClient.CancelSubscription(context.Background(), subscription.PayPalSubscriptionID, "User requested cancellation")
		if err != nil {
			log.Println("Error canceling PayPal subscription:", err)
			// Continue with local cancellation even if PayPal fails
		}
	}

	// Update local subscription
	now := time.Now()
	subscription.Status = "cancelled"
	subscription.CancelledAt = &now
	subscription.AutoRenew = false

	if err := config.GetDB().Save(&subscription).Error; err != nil {
		log.Println("Error updating subscription:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

// HandlePayPalWebhook processes PayPal webhook events
func HandlePayPalWebhook(c *gin.Context) {
	var webhookEvent PayPalWebhookEvent

	if err := c.ShouldBindJSON(&webhookEvent); err != nil {
		log.Println("Error parsing webhook:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
		return
	}

	// Store webhook event for processing
	eventRecord := models.PayPalWebhookEvent{
		EventType:  webhookEvent.EventType,
		ResourceID: webhookEvent.Resource.ID,
		EventData:  "", // You might want to marshal the full event
		Processed:  false,
	}

	if err := config.GetDB().Create(&eventRecord).Error; err != nil {
		log.Println("Error storing webhook event:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	// Process the webhook event
	switch webhookEvent.EventType {
	case "BILLING.SUBSCRIPTION.ACTIVATED":
		handleSubscriptionActivated(webhookEvent)
	case "BILLING.SUBSCRIPTION.CANCELLED":
		handleSubscriptionCancelled(webhookEvent)
	case "PAYMENT.SALE.COMPLETED":
		handlePaymentCompleted(webhookEvent)
	}

	// Mark as processed
	eventRecord.Processed = true
	eventRecord.ProcessedAt = &time.Time{}
	*eventRecord.ProcessedAt = time.Now()

	config.GetDB().Save(&eventRecord)

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// Helper functions for webhook processing
func handleSubscriptionActivated(event PayPalWebhookEvent) {
	subscriptionID := event.Resource.ID

	var subscription models.Subscription
	if err := config.GetDB().Where("paypal_subscription_id = ?", subscriptionID).First(&subscription).Error; err != nil {
		log.Println("Error finding subscription for activation:", err)
		return
	}

	subscription.Status = "active"
	config.GetDB().Save(&subscription)
}

func handleSubscriptionCancelled(event PayPalWebhookEvent) {
	subscriptionID := event.Resource.ID

	var subscription models.Subscription
	if err := config.GetDB().Where("paypal_subscription_id = ?", subscriptionID).First(&subscription).Error; err != nil {
		log.Println("Error finding subscription for cancellation:", err)
		return
	}

	now := time.Now()
	subscription.Status = "cancelled"
	subscription.CancelledAt = &now
	subscription.AutoRenew = false

	config.GetDB().Save(&subscription)
}

func handlePaymentCompleted(event PayPalWebhookEvent) {
	// Create payment record
	payment := models.Payment{
		PayPalPaymentID: event.Resource.ID,
		Amount:          0, // Extract from event data
		Currency:        "USD",
		Status:          "completed",
		PaymentDate:     time.Now(),
		Description:     "Subscription payment",
	}

	// You would need to extract subscription ID from the event and link it
	config.GetDB().Create(&payment)
}
