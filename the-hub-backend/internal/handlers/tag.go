package handlers

import (
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// GetTags godoc
// @Summary      Get all tags
// @Description  Fetch tags for the logged-in user with optional ordering
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_by  query     string  false  "Order by field (name, created_at)"  default(name)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.Tag
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags [get]
func GetTags(c *gin.Context) {
	var tags []models.Tag
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get query parameters for ordering
	orderBy := c.DefaultQuery("order_by", "name")
	sortDir := c.DefaultQuery("sort", "asc")

	// Validate order_by parameter
	validOrderFields := map[string]bool{
		"name":       true,
		"created_at": true,
		"color":      true,
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

	config.Logger.Infof("Fetching tags for user ID: %v with order: %s", userID, orderClause)
	if err := config.GetDB().Where("user_id = ?", userID).Order(orderClause).Find(&tags).Error; err != nil {
		config.Logger.Errorf("Error fetching tags for user %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tags"})
		return
	}

	config.Logger.Infof("Found %d tags for user ID %v", len(tags), userID)
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// GetTag godoc
// @Summary      Get a specific tag
// @Description  Fetch a specific tag by ID for the logged-in user
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Tag ID"
// @Success      200  {object}  map[string]models.Tag
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags/{ID} [get]
func GetTag(c *gin.Context) {
	tagIDStr := c.Param("ID")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid tag ID param: %s", tagIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	config.Logger.Infof("Fetching tag ID: %d for user ID: %v", tagID, userID)
	var tag models.Tag
	// Ensure user can only access their own tags
	if err := config.GetDB().Where("id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
		config.Logger.Errorf("Tag ID %d not found for user %v: %v", tagID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved tag ID %d for user %v", tagID, userID)
	c.JSON(http.StatusOK, gin.H{"tag": tag})
}

// CreateTagRequest represents the request body for creating a tag
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required" example:"Work"`
	Color string `json:"color" example:"#FF5733"`
}

// CreateTag godoc
// @Summary      Create a new tag
// @Description  Create a new tag for the logged-in user
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tag  body      CreateTagRequest  true  "Tag creation data"
// @Success      201  {object}  models.Tag
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags [post]
func CreateTag(c *gin.Context) {
	var input CreateTagRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid tag input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for tag", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during tag creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Check if tag name already exists for this user
	var existingTag models.Tag
	if err := config.GetDB().Where("name = ? AND user_id = ?", input.Name, userIDUint).First(&existingTag).Error; err == nil {
		config.Logger.Warnf("Tag name '%s' already exists for user %d", input.Name, userIDUint)
		c.JSON(http.StatusConflict, gin.H{"error": "Tag with this name already exists"})
		return
	}

	tag := models.Tag{
		Name:   input.Name,
		Color:  input.Color,
		UserID: userIDUint,
	}

	config.Logger.Infof("Creating tag for user %d: %s", userIDUint, input.Name)
	if err := config.GetDB().Create(&tag).Error; err != nil {
		config.Logger.Errorf("Error creating tag for user %d: %v", userIDUint, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create tag"})
		return
	}

	config.Logger.Infof("Successfully created tag ID %d for user %d", tag.ID, userIDUint)
	c.JSON(http.StatusCreated, tag)
}

// UpdateTagRequest represents the request body for updating a tag
type UpdateTagRequest struct {
	Name  *string `json:"name" example:"Updated Work"`
	Color *string `json:"color" example:"#33FF57"`
}

// UpdateTag godoc
// @Summary      Update a tag
// @Description  Update a specific tag by ID for the logged-in user
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int               true  "Tag ID"
// @Param        tag  body      UpdateTagRequest  true  "Tag update data"
// @Success      200  {object}  models.Tag
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags/{ID} [put]
func UpdateTag(c *gin.Context) {
	tagIDStr := c.Param("ID")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid tag ID param for update: %s", tagIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during tag update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var tag models.Tag
	// Ensure user can only update their own tags
	if err := config.GetDB().Where("id = ? AND user_id = ?", tagID, userIDUint).First(&tag).Error; err != nil {
		config.Logger.Warnf("Tag not found for update: ID %d, User %d", tagID, userIDUint)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	var input UpdateTagRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for tag ID %d: %v", tagID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		// Check if the new name conflicts with existing tags for this user
		if *input.Name != tag.Name {
			var existingTag models.Tag
			if err := config.GetDB().Where("name = ? AND user_id = ? AND id != ?", *input.Name, userIDUint, tagID).First(&existingTag).Error; err == nil {
				config.Logger.Warnf("Tag name '%s' already exists for user %d", *input.Name, userIDUint)
				c.JSON(http.StatusConflict, gin.H{"error": "Tag with this name already exists"})
				return
			}
		}
		updates["name"] = *input.Name
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for tag update: ID %d", tagID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating tag ID %d for user %d with data: %+v", tagID, userIDUint, updates)
	if err := config.GetDB().Model(&tag).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update tag ID %d: %v", tagID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
		return
	}

	// Reload the updated tag
	if err := config.GetDB().First(&tag, tag.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated tag ID %d: %v", tag.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated tag"})
		return
	}

	config.Logger.Infof("Successfully updated tag ID %d for user %d", tag.ID, userIDUint)
	c.JSON(http.StatusOK, tag)
}

// DeleteTag godoc
// @Summary      Delete a tag
// @Description  Delete a specific tag by ID for the logged-in user
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Tag ID"
// @Success      200  {object}  models.Tag
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tags/{ID} [delete]
func DeleteTag(c *gin.Context) {
	tagIDStr := c.Param("ID")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid tag ID param for delete: %s", tagIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during tag deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var tag models.Tag
	// Ensure user can only delete their own tags
	if err := config.GetDB().Where("id = ? AND user_id = ?", tagID, userIDUint).First(&tag).Error; err != nil {
		config.Logger.Warnf("Tag not found for delete: ID %d, User %d", tagID, userIDUint)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Check if tag is being used by any tasks
	var taskCount int64
	if err := config.GetDB().Table("task_tags").Where("tag_id = ?", tagID).Count(&taskCount).Error; err != nil {
		config.Logger.Errorf("Error checking tag usage for tag ID %d: %v", tagID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify tag usage"})
		return
	}

	if taskCount > 0 {
		config.Logger.Warnf("Attempted to delete tag ID %d which is used by %d tasks", tagID, taskCount)
		c.JSON(http.StatusConflict, gin.H{
			"error": "Cannot delete tag that is currently being used by tasks",
			"task_count": taskCount,
		})
		return
	}

	config.Logger.Infof("Deleting tag ID %d for user %d", tagID, userIDUint)
	if err := config.GetDB().Delete(&tag).Error; err != nil {
		config.Logger.Errorf("Failed to delete tag ID %d: %v", tagID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	config.Logger.Infof("Successfully deleted tag ID %d for user %d", tagID, userIDUint)
	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully", "tag": tag})
}
