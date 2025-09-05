package handlers

import (
	"net/http"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateLearningPathRequest represents the request body for creating a learning path
type CreateLearningPathRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	TopicIDs    []string `json:"topic_ids" binding:"required,min=1"`
}

// CreateLearningPath godoc
// @Summary      Create a new learning path
// @Description  Create a new learning path with ordered topics for the logged-in user
// @Tags         learning-paths
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        learning_path  body      CreateLearningPathRequest  true  "Learning path creation data"
// @Success      201  {object}  models.LearningPath
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /learning-paths [post]
func CreateLearningPath(c *gin.Context) {
	var input CreateLearningPathRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid learning path input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for learning path", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during learning path creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Validate that all topics exist and belong to the user
	var topics []models.Topic
	for _, topicIDStr := range input.TopicIDs {
		topicID, err := uuid.Parse(topicIDStr)
		if err != nil {
			config.Logger.Warnf("Invalid topic ID in learning path: %s", topicIDStr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID", "topic_id": topicIDStr})
			return
		}

		var topic models.Topic
		if err := config.GetDB().Where("id = ? AND user_id = ?", topicID, userIDUUID).First(&topic).Error; err != nil {
			config.Logger.Warnf("Topic ID %s not found or not owned by user %s", topicID, userIDUUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Topic not found or access denied", "topic_id": topicIDStr})
			return
		}
		topics = append(topics, topic)
	}

	// Create learning path
	learningPath := models.LearningPath{
		UserID:      userIDUUID,
		Title:       input.Title,
		Description: input.Description,
		Topics:      topics,
	}

	config.Logger.Infof("Creating learning path for user %s: %s", userIDUUID, input.Title)
	if err := config.GetDB().Create(&learningPath).Error; err != nil {
		config.Logger.Errorf("Error creating learning path for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create learning path"})
		return
	}

	// Create learning path topics with order
	for i, topic := range topics {
		learningPathTopic := models.LearningPathTopic{
			LearningPathID: learningPath.ID,
			TopicID:        topic.ID,
			OrderIndex:     i,
		}
		if err := config.GetDB().Create(&learningPathTopic).Error; err != nil {
			config.Logger.Errorf("Error creating learning path topic: %v", err)
			// Continue with other topics even if one fails
		}
	}

	config.Logger.Infof("Successfully created learning path ID %s for user %s", learningPath.ID, userIDUUID)
	c.JSON(http.StatusCreated, learningPath)
}

// GetLearningPaths godoc
// @Summary      Get learning paths
// @Description  Fetch learning paths for the logged-in user
// @Tags         learning-paths
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.LearningPath
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /learning-paths [get]
func GetLearningPaths(c *gin.Context) {
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

	var learningPaths []models.LearningPath
	if err := config.GetDB().Preload("Topics").Where("user_id = ?", userIDUUID).Find(&learningPaths).Error; err != nil {
		config.Logger.Errorf("Error fetching learning paths for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch learning paths"})
		return
	}

	config.Logger.Infof("Found %d learning paths for user %s", len(learningPaths), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"learning_paths": learningPaths})
}

// GetLearningPath godoc
// @Summary      Get a specific learning path
// @Description  Fetch a specific learning path by ID for the logged-in user
// @Tags         learning-paths
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Learning Path ID"
// @Success      200  {object}  models.LearningPath
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /learning-paths/{ID} [get]
func GetLearningPath(c *gin.Context) {
	learningPathIDStr := c.Param("ID")
	learningPathID, err := uuid.Parse(learningPathIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid learning path ID param: %s", learningPathIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid learning path ID"})
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

	var learningPath models.LearningPath
	if err := config.GetDB().Preload("Topics").Where("id = ? AND user_id = ?", learningPathID, userIDUUID).First(&learningPath).Error; err != nil {
		config.Logger.Warnf("Learning path ID %s not found for user %s: %v", learningPathID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning path not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved learning path ID %s for user %s", learningPathID, userIDUUID)
	c.JSON(http.StatusOK, learningPath)
}

// DeleteLearningPath godoc
// @Summary      Delete a learning path
// @Description  Delete a specific learning path by ID for the logged-in user
// @Tags         learning-paths
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Learning Path ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /learning-paths/{ID} [delete]
func DeleteLearningPath(c *gin.Context) {
	learningPathIDStr := c.Param("ID")
	learningPathID, err := uuid.Parse(learningPathIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid learning path ID param for delete: %s", learningPathIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid learning path ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during learning path deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var learningPath models.LearningPath
	if err := config.GetDB().Where("id = ? AND user_id = ?", learningPathID, userIDUUID).First(&learningPath).Error; err != nil {
		config.Logger.Warnf("Learning path not found for delete: ID %s, User %s", learningPathID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning path not found"})
		return
	}

	// Delete associated learning path topics first
	if err := config.GetDB().Where("learning_path_id = ?", learningPathID).Delete(&models.LearningPathTopic{}).Error; err != nil {
		config.Logger.Errorf("Failed to delete learning path topics for path ID %s: %v", learningPathID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete learning path topics"})
		return
	}

	// Delete the learning path
	if err := config.GetDB().Delete(&learningPath).Error; err != nil {
		config.Logger.Errorf("Failed to delete learning path ID %s: %v", learningPathID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete learning path"})
		return
	}

	config.Logger.Infof("Successfully deleted learning path ID %s for user %s", learningPathID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Learning path deleted successfully"})
}
