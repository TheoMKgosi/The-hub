package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Get all schedule
func GetSchedule(c *gin.Context) {
	var schedule []models.ScheduledTask
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist",
		})
		return
	}

	result := config.GetDB().Where("user_id = ?", userID).Find(&schedule)

	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule": schedule,
	})

}

// Create a scheduled task
func CreateSchedule(c *gin.Context) {

	var input struct {
		Title string    `json:"title" binding:"required"`
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create Task"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context during schedule creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	schedule := models.ScheduledTask{
		Title:  input.Title,
		Start:  input.Start,
		End:    input.End,
		UserID: userIDUUID,
	}

	if err := config.GetDB().Create(&schedule).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Task"})
		return
	}
	c.JSON(http.StatusCreated, schedule)

}

// Update a specific task
func UpdateSchedule(c *gin.Context) {
	// FIX: Fix the implementation
	var schedule models.ScheduledTask

	scheduleTaskID := c.Param("ID")
	if err := config.GetDB().First(&schedule, scheduleTaskID).Error; err != nil {
		log.Println("Error ID: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	var input struct {
		Start *time.Time `json:"star "`
		End   *time.Time `json:"end"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSchedule := map[string]interface{}{}
	if input.Start != nil {
		updatedSchedule["start"] = *input.Start
	}
	if input.End != nil {
		updatedSchedule["end"] = *input.End
	}

	if err := config.GetDB().Model(&schedule).Updates(updatedSchedule).Error; err != nil {
		log.Println("Error updating task:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update or remove calendar event based on new due date
	c.JSON(http.StatusOK, schedule)

}

// Delete a specific task
func DeleteSchedule(c *gin.Context) {
	var schedule models.ScheduledTask

	// FIX: Make a transaction to make due date null after deleting task
	scheduleTaskID := c.Param("ID")
	if err := config.GetDB().First(&schedule, scheduleTaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Scheduled task not found",
		})
		return
	}

	if err := config.GetDB().Delete(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
