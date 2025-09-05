package handlers

import (
	"net/http"
	"strings"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateResourceRequest represents the request body for creating a resource
type CreateResourceRequest struct {
	TopicID *uuid.UUID `json:"topic_id" binding:"omitempty"`
	TaskID  *uuid.UUID `json:"task_id" binding:"omitempty"`
	Title   string     `json:"title" binding:"required"`
	Link    string     `json:"link"`
	Type    string     `json:"type" binding:"required,oneof=video article document book course"`
	Notes   string     `json:"notes"`
}

// CreateResource godoc
// @Summary      Create a new resource
// @Description  Create a new learning resource for a topic or task
// @Tags         resources
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        resource  body      CreateResourceRequest  true  "Resource creation data"
// @Success      201  {object}  models.Resource
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /resources [post]
func CreateResource(c *gin.Context) {
	var input CreateResourceRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid resource input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for resource", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during resource creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Validate that topic exists and belongs to user (if provided)
	if input.TopicID != nil {
		var topic models.Topic
		if err := config.GetDB().Where("id = ? AND user_id = ?", input.TopicID, userIDUUID).First(&topic).Error; err != nil {
			config.Logger.Warnf("Topic ID %s not found or not owned by user %s", input.TopicID, userIDUUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Topic not found or access denied"})
			return
		}
	}

	// Validate that task exists and belongs to user's topic (if provided)
	if input.TaskID != nil {
		var task models.Task_learning
		if err := config.GetDB().Joins("JOIN topics ON task_learnings.topic_id = topics.id").
			Where("task_learnings.id = ? AND topics.user_id = ?", input.TaskID, userIDUUID).
			First(&task).Error; err != nil {
			config.Logger.Warnf("Task ID %s not found or not owned by user %s", input.TaskID, userIDUUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found or access denied"})
			return
		}
	}

	// Validate link format if provided
	if input.Link != "" {
		if !strings.HasPrefix(input.Link, "http://") && !strings.HasPrefix(input.Link, "https://") {
			config.Logger.Warnf("Invalid link format: %s", input.Link)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Link must be a valid HTTP or HTTPS URL"})
			return
		}
	}

	resource := models.Resource{
		TopicID: input.TopicID,
		TaskID:  input.TaskID,
		Title:   input.Title,
		Link:    input.Link,
		Type:    input.Type,
		Notes:   input.Notes,
	}

	config.Logger.Infof("Creating resource for user %s: %s", userIDUUID, input.Title)
	if err := config.GetDB().Create(&resource).Error; err != nil {
		config.Logger.Errorf("Error creating resource for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create resource"})
		return
	}

	config.Logger.Infof("Successfully created resource ID %s for user %s", resource.ID, userIDUUID)
	c.JSON(http.StatusCreated, resource)
}

// GetResources godoc
// @Summary      Get resources
// @Description  Fetch resources for the logged-in user with optional filtering
// @Tags         resources
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        topic_id  query     string  false  "Filter by topic ID"
// @Param        task_id   query     string  false  "Filter by task ID"
// @Param        type      query     string  false  "Filter by resource type"
// @Param        search    query     string  false  "Search in title and notes"
// @Success      200  {object}  map[string][]models.Resource
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /resources [get]
func GetResources(c *gin.Context) {
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

	// Build query with proper joins to ensure user ownership
	query := config.GetDB().Model(&models.Resource{}).
		Joins("LEFT JOIN topics ON resources.topic_id = topics.id").
		Joins("LEFT JOIN task_learnings ON resources.task_id = task_learnings.id").
		Joins("LEFT JOIN topics AS task_topics ON task_learnings.topic_id = task_topics.id").
		Where("(topics.user_id = ? OR task_topics.user_id = ?)", userIDUUID, userIDUUID)

	// Apply filters
	if topicIDStr := c.Query("topic_id"); topicIDStr != "" {
		if topicID, err := uuid.Parse(topicIDStr); err == nil {
			query = query.Where("resources.topic_id = ?", topicID)
		}
	}

	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		if taskID, err := uuid.Parse(taskIDStr); err == nil {
			query = query.Where("resources.task_id = ?", taskID)
		}
	}

	if resourceType := c.Query("type"); resourceType != "" {
		validTypes := map[string]bool{
			"video": true, "article": true, "document": true, "book": true, "course": true,
		}
		if validTypes[resourceType] {
			query = query.Where("resources.type = ?", resourceType)
		}
	}

	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("(resources.title ILIKE ? OR resources.notes ILIKE ?)", searchTerm, searchTerm)
	}

	var resources []models.Resource
	if err := query.Order("resources.created_at DESC").Find(&resources).Error; err != nil {
		config.Logger.Errorf("Error fetching resources for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch resources"})
		return
	}

	config.Logger.Infof("Found %d resources for user %s", len(resources), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"resources": resources})
}

// GetResource godoc
// @Summary      Get a specific resource
// @Description  Fetch a specific resource by ID for the logged-in user
// @Tags         resources
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Resource ID"
// @Success      200  {object}  models.Resource
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /resources/{ID} [get]
func GetResource(c *gin.Context) {
	resourceIDStr := c.Param("ID")
	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid resource ID param: %s", resourceIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
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

	var resource models.Resource
	if err := config.GetDB().Where("id = ?", resourceID).
		Joins("LEFT JOIN topics ON resources.topic_id = topics.id").
		Joins("LEFT JOIN task_learnings ON resources.task_id = task_learnings.id").
		Joins("LEFT JOIN topics AS task_topics ON task_learnings.topic_id = task_topics.id").
		Where("(topics.user_id = ? OR task_topics.user_id = ?)", userIDUUID, userIDUUID).
		First(&resource).Error; err != nil {
		config.Logger.Warnf("Resource ID %s not found for user %s: %v", resourceID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved resource ID %s for user %s", resourceID, userIDUUID)
	c.JSON(http.StatusOK, resource)
}

// UpdateResourceRequest represents the request body for updating a resource
type UpdateResourceRequest struct {
	Title *string `json:"title"`
	Link  *string `json:"link"`
	Type  *string `json:"type"`
	Notes *string `json:"notes"`
}

// UpdateResource godoc
// @Summary      Update a resource
// @Description  Update a specific resource by ID for the logged-in user
// @Tags         resources
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID        path      int                     true  "Resource ID"
// @Param        resource  body      UpdateResourceRequest  true  "Resource update data"
// @Success      200  {object}  models.Resource
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /resources/{ID} [put]
func UpdateResource(c *gin.Context) {
	resourceIDStr := c.Param("ID")
	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid resource ID param for update: %s", resourceIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during resource update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var resource models.Resource
	if err := config.GetDB().Where("id = ?", resourceID).
		Joins("LEFT JOIN topics ON resources.topic_id = topics.id").
		Joins("LEFT JOIN task_learnings ON resources.task_id = task_learnings.id").
		Joins("LEFT JOIN topics AS task_topics ON task_learnings.topic_id = task_topics.id").
		Where("(topics.user_id = ? OR task_topics.user_id = ?)", userIDUUID, userIDUUID).
		First(&resource).Error; err != nil {
		config.Logger.Warnf("Resource not found for update: ID %s, User %s", resourceID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	var input UpdateResourceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for resource ID %s: %v", resourceID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Link != nil {
		// Validate link format
		if *input.Link != "" && !strings.HasPrefix(*input.Link, "http://") && !strings.HasPrefix(*input.Link, "https://") {
			config.Logger.Warnf("Invalid link format: %s", *input.Link)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Link must be a valid HTTP or HTTPS URL"})
			return
		}
		updates["link"] = *input.Link
	}
	if input.Type != nil {
		validTypes := map[string]bool{
			"video": true, "article": true, "document": true, "book": true, "course": true,
		}
		if !validTypes[*input.Type] {
			config.Logger.Warnf("Invalid resource type: %s", *input.Type)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource type"})
			return
		}
		updates["type"] = *input.Type
	}
	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for resource update: ID %s", resourceID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating resource ID %s for user %s with data: %+v", resourceID, userIDUUID, updates)
	if err := config.GetDB().Model(&resource).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update resource ID %s: %v", resourceID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource"})
		return
	}

	// Reload the updated resource
	if err := config.GetDB().First(&resource, resource.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated resource ID %s: %v", resource.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated resource"})
		return
	}

	config.Logger.Infof("Successfully updated resource ID %s for user %s", resource.ID, userIDUUID)
	c.JSON(http.StatusOK, resource)
}

// DeleteResource godoc
// @Summary      Delete a resource
// @Description  Delete a specific resource by ID for the logged-in user
// @Tags         resources
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Resource ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /resources/{ID} [delete]
func DeleteResource(c *gin.Context) {
	resourceIDStr := c.Param("ID")
	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid resource ID param for delete: %s", resourceIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during resource deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var resource models.Resource
	if err := config.GetDB().Where("id = ?", resourceID).
		Joins("LEFT JOIN topics ON resources.topic_id = topics.id").
		Joins("LEFT JOIN task_learnings ON resources.task_id = task_learnings.id").
		Joins("LEFT JOIN topics AS task_topics ON task_learnings.topic_id = task_topics.id").
		Where("(topics.user_id = ? OR task_topics.user_id = ?)", userIDUUID, userIDUUID).
		First(&resource).Error; err != nil {
		config.Logger.Warnf("Resource not found for delete: ID %s, User %s", resourceID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	config.Logger.Infof("Deleting resource ID %s for user %s", resourceID, userIDUUID)
	if err := config.GetDB().Delete(&resource).Error; err != nil {
		config.Logger.Errorf("Failed to delete resource ID %s: %v", resourceID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
		return
	}

	config.Logger.Infof("Successfully deleted resource ID %s for user %s", resourceID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}
