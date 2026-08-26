package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/ai"
	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// hasTimeConflict checks if a new time slot conflicts with existing scheduled tasks
func hasTimeConflict(db *gorm.DB, userID uuid.UUID, start, end time.Time, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := db.Model(&models.ScheduledTask{}).Where(`user_id = ? AND (("start" < ? AND "end" > ?) OR ("start" < ? AND "end" > ?))`, userID, end, start, start, end)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// validateScheduleInput validates the schedule input data
func validateScheduleInput(input struct {
	Title            string     `json:"title" binding:"required"`
	Start            time.Time  `json:"start" binding:"required"`
	End              time.Time  `json:"end" binding:"required"`
	RecurrenceRuleID *uuid.UUID `json:"recurrence_rule_id"`
}) error {
	if input.Start.After(input.End) || input.Start.Equal(input.End) {
		return fmt.Errorf("start time must be before end time")
	}

	// Check for reasonable duration (max 24 hours)
	if input.End.Sub(input.Start) > 24*time.Hour {
		return fmt.Errorf("event duration cannot exceed 24 hours")
	}

	// Check for past dates (allow events up to 1 hour in the past for flexibility)
	if input.Start.Before(time.Now().Add(-time.Hour)) {
		return fmt.Errorf("cannot schedule events in the past")
	}

	return nil
}

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

	result := config.GetDB().Preload("RecurrenceRule").Where("user_id = ?", userID).Find(&schedule)

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
		RecurrenceRuleID *uuid.UUID `json:"recurrence_rule_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
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

	// Validate input
	if err := validateScheduleInput(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for time conflicts
	hasConflict, err := hasTimeConflict(config.GetDB(), userIDUUID, input.Start, input.End, nil)
	if err != nil {
		log.Println("Error checking for conflicts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check for scheduling conflicts"})
		return
	}

	if hasConflict {
		c.JSON(http.StatusConflict, gin.H{"error": "This time slot conflicts with an existing scheduled event"})
		return
	}

	schedule := models.ScheduledTask{
		Title:            input.Title,
		Start:            input.Start,
		End:              input.End,
		UserID:           userIDUUID,
		RecurrenceRuleID: input.RecurrenceRuleID,
	}

	if err := config.GetDB().Create(&schedule).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create scheduled task"})
		return
	}
	c.JSON(http.StatusCreated, schedule)

}

// Update a specific task
func UpdateSchedule(c *gin.Context) {
	var schedule models.ScheduledTask

	scheduleTaskID := c.Param("ID")
	if err := config.GetDB().Where("id = ?", scheduleTaskID).First(&schedule).Error; err != nil {
		log.Println("Error ID: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	// Check if user owns this scheduled task
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if schedule.UserID != userIDUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this scheduled task"})
		return
	}

	var input struct {
		Title            *string    `json:"title"`
		Start            *time.Time `json:"start"`
		End              *time.Time `json:"end"`
		RecurrenceRuleID *uuid.UUID `json:"recurrence_rule_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate time fields if provided
	if input.Start != nil && input.End != nil {
		if input.Start.After(*input.End) || input.Start.Equal(*input.End) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
			return
		}

		if input.End.Sub(*input.Start) > 24*time.Hour {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event duration cannot exceed 24 hours"})
			return
		}

		// Check for time conflicts if both start and end are being updated
		hasConflict, err := hasTimeConflict(config.GetDB(), userIDUUID, *input.Start, *input.End, &schedule.ID)
		if err != nil {
			log.Println("Error checking for conflicts:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check for scheduling conflicts"})
			return
		}

		if hasConflict {
			c.JSON(http.StatusConflict, gin.H{"error": "This time slot conflicts with an existing scheduled event"})
			return
		}
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

	scheduleTaskID := c.Param("ID")
	if err := config.GetDB().Where("id = ?", scheduleTaskID).First(&schedule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Scheduled task not found",
		})
		return
	}

	// Check if user owns this scheduled task
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if schedule.UserID != userIDUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this scheduled task"})
		return
	}

	// Use transaction for safe deletion
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println("Error starting transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start transaction"})
		return
	}

	if err := tx.Delete(&schedule).Error; err != nil {
		tx.Rollback()
		log.Println("Error deleting scheduled task:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit transaction"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// CreateRecurrenceRule creates a new recurrence rule
func CreateRecurrenceRule(c *gin.Context) {
	var input struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Frequency   string     `json:"frequency" binding:"required"`
		Interval    int        `json:"interval"`
		EndDate     *time.Time `json:"end_date"`
		Count       *int       `json:"count"`
		ByDay       string     `json:"by_day"`
		ByMonthDay  *int       `json:"by_month_day"`
		ByMonth     *int       `json:"by_month"`
		StartDate   *time.Time `json:"start_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	rule := models.RecurrenceRule{
		UserID:      userIDUUID,
		Name:        input.Name,
		Description: input.Description,
		Frequency:   input.Frequency,
		Interval:    input.Interval,
		EndDate:     input.EndDate,
		Count:       input.Count,
		ByDay:       input.ByDay,
		ByMonthDay:  input.ByMonthDay,
		ByMonth:     input.ByMonth,
		StartDate:   input.StartDate,
	}

	if err := config.GetDB().Create(&rule).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create recurrence rule"})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// BulkCreateSchedule creates multiple scheduled tasks at once
func BulkCreateSchedule(c *gin.Context) {
	var inputs []struct {
		Title            string     `json:"title" binding:"required"`
		Start            time.Time  `json:"start" binding:"required"`
		End              time.Time  `json:"end" binding:"required"`
		RecurrenceRuleID *uuid.UUID `json:"recurrence_rule_id"`
	}

	if err := c.ShouldBindJSON(&inputs); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context during bulk schedule creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Validate all inputs first
	for i, input := range inputs {
		if err := validateScheduleInput(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid input at index %d: %s", i, err.Error()),
			})
			return
		}
	}

	// Check for conflicts
	for i, input := range inputs {
		hasConflict, err := hasTimeConflict(config.GetDB(), userIDUUID, input.Start, input.End, nil)
		if err != nil {
			log.Println("Error checking for conflicts:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check for scheduling conflicts"})
			return
		}

		if hasConflict {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Time slot conflict at index %d", i),
			})
			return
		}
	}

	// Create all schedules in a transaction
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println("Error starting transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start transaction"})
		return
	}

	var createdSchedules []models.ScheduledTask
	for _, input := range inputs {
		schedule := models.ScheduledTask{
			Title:            input.Title,
			Start:            input.Start,
			End:              input.End,
			UserID:           userIDUUID,
			RecurrenceRuleID: input.RecurrenceRuleID,
		}

		if err := tx.Create(&schedule).Error; err != nil {
			tx.Rollback()
			log.Println("Error creating scheduled task:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create scheduled tasks"})
			return
		}

		createdSchedules = append(createdSchedules, schedule)
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"scheduled_tasks": createdSchedules})
}

// BulkDeleteSchedule deletes multiple scheduled tasks at once
func BulkDeleteSchedule(c *gin.Context) {
	var input struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Use transaction for safe bulk deletion
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println("Error starting transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start transaction"})
		return
	}

	// Delete the scheduled tasks
	result := tx.Where("id IN ? AND user_id = ?", input.IDs, userIDUUID).Delete(&models.ScheduledTask{})
	if result.Error != nil {
		tx.Rollback()
		log.Println("Error deleting scheduled tasks:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete scheduled tasks"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Successfully deleted %d scheduled tasks", result.RowsAffected),
		"deleted_count": result.RowsAffected,
	})
}

const aiSuggestionsWindowDays = 7

// isSlotInWindow reports whether the given slot falls within the scheduling window.
func isSlotInWindow(start, end, windowStart, windowEnd time.Time) bool {
	return !start.Before(windowStart) && !end.After(windowEnd)
}

// GetScheduleSuggestions uses AI to slot the user's pending tasks into free time windows.
func GetScheduleSuggestions(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	db := config.GetDB()

	// Fetch pending, top-level tasks for the user
	var tasks []models.Task
	if err := db.Where("user_id = ? AND status != 'completed' AND parent_task_id IS NULL", userIDUUID).
		Order("created_at ASC").
		Find(&tasks).Error; err != nil {
		log.Println("Error fetching tasks for schedule suggestions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Fetch existing scheduled events within the scheduling window
	windowStart := time.Now()
	windowEnd := windowStart.AddDate(0, 0, aiSuggestionsWindowDays)
	var existingEvents []models.ScheduledTask
	if err := db.Where("user_id = ? AND start >= ? AND start < ?", userIDUUID, windowStart, windowEnd).
		Find(&existingEvents).Error; err != nil {
		log.Println("Error fetching existing schedule for suggestions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}

	// Skip tasks that already have a matching scheduled event in the window
	existingTitles := make(map[string]bool)
	for _, event := range existingEvents {
		existingTitles[strings.ToLower(strings.TrimSpace(event.Title))] = true
	}

	schedulable := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		if existingTitles[strings.ToLower(strings.TrimSpace(task.Title))] {
			continue
		}
		schedulable = append(schedulable, task)
	}

	if len(schedulable) == 0 {
		c.JSON(http.StatusOK, gin.H{"suggestions": []models.ScheduledTask{}})
		return
	}

	client, err := ai.GetOpenRouterClient()
	if err != nil {
		// Fall back to the algorithmic scheduler when the AI client is unavailable
		log.Println("AI client unavailable, using algorithmic schedule suggestions:", err)
		suggestions, err := ai.GetAISuggestions(userIDUUID)
		if err != nil {
			log.Println("Error generating algorithmic schedule suggestions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate schedule suggestions"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
		return
	}

	aiResponse, err := client.ScheduleTasksIntoSchedule(schedulable, existingEvents, aiSuggestionsWindowDays)
	if err != nil {
		log.Println("Error generating AI schedule suggestions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate schedule suggestions"})
		return
	}

	var slots []ai.ScheduleSlot
	if err := json.Unmarshal([]byte(aiResponse), &slots); err != nil {
		start := strings.Index(aiResponse, "[")
		end := strings.LastIndex(aiResponse, "]")
		if start == -1 || end == -1 || start >= end {
			log.Println("Failed to parse AI schedule response:", aiResponse)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI schedule suggestions"})
			return
		}

		if err := json.Unmarshal([]byte(aiResponse[start:end+1]), &slots); err != nil {
			log.Println("Failed to parse extracted AI schedule response:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI schedule suggestions"})
			return
		}
	}

	suggestions := make([]models.ScheduledTask, 0, len(slots))
	for _, slot := range slots {
		if slot.Start.After(slot.End) || slot.Start.Equal(slot.End) {
			continue
		}

		// Reject slots that exceed reasonable event duration (max 24 hours)
		if slot.End.Sub(slot.Start) > 24*time.Hour {
			continue
		}

		// Reject slots outside the scheduling window
		if !isSlotInWindow(slot.Start, slot.End, windowStart, windowEnd) {
			continue
		}

		// Reject slots that conflict with existing scheduled events
		hasConflict, err := hasTimeConflict(db, userIDUUID, slot.Start, slot.End, nil)
		if err != nil {
			log.Println("Error checking conflicts for AI suggestion:", err)
			continue
		}
		if hasConflict {
			continue
		}

		suggestions = append(suggestions, models.ScheduledTask{
			ID:          uuid.New(),
			Title:       slot.Title,
			Start:       slot.Start,
			End:         slot.End,
			UserID:      userIDUUID,
			CreatedByAI: true,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
	})
}
