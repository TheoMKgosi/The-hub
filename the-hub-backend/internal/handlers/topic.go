package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// Get all topics for the user
func GetTopics(c *gin.Context) {
	var topics []models.Topic
	userID := c.MustGet("userID").(uint)
	log.Printf("Fetching topics for user ID: %d", userID)

	if err := config.GetDB().Where("user_id = ?", userID).Preload("Tags").Find(&topics).Error; err != nil {
		log.Printf("Error fetching topics for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve topics"})
		return
	}

	log.Printf("Fetched %d topics for user ID: %d", len(topics), userID)
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

// Get one topic
func GetTopic(c *gin.Context) {
	topicID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Printf("Invalid topic ID param: %s", c.Param("ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	var topic models.Topic
	log.Printf("Fetching topic ID: %d", topicID)
	if err := config.GetDB().Preload("Tags").First(&topic, topicID).Error; err != nil {
		log.Printf("Topic not found: %d", topicID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	c.JSON(http.StatusOK, topic)
}

// Create a topic
func CreateTopic(c *gin.Context) {
	var input struct {
		Title          string     `json:"title" binding:"required"`
		Description    string     `json:"description"`
		Status         string     `json:"status"`
		EstimatedHours int        `json:"estimated_hours"`
		Deadline       *time.Time `json:"deadline"`
		TagIDs         []uint     `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input on topic creation: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.MustGet("userID").(uint)
	topic := models.Topic{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Deadline:    input.Deadline,
		UserID:      userID,
	}

	if len(input.TagIDs) > 0 {
		var tags []models.Tag
		if err := config.GetDB().Where("id IN ?", input.TagIDs).Find(&tags).Error; err == nil {
			topic.Tags = tags
			log.Printf("Loaded %d tags for new topic", len(tags))
		} else {
			log.Printf("Failed to load tags: %v", err)
		}
	}

	if err := config.GetDB().Create(&topic).Error; err != nil {
		log.Printf("Failed to create topic: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create topic"})
		return
	}

	log.Printf("Created topic ID: %d for user ID: %d", topic.ID, userID)
	c.JSON(http.StatusCreated, topic)
}

// Update topic
func UpdateTopic(c *gin.Context) {
	topicID := c.Param("ID")
	var topic models.Topic
	if err := config.GetDB().First(&topic, topicID).Error; err != nil {
		log.Printf("Topic not found for update, ID: %s", topicID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	var input struct {
		Title          *string    `json:"title"`
		Description    *string    `json:"description"`
		Status         *string    `json:"status"`
		EstimatedHours *int       `json:"estimated_hours"`
		Deadline       *time.Time `json:"deadline"`
		TagIDs         *[]uint    `json:"tag_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input on update for topic ID %s: %v", topicID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.EstimatedHours != nil {
		updates["estimated_hours"] = *input.EstimatedHours
	}
	if input.Deadline != nil {
		updates["deadline"] = *input.Deadline
	}

	if len(updates) > 0 {
		log.Printf("Updating topic ID %s with data: %+v", topicID, updates)
	}

	if err := config.GetDB().Model(&topic).Updates(updates).Error; err != nil {
		log.Printf("Failed to update topic ID %s: %v", topicID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update topic"})
		return
	}

	if input.TagIDs != nil {
		var tags []models.Tag
		if err := config.GetDB().Where("id IN ?", *input.TagIDs).Find(&tags).Error; err == nil {
			config.GetDB().Model(&topic).Association("Tags").Replace(&tags)
			log.Printf("Updated tags for topic ID %s", topicID)
		} else {
			log.Printf("Failed to update tags for topic ID %s: %v", topicID, err)
		}
	}

	c.JSON(http.StatusOK, topic)
}

// Delete topic
func DeleteTopic(c *gin.Context) {
	topicID := c.Param("ID")
	var topic models.Topic
	if err := config.GetDB().First(&topic, topicID).Error; err != nil {
		log.Printf("Topic not found for delete, ID: %s", topicID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	if err := config.GetDB().Delete(&topic).Error; err != nil {
		log.Printf("Failed to delete topic ID %s: %v", topicID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete topic"})
		return
	}

	log.Printf("Deleted topic ID: %s", topicID)
	c.JSON(http.StatusOK, topic)
}
