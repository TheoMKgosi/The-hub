package handlers

import (
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateFeedback godoc
// @Summary      Submit new feedback
// @Description  Create a new feedback entry from a user
// @Tags         feedback
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        feedback body      models.Feedback  true  "Feedback data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /feedback [post]
func CreateFeedback(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context during feedback creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var feedback models.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		config.Logger.Warnf("Invalid feedback data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback data"})
		return
	}

	// Set the user ID from context
	feedback.UserID = userIDUUID

	// Validate required fields
	if feedback.Type == "" || feedback.Subject == "" || feedback.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type, subject, and description are required"})
		return
	}

	// Validate type
	validTypes := map[string]bool{"bug": true, "feature": true, "improvement": true, "general": true}
	if !validTypes[feedback.Type] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback type"})
		return
	}

	// Validate rating if provided
	if feedback.Rating != nil && (*feedback.Rating < 1 || *feedback.Rating > 5) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	// Create feedback in database
	if err := config.GetDB().Create(&feedback).Error; err != nil {
		config.Logger.Errorf("Failed to create feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit feedback"})
		return
	}

	config.Logger.Infof("Feedback created successfully: ID %s, User %s", feedback.ID, userIDUUID)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Feedback submitted successfully",
		"feedback": feedback,
	})
}

// GetUserFeedback godoc
// @Summary      Get user's feedback history
// @Description  Retrieve all feedback submitted by the authenticated user
// @Tags         feedback
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page     query     int  false  "Page number (default: 1)"
// @Param        limit    query     int  false  "Items per page (default: 10)"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /feedback [get]
func GetUserFeedback(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		config.Logger.Warn("userID not found in context during feedback retrieval")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var feedback []models.Feedback
	var total int64

	// Get total count
	if err := config.GetDB().Model(&models.Feedback{}).Where("user_id = ?", userIDUUID).Count(&total).Error; err != nil {
		config.Logger.Errorf("Failed to count feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}

	// Get feedback with pagination
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedback).Error; err != nil {
		config.Logger.Errorf("Failed to retrieve feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}

	config.Logger.Infof("User feedback retrieved: User %s, Count %d", userIDUUID, len(feedback))
	c.JSON(http.StatusOK, gin.H{
		"feedback": feedback,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetAllFeedback godoc
// @Summary      Get all feedback (Admin only)
// @Description  Retrieve all feedback entries for admin review
// @Tags         feedback, admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page     query     int     false  "Page number (default: 1)"
// @Param        limit    query     int     false  "Items per page (default: 10)"
// @Param        status   query     string  false  "Filter by status"
// @Param        type     query     string  false  "Filter by type"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /admin/feedback [get]
func GetAllFeedback(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var feedback []models.Feedback
	var total int64

	// Build query
	query := config.GetDB().Model(&models.Feedback{}).Preload("User")

	// Apply filters
	if status := c.Query("status"); status != "" {
		validStatuses := map[string]bool{"pending": true, "reviewed": true, "implemented": true, "declined": true}
		if validStatuses[status] {
			query = query.Where("status = ?", status)
		}
	}

	if feedbackType := c.Query("type"); feedbackType != "" {
		validTypes := map[string]bool{"bug": true, "feature": true, "improvement": true, "general": true}
		if validTypes[feedbackType] {
			query = query.Where("type = ?", feedbackType)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		config.Logger.Errorf("Failed to count all feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}

	// Get feedback with pagination
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&feedback).Error; err != nil {
		config.Logger.Errorf("Failed to retrieve all feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"})
		return
	}

	config.Logger.Infof("All feedback retrieved: Count %d", len(feedback))
	c.JSON(http.StatusOK, gin.H{
		"feedback": feedback,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateFeedbackStatus godoc
// @Summary      Update feedback status (Admin only)
// @Description  Update the status and admin response for a feedback entry
// @Tags         feedback, admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                      true  "Feedback ID"
// @Param        update   body      map[string]interface{}   true  "Status update data"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /admin/feedback/{id} [patch]
func UpdateFeedbackStatus(c *gin.Context) {
	feedbackIDStr := c.Param("id")
	feedbackID, err := uuid.Parse(feedbackIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid feedback ID param: %s", feedbackIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	var updateData struct {
		Status        string `json:"status"`
		AdminResponse string `json:"admin_response"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		config.Logger.Warnf("Invalid update data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	// Validate status
	if updateData.Status != "" {
		validStatuses := map[string]bool{"pending": true, "reviewed": true, "implemented": true, "declined": true}
		if !validStatuses[updateData.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
	}

	// Find and update feedback
	var feedback models.Feedback
	if err := config.GetDB().First(&feedback, feedbackID).Error; err != nil {
		config.Logger.Warnf("Feedback not found: ID %s", feedbackID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}

	updates := make(map[string]interface{})
	if updateData.Status != "" {
		updates["status"] = updateData.Status
	}
	if updateData.AdminResponse != "" {
		updates["admin_response"] = updateData.AdminResponse
	}

	if err := config.GetDB().Model(&feedback).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feedback"})
		return
	}

	config.Logger.Infof("Feedback updated successfully: ID %s, Status %s", feedbackID, updateData.Status)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Feedback updated successfully",
		"feedback": feedback,
	})
}
