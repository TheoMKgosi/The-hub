package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// Get all tags for the user
func GetTags(c *gin.Context) {
	var tags []models.Tag
	userID := c.MustGet("userID").(uint)
	log.Printf("Fetching tags for user ID: %d", userID)

	if err := config.GetDB().Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		log.Printf("Error fetching tags for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tags"})
		return
	}

	log.Printf("Fetched %d tags for user ID: %d", len(tags), userID)
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// Get one tag
func GetTag(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Printf("Invalid tag ID param: %s", c.Param("ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	log.Printf("Fetching tag ID: %d", tagID)
	var tag models.Tag
	if err := config.GetDB().First(&tag, tagID).Error; err != nil {
		log.Printf("Tag not found: %d", tagID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// Create a tag
func CreateTag(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input on tag creation: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.MustGet("userID").(uint)
	tag := models.Tag{
		Name:   input.Name,
		Color:  input.Color,
		UserID: userID,
	}

	if err := config.GetDB().Create(&tag).Error; err != nil {
		log.Printf("Error creating tag for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create tag"})
		return
	}

	log.Printf("Created tag ID: %d for user ID: %d", tag.ID, userID)
	c.JSON(http.StatusCreated, tag)
}

// Update a tag
func UpdateTag(c *gin.Context) {
	tagID := c.Param("ID")
	var tag models.Tag
	if err := config.GetDB().First(&tag, tagID).Error; err != nil {
		log.Printf("Tag not found for update, ID: %s", tagID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	var input struct {
		Name  *string `json:"name"`
		Color *string `json:"color"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input on tag update ID %s: %v", tagID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}

	log.Printf("Updating tag ID %s with: %+v", tagID, updates)
	if err := config.GetDB().Model(&tag).Updates(updates).Error; err != nil {
		log.Printf("Failed to update tag ID %s: %v", tagID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update tag"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// Delete a tag
func DeleteTag(c *gin.Context) {
	tagID := c.Param("ID")
	var tag models.Tag
	if err := config.GetDB().First(&tag, tagID).Error; err != nil {
		log.Printf("Tag not found for delete, ID: %s", tagID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	if err := config.GetDB().Delete(&tag).Error; err != nil {
		log.Printf("Failed to delete tag ID %s: %v", tagID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete tag"})
		return
	}

	log.Printf("Deleted tag ID: %s", tagID)
	c.JSON(http.StatusOK, tag)
}

