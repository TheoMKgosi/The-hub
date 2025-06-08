package main

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/handlers"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, 
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// Testing route
	router.GET("/ping", ping)

	// Goal routes
	router.GET("/goals", handlers.GetGoals)
	router.GET("/goals/:ID", handlers.GetGoal)
	router.POST("/goals", handlers.CreateGoal)
	router.PUT("/goals/:ID", handlers.UpdateGoal)
	router.DELETE("/goals/:ID", handlers.DeleteGoal)

	// Task routes
	router.GET("/tasks", handlers.GetTasks)
	router.GET("/tasks/:ID", handlers.GetTask)
	router.POST("/tasks", handlers.CreateTask)
	router.PUT("/tasks/:ID", handlers.UpdateTask)
	router.DELETE("/tasks/:ID", handlers.DeleteTask)

	router.Run(":8080")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
