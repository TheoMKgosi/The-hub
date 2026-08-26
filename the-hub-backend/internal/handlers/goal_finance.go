package handlers

import (
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetFinancialGoals(c *gin.Context) {
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

	var goals []models.FinancialGoal
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Order("created_at desc").Find(&goals).Error; err != nil {
		config.Logger.Errorf("Error fetching financial goals for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch financial goals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goals": goals})
}

type CreateFinancialGoalRequest struct {
	Title         string     `json:"title" binding:"required"`
	Description   string     `json:"description"`
	TargetAmount  *float64   `json:"target_amount"`
	CurrentAmount float64    `json:"current_amount"`
	Type          string     `json:"type" binding:"required"`
	Priority      *int       `json:"priority"`
	TargetDate    *time.Time `json:"target_date"`
	CategoryID    *uuid.UUID `json:"category_id"`
	Color         string     `json:"color"`
}

func CreateFinancialGoal(c *gin.Context) {
	var input CreateFinancialGoalRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid financial goal input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for financial goal", "details": err.Error()})
		return
	}

	if input.Type != "savings" && input.Type != "checklist" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'savings' or 'checklist'"})
		return
	}

	if input.Type == "checklist" {
		input.TargetAmount = nil
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during financial goal creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	color := input.Color
	if color == "" {
		color = "#3B82F6"
	}

	goal := models.FinancialGoal{
		Title:         input.Title,
		Description:   input.Description,
		TargetAmount:  input.TargetAmount,
		CurrentAmount: input.CurrentAmount,
		Type:          input.Type,
		Status:        "active",
		Priority:      input.Priority,
		TargetDate:    input.TargetDate,
		CategoryID:    input.CategoryID,
		Color:         color,
		UserID:        userIDUUID,
	}

	config.Logger.Infof("Creating financial goal for user %s: %s", userIDUUID, input.Title)
	if err := config.GetDB().Create(&goal).Error; err != nil {
		config.Logger.Errorf("Error creating financial goal for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create financial goal"})
		return
	}

	config.Logger.Infof("Successfully created financial goal ID %s for user %s", goal.ID, userIDUUID)
	c.JSON(http.StatusCreated, goal)
}

type UpdateFinancialGoalRequest struct {
	Title         *string    `json:"title"`
	Description   *string    `json:"description"`
	TargetAmount  *float64   `json:"target_amount"`
	CurrentAmount *float64   `json:"current_amount"`
	Type          *string    `json:"type"`
	Status        *string    `json:"status"`
	Priority      *int       `json:"priority"`
	TargetDate    *time.Time `json:"target_date"`
	CategoryID    *uuid.UUID `json:"category_id"`
	Color         *string    `json:"color"`
}

func UpdateFinancialGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param for update: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during financial goal update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var goal models.FinancialGoal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal not found for update: ID %s, User %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Financial goal not found"})
		return
	}

	var input UpdateFinancialGoalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for goal ID %s: %v", goalID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Type != nil {
		if *input.Type != "savings" && *input.Type != "checklist" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'savings' or 'checklist'"})
			return
		}
		updates["type"] = *input.Type
	}
	if input.Status != nil {
		if *input.Status != "active" && *input.Status != "completed" && *input.Status != "cancelled" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'active', 'completed', or 'cancelled'"})
			return
		}
		updates["status"] = *input.Status
	}
	if input.Priority != nil {
		if *input.Priority < 1 || *input.Priority > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Priority must be between 1 and 5"})
			return
		}
		updates["priority"] = *input.Priority
	}
	if input.TargetAmount != nil {
		updates["target_amount"] = *input.TargetAmount
	}
	if input.CurrentAmount != nil {
		updates["current_amount"] = *input.CurrentAmount
	}
	if input.TargetDate != nil {
		updates["target_date"] = *input.TargetDate
	}
	if input.CategoryID != nil {
		updates["category_id"] = *input.CategoryID
	}
	if input.Color != nil {
		updates["color"] = *input.Color
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for goal update: ID %s", goalID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating financial goal ID %s for user %s", goalID, userIDUUID)
	if err := config.GetDB().Model(&goal).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update financial goal ID %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update financial goal"})
		return
	}

	if err := config.GetDB().First(&goal, goal.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated goal ID %s: %v", goal.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated goal"})
		return
	}

	config.Logger.Infof("Successfully updated financial goal ID %s for user %s", goal.ID, userIDUUID)
	c.JSON(http.StatusOK, goal)
}

func DeleteFinancialGoal(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param for delete: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during financial goal deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var goal models.FinancialGoal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal not found for delete: ID %s, User %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Financial goal not found"})
		return
	}

	config.Logger.Infof("Deleting financial goal ID %s for user %s", goalID, userIDUUID)
	if err := config.GetDB().Delete(&goal).Error; err != nil {
		config.Logger.Errorf("Failed to delete financial goal ID %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete financial goal"})
		return
	}

	config.Logger.Infof("Successfully deleted financial goal ID %s for user %s", goalID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Financial goal deleted successfully"})
}

type UpdateGoalProgressRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

func UpdateGoalProgress(c *gin.Context) {
	goalIDStr := c.Param("ID")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid goal ID param for progress update: %s", goalIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during progress update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var input UpdateGoalProgressRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid progress input for goal ID %s: %v", goalID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	var goal models.FinancialGoal
	if err := config.GetDB().Where("id = ? AND user_id = ?", goalID, userIDUUID).First(&goal).Error; err != nil {
		config.Logger.Warnf("Goal not found for progress update: ID %s, User %s", goalID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Financial goal not found"})
		return
	}

	if goal.Type == "checklist" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot track progress on checklist items. Mark as completed instead."})
		return
	}

	newAmount := goal.CurrentAmount + input.Amount
	updates := map[string]interface{}{
		"current_amount": newAmount,
	}

	if goal.TargetAmount != nil && newAmount >= *goal.TargetAmount {
		updates["status"] = "completed"
	}

	if err := config.GetDB().Model(&goal).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update goal progress ID %s: %v", goalID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goal progress"})
		return
	}

	if err := config.GetDB().First(&goal, goal.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated goal ID %s: %v", goal.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated goal"})
		return
	}

	config.Logger.Infof("Updated progress for goal %s: added %.2f, new total: %.2f", goalID, input.Amount, goal.CurrentAmount)
	c.JSON(http.StatusOK, goal)
}
