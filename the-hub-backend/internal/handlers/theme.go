package handlers

import (
	"net/http"
	"strings"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// validThemeDays holds the allowed weekday identifiers.
var validThemeDays = map[string]bool{
	"monday":    true,
	"tuesday":   true,
	"wednesday": true,
	"thursday":  true,
	"friday":    true,
	"saturday":  true,
	"sunday":    true,
}

// normalizeThemeDays lowercases and de-duplicates weekday names.
func normalizeThemeDays(days []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(days))
	for _, day := range days {
		day = strings.ToLower(strings.TrimSpace(day))
		if day == "" || seen[day] {
			continue
		}
		seen[day] = true
		result = append(result, day)
	}
	return result
}

// validateThemeDays returns an error when any day is not a valid weekday name.
func validateThemeDays(days []string) error {
	for _, day := range days {
		if !validThemeDays[day] {
			return &invalidDayError{day: day}
		}
	}
	return nil
}

type invalidDayError struct {
	day string
}

func (e *invalidDayError) Error() string {
	return "invalid weekday: '" + e.day + "'. Expected one of monday..sunday"
}

// themeDayConflict returns the names of existing themes that already use any of the given days.
func themeDayConflict(db *gorm.DB, userID uuid.UUID, days []string, excludeID *uuid.UUID) ([]models.Theme, error) {
	if len(days) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(days))
	args := make([]interface{}, 0, len(days)+1)
	args = append(args, userID)
	for i, day := range days {
		placeholders[i] = "?"
		args = append(args, day)
	}
	query := db.Where("user_id = ? AND days && ARRAY["+strings.Join(placeholders, ",")+"]", args...)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var conflicting []models.Theme
	if err := query.Find(&conflicting).Error; err != nil {
		return nil, err
	}
	return conflicting, nil
}

// themeBelongsToUser reports whether the theme exists and belongs to the user.
func themeBelongsToUser(db *gorm.DB, userID uuid.UUID, themeID uuid.UUID) bool {
	var count int64
	if err := db.Model(&models.Theme{}).Where("id = ? AND user_id = ?", themeID, userID).Count(&count).Error; err != nil {
		config.Logger.Errorf("Failed to validate theme ownership: %v", err)
		return false
	}
	return count > 0
}

// GetThemes godoc
// @Summary      Get all themes
// @Description  Fetch the logged-in user's themes
// @Tags         themes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.Theme
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /themes [get]
func GetThemes(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var themes []models.Theme
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Order("created_at ASC").Find(&themes).Error; err != nil {
		config.Logger.Errorf("Failed to fetch themes for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch themes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"themes": themes})
}

// GetTheme godoc
// @Summary      Get a specific theme
// @Description  Fetch a specific theme by ID for the logged-in user
// @Tags         themes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      string  true  "Theme ID"
// @Success      200  {object}  models.Theme
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /themes/{ID} [get]
func GetTheme(c *gin.Context) {
	themeID, err := uuid.Parse(c.Param("ID"))
	if err != nil {
		config.Logger.Warnf("Invalid theme ID param: %s", c.Param("ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theme ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var theme models.Theme
	if err := config.GetDB().Where("id = ? AND user_id = ?", themeID, userIDUUID).First(&theme).Error; err != nil {
		config.Logger.Warnf("Theme ID %s not found for user %s: %v", themeID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}

	c.JSON(http.StatusOK, theme)
}

// CreateTheme godoc
// @Summary      Create a new theme
// @Description  Create a new theme for the logged-in user with assigned days
// @Tags         themes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        theme  body      models.CreateThemeRequest  true  "Theme creation data"
// @Success      201    {object}  models.Theme
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /themes [post]
func CreateTheme(c *gin.Context) {
	var input models.CreateThemeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid theme input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for theme", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	days := normalizeThemeDays(input.Days)
	if err := validateThemeDays(days); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conflicts, err := themeDayConflict(config.GetDB(), userIDUUID, days, nil)
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

	theme := models.Theme{
		UserID:      userIDUUID,
		Name:        strings.TrimSpace(input.Name),
		Description: input.Description,
		Color:       input.Color,
		Icon:        input.Icon,
		Days:        pq.StringArray(days),
	}

	if err := config.GetDB().Create(&theme).Error; err != nil {
		config.Logger.Errorf("Failed to create theme for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create theme"})
		return
	}

	c.JSON(http.StatusCreated, theme)
}

// UpdateTheme godoc
// @Summary      Update a theme
// @Description  Update a specific theme by ID for the logged-in user
// @Tags         themes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID     path      string                    true  "Theme ID"
// @Param        theme  body      models.UpdateThemeRequest  true  "Theme update data"
// @Success      200   {object}  models.Theme
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /themes/{ID} [patch]
func UpdateTheme(c *gin.Context) {
	themeID, err := uuid.Parse(c.Param("ID"))
	if err != nil {
		config.Logger.Warnf("Invalid theme ID param for update: %s", c.Param("ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theme ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var theme models.Theme
	if err := config.GetDB().Where("id = ? AND user_id = ?", themeID, userIDUUID).First(&theme).Error; err != nil {
		config.Logger.Warnf("Theme not found for update: ID %s, User %s", themeID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}

	var input models.UpdateThemeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for theme ID %s: %v", themeID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var newDays []string
	if input.Days != nil {
		newDays = normalizeThemeDays(input.Days)
		if err := validateThemeDays(newDays); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		conflicts, err := themeDayConflict(config.GetDB(), userIDUUID, newDays, &themeID)
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

	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}
	if input.Icon != nil {
		updates["icon"] = *input.Icon
	}
	if input.Days != nil {
		updates["days"] = pq.StringArray(newDays)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := config.GetDB().Model(&theme).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update theme ID %s: %v", themeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update theme"})
		return
	}

	if err := config.GetDB().First(&theme, themeID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated theme ID %s: %v", themeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated theme"})
		return
	}

	c.JSON(http.StatusOK, theme)
}

// DeleteTheme godoc
// @Summary      Delete a theme
// @Description  Delete a specific theme by ID for the logged-in user. Tasks assigned to the theme are kept but unassigned.
// @Tags         themes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      string  true  "Theme ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /themes/{ID} [delete]
func DeleteTheme(c *gin.Context) {
	themeID, err := uuid.Parse(c.Param("ID"))
	if err != nil {
		config.Logger.Warnf("Invalid theme ID param for delete: %s", c.Param("ID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theme ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var theme models.Theme
	if err := config.GetDB().Where("id = ? AND user_id = ?", themeID, userIDUUID).First(&theme).Error; err != nil {
		config.Logger.Warnf("Theme not found for delete: ID %s, User %s", themeID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}

	if err := config.GetDB().Delete(&theme).Error; err != nil {
		config.Logger.Errorf("Failed to delete theme ID %s: %v", themeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete theme"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theme deleted successfully"})
}
