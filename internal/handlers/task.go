package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// Task handlers

// Get all tasks
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	result := config.GetDB().Find(&tasks)

	if result.Error != nil {
		log.Fatal(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})

}

// Get a specific task
func GetTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Fatal("GetTask error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting task",
		})
		return
	}

	var task models.Task
	result := config.GetDB().First(&task, taskID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Task does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

// Create a task
func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// Update a specific task
func UpdateTask(c *gin.Context) {
	var task models.Task

	taskID := c.Param("ID")
	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     time.Time `json:"due_date"`
		Priority    int       `json:"priority"`
		Status      string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTask := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Priority:    &input.Priority,
		Status:      input.Status,
	}

	if err := config.GetDB().Model(&task).Updates(updatedTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)

}

// Delete a specific task
func DeleteTask(c *gin.Context) {
	var task models.Task

	taskID := c.Param("ID")
	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	if err := config.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)

}
