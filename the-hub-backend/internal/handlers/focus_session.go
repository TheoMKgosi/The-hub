package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StartFocusSessionRequest struct {
	TaskID          *uuid.UUID `json:"task_id" binding:"omitempty"`
	GoalID          *uuid.UUID `json:"goal_id" binding:"omitempty"`
	SessionType     string     `json:"session_type" binding:"omitempty"`      // focus, break
	PlannedDuration int        `json:"planned_duration" binding:"omitempty"`  // planned duration in minutes
	Notes           string     `json:"notes" binding:"omitempty"`
	CycleGroup      *uuid.UUID `json:"cycle_group" binding:"omitempty"`
	CycleOrder      int        `json:"cycle_order" binding:"omitempty"`
}

type StopFocusSessionRequest struct {
	Status string `json:"status" binding:"required,oneof=completed interrupted"`
}

func StartFocusSession(c *gin.Context) {
	var input StartFocusSessionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid focus session input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
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
		config.Logger.Errorf("Invalid userID type: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Auto-stop any currently active session for this user
	if err := config.GetDB().Model(&models.FocusSession{}).
		Where("user_id = ? AND status = ?", userIDUUID, "active").
		Updates(map[string]interface{}{
			"status":    "interrupted",
			"ended_at":  time.Now(),
			"updated_at": time.Now(),
		}).Error; err != nil {
		config.Logger.Errorf("Failed to stop active session for user %s: %v", userIDUUID, err)
	}

	if input.TaskID != nil {
		var task models.Task
		if err := config.GetDB().Where("id = ? AND user_id = ?", input.TaskID, userIDUUID).First(&task).Error; err != nil {
			config.Logger.Warnf("Task %s not found for user %s", *input.TaskID, userIDUUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found or access denied"})
			return
		}
	}

	if input.GoalID != nil {
		var goal models.Goal
		if err := config.GetDB().Where("id = ? AND user_id = ?", input.GoalID, userIDUUID).First(&goal).Error; err != nil {
			config.Logger.Warnf("Goal %s not found for user %s", *input.GoalID, userIDUUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Goal not found or access denied"})
			return
		}
	}

	sessionType := input.SessionType
	if sessionType == "" {
		sessionType = "focus"
	}

	session := models.FocusSession{
		UserID:      userIDUUID,
		TaskID:      input.TaskID,
		GoalID:      input.GoalID,
		DurationMin: 0,
		StartedAt:   time.Now(),
		Status:      "active",
		SessionType: sessionType,
		CycleGroup:  input.CycleGroup,
		CycleOrder:  input.CycleOrder,
		Notes:       input.Notes,
	}

	config.Logger.Infof("Starting %s session for user %s", sessionType, userIDUUID)
	if err := config.GetDB().Create(&session).Error; err != nil {
		config.Logger.Errorf("Failed to create focus session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create focus session"})
		return
	}

	config.Logger.Infof("Created focus session %s for user %s", session.ID, userIDUUID)
	c.JSON(http.StatusCreated, session)
}

func StopFocusSession(c *gin.Context) {
	sessionIDStr := c.Param("ID")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid session ID: %s", sessionIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
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
		config.Logger.Errorf("Invalid userID type: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var input StopFocusSessionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid stop session input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var session models.FocusSession
	if err := config.GetDB().Where("id = ? AND user_id = ?", sessionID, userIDUUID).First(&session).Error; err != nil {
		config.Logger.Warnf("Focus session %s not found for user %s", sessionID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.Status != "active" {
		config.Logger.Warnf("Session %s is not active (status: %s)", sessionID, session.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is not active"})
		return
	}

	now := time.Now()
	duration := int(now.Sub(session.StartedAt).Minutes())

	updates := map[string]interface{}{
		"status":       input.Status,
		"ended_at":     now,
		"duration_min": duration,
		"updated_at":   now,
	}

	config.Logger.Infof("Stopping focus session %s for user %s: status=%s, duration=%dmin", sessionID, userIDUUID, input.Status, duration)
	if err := config.GetDB().Model(&session).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to stop focus session %s: %v", sessionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not stop focus session"})
		return
	}

	session.Status = input.Status
	session.EndedAt = &now
	session.DurationMin = duration

	config.Logger.Infof("Successfully stopped focus session %s", sessionID)
	c.JSON(http.StatusOK, session)
}

func GetFocusSessions(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	query := config.GetDB().Where("user_id = ?", userIDUUID)

	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		if taskID, err := uuid.Parse(taskIDStr); err == nil {
			query = query.Where("task_id = ?", taskID)
		}
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if sessionType := c.Query("session_type"); sessionType != "" {
		query = query.Where("session_type = ?", sessionType)
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query = query.Where("started_at >= ?", parsedDate)
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			query = query.Where("started_at <= ?", parsedDate.Add(24*time.Hour-time.Second))
		}
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	var sessions []models.FocusSession
	if err := query.Order("started_at DESC").Limit(limit).Find(&sessions).Error; err != nil {
		config.Logger.Errorf("Failed to fetch focus sessions for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch sessions"})
		return
	}

	config.Logger.Infof("Found %d focus sessions for user %s", len(sessions), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"focus_sessions": sessions})
}

func GetFocusSessionStats(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	sessionType := c.DefaultQuery("session_type", "focus")

	startDate := time.Now().AddDate(0, 0, -days)

	var totalMinutes int64
	if err := config.GetDB().Model(&models.FocusSession{}).
		Where("user_id = ? AND started_at >= ? AND session_type = ? AND status IN ('completed', 'interrupted')", userIDUUID, startDate, sessionType).
		Select("COALESCE(SUM(duration_min), 0)").
		Scan(&totalMinutes).Error; err != nil {
		config.Logger.Errorf("Failed to calculate total focus time: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}

	type DailyStats struct {
		Date          time.Time `json:"date"`
		Minutes       int       `json:"minutes"`
		SessionsCount int       `json:"sessions_count"`
	}

	var dailyStats []DailyStats
	if err := config.GetDB().Model(&models.FocusSession{}).
		Select("DATE(started_at) as date, COALESCE(SUM(duration_min), 0) as minutes, COUNT(*) as sessions_count").
		Where("user_id = ? AND started_at >= ? AND session_type = ? AND status IN ('completed', 'interrupted')", userIDUUID, startDate, sessionType).
		Group("DATE(started_at)").
		Order("date").
		Scan(&dailyStats).Error; err != nil {
		config.Logger.Errorf("Failed to fetch daily stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch daily stats"})
		return
	}

	type TaskBreakdown struct {
		TaskID        *uuid.UUID `json:"task_id"`
		TaskTitle     string     `json:"task_title"`
		Minutes       int        `json:"minutes"`
		SessionsCount int        `json:"sessions_count"`
	}

	var taskBreakdown []TaskBreakdown
	if err := config.GetDB().Model(&models.FocusSession{}).
		Select("focus_sessions.task_id, COALESCE(tasks.title, 'No Task') as task_title, COALESCE(SUM(focus_sessions.duration_min), 0) as minutes, COUNT(*) as sessions_count").
		Joins("LEFT JOIN tasks ON focus_sessions.task_id = tasks.id").
		Where("focus_sessions.user_id = ? AND focus_sessions.started_at >= ? AND focus_sessions.session_type = ? AND focus_sessions.status IN ('completed', 'interrupted')", userIDUUID, startDate, sessionType).
		Group("focus_sessions.task_id, tasks.title").
		Order("minutes DESC").
		Scan(&taskBreakdown).Error; err != nil {
		config.Logger.Errorf("Failed to fetch task breakdown: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch task breakdown"})
		return
	}

	type WeeklyStats struct {
		WeekStart     string `json:"week_start"`
		Minutes       int    `json:"minutes"`
		SessionsCount int    `json:"sessions_count"`
	}

	var weeklyStats []WeeklyStats
	if err := config.GetDB().Model(&models.FocusSession{}).
		Select("DATE_TRUNC('week', started_at)::date as week_start, COALESCE(SUM(duration_min), 0) as minutes, COUNT(*) as sessions_count").
		Where("user_id = ? AND started_at >= ? AND session_type = ? AND status IN ('completed', 'interrupted')", userIDUUID, startDate, sessionType).
		Group("DATE_TRUNC('week', started_at)").
		Order("week_start").
		Scan(&weeklyStats).Error; err != nil {
		config.Logger.Errorf("Failed to fetch weekly stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch weekly stats"})
		return
	}

	var totalSessions int64
	if err := config.GetDB().Model(&models.FocusSession{}).
		Where("user_id = ? AND started_at >= ? AND session_type = ? AND status IN ('completed', 'interrupted')", userIDUUID, startDate, sessionType).
		Select("COUNT(*)").
		Scan(&totalSessions).Error; err != nil {
		config.Logger.Errorf("Failed to count sessions: %v", err)
	}

	stats := map[string]interface{}{
		"total_minutes":   totalMinutes,
		"total_hours":     float64(totalMinutes) / 60,
		"total_sessions":  totalSessions,
		"days":            days,
		"daily_stats":     dailyStats,
		"task_breakdown":  taskBreakdown,
		"weekly_stats":    weeklyStats,
		"average_daily":   float64(totalMinutes) / float64(days),
	}

	config.Logger.Infof("Calculated focus session stats for user %s over %d days", userIDUUID, days)
	c.JSON(http.StatusOK, stats)
}
