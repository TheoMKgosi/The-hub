package routes

import (
	"github.com/TheoMKgosi/The-hub/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupSubscriptionRoutes sets up subscription-related routes
func SetupSubscriptionRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	subscriptionGroup := router.Group("/api/v1/subscriptions")
	subscriptionGroup.Use(authMiddleware)
	{
		subscriptionGroup.GET("/plans", handlers.GetSubscriptionPlans)
		subscriptionGroup.GET("/user", handlers.GetUserSubscription)
		subscriptionGroup.POST("", handlers.CreateSubscription)
		subscriptionGroup.DELETE("", handlers.CancelSubscription)
	}

	// Webhook endpoint (no auth required)
	router.POST("/api/v1/webhooks/paypal", handlers.HandlePayPalWebhook)
}
