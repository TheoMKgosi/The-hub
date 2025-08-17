package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// GetTasks godoc
// @Summary      Get all tasks
// @Description  Fetch tasks for the logged-in user with optional ordering
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_by  query     string  false  "Order by field (order, priority, due_date, created_at)"  default(order)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks [get]
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get query parameters for ordering
	orderBy := c.DefaultQuery("order_by", "order_index")
	sortDir := c.DefaultQuery("sort", "asc")

	// Validate order_by parameter
	validOrderFields := map[string]bool{
		"order_index": true,
		"priority":    true,
		"due_date":    true,
		"created_at":  true,
		"title":       true,
		"status":      true,
	}

	if !validOrderFields[orderBy] {
		config.Logger.Warnf("Invalid order_by parameter: %s", orderBy)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order_by parameter"})
		return
	}

	// Validate sort direction
	if sortDir != "asc" && sortDir != "desc" {
		config.Logger.Warnf("Invalid sort direction: %s", sortDir)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort direction. Use 'asc' or 'desc'"})
		return
	}

	orderClause := orderBy + " " + sortDir

	config.Logger.Infof("Fetching tasks for user ID: %v with order: %s", userID, orderClause)
	if err := config.GetDB().Where("user_id = ?", userID).Order(orderClause).Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Error fetching tasks for user %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	config.Logger.Infof("Found %d tasks for user ID %v", len(tasks), userID)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTask godoc
// @Summary      Get a specific task
// @Description  Fetch a specific task by ID for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Task ID"
// @Success      200  {object}  map[string]models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{ID} [get]
func GetTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	config.Logger.Infof("Fetching task ID: %d for user ID: %v", taskID, userID)
	var task models.Task
	// Ensure user can only access their own tasks
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		config.Logger.Errorf("Task ID %d not found for user %v: %v", taskID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved task ID %d for user %v", taskID, userID)
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required" example:"Complete project proposal"`
	Description string     `json:"description" example:"Finish the quarterly project proposal document"`
	Priority    *int       `json:"priority" example:"3"`
	DueDate     *time.Time `json:"due_date" example:"2024-12-31T23:59:59Z"`
	OrderIndex  *int       `json:"order" example:"1"`
	GoalID      *uint      `json:"goal_id" example:"1"`
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        task  body      CreateTaskRequest  true  "Task creation data"
// @Success      201   {object}  models.Task
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tasks [post]
func CreateTask(c *gin.Context) {
	var input CreateTaskRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid task input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for task", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// If no order is specified, set it to the next available position
	order := 0
	if input.OrderIndex != nil {
		order = *input.OrderIndex
	} else {
		// Get the highest order number and add 1
		var maxOrder int
		if err := config.GetDB().Model(&models.Task{}).Where("user_id = ?", userIDUint).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder).Error; err != nil {
			config.Logger.Warnf("Failed to get max order for user %d: %v", userIDUint, err)
		}
		order = maxOrder + 1
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		OrderIndex:  order,
		GoalID:      input.GoalID,
		UserID:      userIDUint,
	}

	config.Logger.Infof("Creating task for user %d: %s with order %d", userIDUint, input.Title, order)
	if err := config.GetDB().Create(&task).Error; err != nil {
		config.Logger.Errorf("Error creating task for user %d: %v", userIDUint, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	// Handle scheduled task if due date is provided
	if task.DueDate != nil {
		if err := UpsertScheduledTask(task); err != nil {
			config.Logger.Warnf("Failed to create scheduled task for task ID %d: %v", task.ID, err)
			// Don't return error as the main task was created successfully
		}
	}

	config.Logger.Infof("Successfully created task ID %d for user %d", task.ID, userIDUint)
	c.JSON(http.StatusCreated, task)
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title       *string    `json:"title" example:"Updated task title"`
	Description *string    `json:"description" example:"Updated task description"`
	Priority    *int       `json:"priority" example:"2"`
	Status      *string    `json:"status" example:"completed"`
	DueDate     *time.Time `json:"due_date" example:"2024-12-31T23:59:59Z"`
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Update a specific task by ID for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID    path      int                true  "Task ID"
// @Param        task  body      UpdateTaskRequest  true  "Task update data"
// @Success      200   {object}  models.Task
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tasks/{ID} [patch]
func UpdateTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param for update: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var task models.Task
	// Ensure user can only update their own tasks
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task not found for update: ID %d, User %v", taskID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input UpdateTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for task ID %d: %v", taskID, err)
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
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.DueDate != nil {
		updates["due_date"] = *input.DueDate
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for task update: ID %d", taskID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating task ID %d for user %v with data: %+v", taskID, userID, updates)
	if err := config.GetDB().Model(&task).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update task ID %d: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Reload the updated task
	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated task ID %d: %v", task.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	config.Logger.Infof("Successfully updated task ID %d for user %v", task.ID, userID)
	c.JSON(http.StatusOK, task)
}

// ReorderTasksRequest represents the request body for reordering tasks
type ReorderTasksRequest struct {
	TaskOrders []TaskOrderItem `json:"task_orders" binding:"required"`
}

type TaskOrderItem struct {
	TaskID uint `json:"task_id" binding:"required"`
	Order  int  `json:"order" binding:"required"`
}

// ReorderTasks godoc
// @Summary      Reorder multiple tasks
// @Description  Update the order of multiple tasks at once
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        reorder  body      ReorderTasksRequest  true  "Task reorder data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /tasks/reorder [put]
func ReorderTasks(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task reorder")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input ReorderTasksRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid reorder input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if len(input.TaskOrders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tasks to reorder"})
		return
	}

	// Start transaction
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	config.Logger.Infof("Reordering %d tasks for user %v", len(input.TaskOrders), userID)

	// Update each task's order
	for _, item := range input.TaskOrders {
		var task models.Task
		// Ensure user owns the task
		if err := tx.Where("id = ? AND user_id = ?", item.TaskID, userID).First(&task).Error; err != nil {
			tx.Rollback()
			config.Logger.Warnf("Task ID %d not found for user %v during reorder", item.TaskID, userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "One or more tasks not found"})
			return
		}

		if err := tx.Model(&task).Update("order_index", item.Order).Error; err != nil {
			tx.Rollback()
			config.Logger.Errorf("Failed to update order for task ID %d: %v", item.TaskID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder tasks"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		config.Logger.Errorf("Failed to commit task reorder transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task order"})
		return
	}

	config.Logger.Infof("Successfully reordered tasks for user %v", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Tasks reordered successfully"})
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a specific task by ID for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Task ID"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{ID} [delete]
func DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param for delete: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var task models.Task
	// Ensure user can only delete their own tasks
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task not found for delete: ID %d, User %v", taskID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	config.Logger.Infof("Deleting task ID %d for user %v", taskID, userID)
	if err := config.GetDB().Delete(&task).Error; err != nil {
		config.Logger.Errorf("Failed to delete task ID %d: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// Clean up scheduled task if it exists
	if err := config.GetDB().Where("task_id = ?", task.ID).Delete(&models.ScheduledTask{}).Error; err != nil {
		config.Logger.Warnf("Failed to delete scheduled task for task ID %d: %v", task.ID, err)
		// Don't return error as the main task was deleted successfully
	}

	config.Logger.Infof("Successfully deleted task ID %d for user %v", taskID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task": task})
}

// UpsertScheduledTask creates or updates a scheduled task based on the task's due date
func UpsertScheduledTask(task models.Task) error {
	db := config.GetDB()

	if task.DueDate == nil {
		config.Logger.Infof("Removing scheduled task for task ID %d (no due date)", task.ID)
		return db.Where("task_id = ?", task.ID).Delete(&models.ScheduledTask{}).Error
	}

	start := *task.DueDate
	end := start.Add(time.Hour)

	var scheduled models.ScheduledTask
	err := db.Where("task_id = ?", task.ID).First(&scheduled).Error

	if err != nil {
		// Create new scheduled task
		config.Logger.Infof("Creating new scheduled task for task ID %d", task.ID)
		scheduled = models.ScheduledTask{
			Title:  task.Title,
			Start:  start,
			End:    end,
			UserID: task.UserID,
		}
		if createErr := db.Create(&scheduled).Error; createErr != nil {
			config.Logger.Errorf("Failed to create scheduled task for task ID %d: %v", task.ID, createErr)
			return createErr
		}
		config.Logger.Infof("Successfully created scheduled task for task ID %d", task.ID)
		return nil
	}

	// Update existing scheduled task
	config.Logger.Infof("Updating existing scheduled task for task ID %d", task.ID)
	updateData := models.ScheduledTask{
		Title: task.Title,
		Start: start,
		End:   end,
	}
	if updateErr := db.Model(&scheduled).Updates(updateData).Error; updateErr != nil {
		config.Logger.Errorf("Failed to update scheduled task for task ID %d: %v", task.ID, updateErr)
		return updateErr
	}

	config.Logger.Infof("Successfully updated scheduled task for task ID %d", task.ID)
	return nil
}
