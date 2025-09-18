package handlers

import (
	"fmt"
	"log"
	"net/http"
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
	TaskID           *uuid.UUID `json:"task_id"`
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
		TaskID:           input.TaskID,
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
	if err := config.GetDB().First(&schedule, scheduleTaskID).Error; err != nil {
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
		TaskID           *uuid.UUID `json:"task_id"`
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

	scheduleTaskID := c.Param("ID")
	if err := config.GetDB().First(&schedule, scheduleTaskID).Error; err != nil {
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

	// If this scheduled task is linked to a task, update the task's due date
	if schedule.TaskID != nil {
		if err := tx.Model(&models.Task{}).Where("id = ?", *schedule.TaskID).Update("due_date", nil).Error; err != nil {
			tx.Rollback()
			log.Println("Error updating task due date:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update associated task"})
			return
		}
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
		TaskID           *uuid.UUID `json:"task_id"`
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
			TaskID:           input.TaskID,
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

	// Delete associated task due dates
	for _, id := range input.IDs {
		var schedule models.ScheduledTask
		if err := tx.Where("id = ? AND user_id = ?", id, userIDUUID).First(&schedule).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue // Skip if not found or not owned by user
			}
			tx.Rollback()
			log.Println("Error finding scheduled task:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find scheduled tasks"})
			return
		}

		if schedule.TaskID != nil {
			if err := tx.Model(&models.Task{}).Where("id = ?", *schedule.TaskID).Update("due_date", nil).Error; err != nil {
				tx.Rollback()
				log.Println("Error updating task due date:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update associated tasks"})
				return
			}
		}
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

// CreateCalendarZone creates a new calendar zone
func CreateCalendarZone(c *gin.Context) {
	var input struct {
		Name            string     `json:"name" binding:"required"`
		Description     string     `json:"description"`
		Category        string     `json:"category" binding:"required"`
		Color           string     `json:"color"`
		StartTime       time.Time  `json:"start_time" binding:"required"`
		EndTime         time.Time  `json:"end_time" binding:"required"`
		DaysOfWeek      string     `json:"days_of_week"`
		Priority        int        `json:"priority"`
		IsActive        *bool      `json:"is_active"`
		AllowScheduling *bool      `json:"allow_scheduling"`
		MaxEventsPerDay *int       `json:"max_events_per_day"`
		IsRecurring     *bool      `json:"is_recurring"`
		RecurrenceStart *time.Time `json:"recurrence_start"`
		RecurrenceEnd   *time.Time `json:"recurrence_end"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid calendar zone input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Validate time range
	if input.StartTime.After(input.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
		return
	}

	// Set defaults
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	isRecurring := false
	if input.IsRecurring != nil {
		isRecurring = *input.IsRecurring
	}

	priority := 5
	if input.Priority > 0 && input.Priority <= 10 {
		priority = input.Priority
	}

	color := "#3b82f6"
	if input.Color != "" {
		color = input.Color
	}

	allowScheduling := false
	config.Logger.Debug(allowScheduling)
	if input.AllowScheduling != nil {
		config.Logger.Debug("Entered allowing")
		allowScheduling = *input.AllowScheduling
		config.Logger.Debug(allowScheduling)
	}

	zone := models.CalendarZone{
		UserID:          userIDUUID,
		Name:            input.Name,
		Description:     input.Description,
		Category:        input.Category,
		Color:           color,
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
		DaysOfWeek:      input.DaysOfWeek,
		Priority:        priority,
		IsActive:        isActive,
		AllowScheduling: allowScheduling,
		MaxEventsPerDay: input.MaxEventsPerDay,
		IsRecurring:     isRecurring,
		RecurrenceStart: input.RecurrenceStart,
		RecurrenceEnd:   input.RecurrenceEnd,
	}

	if err := config.GetDB().Create(&zone).Error; err != nil {
		config.Logger.Errorf("Error creating calendar zone: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create calendar zone"})
		return
	}

	config.Logger.Infof("Created calendar zone %s for user %s", zone.ID, userIDUUID)
	c.JSON(http.StatusCreated, zone)
}

// GetCalendarZones gets all calendar zones for the user
func GetCalendarZones(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var zones []models.CalendarZone
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Order("created_at DESC").Find(&zones).Error; err != nil {
		config.Logger.Errorf("Error fetching calendar zones for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch calendar zones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"zones": zones})
}

// UpdateCalendarZone updates a calendar zone
func UpdateCalendarZone(c *gin.Context) {
	zoneIDStr := c.Param("zoneID")
	zoneID, err := uuid.Parse(zoneIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid zone ID param: %s", zoneIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid zone ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify zone exists and belongs to user
	var zone models.CalendarZone
	if err := config.GetDB().Where("id = ? AND user_id = ?", zoneID, userIDUUID).First(&zone).Error; err != nil {
		config.Logger.Warnf("Calendar zone ID %s not found for user %s", zoneID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar zone not found"})
		return
	}

	var input struct {
		Name            *string    `json:"name"`
		Description     *string    `json:"description"`
		Category        *string    `json:"category"`
		Color           *string    `json:"color"`
		StartTime       *time.Time `json:"start_time"`
		EndTime         *time.Time `json:"end_time"`
		DaysOfWeek      *string    `json:"days_of_week"`
		Priority        *int       `json:"priority"`
		IsActive        *bool      `json:"is_active"`
		AllowScheduling *bool      `json:"allow_scheduling"`
		MaxEventsPerDay *int       `json:"max_events_per_day"`
		IsRecurring     *bool      `json:"is_recurring"`
		RecurrenceStart *time.Time `json:"recurrence_start"`
		RecurrenceEnd   *time.Time `json:"recurrence_end"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid calendar zone update input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate time range if both times are provided
	if input.StartTime != nil && input.EndTime != nil && input.StartTime.After(*input.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
		return
	}

	// Update fields
	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}
	if input.StartTime != nil {
		updates["start_time"] = *input.StartTime
	}
	if input.EndTime != nil {
		updates["end_time"] = *input.EndTime
	}
	if input.DaysOfWeek != nil {
		updates["days_of_week"] = *input.DaysOfWeek
	}
	if input.Priority != nil && *input.Priority > 0 && *input.Priority <= 10 {
		updates["priority"] = *input.Priority
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.AllowScheduling != nil {
		updates["allow_scheduling"] = *input.AllowScheduling
	}
	if input.MaxEventsPerDay != nil {
		updates["max_events_per_day"] = *input.MaxEventsPerDay
	}
	if input.IsRecurring != nil {
		updates["is_recurring"] = *input.IsRecurring
	}
	if input.RecurrenceStart != nil {
		updates["recurrence_start"] = *input.RecurrenceStart
	}
	if input.RecurrenceEnd != nil {
		updates["recurrence_end"] = *input.RecurrenceEnd
	}

	if len(updates) > 0 {
		if err := config.GetDB().Model(&zone).Updates(updates).Error; err != nil {
			config.Logger.Errorf("Error updating calendar zone: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update calendar zone"})
			return
		}
	}

	// Fetch updated zone
	if err := config.GetDB().Where("id = ?", zoneID).First(&zone).Error; err != nil {
		config.Logger.Errorf("Error fetching updated calendar zone: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch updated calendar zone"})
		return
	}

	config.Logger.Infof("Updated calendar zone %s for user %s", zoneID, userIDUUID)
	c.JSON(http.StatusOK, zone)
}

// DeleteCalendarZone deletes a calendar zone
func DeleteCalendarZone(c *gin.Context) {
	zoneIDStr := c.Param("zoneID")
	zoneID, err := uuid.Parse(zoneIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid zone ID param: %s", zoneIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid zone ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify zone exists and belongs to user
	var zone models.CalendarZone
	if err := config.GetDB().Where("id = ? AND user_id = ?", zoneID, userIDUUID).First(&zone).Error; err != nil {
		config.Logger.Warnf("Calendar zone ID %s not found for user %s", zoneID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Calendar zone not found"})
		return
	}

	if err := config.GetDB().Delete(&zone).Error; err != nil {
		config.Logger.Errorf("Error deleting calendar zone: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete calendar zone"})
		return
	}

	config.Logger.Infof("Deleted calendar zone %s for user %s", zoneID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Calendar zone deleted successfully"})
}

// GetZoneCategories returns available zone categories
func GetZoneCategories(c *gin.Context) {
	categories := models.GetDefaultZoneCategories()
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
