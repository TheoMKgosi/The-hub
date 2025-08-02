package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

func GetTaskLearnings(c *gin.Context) {
	var tasks []models.Task_learning
	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Printf("Invalid ID param: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	log.Printf("Fetching task learnings for topic_id=%d", id)
	if err := config.GetDB().Where("topic_id = ?", id).Find(&tasks).Error; err != nil {
		log.Printf("Error fetching task learnings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task learnings"})
		return
	}

	log.Printf("Successfully fetched %d task learnings for topic_id=%d", len(tasks), id)
	c.JSON(http.StatusOK, gin.H{"task_learnings": tasks})
}

func GetTaskLearning(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Printf("Invalid ID param: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var task models.Task_learning
	log.Printf("Fetching task learning with ID=%d", id)
	if err := config.GetDB().Preload("Resources").First(&task, id).Error; err != nil {
		log.Printf("Task_learning with ID=%d not found: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task_learning not found"})
		return
	}

	log.Printf("Successfully fetched task learning with ID=%d", id)
	c.JSON(http.StatusOK, task)
}

func CreateTaskLearning(c *gin.Context) {
	var input struct {
		TopicID uint   `json:"topic_id"`
		Title   string `json:"title"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input on CreateTaskLearning: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	task := models.Task_learning{
		TopicID: input.TopicID,
		Title:   input.Title,
	}

	log.Printf("Creating task learning: %+v", input)
	if err := config.GetDB().Create(&task).Error; err != nil {
		log.Printf("Failed to create Task_learning: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Task_learning"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func UpdateTaskLearning(c *gin.Context) {
	id := c.Param("ID")
	var task models.Task_learning

	log.Printf("Updating task learning with ID=%s", id)
	if err := config.GetDB().First(&task, id).Error; err != nil {
		log.Printf("Task_learning with ID=%s not found: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task_learning not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input on update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.GetDB().Model(&task).Updates(input).Error; err != nil {
		log.Printf("Failed to update Task_learning with ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Task_learning"})
		return
	}

	log.Printf("Successfully updated Task_learning with ID=%s", id)
	c.JSON(http.StatusOK, task)
}

func DeleteTaskLearning(c *gin.Context) {
	id := c.Param("ID")
	var task models.Task_learning

	log.Printf("Deleting task learning with ID=%s", id)
	if err := config.GetDB().First(&task, id).Error; err != nil {
		log.Printf("Task_learning with ID=%s not found: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task_learning not found"})
		return
	}

	if err := config.GetDB().Delete(&task).Error; err != nil {
		log.Printf("Failed to delete Task_learning with ID=%s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Task_learning"})
		return
	}

	log.Printf("Successfully deleted Task_learning with ID=%s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Task_learning deleted"})
}
