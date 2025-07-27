package handlers

import (
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

func GetTaskLearnings(c *gin.Context) {
	var tasks []models.Task_learning
	if err := config.GetDB().Preload("Resources").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task learnings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task_learnings": tasks})
}

func GetTaskLearning(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var task models.Task_learning
	if err := config.GetDB().Preload("Resources").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task_learning not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func CreateTaskLearning(c *gin.Context) {
	var input models.Task_learning
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.GetDB().Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Task_learning"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func UpdateTaskLearning(c *gin.Context) {
	id := c.Param("ID")
	var task models.Task_learning

	if err := config.GetDB().First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task_learning not found"})
		return
	}

	var input models.Task_learning
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.GetDB().Model(&task).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Task_learning"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTaskLearning(c *gin.Context) {
	id := c.Param("ID")
	var task models.Task_learning

	if err := config.GetDB().First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task_learning not found"})
		return
	}

	if err := config.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Task_learning"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task_learning deleted"})
}

