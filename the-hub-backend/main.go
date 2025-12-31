package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/TheoMKgosi/The-hub/docs"
	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitLogger()
	defer config.Logger.Sync()

	if err := config.InitDBManager(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Add health check
	if err := config.GetDBManager().HealthCheck(context.Background()); err != nil {
		log.Fatal("Database health check failed:", err)
	}

	router := gin.Default()

	if os.Getenv("GIN_MODE") == "release" {
		log.Println("Entered production mode")
		gin.SetMode(gin.ReleaseMode)
	}

	rawUrl := os.Getenv("ALLOWED_URL")
	allowedOrigins := strings.Split(rawUrl, ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
