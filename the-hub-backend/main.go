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

	router.Run(":8080")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
