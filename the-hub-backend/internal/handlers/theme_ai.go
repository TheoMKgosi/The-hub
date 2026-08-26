package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/TheoMKgosi/The-hub/internal/ai"
	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const maxTasksForThemeAI = 50

// GetAIThemeSuggestions godoc
// @Summary      Suggest themes from tasks
// @Description  Ask the AI to group the user's pending tasks into themes and propose a weekday assignment (one theme per day)
// @Tags         themes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]ai.ThemeSuggestion
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /themes/ai/suggest [post]
func GetAIThemeSuggestions(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var tasks []models.Task
	if err := config.GetDB().
		Where("user_id = ? AND status != 'completed' AND parent_task_id IS NULL", userIDUUID).
		Order("created_at ASC").
		Limit(maxTasksForThemeAI).
		Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Failed to fetch tasks for theme suggestions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"suggestions": []ai.ThemeSuggestion{}})
		return
	}

	client, err := ai.GetOpenRouterClient()
	if err != nil {
		config.Logger.Errorf("Failed to get AI client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable"})
		return
	}

	aiResponse, err := client.SuggestThemes(tasks)
	if err != nil {
		config.Logger.Errorf("Failed to suggest themes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI theme suggestions"})
		return
	}

	var suggestions []ai.ThemeSuggestion
	if err := json.Unmarshal([]byte(aiResponse), &suggestions); err != nil {
		start := strings.Index(aiResponse, "[")
		end := strings.LastIndex(aiResponse, "]")
		if start == -1 || end == -1 || start >= end {
			config.Logger.Debug(aiResponse)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI theme suggestions"})
			return
		}
		if err := json.Unmarshal([]byte(aiResponse[start:end+1]), &suggestions); err != nil {
			config.Logger.Errorf("Failed to parse AI theme suggestions: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI theme suggestions"})
			return
		}
	}

	// Validate and clean up suggestions
	validTaskIDs := make(map[string]bool, len(tasks))
	for _, task := range tasks {
		validTaskIDs[task.ID.String()] = true
	}

	usedDays := make(map[string]bool)
	cleaned := make([]ai.ThemeSuggestion, 0, len(suggestions))
	for _, s := range suggestions {
		s.Name = strings.TrimSpace(s.Name)
		s.Days = normalizeThemeDays(s.Days)

		if s.Name == "" || len(s.Days) == 0 {
			continue
		}
		if err := validateThemeDays(s.Days); err != nil {
			config.Logger.Warnf("AI theme %s returned invalid days: %v", s.Name, err)
			continue
		}

		// Enforce one theme per day across suggestions
		dayOverlap := false
		for _, day := range s.Days {
			if usedDays[day] {
				dayOverlap = true
				break
			}
		}
		if dayOverlap {
			continue
		}
		for _, day := range s.Days {
			usedDays[day] = true
		}

		// Keep only known task IDs
		taskIDs := make([]string, 0, len(s.Tasks))
		for _, taskID := range s.Tasks {
			if validTaskIDs[taskID] {
				taskIDs = append(taskIDs, taskID)
			}
		}
		s.Tasks = taskIDs

		cleaned = append(cleaned, s)
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": cleaned})
}

// ApplyThemeRequest represents the request body for applying AI theme suggestions.
type ApplyThemeRequest struct {
	Themes      []ApplyTheme      `json:"themes"`
	Assignments []ThemeAssignment `json:"assignments"`
}

// ApplyTheme describes a theme to create.
type ApplyTheme struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Color       string   `json:"color"`
	Days        []string `json:"days"`
}

// ThemeAssignment links a task to a theme by name.
type ThemeAssignment struct {
	TaskID    string `json:"task_id" binding:"required"`
	ThemeName string `json:"theme_name" binding:"required"`
}

// ApplyAIThemes godoc
// @Summary      Apply AI theme suggestions
// @Description  Create the selected themes and assign tasks to them
// @Tags         themes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      ApplyThemeRequest  true  "Themes to create and task assignments"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /themes/ai/apply [post]
func ApplyAIThemes(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var req ApplyThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		config.Logger.Warnf("Invalid apply themes input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if len(req.Themes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No themes to apply"})
		return
	}

	db := config.GetDB()

	// Validate the requested themes and enforce one theme per day
	usedDays := make(map[string]bool)
	themeNames := make(map[string]bool)
	for _, theme := range req.Themes {
		name := strings.TrimSpace(theme.Name)
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Theme name is required"})
			return
		}
		if themeNames[name] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Duplicate theme name in request: %s", name)})
			return
		}
		themeNames[name] = true

		days := normalizeThemeDays(theme.Days)
		if err := validateThemeDays(days); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, day := range days {
			if usedDays[day] {
				c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Day '%s' is assigned to more than one theme", day)})
				return
			}
			usedDays[day] = true
		}
	}

	// Validate against existing themes
	for _, theme := range req.Themes {
		days := normalizeThemeDays(theme.Days)
		conflicts, err := themeDayConflict(db, userIDUUID, days, nil)
		if err != nil {
			config.Logger.Errorf("Failed to check theme day conflicts for user %s: %v", userIDUUID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not validate theme days"})
			return
		}
		if len(conflicts) > 0 {
			conflictNames := make([]string, len(conflicts))
			for i, conflict := range conflicts {
				conflictNames[i] = conflict.Name
			}
			c.JSON(http.StatusConflict, gin.H{
				"error":     "One or more days already belong to another theme",
				"conflicts": conflictNames,
			})
			return
		}
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	themeIDByName := make(map[string]uuid.UUID)
	for _, theme := range req.Themes {
		// Reuse an existing theme with the same name when present
		var existing models.Theme
		if err := tx.Where("user_id = ? AND name = ?", userIDUUID, strings.TrimSpace(theme.Name)).First(&existing).Error; err == nil {
			themeIDByName[existing.Name] = existing.ID
			continue
		}

		newTheme := models.Theme{
			UserID:      userIDUUID,
			Name:        strings.TrimSpace(theme.Name),
			Description: theme.Description,
			Color:       theme.Color,
			Days:        pq.StringArray(normalizeThemeDays(theme.Days)),
		}
		if err := tx.Create(&newTheme).Error; err != nil {
			tx.Rollback()
			config.Logger.Errorf("Failed to create theme '%s': %v", newTheme.Name, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create themes"})
			return
		}
		themeIDByName[newTheme.Name] = newTheme.ID
	}

	assigned := 0
	for _, assignment := range req.Assignments {
		themeName := strings.TrimSpace(assignment.ThemeName)
		themeID, ok := themeIDByName[themeName]
		if !ok {
			continue
		}

		taskID, err := uuid.Parse(assignment.TaskID)
		if err != nil {
			continue
		}

		result := tx.Model(&models.Task{}).
			Where("id = ? AND user_id = ?", taskID, userIDUUID).
			Update("theme_id", themeID)
		if result.Error != nil {
			config.Logger.Warnf("Failed to assign task %s to theme %s: %v", taskID, themeName, result.Error)
			continue
		}
		if result.RowsAffected > 0 {
			assigned++
		}
	}

	if err := tx.Commit().Error; err != nil {
		config.Logger.Errorf("Failed to commit theme apply transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save themes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        fmt.Sprintf("Created %d themes and assigned %d tasks", len(themeIDByName), assigned),
		"themes_created": len(themeIDByName),
		"tasks_assigned": assigned,
	})
}
