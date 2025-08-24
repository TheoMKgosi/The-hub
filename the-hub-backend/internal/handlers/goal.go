package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Goal handlers
func GetGoals(c *gin.Context) {
	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var goals []models.Goal
	if err := config.GetDB().Preload("Tasks").Where("user_id = ?", userIDUUID).Find(&goals).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"goals": goals,
	})

}

func GetGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid goal ID",
		})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var goal models.Goal
	if err := config.GetDB().Preload("Tasks").Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"goal": goal,
	})

}

func CreateGoal(c *gin.Context) {
	var goal models.Goal

	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set the user ID from authentication
	goal.UserID = userID.(uuid.UUID)

	if err := config.GetDB().Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

func UpdateGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := config.GetDB().Model(&goal).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)

}

func DeleteGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	if err := config.GetDB().Delete(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully", "goal": goal})

}

// AddTaskToGoalRequest represents the request body for adding a task to a goal
type AddTaskToGoalRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    *int       `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	OrderIndex  *int       `json:"order"`
}

// AddTaskToGoal adds a new task to a specific goal
func AddTaskToGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	var input AddTaskToGoalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If no order is specified, set it to the next available position for this goal
	order := 0
	if input.OrderIndex != nil {
		order = *input.OrderIndex
	} else {
		// Get the highest order number for tasks in this goal and add 1
		var maxOrder int
		if err := config.GetDB().Model(&models.Task{}).Where("goal_id = ?", goalID).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder).Error; err != nil {
			log.Printf("Failed to get max order for goal %s: %v", goalID, err)
		}
		order = maxOrder + 1
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		OrderIndex:  order,
		GoalID:      &goalID,
		UserID:      userIDUUID,
	}

	if err := config.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetGoalTasks retrieves all tasks for a specific goal
func GetGoalTasks(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	var tasks []models.Task
	if err := config.GetDB().Where("goal_id = ? AND user_id = ?", goalID, userIDUUID).Order("order_index").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// UpdateGoalTaskRequest represents the request body for updating a task in a goal
type UpdateGoalTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    *int       `json:"priority"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateGoalTask updates a specific task within a goal
func UpdateGoalTask(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	taskIDStr := c.Param("taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Verify the task exists, belongs to the goal, and belongs to the user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND goal_id = ? AND user_id = ?", taskID, goalID, userIDUUID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this goal"})
		return
	}

	var input UpdateGoalTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.DueDate != nil {
		updates["due_date"] = *input.DueDate
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := config.GetDB().Model(&task).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload the updated task
	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteGoalTask deletes a specific task from a goal
func DeleteGoalTask(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	taskIDStr := c.Param("taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Verify the task exists, belongs to the goal, and belongs to the user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND goal_id = ? AND user_id = ?", taskID, goalID, userIDUUID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this goal"})
		return
	}

	if err := config.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task": task})
}

// CompleteGoalTask marks a task as completed
func CompleteGoalTask(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	taskIDStr := c.Param("taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Verify the task exists, belongs to the goal, and belongs to the user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND goal_id = ? AND user_id = ?", taskID, goalID, userIDUUID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this goal"})
		return
	}

	// Toggle completion status
	newStatus := "completed"
	if task.Status == "completed" {
		newStatus = "pending"
	}

	if err := config.GetDB().Model(&task).Update("status", newStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload the updated task
	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	c.JSON(http.StatusOK, task)
}
