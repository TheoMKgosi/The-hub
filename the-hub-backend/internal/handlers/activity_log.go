package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CreateActivityRequest struct {
	Title        string   `json:"title" binding:"required"`
	Notes        string   `json:"notes" binding:"omitempty"`
	Category     string   `json:"category" binding:"omitempty"`
	Tags         []string `json:"tags" binding:"omitempty"`
	ActivityDate string   `json:"activity_date" binding:"omitempty"` // YYYY-MM-DD, defaults to today
	DurationMin  int      `json:"duration_min" binding:"omitempty,min=0"`
}

type UpdateActivityRequest struct {
	Title        *string   `json:"title" binding:"omitempty"`
	Notes        *string   `json:"notes" binding:"omitempty"`
	Category     *string   `json:"category" binding:"omitempty"`
	Tags         *[]string `json:"tags" binding:"omitempty"`
	ActivityDate *string   `json:"activity_date" binding:"omitempty"`
	DurationMin  *int      `json:"duration_min" binding:"omitempty,min=0"`
}

// NormalizeTags trims, drops empty and duplicate entries (case-insensitive,
// keeping the first occurrence's casing).
func NormalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if seen[key] {
			continue
		}
		seen[key] = true
		normalized = append(normalized, trimmed)
	}
	return normalized
}

// ComputeStreak returns the number of consecutive days ending today (or
// yesterday, if nothing was logged today) present in the given dates.
// Dates are compared by calendar day in the caller's timezone, matching
// how activity_date is stored (a local calendar date, not a UTC instant).
func ComputeStreak(dates []time.Time, today time.Time) int {
	set := make(map[string]bool, len(dates))
	for _, d := range dates {
		set[d.Format("2006-01-02")] = true
	}

	offset := 0
	if !set[today.Format("2006-01-02")] {
		// Nothing logged today yet: a streak ending yesterday still counts.
		offset = 1
	}

	streak := 0
	for {
		day := today.AddDate(0, 0, -offset).Format("2006-01-02")
		if !set[day] {
			break
		}
		streak++
		offset++
	}
	return streak
}

func parseActivityDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return time.Parse("2006-01-02", dateStr)
}

func CreateActivity(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	var input CreateActivityRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Sugar().Warnf("Invalid activity input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	activityDate, err := parseActivityDate(input.ActivityDate)
	if err != nil {
		config.Logger.Sugar().Warnf("Invalid activity date %q: %v", input.ActivityDate, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity_date, expected YYYY-MM-DD"})
		return
	}

	entry := models.ActivityEntry{
		UserID:       userID,
		Title:        input.Title,
		Notes:        input.Notes,
		Category:     strings.TrimSpace(input.Category),
		Tags:         pq.StringArray(NormalizeTags(input.Tags)),
		ActivityDate: activityDate,
		DurationMin:  input.DurationMin,
	}

	if err := config.GetDB().Create(&entry).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to create activity entry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create activity entry"})
		return
	}

	config.Logger.Sugar().Infof("Created activity entry %s for user %s", entry.ID, userID)
	c.JSON(http.StatusCreated, entry)
}

func GetActivities(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	query := config.GetDB().Where("user_id = ?", userID)

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			query = query.Where("activity_date >= ?", parsedDate)
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			query = query.Where("activity_date <= ?", parsedDate)
		}
	}

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tags @> ARRAY[?]", tag)
	}

	if search := c.Query("search"); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR notes ILIKE ?", pattern, pattern)
	}

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	var entries []models.ActivityEntry
	if err := query.Order("activity_date DESC, created_at DESC").Limit(limit).Offset(offset).
		Find(&entries).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to fetch activities for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch activities"})
		return
	}

	var total int64
	if err := config.GetDB().Model(&models.ActivityEntry{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to count activities for user %s: %v", userID, err)
	}

	config.Logger.Sugar().Infof("Found %d activities for user %s", len(entries), userID)
	c.JSON(http.StatusOK, gin.H{"activities": entries, "total": total})
}

func UpdateActivity(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	entryID, err := uuid.Parse(c.Param("ID"))
	if err != nil {
		config.Logger.Sugar().Warnf("Invalid activity ID: %s", c.Param("ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	var entry models.ActivityEntry
	if err := config.GetDB().Where("id = ? AND user_id = ?", entryID, userID).First(&entry).Error; err != nil {
		config.Logger.Sugar().Warnf("Activity %s not found for user %s", entryID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	var input UpdateActivityRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Sugar().Warnf("Invalid activity update input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}
	if input.Category != nil {
		updates["category"] = strings.TrimSpace(*input.Category)
	}
	if input.DurationMin != nil {
		updates["duration_min"] = *input.DurationMin
	}
	if input.ActivityDate != nil {
		activityDate, err := time.Parse("2006-01-02", *input.ActivityDate)
		if err != nil {
			config.Logger.Sugar().Warnf("Invalid activity date %q: %v", *input.ActivityDate, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity_date, expected YYYY-MM-DD"})
			return
		}
		updates["activity_date"] = activityDate
	}

	if err := config.GetDB().Model(&entry).Updates(updates).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to update activity %s: %v", entryID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update activity"})
		return
	}

	// Tags go through the model field (not a map update) so GORM binds the
	// text[] column type correctly.
	if input.Tags != nil {
		entry.Tags = pq.StringArray(NormalizeTags(*input.Tags))
		if err := config.GetDB().Model(&entry).Select("tags").Updates(&entry).Error; err != nil {
			config.Logger.Sugar().Errorf("Failed to update activity tags %s: %v", entryID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update activity"})
			return
		}
	}

	if err := config.GetDB().Where("id = ? AND user_id = ?", entryID, userID).First(&entry).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to reload activity %s: %v", entryID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update activity"})
		return
	}

	config.Logger.Sugar().Infof("Updated activity %s for user %s", entryID, userID)
	c.JSON(http.StatusOK, entry)
}

func DeleteActivity(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	entryID, err := uuid.Parse(c.Param("ID"))
	if err != nil {
		config.Logger.Sugar().Warnf("Invalid activity ID: %s", c.Param("ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	result := config.GetDB().Where("id = ? AND user_id = ?", entryID, userID).Delete(&models.ActivityEntry{})
	if result.Error != nil {
		config.Logger.Sugar().Errorf("Failed to delete activity %s: %v", entryID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete activity"})
		return
	}
	if result.RowsAffected == 0 {
		config.Logger.Sugar().Warnf("Activity %s not found for user %s", entryID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	config.Logger.Sugar().Infof("Deleted activity %s for user %s", entryID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted"})
}

func GetActivityStats(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	startDate := time.Now().AddDate(0, 0, -days)
	// Each aggregate builds a fresh query: reusing one *gorm.DB across
	// multiple Select/Scan calls leaks clauses between them.
	baseScope := func() *gorm.DB {
		return config.GetDB().Model(&models.ActivityEntry{}).
			Where("user_id = ? AND activity_date >= ?", userID, startDate)
	}

	var totalEntries, totalMinutes int64
	if err := baseScope().Select("COUNT(*)").Scan(&totalEntries).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to count activities: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}
	if err := baseScope().Select("COALESCE(SUM(duration_min), 0)").Scan(&totalMinutes).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to sum activity minutes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}

	type CategoryBreakdown struct {
		Category string `json:"category"`
		Minutes  int    `json:"minutes"`
		Entries  int    `json:"entries"`
	}
	var categories []CategoryBreakdown
	if err := baseScope().Select("COALESCE(NULLIF(category, ''), 'Uncategorized') as category, COALESCE(SUM(duration_min), 0) as minutes, COUNT(*) as entries").
		Group("COALESCE(NULLIF(category, ''), 'Uncategorized')").
		Order("minutes DESC").
		Scan(&categories).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to fetch category breakdown: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}

	type DailyStats struct {
		Date    string `json:"date"`
		Minutes int    `json:"minutes"`
		Entries int    `json:"entries"`
	}
	var dailyStats []DailyStats
	if err := baseScope().Select("TO_CHAR(activity_date, 'YYYY-MM-DD') as date, COALESCE(SUM(duration_min), 0) as minutes, COUNT(*) as entries").
		Group("activity_date").
		Order("activity_date").
		Scan(&dailyStats).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to fetch daily stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}

	type WeeklyStats struct {
		WeekStart string `json:"week_start"`
		Minutes   int    `json:"minutes"`
		Entries   int    `json:"entries"`
	}
	var weeklyStats []WeeklyStats
	if err := baseScope().Select("TO_CHAR(DATE_TRUNC('week', activity_date), 'YYYY-MM-DD') as week_start, COALESCE(SUM(duration_min), 0) as minutes, COUNT(*) as entries").
		Group("DATE_TRUNC('week', activity_date)").
		Order("week_start").
		Scan(&weeklyStats).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to fetch weekly stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}

	type TagUsage struct {
		Tag     string `json:"tag"`
		Minutes int    `json:"minutes"`
		Count   int    `json:"count"`
	}
	var topTags []TagUsage
	if err := config.GetDB().Model(&models.ActivityEntry{}).
		Select("tag as tag, COALESCE(SUM(duration_min), 0) as minutes, COUNT(*) as count").
		Joins(", unnest(tags) AS tag").
		Where("user_id = ? AND activity_date >= ?", userID, startDate).
		Group("tag").
		Order("count DESC").
		Limit(10).
		Scan(&topTags).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to fetch top tags: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}

	var loggedDates []time.Time
	if err := baseScope().Distinct().Select("activity_date").Scan(&loggedDates).Error; err != nil {
		config.Logger.Sugar().Errorf("Failed to fetch logged dates: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate stats"})
		return
	}

	stats := map[string]interface{}{
		"days":           days,
		"total_entries":  totalEntries,
		"total_minutes":  totalMinutes,
		"total_hours":    float64(totalMinutes) / 60,
		"categories":     categories,
		"daily_stats":    dailyStats,
		"weekly_stats":   weeklyStats,
		"top_tags":       topTags,
		"current_streak": ComputeStreak(loggedDates, time.Now()),
	}

	config.Logger.Sugar().Infof("Calculated activity stats for user %s over %d days", userID, days)
	c.JSON(http.StatusOK, stats)
}

// requireUserID extracts the authenticated userID from the request context,
// writing the error response itself when missing or malformed.
func requireUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Sugar().Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return uuid.Nil, false
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Sugar().Errorf("Invalid userID type: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return uuid.Nil, false
	}
	return userIDUUID, true
}
