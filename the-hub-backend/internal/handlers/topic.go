package handlers

import (
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetTopics godoc
// @Summary      Get all topics
// @Description  Fetch topics for the logged-in user with optional ordering and filtering
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_by  query     string  false  "Order by field (title, created_at, updated_at, deadline, status)"  default(created_at)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(desc)
// @Param        status    query     string  false  "Filter by status"
// @Param        tag       query     string  false  "Filter by tag name"
// @Success      200  {object}  map[string][]models.Topic
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /topics [get]
func GetTopics(c *gin.Context) {
	var topics []models.Topic
	userID, exist := c.Get("userID")
	if !exist {
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

	// Get query parameters
	orderBy := c.DefaultQuery("order_by", "created_at")
	sortDir := c.DefaultQuery("sort", "desc")
	statusFilter := c.Query("status")
	tagFilter := c.Query("tag")

	// Validate order_by parameter
	validOrderFields := map[string]bool{
		"title":      true,
		"created_at": true,
		"updated_at": true,
		"deadline":   true,
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

	// Build query
	query := config.GetDB().Where("user_id = ?", userIDUUID).Preload("Tags")

	// Apply status filter
	if statusFilter != "" {
		validStatuses := map[string]bool{
			"not_started": true,
			"in_progress": true,
			"completed":   true,
			"on_hold":     true,
		}
		if !validStatuses[statusFilter] {
			config.Logger.Warnf("Invalid status filter: %s", statusFilter)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status filter"})
			return
		}
		query = query.Where("status = ?", statusFilter)
	}

	// Apply tag filter
	if tagFilter != "" {
		query = query.Joins("JOIN topic_tags ON topics.id = topic_tags.topic_id").
			Joins("JOIN tags ON topic_tags.tag_id = tags.id").
			Where("tags.name = ? AND tags.user_id = ?", tagFilter, userIDUUID)
	}

	config.Logger.Infof("Fetching topics for user ID: %s with order: %s, status: %s, tag: %s", userIDUUID, orderClause, statusFilter, tagFilter)

	if err := query.Order(orderClause).Find(&topics).Error; err != nil {
		config.Logger.Errorf("Error fetching topics for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve topics"})
		return
	}

	config.Logger.Infof("Found %d topics for user ID %s", len(topics), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

// GetTopic godoc
// @Summary      Get a specific topic
// @Description  Fetch a specific topic by ID for the logged-in user
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Topic ID"
// @Success      200  {object}  models.Topic
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /topics/{ID} [get]
func GetTopic(c *gin.Context) {
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

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	config.Logger.Infof("Fetching topic ID: %s for user ID: %s", topicID, userIDUUID)
	var topic models.Topic
	// Ensure user can only access their own topics
	if err := config.GetDB().Where("id = ? AND user_id = ?", topicID, userIDUUID).Preload("Tags").First(&topic).Error; err != nil {
		config.Logger.Errorf("Topic ID %s not found for user %s: %v", topicID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved topic ID %s for user %s", topicID, userIDUUID)
	c.JSON(http.StatusOK, topic)
}

// CreateTopicRequest represents the request body for creating a topic
type CreateTopicRequest struct {
	Title          string      `json:"title" binding:"required" example:"Learn Go Programming"`
	Description    string      `json:"description" example:"Comprehensive study of Go programming language"`
	Status         string      `json:"status" example:"not_started"`
	EstimatedHours *int        `json:"estimated_hours" example:"40"`
	Deadline       *time.Time  `json:"deadline" example:"2024-12-31T23:59:59Z"`
	TagIDs         []uuid.UUID `json:"tag_ids" example:"[\"550e8400-e29b-41d4-a716-446655440000\",\"550e8400-e29b-41d4-a716-446655440001\"]"`
}

// CreateTopic godoc
// @Summary      Create a new topic
// @Description  Create a new topic for the logged-in user
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        topic  body      CreateTopicRequest  true  "Topic creation data"
// @Success      201  {object}  models.Topic
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /topics [post]
func CreateTopic(c *gin.Context) {
	var input CreateTopicRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid topic input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for topic", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during topic creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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

	// Check if topic title already exists for this user
	var existingTopic models.Topic
	if err := config.GetDB().Where("title = ? AND user_id = ?", input.Title, userIDUUID).First(&existingTopic).Error; err == nil {
		config.Logger.Warnf("Topic title '%s' already exists for user %s", input.Title, userIDUUID)
		c.JSON(http.StatusConflict, gin.H{"error": "Topic with this title already exists"})
		return
	}

	topic := models.Topic{
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		Deadline:    input.Deadline,
		UserID:      userIDUUID,
	}

	// if input.EstimatedHours != nil {
	// 	if *input.EstimatedHours < 0 {
	// 		config.Logger.Warnf("Invalid estimated hours: %d", *input.EstimatedHours)
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Estimated hours must be non-negative"})
	// 		return
	// 	}
	// 	topic.EstimatedHours = *input.EstimatedHours
	// }

	// Handle tags if provided
	if len(input.TagIDs) > 0 {
		var tags []models.Tag
		// Ensure user owns all the tags they're trying to assign
		if err := config.GetDB().Where("id IN ? AND user_id = ?", input.TagIDs, userIDUUID).Find(&tags).Error; err != nil {
			config.Logger.Errorf("Error loading tags for user %s: %v", userIDUUID, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag IDs or tags not found"})
			return
		}

		if len(tags) != len(input.TagIDs) {
			config.Logger.Warnf("Not all tags found or owned by user %s. Expected %d, found %d", userIDUUID, len(input.TagIDs), len(tags))
			c.JSON(http.StatusForbidden, gin.H{"error": "Some tags not found or access denied"})
			return
		}

		topic.Tags = tags
		config.Logger.Infof("Loaded %d tags for new topic", len(tags))
	}

	config.Logger.Infof("Creating topic for user %s: %s", userIDUUID, input.Title)
	if err := config.GetDB().Create(&topic).Error; err != nil {
		config.Logger.Errorf("Error creating topic for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create topic"})
		return
	}

	// Reload topic with tags for response
	if err := config.GetDB().Preload("Tags").First(&topic, topic.ID).Error; err != nil {
		config.Logger.Warnf("Failed to reload topic with tags: %v", err)
	}

	config.Logger.Infof("Successfully created topic ID %s for user %s", topic.ID, userIDUUID)
	c.JSON(http.StatusCreated, topic)
}

// UpdateTopicRequest represents the request body for updating a topic
type UpdateTopicRequest struct {
	Title          *string      `json:"title" example:"Updated Topic Title"`
	Description    *string      `json:"description" example:"Updated description"`
	Status         *string      `json:"status" example:"in_progress"`
	EstimatedHours *int         `json:"estimated_hours" example:"50"`
	Deadline       *time.Time   `json:"deadline" example:"2024-12-31T23:59:59Z"`
	TagIDs         *[]uuid.UUID `json:"tag_ids" example:"[\"550e8400-e29b-41d4-a716-446655440000\",\"550e8400-e29b-41d4-a716-446655440001\"]"`
}

// UpdateTopic godoc
// @Summary      Update a topic
// @Description  Update a specific topic by ID for the logged-in user
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID    path      int                 true  "Topic ID"
// @Param        topic body      UpdateTopicRequest  true  "Topic update data"
// @Success      200  {object}  models.Topic
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /topics/{ID} [put]
func UpdateTopic(c *gin.Context) {
	topicIDStr := c.Param("ID")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid topic ID param for update: %s", topicIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during topic update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var topic models.Topic
	// Ensure user can only update their own topics
	if err := config.GetDB().Where("id = ? AND user_id = ?", topicID, userIDUUID).First(&topic).Error; err != nil {
		config.Logger.Warnf("Topic not found for update: ID %s, User %s", topicID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	var input UpdateTopicRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for topic ID %d: %v", topicID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		// Check for duplicate title (excluding current topic)
		if *input.Title != topic.Title {
			var existingTopic models.Topic
			if err := config.GetDB().Where("title = ? AND user_id = ? AND id != ?", *input.Title, userIDUUID, topicID).First(&existingTopic).Error; err == nil {
				config.Logger.Warnf("Topic title '%s' already exists for user %s", *input.Title, userIDUUID)
				c.JSON(http.StatusConflict, gin.H{"error": "Topic with this title already exists"})
				return
			}
		}
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
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
	if input.EstimatedHours != nil {
		if *input.EstimatedHours < 0 {
			config.Logger.Warnf("Invalid estimated hours: %d", *input.EstimatedHours)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Estimated hours must be non-negative"})
			return
		}
		updates["estimated_hours"] = *input.EstimatedHours
	}
	if input.Deadline != nil {
		updates["deadline"] = *input.Deadline
	}

	// Start transaction for atomic updates
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update basic fields
	if len(updates) > 0 {
		config.Logger.Infof("Updating topic ID %s for user %s with data: %+v", topicID, userIDUUID, updates)
		if err := tx.Model(&topic).Updates(updates).Error; err != nil {
			tx.Rollback()
			config.Logger.Errorf("Failed to update topic ID %d: %v", topicID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update topic"})
			return
		}
	}

	// Handle tag updates
	if input.TagIDs != nil {
		if len(*input.TagIDs) > 0 {
			var tags []models.Tag
			// Ensure user owns all the tags they're trying to assign
			if err := tx.Where("id IN ? AND user_id = ?", *input.TagIDs, userIDUUID).Find(&tags).Error; err != nil {
				tx.Rollback()
				config.Logger.Errorf("Error loading tags for user %s: %v", userIDUUID, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag IDs or tags not found"})
				return
			}

			if len(tags) != len(*input.TagIDs) {
				tx.Rollback()
				config.Logger.Warnf("Not all tags found or owned by user %s. Expected %d, found %d", userIDUUID, len(*input.TagIDs), len(tags))
				c.JSON(http.StatusForbidden, gin.H{"error": "Some tags not found or access denied"})
				return
			}

			if err := tx.Model(&topic).Association("Tags").Replace(&tags); err != nil {
				tx.Rollback()
				config.Logger.Errorf("Failed to update tags for topic ID %d: %v", topicID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update topic tags"})
				return
			}
			config.Logger.Infof("Updated tags for topic ID %d", topicID)
		} else {
			// Clear all tags
			if err := tx.Model(&topic).Association("Tags").Clear(); err != nil {
				tx.Rollback()
				config.Logger.Errorf("Failed to clear tags for topic ID %d: %v", topicID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear topic tags"})
				return
			}
			config.Logger.Infof("Cleared all tags for topic ID %d", topicID)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		config.Logger.Errorf("Failed to commit topic update transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save topic updates"})
		return
	}

	// Reload topic with tags for response
	if err := config.GetDB().Preload("Tags").First(&topic, topic.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated topic ID %d: %v", topic.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated topic"})
		return
	}

	config.Logger.Infof("Successfully updated topic ID %s for user %s", topic.ID, userIDUUID)
	c.JSON(http.StatusOK, topic)
}

// DeleteTopic godoc
// @Summary      Delete a topic
// @Description  Delete a specific topic by ID for the logged-in user
// @Tags         topics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Topic ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /topics/{ID} [delete]
func DeleteTopic(c *gin.Context) {
	topicIDStr := c.Param("ID")
	topicID, err := uuid.Parse(topicIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid topic ID param for delete: %s", topicIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during topic deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var topic models.Topic
	// Ensure user can only delete their own topics
	if err := config.GetDB().Where("id = ? AND user_id = ?", topicID, userIDUUID).First(&topic).Error; err != nil {
		config.Logger.Warnf("Topic not found for delete: ID %s, User %s", topicID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	// Check if topic has associated task learnings
	var taskLearningCount int64
	if err := config.GetDB().Model(&models.Task_learning{}).Where("topic_id = ?", topicID).Count(&taskLearningCount).Error; err != nil {
		config.Logger.Errorf("Error checking task learnings for topic ID %d: %v", topicID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify topic usage"})
		return
	}

	if taskLearningCount > 0 {
		config.Logger.Warnf("Attempted to delete topic ID %d which has %d task learnings", topicID, taskLearningCount)
		c.JSON(http.StatusConflict, gin.H{
			"error":               "Cannot delete topic that has associated task learnings",
			"task_learning_count": taskLearningCount,
		})
		return
	}

	config.Logger.Infof("Deleting topic ID %s for user %s", topicID, userIDUUID)

	// Start transaction for cascading deletes
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Clear tag associations
	if err := tx.Model(&topic).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		config.Logger.Errorf("Failed to clear tag associations for topic ID %d: %v", topicID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear topic associations"})
		return
	}

	// Delete the topic
	if err := tx.Delete(&topic).Error; err != nil {
		tx.Rollback()
		config.Logger.Errorf("Failed to delete topic ID %d: %v", topicID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete topic"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		config.Logger.Errorf("Failed to commit topic deletion transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete deletion"})
		return
	}

	config.Logger.Infof("Successfully deleted topic ID %s for user %s", topicID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Topic deleted successfully",
		"topic":   topic,
	})
}
