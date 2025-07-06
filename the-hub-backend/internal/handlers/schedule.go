package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		"tasks": schedule,
	})

}

// Create a scheduled task
func CreateSchedule(c *gin.Context) {

	var input struct {
		TaskID uint      `json:"task_id"`
		Title  string    `json:"title"`
		Start  time.Time `json:"start"`
		End    time.Time `json:"end"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create Task"})
		return
	}

	userID := c.MustGet("userID").(uint)

	// TODO: Make a transaction to also change the due date in task
	err := config.GetDB().Transaction(func(tx *gorm.DB) error {
		// 1. Create the schedule
		schedule := models.ScheduledTask{
			Title:  input.Title,
			TaskID: input.TaskID,
			Start:  input.Start,
			End:    input.End,
			UserID: userID,
		}

		if err := tx.Create(&schedule).Error; err != nil {
			return err // rollback
		}

		// 2. Update the task's due_date to match the schedule's end
		if err := tx.Model(&models.Task{}).
			Where("task_id = ? AND user_id = ?", input.TaskID, userID).
			Update("due_date", input.End).Error; err != nil {
			return err // rollback
		}

		// All good
		c.JSON(http.StatusCreated, schedule)
		return nil
	})

	if err != nil {
		log.Println("Transaction failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create schedule"})
	}

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
