package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/handlers"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDBSQLite()

	router := gin.Default()

	if os.Getenv("GIN_MODE") == "release" {
		log.Println("Entered production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("ALLOWED_URL")},
		AllowMethods:     []string{os.Getenv("ALLOWED_URL")},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Testing route
	router.GET("/ping", ping)

	// User routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	protected := router.Group("/")
	protected.Use(util.JWTAuthMiddleware())

	// Goal routes
	protected.GET("/goals", handlers.GetGoals)
	protected.GET("/goals/:ID", handlers.GetGoal)
	protected.POST("/goals", handlers.CreateGoal)
	protected.PUT("/goals/:ID", handlers.UpdateGoal)
	protected.DELETE("/goals/:ID", handlers.DeleteGoal)

	// Task routes
	protected.GET("/tasks", handlers.GetTasks)
	protected.GET("/tasks/:ID", handlers.GetTask)
	protected.POST("/tasks", handlers.CreateTask)
	protected.PATCH("/tasks/:ID", handlers.UpdateTask)
	protected.DELETE("/tasks/:ID", handlers.DeleteTask)

	// Schedule routes
	protected.GET("/schedule", handlers.GetSchedule)
	protected.POST("/schedule", handlers.CreateSchedule)
	protected.DELETE("/schedule/:ID", handlers.DeleteSchedule)

	// Learning routes
	protected.GET("/decks", handlers.GetDecks)
	protected.GET("/decks/:ID", handlers.GetDeck)
	protected.POST("/decks", handlers.CreateDeck)
	protected.PATCH("/decks/:ID", handlers.UpdateDeck)
	protected.DELETE("/decks/:ID", handlers.DeleteDeck)

	protected.GET("decks/cards/:deckID", handlers.GetCards)
	protected.GET("/cards/:ID", handlers.GetCard)
	protected.POST("/cards", handlers.CreateCard)
	protected.PATCH("/cards/:ID", handlers.UpdateCard)
	protected.DELETE("/cards/:ID", handlers.DeleteCard)

	protected.GET("/cards/review/:ID", handlers.ReviewCard)
	protected.GET("/cards/due/:deckID", handlers.GetDueCards)

	// Finance routes 
	protected.GET("/categories", handlers.GetBudgetCategories)
	protected.POST("/categories", handlers.CreateBudgetCategory)
	protected.PATCH("/categories/:ID", handlers.UpdateBudgetCategory)
	protected.DELETE("/categories/:ID", handlers.DeleteBudgetCategory)

	protected.GET("/budgets", handlers.GetBudgets)
	protected.POST("/budgets", handlers.CreateBudget)
	protected.PATCH("/budgets/:ID", handlers.UpdateBudget)
	protected.DELETE("/budgets/:ID", handlers.DeleteBudget)


	protected.GET("/incomes", handlers.GetIncomes)
	protected.POST("/incomes", handlers.CreateIncome)
	protected.PATCH("/incomes/:ID", handlers.UpdateIncome)
	protected.DELETE("/incomes/:ID", handlers.DeleteIncome)

	router.Run(os.Getenv("PORT"))
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
