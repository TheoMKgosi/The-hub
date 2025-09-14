package routes

import (
	"github.com/TheoMKgosi/The-hub/internal/handlers"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(router *gin.Engine) {
	// Public routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)
	router.POST("/auth/forgot-password", handlers.RequestPasswordReset)
	router.POST("/auth/reset-password", handlers.ResetPassword)
	router.POST("/auth/refresh", handlers.RefreshToken)
	router.POST("/auth/logout", handlers.Logout)

	protected := router.Group("/")
	protected.Use(util.JWTAuthMiddleware())

	// User management routes
	protected.GET("/users/:ID", handlers.GetUser)
	protected.PUT("/users/:ID", handlers.UpdateUser)
	protected.DELETE("/users/:ID", handlers.DeleteUser)

	// User settings routes
	protected.GET("/users/:ID/settings", handlers.GetUserSettings)
	protected.PUT("/users/:ID/settings", handlers.UpdateUserSettings)
	protected.PATCH("/users/:ID/settings", handlers.PatchUserSettings)

	// Plan routes
	// -- Goal routes
	protected.GET("/goals", handlers.GetGoals)
	protected.GET("/goals/:ID", handlers.GetGoal)
	protected.POST("/goals", handlers.CreateGoal)
	protected.PUT("/goals/:ID", handlers.UpdateGoal)
	protected.DELETE("/goals/:ID", handlers.DeleteGoal)

	// -- Goal Task routes
	protected.POST("/goals/:ID/tasks", handlers.AddTaskToGoal)
	protected.GET("/goals/:ID/tasks", handlers.GetGoalTasks)
	protected.PATCH("/goals/:ID/tasks/:taskID", handlers.UpdateGoalTask)
	protected.DELETE("/goals/:ID/tasks/:taskID", handlers.DeleteGoalTask)
	protected.PATCH("/goals/:ID/tasks/:taskID/complete", handlers.CompleteGoalTask)

	// -- Task routes
	protected.GET("/tasks", handlers.GetTasks)
	protected.GET("/tasks/:ID", handlers.GetTask)
	protected.POST("/tasks", handlers.CreateTask)
	protected.PUT("/tasks/reorder", handlers.ReorderTasks)
	protected.PATCH("/tasks/:ID", handlers.UpdateTask)
	protected.DELETE("/tasks/:ID", handlers.DeleteTask)

	// -- Task Statistics routes
	protected.GET("/stats/tasks", handlers.GetTaskStats)
	protected.GET("/stats/tasks/trends", handlers.GetTaskStatsTrends)

	// Time routes
	// -- Schedule routes
	protected.GET("/schedule", handlers.GetSchedule)
	protected.POST("/schedule", handlers.CreateSchedule)
	protected.PUT("/schedule/:ID", handlers.UpdateSchedule)
	protected.DELETE("/schedule/:ID", handlers.DeleteSchedule)
	protected.POST("/schedule/bulk", handlers.BulkCreateSchedule)
	protected.DELETE("/schedule/bulk", handlers.BulkDeleteSchedule)

	// -- Recurrence rule routes
	protected.POST("/recurrence-rules", handlers.CreateRecurrenceRule)

	// -- AI routes
	protected.GET("/ai/suggestions", handlers.GetAISuggestions)

	// Calendar integration routes
	protected.POST("/calendar/google/auth", handlers.InitiateGoogleCalendarAuth)
	protected.GET("/calendar/google/callback", handlers.HandleGoogleCalendarCallback)
	protected.GET("/calendar/integrations", handlers.GetCalendarIntegrations)
	protected.POST("/calendar/integrations/:integrationID/sync", handlers.SyncCalendarEvents)
	protected.DELETE("/calendar/integrations/:integrationID", handlers.DeleteCalendarIntegration)

	// Learning routes
	// -- Deck routes
	protected.GET("/decks", handlers.GetDecks)
	protected.GET("/decks/:ID", handlers.GetDeck)
	protected.POST("/decks", handlers.CreateDeck)
	protected.PATCH("/decks/:ID", handlers.UpdateDeck)
	protected.DELETE("/decks/:ID", handlers.DeleteDeck)

	// -- Card routes
	protected.GET("/decks/cards/:deckID", handlers.GetCards)
	protected.GET("/cards/:ID", handlers.GetCard)
	protected.POST("/cards", handlers.CreateCard)
	protected.PATCH("/cards/:ID", handlers.UpdateCard)
	protected.DELETE("/cards/:ID", handlers.DeleteCard)

	protected.POST("/cards/review/:ID", handlers.ReviewCard)
	protected.GET("/cards/due/:deckID", handlers.GetDueCards)

	// -- Topic routes
	protected.GET("/topics", handlers.GetTopics)
	protected.GET("/topics/:ID", handlers.GetTopic)
	protected.POST("/topics", handlers.CreateTopic)
	protected.PATCH("/topics/:ID", handlers.UpdateTopic)
	protected.DELETE("/topics/:ID", handlers.DeleteTopic)
	protected.GET("/topics/public", handlers.GetPublicTopics)

	// -- Task Learning
	protected.GET("/task-learning/:ID", handlers.GetTaskLearnings)
	protected.POST("/task-learning", handlers.CreateTaskLearning)
	protected.PATCH("/task-learning/:ID", handlers.UpdateTaskLearning)
	protected.DELETE("/task-learning/:ID", handlers.DeleteTaskLearning)

	// -- Study Session routes
	protected.GET("/study-sessions", handlers.GetStudySessions)
	protected.POST("/study-sessions", handlers.CreateStudySession)
	protected.GET("/study-sessions/stats", handlers.GetStudySessionStats)

	// -- Resource routes
	protected.GET("/resources", handlers.GetResources)
	protected.GET("/resources/:ID", handlers.GetResource)
	protected.POST("/resources", handlers.CreateResource)
	protected.PUT("/resources/:ID", handlers.UpdateResource)
	protected.DELETE("/resources/:ID", handlers.DeleteResource)

	// -- Learning Path routes
	protected.GET("/learning-paths", handlers.GetLearningPaths)
	protected.GET("/learning-paths/:ID", handlers.GetLearningPath)
	protected.POST("/learning-paths", handlers.CreateLearningPath)
	protected.DELETE("/learning-paths/:ID", handlers.DeleteLearningPath)

	// -- Tag routes
	protected.GET("/tags", handlers.GetTags)
	protected.POST("/tags", handlers.CreateTag)
	protected.PATCH("/tags/:ID", handlers.UpdateTag)
	protected.DELETE("/tags/:ID", handlers.DeleteTag)

	// Finance routes
	// -- Category routes
	protected.GET("/categories", handlers.GetBudgetCategories)
	protected.POST("/categories", handlers.CreateBudgetCategory)
	protected.PATCH("/categories/:ID", handlers.UpdateBudgetCategory)
	protected.DELETE("/categories/:ID", handlers.DeleteBudgetCategory)

	// -- Budget routes
	protected.GET("/budgets", handlers.GetBudgets)
	protected.POST("/budgets", handlers.CreateBudget)
	protected.PATCH("/budgets/:ID", handlers.UpdateBudget)
	protected.DELETE("/budgets/:ID/:incomeID", handlers.DeleteBudget)
	protected.GET("/budgets/analytics", handlers.GetBudgetAnalytics)
	protected.GET("/budgets/suggestions", handlers.GetBudgetSuggestions)
	protected.GET("/budgets/alerts", handlers.GetBudgetAlerts)

	// -- Income routes
	protected.GET("/incomes", handlers.GetIncomes)
	protected.POST("/incomes", handlers.CreateIncome)
	protected.PATCH("/incomes/:ID", handlers.UpdateIncome)
	protected.DELETE("/incomes/:ID", handlers.DeleteIncome)

	// -- Transaction routes
	protected.GET("/transactions", handlers.GetTransactions)
	protected.POST("/transactions", handlers.CreateTransaction)
	protected.PATCH("/transactions/:ID", handlers.UpdateTransaction)
	protected.DELETE("/transactions/:ID", handlers.DeleteTransaction)

	// Push notification routes
	protected.POST("/push/subscription", handlers.SubscribePush)
	protected.DELETE("/push/subscription", handlers.UnsubscribePush)
	protected.GET("/push/subscriptions", handlers.GetPushSubscriptions)
	protected.POST("/push/notification", handlers.SendPushNotification)
	protected.POST("/push/task-reminder/:task_id", handlers.SendTaskReminder)
	protected.POST("/push/goal-reminder/:goal_id", handlers.SendGoalReminder)
}
