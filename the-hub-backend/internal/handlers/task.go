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
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist",
		})
		return
	}

	result := config.GetDB().Where("user_id = ?", userID).Find(&tasks)

	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
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
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist",
		})
		return
	}

	id := userID.(uint)

	var input struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description"`
		Priority    *int       `json:"priority"`
		DueDate     *time.Time `json:"due_date"`
		GoalID      *uint      `json:"goal_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create Task"})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		GoalID:      input.GoalID,
		UserID:      id,
	}

	if err := config.GetDB().Create(&task).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Task"})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// Update a specific task
func UpdateTask(c *gin.Context) {
	var task models.Task

	taskID := c.Param("ID")
	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		log.Println("Error ID: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Priority    *int    `json:"priority"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTask := map[string]interface{}{}
	if input.Title != nil {
		updatedTask["title"] = *input.Title
	}
	if input.Description != nil {
		updatedTask["description"] = *input.Description
	}
	if input.Priority != nil {
		updatedTask["priority"] = *input.Priority
	}
	if input.Status != nil {
		updatedTask["status"] = *input.Status
	}

	if err := config.GetDB().Model(&task).Updates(updatedTask).Error; err != nil {
		log.Println("Error updating task:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// reload the task to get updated due date
	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		log.Println("Error retrieving updated task:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving updated task"})
		return
	}

	// update or remove calendar event based on new due date
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
