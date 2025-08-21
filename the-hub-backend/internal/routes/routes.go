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

	protected := router.Group("/")
	protected.Use(util.JWTAuthMiddleware())

	// Plan routes
	// -- Goal routes
	protected.GET("/goals", handlers.GetGoals)
	protected.GET("/goals/:ID", handlers.GetGoal)
	protected.POST("/goals", handlers.CreateGoal)
	protected.PUT("/goals/:ID", handlers.UpdateGoal)
	protected.DELETE("/goals/:ID", handlers.DeleteGoal)

	// -- Task routes
	protected.GET("/tasks", handlers.GetTasks)
	protected.GET("/tasks/:ID", handlers.GetTask)
	protected.POST("/tasks", handlers.CreateTask)
	protected.PUT("/tasks/reorder", handlers.ReorderTasks)
	protected.PATCH("/tasks/:ID", handlers.UpdateTask)
	protected.DELETE("/tasks/:ID", handlers.DeleteTask)

	// Time routes
	// -- Schedule routes
	protected.GET("/schedule", handlers.GetSchedule)
	protected.POST("/schedule", handlers.CreateSchedule)
	protected.DELETE("/schedule/:ID", handlers.DeleteSchedule)

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

	// -- Task Learning
	protected.GET("/task-learning/:ID", handlers.GetTaskLearnings)
	protected.POST("/task-learning", handlers.CreateTaskLearning)
	protected.PATCH("/task-learning/:ID", handlers.UpdateTaskLearning)
	protected.DELETE("/task-learning/:ID", handlers.DeleteTaskLearning)

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

	// -- Income routes
	protected.GET("/incomes", handlers.GetIncomes)
	protected.POST("/incomes", handlers.CreateIncome)
	protected.PATCH("/incomes/:ID", handlers.UpdateIncome)
	protected.DELETE("/incomes/:ID", handlers.DeleteIncome)
}

