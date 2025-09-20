package handlers

import (
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/ai"
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
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	config.Logger.Infof("Fetching goals for user ID: %s", userIDUUID)

	var goals []models.Goal
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Find(&goals).Error; err != nil {
		config.Logger.Errorf("Error fetching goals for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch goals",
		})
		return
	}

	// Calculate progress for each goal
	for i := range goals {
		if err := goals[i].CalculateProgress(config.GetDB()); err != nil {
			config.Logger.Warnf("Failed to calculate progress for goal %s: %v", goals[i].ID, err)
		}
	}

	config.Logger.Infof("Found %d goals for user ID %s", len(goals), userIDUUID)
	c.JSON(http.StatusOK, gin.H{
		"goals": goals,
	})

}

func GetGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid goal ID",
		})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	config.Logger.Infof("Fetching goal ID: %s for user ID: %s", goalID, userIDUUID)

	var goal models.Goal
	if err := config.GetDB().Preload("Tasks").Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Errorf("Goal ID %s not found for user %s: %v", goalID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	config.Logger.Infof("Successfully retrieved goal ID %s for user %s", goalID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{
		"goal": goal,
	})

}

// CreateGoalRequest represents the request body for creating a goal
type CreateGoalRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Priority    *int       `json:"priority"`
	Category    string     `json:"category"`
	Color       string     `json:"color"`
}

func CreateGoal(c *gin.Context) {
	var input CreateGoalRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid goal input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for goal", "details": err.Error()})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context during goal creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	goal := models.Goal{
		UserID:      userIDUUID,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Priority:    input.Priority,
		Category:    input.Category,
		Color:       input.Color,
		Status:      "active", // Default status
	}

	config.Logger.Infof("Creating goal for user %s: %s", userIDUUID, input.Title)
	if err := config.GetDB().Create(&goal).Error; err != nil {
		config.Logger.Errorf("Error creating goal for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create goal"})
		return
	}

	// Calculate initial progress (should be 0 for new goal)
	if err := goal.CalculateProgress(config.GetDB()); err != nil {
		config.Logger.Warnf("Failed to calculate progress for new goal %s: %v", goal.ID, err)
	}

	config.Logger.Infof("Successfully created goal ID %s for user %s", goal.ID, userIDUUID)
	c.JSON(http.StatusCreated, goal)
}

// UpdateGoalRequest represents the request body for updating a goal
type UpdateGoalRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Priority    *int       `json:"priority"`
	Status      *string    `json:"status"`
	Category    *string    `json:"category"`
	Color       *string    `json:"color"`
}

func UpdateGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param for update: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context during goal update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal not found for update: ID %s, User %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	var input UpdateGoalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for goal ID %s: %v", goalID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.DueDate != nil {
		updates["due_date"] = *input.DueDate
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for goal update: ID %s", goalID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating goal ID %s for user %s with data: %+v", goalID, userIDUUID, updates)
	if err := config.GetDB().Model(&goal).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update goal ID %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goal"})
		return
	}

	// Reload the updated goal
	if err := config.GetDB().First(&goal, goal.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated goal ID %s: %v", goal.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated goal"})
		return
	}

	// Recalculate progress after update
	if err := goal.CalculateProgress(config.GetDB()); err != nil {
		config.Logger.Warnf("Failed to calculate progress for updated goal %s: %v", goal.ID, err)
	}

	config.Logger.Infof("Successfully updated goal ID %s for user %s", goal.ID, userIDUUID)
	c.JSON(http.StatusOK, goal)

}

func DeleteGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param for delete: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context during goal deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal not found for delete: ID %s, User %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	config.Logger.Infof("Deleting goal ID %s for user %s", goalID, userIDUUID)
	if err := config.GetDB().Delete(&goal).Error; err != nil {
		config.Logger.Errorf("Failed to delete goal ID %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goal"})
		return
	}

	config.Logger.Infof("Successfully deleted goal ID %s for user %s", goalID, userIDUUID)
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
		config.Logger.Warnf("Invalid goal ID param: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal ID %s not found for user %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	var input AddTaskToGoalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid task input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
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
			config.Logger.Warnf("Failed to get max order for goal %s: %v", goalID, err)
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

	config.Logger.Infof("Creating task for goal %s by user %s: %s", goalID, userIDUUID, input.Title)
	if err := config.GetDB().Create(&task).Error; err != nil {
		config.Logger.Errorf("Error creating task for goal %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	// Recalculate goal progress after adding task
	if err := goal.CalculateProgress(config.GetDB()); err != nil {
		config.Logger.Warnf("Failed to recalculate progress for goal %s: %v", goal.ID, err)
	} else {
		// Update goal in database
		config.GetDB().Model(&goal).Updates(map[string]interface{}{
			"progress":        goal.Progress,
			"total_tasks":     goal.TotalTasks,
			"completed_tasks": goal.CompletedTasks,
			"status":          goal.Status,
		})
	}

	config.Logger.Infof("Successfully created task ID %s for goal %s", task.ID, goalID)
	c.JSON(http.StatusCreated, task)
}

// GetGoalTasks retrieves all tasks for a specific goal
func GetGoalTasks(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal ID %s not found for user %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	config.Logger.Infof("Fetching tasks for goal ID: %s for user ID: %s", goalID, userIDUUID)

	var tasks []models.Task
	if err := config.GetDB().Where("goal_id = ? AND user_id = ?", goalID, userIDUUID).Order("order_index").Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Error fetching tasks for goal %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	config.Logger.Infof("Found %d tasks for goal ID %s", len(tasks), goalID)
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
		config.Logger.Warnf("Invalid goal ID param: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	taskIDStr := c.Param("taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal ID %s not found for user %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Verify the task exists, belongs to the goal, and belongs to the user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND goal_id = ? AND user_id = ?", taskID, goalID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found in goal %s for user %s", taskID, goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this goal"})
		return
	}

	var input UpdateGoalTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for task ID %s: %v", taskID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
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
		config.Logger.Warnf("No valid fields provided for task update: ID %s", taskID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating task ID %s in goal %s for user %s", taskID, goalID, userIDUUID)
	if err := config.GetDB().Model(&task).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update task ID %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Reload the updated task
	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated task ID %s: %v", task.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	config.Logger.Infof("Successfully updated task ID %s for user %s", task.ID, userIDUUID)
	c.JSON(http.StatusOK, task)
}

// DeleteGoalTask deletes a specific task from a goal
func DeleteGoalTask(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	taskIDStr := c.Param("taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal ID %s not found for user %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Verify the task exists, belongs to the goal, and belongs to the user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND goal_id = ? AND user_id = ?", taskID, goalID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found in goal %s for user %s", taskID, goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this goal"})
		return
	}

	config.Logger.Infof("Deleting task ID %s from goal %s for user %s", taskID, goalID, userIDUUID)
	if err := config.GetDB().Delete(&task).Error; err != nil {
		config.Logger.Errorf("Failed to delete task ID %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	config.Logger.Infof("Successfully deleted task ID %s for user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task": task})
}

// CompleteGoalTask marks a task as completed
func CompleteGoalTask(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	taskIDStr := c.Param("taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal ID %s not found for user %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	// Verify the task exists, belongs to the goal, and belongs to the user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND goal_id = ? AND user_id = ?", taskID, goalID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found in goal %s for user %s", taskID, goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this goal"})
		return
	}

	// Toggle completion status
	newStatus := "completed"
	if task.Status == "completed" {
		newStatus = "pending"
	}

	config.Logger.Infof("Toggling task ID %s status to %s for user %s", taskID, newStatus, userIDUUID)
	if err := config.GetDB().Model(&task).Update("status", newStatus).Error; err != nil {
		config.Logger.Errorf("Failed to update task status for task ID %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	// Reload the updated task
	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated task ID %s: %v", task.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	// Recalculate goal progress after task status change
	if err := goal.CalculateProgress(config.GetDB()); err != nil {
		config.Logger.Warnf("Failed to recalculate progress for goal %s: %v", goal.ID, err)
	} else {
		// Update goal in database
		config.GetDB().Model(&goal).Updates(map[string]interface{}{
			"progress":        goal.Progress,
			"total_tasks":     goal.TotalTasks,
			"completed_tasks": goal.CompletedTasks,
			"status":          goal.Status,
		})
	}

	config.Logger.Infof("Successfully updated task ID %s status for user %s", task.ID, userIDUUID)
	c.JSON(http.StatusOK, task)
}

// GetGoalTaskRecommendations returns AI-generated task recommendations for a goal
func GetGoalTaskRecommendations(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify the goal exists and belongs to the user
	var goal models.Goal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal ID %s not found for user %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	config.Logger.Infof("Generating AI recommendations for goal ID: %s for user ID: %s", goalID, userIDUUID)

	// Generate AI recommendations
	recommendations, err := ai.GenerateGoalTaskRecommendations(goalID, userIDUUID)
	if err != nil {
		config.Logger.Errorf("Error generating AI recommendations for goal %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate recommendations"})
		return
	}

	config.Logger.Infof("Successfully generated %d recommendations for goal ID %s", len(recommendations), goalID)
	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"goal":            goal,
	})
}
