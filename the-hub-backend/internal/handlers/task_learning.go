package handlers

import (
	"net/http"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetTaskLearnings godoc
// @Summary      Get all task learnings for a topic
// @Description  Fetch task learnings for a specific topic ID with optional ordering
// @Tags         task-learnings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID        path      int     true   "Topic ID"
// @Param        order_by  query     string  false  "Order by field (title, created_at, updated_at)"  default(created_at)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.Task_learning
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /topics/{ID}/task-learnings [get]
func GetTaskLearnings(c *gin.Context) {
	var tasks []models.Task_learning
	topicIDStr := c.Param("ID")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid topic ID param: %s", topicIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify user owns the topic
	var topic models.Topic
	if err := config.GetDB().Where("id = ? AND user_id = ?", topicID, userID).First(&topic).Error; err != nil {
		config.Logger.Warnf("Topic ID %d not found or not owned by user %v", topicID, userID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Topic not found or access denied"})
		return
	}

	// Get query parameters for ordering
	orderBy := c.DefaultQuery("order_by", "created_at")
	sortDir := c.DefaultQuery("sort", "asc")

	// Validate order_by parameter
	validOrderFields := map[string]bool{
		"title":      true,
		"created_at": true,
		"updated_at": true,
		"status":     true,
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

	config.Logger.Infof("Fetching task learnings for topic ID: %d with order: %s", topicID, orderClause)
	if err := config.GetDB().Where("topic_id = ?", topicID).Order(orderClause).Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Error fetching task learnings for topic %d: %v", topicID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task learnings"})
		return
	}

	config.Logger.Infof("Found %d task learnings for topic ID %d", len(tasks), topicID)
	c.JSON(http.StatusOK, gin.H{"task_learnings": tasks})
}

// GetTaskLearning godoc
// @Summary      Get a specific task learning
// @Description  Fetch a specific task learning by ID for the logged-in user
// @Tags         task-learnings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Task Learning ID"
// @Success      200  {object}  models.Task_learning
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task-learnings/{ID} [get]
func GetTaskLearning(c *gin.Context) {
	taskLearningIDStr := c.Param("ID")
	taskLearningID, err := uuid.Parse(taskLearningIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task learning ID param: %s", taskLearningIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task learning ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var taskLearning models.Task_learning
	config.Logger.Infof("Fetching task learning ID: %d", taskLearningID)

	// First get the task learning with its topic to verify ownership
	if err := config.GetDB().Preload("Resources").Preload("Topic").First(&taskLearning, taskLearningID).Error; err != nil {
		config.Logger.Errorf("Task learning ID %d not found: %v", taskLearningID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task learning not found"})
		return
	}

	// Verify user owns the topic that contains this task learning
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if taskLearning.Topic.UserID != userIDUUID {
		config.Logger.Warnf("User %s attempted to access task learning %s owned by user %s", userIDUUID, taskLearningID, taskLearning.Topic.UserID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	config.Logger.Infof("Successfully retrieved task learning ID %d", taskLearningID)
	c.JSON(http.StatusOK, taskLearning)
}

// CreateTaskLearningRequest represents the request body for creating a task learning
type CreateTaskLearningRequest struct {
	TopicID uuid.UUID `json:"topic_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title   string    `json:"title" binding:"required" example:"Learn Go interfaces"`
	Status  string    `json:"status" example:"not_started"`
}

// CreateTaskLearning godoc
// @Summary      Create a new task learning
// @Description  Create a new task learning for a topic owned by the logged-in user
// @Tags         task-learnings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        task_learning  body      CreateTaskLearningRequest  true  "Task learning creation data"
// @Success      201  {object}  models.Task_learning
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task-learnings [post]
func CreateTaskLearning(c *gin.Context) {
	var input CreateTaskLearningRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid task learning input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for task learning", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task learning creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify user owns the topic
	var topic models.Topic
	if err := config.GetDB().Where("id = ? AND user_id = ?", input.TopicID, userIDUint).First(&topic).Error; err != nil {
		config.Logger.Warnf("Topic ID %d not found or not owned by user %d", input.TopicID, userIDUint)
		c.JSON(http.StatusForbidden, gin.H{"error": "Topic not found or access denied"})
		return
	}

	// Set default status if not provided
	status := "not_started"
	if input.Status != "" {
		// Validate status
		validStatuses := map[string]bool{
			"not_started": true,
			"in_progress": true,
			"completed":   true,
			"on_hold":     true,
		}
		if !validStatuses[input.Status] {
			config.Logger.Warnf("Invalid status: %s", input.Status)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid options: not_started, in_progress, completed, on_hold"})
			return
		}
		status = input.Status
	}

	taskLearning := models.Task_learning{
		TopicID: input.TopicID,
		Title:   input.Title,
		Status:  status,
	}

	config.Logger.Infof("Creating task learning for topic %d by user %d: %s", input.TopicID, userIDUint, input.Title)
	if err := config.GetDB().Create(&taskLearning).Error; err != nil {
		config.Logger.Errorf("Error creating task learning for user %d: %v", userIDUint, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task learning"})
		return
	}

	config.Logger.Infof("Successfully created task learning ID %d for user %d", taskLearning.ID, userIDUint)
	c.JSON(http.StatusCreated, taskLearning)
}

// UpdateTaskLearningRequest represents the request body for updating a task learning
type UpdateTaskLearningRequest struct {
	Title  *string `json:"title" example:"Updated task learning title"`
	Status *string `json:"status" example:"completed"`
}

// UpdateTaskLearning godoc
// @Summary      Update a task learning
// @Description  Update a specific task learning by ID for the logged-in user
// @Tags         task-learnings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID            path      int                        true  "Task Learning ID"
// @Param        task_learning body      UpdateTaskLearningRequest  true  "Task learning update data"
// @Success      200  {object}  models.Task_learning
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task-learnings/{ID} [put]
func UpdateTaskLearning(c *gin.Context) {
	taskLearningIDStr := c.Param("ID")
	taskLearningID, err := uuid.Parse(taskLearningIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task learning ID param for update: %s", taskLearningIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task learning ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task learning update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var taskLearning models.Task_learning
	// Get task learning with its topic to verify ownership
	if err := config.GetDB().Preload("Topic").First(&taskLearning, taskLearningID).Error; err != nil {
		config.Logger.Warnf("Task learning not found for update: ID %d", taskLearningID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task learning not found"})
		return
	}

	// Verify user owns the topic that contains this task learning
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if taskLearning.Topic.UserID != userIDUUID {
		config.Logger.Warnf("User %s attempted to update task learning %s owned by user %s", userIDUUID, taskLearningID, taskLearning.Topic.UserID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var input UpdateTaskLearningRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for task learning ID %d: %v", taskLearningID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Status != nil {
		// Validate status
		validStatuses := map[string]bool{
			"not_started": true,
			"in_progress": true,
			"completed":   true,
			"on_hold":     true,
		}
		if !validStatuses[*input.Status] {
			config.Logger.Warnf("Invalid status: %s", *input.Status)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid options: not_started, in_progress, completed, on_hold"})
			return
		}
		updates["status"] = *input.Status
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for task learning update: ID %d", taskLearningID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating task learning ID %d for user %d with data: %+v", taskLearningID, userIDUint, updates)
	if err := config.GetDB().Model(&taskLearning).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update task learning ID %d: %v", taskLearningID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task learning"})
		return
	}

	// Reload the updated task learning
	if err := config.GetDB().Preload("Resources").First(&taskLearning, taskLearning.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated task learning ID %d: %v", taskLearning.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task learning"})
		return
	}

	config.Logger.Infof("Successfully updated task learning ID %d for user %d", taskLearning.ID, userIDUint)
	c.JSON(http.StatusOK, taskLearning)
}

// DeleteTaskLearning godoc
// @Summary      Delete a task learning
// @Description  Delete a specific task learning by ID for the logged-in user
// @Tags         task-learnings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Task Learning ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task-learnings/{ID} [delete]
func DeleteTaskLearning(c *gin.Context) {
	taskLearningIDStr := c.Param("ID")
	taskLearningID, err := uuid.Parse(taskLearningIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task learning ID param for delete: %s", taskLearningIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task learning ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task learning deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var taskLearning models.Task_learning
	// Get task learning with its topic to verify ownership
	if err := config.GetDB().Preload("Topic").First(&taskLearning, taskLearningID).Error; err != nil {
		config.Logger.Warnf("Task learning not found for delete: ID %d", taskLearningID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task learning not found"})
		return
	}

	// Verify user owns the topic that contains this task learning
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if taskLearning.Topic.UserID != userIDUUID {
		config.Logger.Warnf("User %s attempted to delete task learning %s owned by user %s", userIDUUID, taskLearningID, taskLearning.Topic.UserID)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	config.Logger.Infof("Deleting task learning ID %s for user %s", taskLearningID, userIDUUID)

	// Start transaction to handle cascading deletes
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete associated resources first (if any)
	if err := tx.Where("task_learning_id = ?", taskLearningID).Delete(&models.Resource{}).Error; err != nil {
		tx.Rollback()
		config.Logger.Errorf("Failed to delete resources for task learning ID %d: %v", taskLearningID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task learning resources"})
		return
	}

	// Delete the task learning
	if err := tx.Delete(&taskLearning).Error; err != nil {
		tx.Rollback()
		config.Logger.Errorf("Failed to delete task learning ID %d: %v", taskLearningID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task learning"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		config.Logger.Errorf("Failed to commit task learning deletion transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete deletion"})
		return
	}

	config.Logger.Infof("Successfully deleted task learning ID %d for user %d", taskLearningID, userIDUint)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Task learning deleted successfully",
		"task_learning": taskLearning,
	})
}
