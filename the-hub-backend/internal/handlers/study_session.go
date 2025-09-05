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

// CreateStudySessionRequest represents the request body for creating a study session
type CreateStudySessionRequest struct {
	TopicID     *uuid.UUID `json:"topic_id" binding:"omitempty"`
	TaskID      *uuid.UUID `json:"task_id" binding:"omitempty"`
	DurationMin int        `json:"duration_min" binding:"required,min=1"`
	Notes       string     `json:"notes"`
}

// CreateStudySession godoc
// @Summary      Create a new study session
// @Description  Create a new study session for the logged-in user
// @Tags         study-sessions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        study_session  body      CreateStudySessionRequest  true  "Study session creation data"
// @Success      201  {object}  models.StudySession
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /study-sessions [post]
func CreateStudySession(c *gin.Context) {
	var input CreateStudySessionRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid study session input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for study session", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during study session creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Validate that topic exists and belongs to user (if provided)
	if input.TopicID != nil {
		var topic models.Topic
		if err := config.GetDB().Where("id = ? AND user_id = ?", input.TopicID, userIDUUID).First(&topic).Error; err != nil {
			config.Logger.Warnf("Topic ID %s not found or not owned by user %s", input.TopicID, userIDUUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Topic not found or access denied"})
			return
		}
	}

	// Validate that task exists and belongs to user's topic (if provided)
	if input.TaskID != nil {
		var task models.Task_learning
		if err := config.GetDB().Joins("JOIN topics ON task_learnings.topic_id = topics.id").
			Where("task_learnings.id = ? AND topics.user_id = ?", input.TaskID, userIDUUID).
			First(&task).Error; err != nil {
			config.Logger.Warnf("Task ID %s not found or not owned by user %s", input.TaskID, userIDUUID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found or access denied"})
			return
		}
	}

	studySession := models.StudySession{
		UserID:      userIDUUID,
		TopicID:     input.TopicID,
		TaskID:      input.TaskID,
		DurationMin: input.DurationMin,
		StartedAt:   time.Now(),
		EndedAt:     time.Now(), // Will be updated when session ends
	}

	config.Logger.Infof("Creating study session for user %s: %d minutes", userIDUUID, input.DurationMin)
	if err := config.GetDB().Create(&studySession).Error; err != nil {
		config.Logger.Errorf("Error creating study session for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create study session"})
		return
	}

	config.Logger.Infof("Successfully created study session ID %s for user %s", studySession.ID, userIDUUID)
	c.JSON(http.StatusCreated, studySession)
}

// GetStudySessions godoc
// @Summary      Get study sessions
// @Description  Fetch study sessions for the logged-in user with optional filtering
// @Tags         study-sessions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        topic_id    query     string  false  "Filter by topic ID"
// @Param        task_id     query     string  false  "Filter by task ID"
// @Param        date_from   query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to     query     string  false  "Filter to date (YYYY-MM-DD)"
// @Param        limit       query     int     false  "Limit number of results"  default(50)
// @Success      200  {object}  map[string][]models.StudySession
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /study-sessions [get]
func GetStudySessions(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Build query
	query := config.GetDB().Where("user_id = ?", userIDUUID)

	// Apply filters
	if topicIDStr := c.Query("topic_id"); topicIDStr != "" {
		if topicID, err := uuid.Parse(topicIDStr); err == nil {
			query = query.Where("topic_id = ?", topicID)
		}
	}

	if taskIDStr := c.Query("task_id"); taskIDStr != "" {
		if taskID, err := uuid.Parse(taskIDStr); err == nil {
			query = query.Where("task_id = ?", taskID)
		}
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

	// Apply limit
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	var studySessions []models.StudySession
	if err := query.Order("started_at DESC").Limit(limit).Find(&studySessions).Error; err != nil {
		config.Logger.Errorf("Error fetching study sessions for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch study sessions"})
		return
	}

	config.Logger.Infof("Found %d study sessions for user %s", len(studySessions), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"study_sessions": studySessions})
}

// GetStudySessionStats godoc
// @Summary      Get study session statistics
// @Description  Get aggregated statistics for the user's study sessions
// @Tags         study-sessions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  false  "Number of days to look back"  default(30)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /study-sessions/stats [get]
func GetStudySessionStats(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	startDate := time.Now().AddDate(0, 0, -days)

	// Get total study time
	var totalMinutes int64
	if err := config.GetDB().Model(&models.StudySession{}).
		Where("user_id = ? AND started_at >= ?", userIDUUID, startDate).
		Select("COALESCE(SUM(duration_min), 0)").
		Scan(&totalMinutes).Error; err != nil {
		config.Logger.Errorf("Error calculating total study time for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate study statistics"})
		return
	}

	// Get daily study time for the last N days
	type DailyStats struct {
		Date    time.Time `json:"date"`
		Minutes int       `json:"minutes"`
	}

	var dailyStats []DailyStats
	if err := config.GetDB().Model(&models.StudySession{}).
		Select("DATE(started_at) as date, COALESCE(SUM(duration_min), 0) as minutes").
		Where("user_id = ? AND started_at >= ?", userIDUUID, startDate).
		Group("DATE(started_at)").
		Order("date").
		Scan(&dailyStats).Error; err != nil {
		config.Logger.Errorf("Error fetching daily stats for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch daily statistics"})
		return
	}

	// Get sessions by topic
	type TopicStats struct {
		TopicID    *uuid.UUID `json:"topic_id"`
		TopicTitle string     `json:"topic_title"`
		Minutes    int        `json:"minutes"`
		Sessions   int        `json:"sessions"`
	}

	var topicStats []TopicStats
	if err := config.GetDB().Model(&models.StudySession{}).
		Select("study_sessions.topic_id, topics.title as topic_title, COALESCE(SUM(study_sessions.duration_min), 0) as minutes, COUNT(*) as sessions").
		Joins("LEFT JOIN topics ON study_sessions.topic_id = topics.id").
		Where("study_sessions.user_id = ? AND study_sessions.started_at >= ?", userIDUUID, startDate).
		Group("study_sessions.topic_id, topics.title").
		Order("minutes DESC").
		Scan(&topicStats).Error; err != nil {
		config.Logger.Errorf("Error fetching topic stats for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch topic statistics"})
		return
	}

	stats := map[string]interface{}{
		"total_minutes": totalMinutes,
		"total_hours":   float64(totalMinutes) / 60,
		"days":          days,
		"daily_stats":   dailyStats,
		"topic_stats":   topicStats,
		"average_daily": float64(totalMinutes) / float64(days),
	}

	config.Logger.Infof("Calculated study statistics for user %s over %d days", userIDUUID, days)
	c.JSON(http.StatusOK, stats)
}
