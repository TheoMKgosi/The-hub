package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/ai"
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

	result := config.GetDB().Preload("Task").Preload("RecurrenceRule").Where("user_id = ?", userID).Find(&schedule)

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
		Title            string     `json:"title" binding:"required"`
		Start            time.Time  `json:"start" binding:"required"`
		End              time.Time  `json:"end" binding:"required"`
		TaskID           *uuid.UUID `json:"task_id"`
		RecurrenceRuleID *uuid.UUID `json:"recurrence_rule_id"`
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
		Title:            input.Title,
		Start:            input.Start,
		End:              input.End,
		UserID:           userIDUUID,
		TaskID:           input.TaskID,
		RecurrenceRuleID: input.RecurrenceRuleID,
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
		Title            *string    `json:"title"`
		Start            *time.Time `json:"start"`
		End              *time.Time `json:"end"`
		TaskID           *uuid.UUID `json:"task_id"`
		RecurrenceRuleID *uuid.UUID `json:"recurrence_rule_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSchedule := map[string]interface{}{}
	if input.Title != nil {
		updatedSchedule["title"] = *input.Title
	}
	if input.Start != nil {
		updatedSchedule["start"] = *input.Start
	}
	if input.End != nil {
		updatedSchedule["end"] = *input.End
	}
	if input.TaskID != nil {
		updatedSchedule["task_id"] = *input.TaskID
	}
	if input.RecurrenceRuleID != nil {
		updatedSchedule["recurrence_rule_id"] = *input.RecurrenceRuleID
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

// CreateRecurrenceRule creates a new recurrence rule
func CreateRecurrenceRule(c *gin.Context) {
	var input struct {
		Frequency  string     `json:"frequency" binding:"required"`
		Interval   int        `json:"interval"`
		EndDate    *time.Time `json:"end_date"`
		Count      *int       `json:"count"`
		ByDay      string     `json:"by_day"`
		ByMonthDay *int       `json:"by_month_day"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	rule := models.RecurrenceRule{
		Frequency:  input.Frequency,
		Interval:   input.Interval,
		EndDate:    input.EndDate,
		Count:      input.Count,
		ByDay:      input.ByDay,
		ByMonthDay: input.ByMonthDay,
	}

	if err := config.GetDB().Create(&rule).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create recurrence rule"})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// GetAISuggestions returns AI-generated schedule suggestions
func GetAISuggestions(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist",
		})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	suggestions, err := ai.GetAISuggestions(userIDUUID)
	if err != nil {
		log.Println("Error getting AI suggestions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate suggestions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
	})
}
