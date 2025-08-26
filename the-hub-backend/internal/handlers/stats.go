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

// GetTaskStats godoc
// @Summary      Get task statistics for the logged-in user
// @Description  Fetch aggregated task statistics for the logged-in user with optional date range
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period  query     string  false  "Time period (today, week, month, year)"  default(today)
// @Param        start_date  query     string  false  "Start date in YYYY-MM-DD format"
// @Param        end_date  query     string  false  "End date in YYYY-MM-DD format"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /stats/tasks [get]
func GetTaskStats(c *gin.Context) {
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

	// Parse query parameters
	period := c.DefaultQuery("period", "today")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	now := time.Now()

	// Calculate date range based on period or explicit dates
	if startDateStr != "" && endDateStr != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
	} else {
		// Use period-based date ranges
		switch period {
		case "today":
			startDate = now.Truncate(24 * time.Hour)
			endDate = startDate.Add(24 * time.Hour)
		case "week":
			// Start of week (Monday)
			days := int(now.Weekday()-time.Monday) % 7
			startDate = now.AddDate(0, 0, -days).Truncate(24 * time.Hour)
			endDate = startDate.AddDate(0, 0, 7)
		case "month":
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(0, 1, 0)
		case "year":
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(1, 0, 0)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period. Use 'today', 'week', 'month', or 'year'"})
			return
		}
	}

	config.Logger.Infof("Fetching task stats for user %s from %s to %s", userIDUUID, startDate, endDate)

	// Calculate real-time statistics
	stats, err := calculateTaskStats(userIDUUID, startDate, endDate)
	if err != nil {
		config.Logger.Errorf("Error calculating task stats for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate task statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// calculateTaskStats computes real-time statistics for the given user and date range
func calculateTaskStats(userID uuid.UUID, startDate, endDate time.Time) (map[string]interface{}, error) {
	db := config.GetDB()

	// Get all tasks for the user within the date range
	var tasks []models.Task
	if err := db.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startDate, endDate).Find(&tasks).Error; err != nil {
		return nil, err
	}

	// Get all tasks for completion analysis (regardless of creation date)
	var allTasks []models.Task
	if err := db.Where("user_id = ?", userID).Find(&allTasks).Error; err != nil {
		return nil, err
	}

	// Initialize stats
	stats := map[string]interface{}{
		"period": map[string]string{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
		"summary": map[string]interface{}{
			"total_tasks":     len(tasks),
			"completed_tasks": 0,
			"pending_tasks":   0,
			"overdue_tasks":   0,
		},
		"priority_distribution": map[string]int{
			"1": 0, "2": 0, "3": 0, "4": 0, "5": 0,
		},
		"priority_completion": map[string]int{
			"1": 0, "2": 0, "3": 0, "4": 0, "5": 0,
		},
		"goal_distribution": map[string]int{
			"with_goals":    0,
			"without_goals": 0,
		},
		"time_based": map[string]int{
			"due_today":     0,
			"due_tomorrow":  0,
			"due_this_week": 0,
		},
	}

	// Calculate current date references
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	tomorrow := today.AddDate(0, 0, 1)
	weekFromNow := today.AddDate(0, 0, 7)

	// Analyze tasks
	for _, task := range allTasks {
		// Count by status
		if task.Status == "complete" {
			stats["summary"].(map[string]interface{})["completed_tasks"] = stats["summary"].(map[string]interface{})["completed_tasks"].(int) + 1
		} else {
			stats["summary"].(map[string]interface{})["pending_tasks"] = stats["summary"].(map[string]interface{})["pending_tasks"].(int) + 1

			// Check if overdue
			if task.DueDate != nil && task.DueDate.Before(now) {
				stats["summary"].(map[string]interface{})["overdue_tasks"] = stats["summary"].(map[string]interface{})["overdue_tasks"].(int) + 1
			}
		}

		// Priority distribution
		if task.Priority != nil {
			priority := *task.Priority
			if priority >= 1 && priority <= 5 {
				priorityStr := string(rune('0' + priority))
				stats["priority_distribution"].(map[string]int)[priorityStr]++

				// Priority completion
				if task.Status == "complete" {
					stats["priority_completion"].(map[string]int)[priorityStr]++
				}
			}
		}

		// Goal distribution
		if task.GoalID != nil {
			stats["goal_distribution"].(map[string]int)["with_goals"]++
		} else {
			stats["goal_distribution"].(map[string]int)["without_goals"]++
		}

		// Time-based analysis
		if task.DueDate != nil {
			dueDate := *task.DueDate
			if dueDate.Truncate(24 * time.Hour).Equal(today) {
				stats["time_based"].(map[string]int)["due_today"]++
			} else if dueDate.Truncate(24 * time.Hour).Equal(tomorrow) {
				stats["time_based"].(map[string]int)["due_tomorrow"]++
			} else if dueDate.After(today) && dueDate.Before(weekFromNow) {
				stats["time_based"].(map[string]int)["due_this_week"]++
			}
		}
	}

	// Calculate completion rate
	totalTasks := stats["summary"].(map[string]interface{})["total_tasks"].(int)
	completedTasks := stats["summary"].(map[string]interface{})["completed_tasks"].(int)
	if totalTasks > 0 {
		completionRate := float64(completedTasks) / float64(totalTasks) * 100
		stats["summary"].(map[string]interface{})["completion_rate"] = completionRate
	} else {
		stats["summary"].(map[string]interface{})["completion_rate"] = 0.0
	}

	// Calculate priority completion rates
	for i := 1; i <= 5; i++ {
		priorityStr := string(rune('0' + i))
		total := stats["priority_distribution"].(map[string]int)[priorityStr]
		completed := stats["priority_completion"].(map[string]int)[priorityStr]

		if total > 0 {
			rate := float64(completed) / float64(total) * 100
			stats["priority_completion"].(map[string]int)[priorityStr] = int(rate) // Store as percentage
		}
	}

	return stats, nil
}

// GetTaskStatsTrends godoc
// @Summary      Get task statistics trends over time
// @Description  Fetch task statistics trends for the logged-in user over a specified period
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        days  query     int  false  "Number of days to look back"  default(30)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /stats/tasks/trends [get]
func GetTaskStatsTrends(c *gin.Context) {
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

	days := 30 // Default to 30 days
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	// Calculate trends for each day in the period
	trends := make([]map[string]interface{}, 0, days)
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startDate := date.Truncate(24 * time.Hour)
		endDate := startDate.Add(24 * time.Hour)

		stats, err := calculateTaskStats(userIDUUID, startDate, endDate)
		if err != nil {
			config.Logger.Errorf("Error calculating trends for date %s: %v", date.Format("2006-01-02"), err)
			continue
		}

		// Add date to the stats
		stats["date"] = date.Format("2006-01-02")
		trends = append(trends, stats)
	}

	c.JSON(http.StatusOK, gin.H{"trends": trends})
}
