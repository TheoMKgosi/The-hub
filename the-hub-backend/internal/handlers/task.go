package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	config.Logger.Infof("Fetching tasks for user ID: %v", userID)
	if err := config.GetDB().Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Error fetching tasks for user %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	config.Logger.Infof("Found %d tasks for user ID %v", len(tasks), userID)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	config.Logger.Infof("Fetching task ID: %d", taskID)
	var task models.Task
	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		config.Logger.Errorf("Task ID %d not found: %v", taskID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func CreateTask(c *gin.Context) {
	var input struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description"`
		Priority    *int       `json:"priority"`
		DueDate     *time.Time `json:"due_date"`
		GoalID      *uint      `json:"goal_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid task input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for task"})
		return
	}

	userID := c.MustGet("userID").(uint)
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		GoalID:      input.GoalID,
		UserID:      userID,
	}

	if err := config.GetDB().Create(&task).Error; err != nil {
		config.Logger.Errorf("Error creating task for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	config.Logger.Infof("Created task ID %d for user %d", task.ID, userID)
	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	taskID := c.Param("ID")
	var task models.Task

	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		config.Logger.Warnf("Task not found for update: ID %s", taskID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Priority    *int    `json:"priority"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for task ID %s: %v", taskID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	config.Logger.Infof("Updating task ID %s with data: %+v", taskID, updates)
	if err := config.GetDB().Model(&task).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update task ID %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated task ID %d: %v", task.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	config.Logger.Infof("Successfully updated task ID %d", task.ID)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	taskID := c.Param("ID")
	var task models.Task

	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		config.Logger.Warnf("Task not found for delete: ID %s", taskID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := config.GetDB().Delete(&task).Error; err != nil {
		config.Logger.Errorf("Failed to delete task ID %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	config.Logger.Infof("Deleted task ID %s", taskID)
	c.JSON(http.StatusOK, task)
}

func UpsertScheduledTask(task models.Task) error {
	db := config.GetDB()

	if task.DueDate == nil {
		config.Logger.Infof("Removing scheduled task for task ID %d", task.ID)
		return db.Where("task_id = ?", task.ID).Delete(&models.ScheduledTask{}).Error
	}

	start := *task.DueDate
	end := start.Add(time.Hour)

	var scheduled models.ScheduledTask
	err := db.Where("task_id = ?", task.ID).First(&scheduled).Error

	if err != nil {
		config.Logger.Infof("Creating new scheduled task for task ID %d", task.ID)
		scheduled = models.ScheduledTask{
			TaskID: task.ID,
			Title:  task.Title,
			Start:  start,
			End:    end,
			UserID: task.UserID,
		}
		return db.Create(&scheduled).Error
	}

	config.Logger.Infof("Updating scheduled task for task ID %d", task.ID)
	return db.Model(&scheduled).Updates(models.ScheduledTask{
		Title: task.Title,
		Start: start,
		End:   end,
	}).Error
}

