package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetBudgets godoc
// @Summary      Get all budgets
// @Description  Fetch budgets for the logged-in user with optional ordering
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_by  query     string  false  "Order by field (amount, start_date, end_date, created_at)"  default(created_at)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.Budget
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budgets [get]
func GetBudgets(c *gin.Context) {
	var budgets []models.Budget
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

	// Get query parameters for ordering
	orderBy := c.DefaultQuery("order_by", "created_at")
	sortDir := c.DefaultQuery("sort", "asc")

	// Validate order_by parameter
	validOrderFields := map[string]bool{
		"amount":     true,
		"start_date": true,
		"end_date":   true,
		"created_at": true,
	}

	if !validOrderFields[orderBy] {
		config.Logger.Warnf("Invalid order_by parameter: %s", orderBy)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order_by parameter"})
		return
	}

	// Validate sort direction
	if sortDir != "asc" && sortDir != "desc" {
		config.Logger.Warnf("Invalid sort direction: %s", sortDir)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort direction. Use 'asc' or 'desc'"})
		return
	}

	orderClause := orderBy + " " + sortDir

	config.Logger.Infof("Fetching budgets for user ID: %s with order: %s", userIDUUID, orderClause)
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Preload("Category").Order(orderClause).Find(&budgets).Error; err != nil {
		config.Logger.Errorf("Error fetching budgets for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch budgets"})
		return
	}

	config.Logger.Infof("Found %d budgets for user ID %s", len(budgets), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"budgets": budgets})
}

// GetBudget godoc
// @Summary      Get a specific budget
// @Description  Fetch a specific budget by ID for the logged-in user
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Budget ID"
// @Success      200  {object}  map[string]models.Budget
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budgets/{ID} [get]
func GetBudget(c *gin.Context) {
	budgetIDStr := c.Param("ID")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid budget ID param: %s", budgetIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
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
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	config.Logger.Infof("Fetching budget ID: %s for user ID: %s", budgetID, userIDUUID)
	var budget models.Budget
	// Ensure user can only access their own budgets
	if err := config.GetDB().Where("id = ? AND user_id = ?", budgetID, userIDUUID).Preload("Category").First(&budget).Error; err != nil {
		config.Logger.Errorf("Budget ID %s not found for user %s: %v", budgetID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved budget ID %s for user %s", budgetID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"budget": budget})
}

// CreateBudgetRequest represents the request body for creating a budget
type CreateBudgetRequest struct {
	CategoryID uuid.UUID  `json:"category_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	IncomeID   *uuid.UUID `json:"income_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount     float64    `json:"amount" binding:"required" example:"1500.00"`
	StartDate  string     `json:"start_date" binding:"required" example:"2024-01-01"`
	EndDate    string     `json:"end_date" binding:"required" example:"2024-12-31"`
}

// CreateBudget godoc
// @Summary      Create a new budget
// @Description  Create a new budget for the logged-in user
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        budget  body      CreateBudgetRequest  true  "Budget creation data"
// @Success      201     {object}  models.Budget
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /budgets [post]
func CreateBudget(c *gin.Context) {
	var input CreateBudgetRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid budget input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for budget", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during budget creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Parse dates
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		config.Logger.Warnf("Invalid start date format: %s", input.StartDate)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		config.Logger.Warnf("Invalid end date format: %s", input.EndDate)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Validate that end date is after start date
	if endDate.Before(startDate) {
		config.Logger.Warnf("End date %s is before start date %s", input.EndDate, input.StartDate)
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	budget := models.Budget{
		CategoryID: input.CategoryID,
		IncomeID:   input.IncomeID,
		Amount:     input.Amount,
		StartDate:  startDate,
		EndDate:    endDate,
		UserID:     userIDUUID,
	}

	config.Logger.Infof("Creating budget for user %s: category %s, amount %.2f", userIDUUID, input.CategoryID, input.Amount)
	if err := config.GetDB().Create(&budget).Error; err != nil {
		config.Logger.Errorf("Error creating budget for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create budget"})
		return
	}

	config.Logger.Infof("Successfully created budget ID %s for user %s", budget.ID, userIDUUID)
	c.JSON(http.StatusCreated, budget)
}

// UpdateBudgetRequest represents the request body for updating a budget
type UpdateBudgetRequest struct {
	Amount    *float64 `json:"amount" example:"2000.00"`
	StartDate *string  `json:"start_date" example:"2024-01-01"`
	EndDate   *string  `json:"end_date" example:"2024-12-31"`
}

// UpdateBudget godoc
// @Summary      Update a budget
// @Description  Update a specific budget by ID for the logged-in user
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID      path      int                  true  "Budget ID"
// @Param        budget  body      UpdateBudgetRequest  true  "Budget update data"
// @Success      200     {object}  models.Budget
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /budgets/{ID} [put]
func UpdateBudget(c *gin.Context) {
	budgetIDStr := c.Param("ID")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid budget ID param for update: %s", budgetIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during budget update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var budget models.Budget
	// Ensure user can only update their own budgets
	if err := config.GetDB().Where("id = ? AND user_id = ?", budgetID, userIDUUID).First(&budget).Error; err != nil {
		config.Logger.Warnf("Budget not found for update: ID %s, User %s", budgetID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	var input UpdateBudgetRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for budget ID %d: %v", budgetID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	layout := "2006-01-02"

	if input.Amount != nil {
		updates["amount"] = *input.Amount
	}

	if input.StartDate != nil {
		startDate, err := time.Parse(layout, *input.StartDate)
		if err != nil {
			config.Logger.Warnf("Invalid start date format in update: %s", *input.StartDate)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}
		updates["start_date"] = startDate
	}

	if input.EndDate != nil {
		endDate, err := time.Parse(layout, *input.EndDate)
		if err != nil {
			config.Logger.Warnf("Invalid end date format in update: %s", *input.EndDate)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}
		updates["end_date"] = endDate
	}

	// Validate date range if both dates are being updated
	if input.StartDate != nil && input.EndDate != nil {
		startDate, _ := time.Parse(layout, *input.StartDate)
		endDate, _ := time.Parse(layout, *input.EndDate)
		if endDate.Before(startDate) {
			config.Logger.Warnf("End date %s is before start date %s in update", *input.EndDate, *input.StartDate)
			c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
			return
		}
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for budget update: ID %d", budgetID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating budget ID %d for user %v with data: %+v", budgetID, userID, updates)
	if err := config.GetDB().Model(&budget).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update budget ID %d: %v", budgetID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	// Reload the updated budget
	if err := config.GetDB().First(&budget, budget.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated budget ID %d: %v", budget.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated budget"})
		return
	}

	config.Logger.Infof("Successfully updated budget ID %d for user %v", budget.ID, userID)
	c.JSON(http.StatusOK, budget)
}

// DeleteBudget godoc
// @Summary      Delete a budget
// @Description  Delete a specific budget by ID for the logged-in user
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Budget ID"
// @Success      200  {object}  models.Budget
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budgets/{ID} [delete]
func DeleteBudget(c *gin.Context) {
	budgetIDStr := c.Param("ID")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid budget ID param for delete: %s", budgetIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during budget deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var budget models.Budget
	// Ensure user can only delete their own budgets
	if err := config.GetDB().Where("id = ? AND user_id = ?", budgetID, userIDUUID).First(&budget).Error; err != nil {
		config.Logger.Warnf("Budget not found for delete: ID %s, User %s", budgetID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	config.Logger.Infof("Deleting budget ID %s for user %s", budgetID, userIDUUID)
	if err := config.GetDB().Delete(&budget).Error; err != nil {
		config.Logger.Errorf("Failed to delete budget ID %s: %v", budgetID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget"})
		return
	}

	config.Logger.Infof("Successfully deleted budget ID %s for user %s", budgetID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully", "budget": budget})
}

// BudgetAnalytics represents budget performance analytics
type BudgetAnalytics struct {
	BudgetID          uuid.UUID `json:"budget_id"`
	CategoryName      string    `json:"category_name"`
	BudgetAmount      float64   `json:"budget_amount"`
	SpentAmount       float64   `json:"spent_amount"`
	RemainingAmount   float64   `json:"remaining_amount"`
	UtilizationRate   float64   `json:"utilization_rate"`
	DaysRemaining     int       `json:"days_remaining"`
	DailySpendingRate float64   `json:"daily_spending_rate"`
	Status            string    `json:"status"` // "on_track", "warning", "over_budget"
}

// GetBudgetAnalytics godoc
// @Summary      Get budget analytics and spending vs budget tracking
// @Description  Fetch budget performance analytics for the logged-in user
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period  query     string  false  "Analysis period (current, last_month, last_3_months)"  default(current)
// @Success      200  {object}  map[string][]BudgetAnalytics
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budgets/analytics [get]
func GetBudgetAnalytics(c *gin.Context) {
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

	period := c.DefaultQuery("period", "current")

	// Get all budgets for the user
	var budgets []models.Budget
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Preload("Category").Find(&budgets).Error; err != nil {
		config.Logger.Errorf("Error fetching budgets for analytics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch budgets"})
		return
	}

	var analytics []BudgetAnalytics

	for _, budget := range budgets {
		// Calculate spending for this budget's category and period
		spentAmount, err := calculateBudgetSpending(budget, period)
		if err != nil {
			config.Logger.Errorf("Error calculating spending for budget %s: %v", budget.ID, err)
			continue
		}

		// Calculate analytics
		remainingAmount := budget.Amount - spentAmount
		utilizationRate := (spentAmount / budget.Amount) * 100

		// Calculate days remaining in budget period
		now := time.Now()
		daysRemaining := int(budget.EndDate.Sub(now).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Calculate daily spending rate
		daysElapsed := int(now.Sub(budget.StartDate).Hours() / 24)
		if daysElapsed < 1 {
			daysElapsed = 1
		}
		dailySpendingRate := spentAmount / float64(daysElapsed)

		// Determine status with sophisticated logic
		status := CalculateBudgetStatus(budget.Amount, spentAmount, daysRemaining, budget.Category.Name)

		analytics = append(analytics, BudgetAnalytics{
			BudgetID:          budget.ID,
			CategoryName:      budget.Category.Name,
			BudgetAmount:      budget.Amount,
			SpentAmount:       spentAmount,
			RemainingAmount:   remainingAmount,
			UtilizationRate:   utilizationRate,
			DaysRemaining:     daysRemaining,
			DailySpendingRate: dailySpendingRate,
			Status:            status,
		})
	}

	config.Logger.Infof("Generated analytics for %d budgets for user %s", len(analytics), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

// calculateBudgetSpending calculates total spending for a budget category within the budget period
func calculateBudgetSpending(budget models.Budget, period string) (float64, error) {
	var startDate, endDate time.Time

	switch period {
	case "current":
		startDate = budget.StartDate
		endDate = budget.EndDate
	case "last_month":
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = endDate.AddDate(0, -1, 0)
	case "last_3_months":
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startDate = endDate.AddDate(0, -3, 0)
	default:
		startDate = budget.StartDate
		endDate = budget.EndDate
	}

	var totalSpent float64
	if err := config.GetDB().Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ? AND type = ? AND date >= ? AND date <= ?",
			budget.UserID, budget.CategoryID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpent).Error; err != nil {
		return 0, err
	}

	return totalSpent, nil
}

// GetBudgetSuggestions godoc
// @Summary      Get AI-powered budget suggestions
// @Description  Generate budget amount suggestions based on spending patterns
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budgets/suggestions [get]
func GetBudgetSuggestions(c *gin.Context) {
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

	suggestions, err := generateBudgetSuggestions(userIDUUID)
	if err != nil {
		config.Logger.Errorf("Error generating budget suggestions for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate suggestions"})
		return
	}

	config.Logger.Infof("Generated budget suggestions for user %s", userIDUUID)
	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}

// generateBudgetSuggestions creates AI-powered budget suggestions
func generateBudgetSuggestions(userID uuid.UUID) (map[string]interface{}, error) {
	// Get user's transaction history for the last 3 months
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	var transactions []models.Transaction
	if err := config.GetDB().Where("user_id = ? AND date >= ?", userID, threeMonthsAgo).
		Preload("Category").Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Analyze spending patterns by category
	categorySpending := make(map[string][]float64)
	categoryNames := make(map[string]string)

	for _, transaction := range transactions {
		if transaction.Type == "expense" && transaction.CategoryID != nil {
			categoryID := transaction.CategoryID.String()
			categorySpending[categoryID] = append(categorySpending[categoryID], transaction.Amount)
			if transaction.Category.Name != "" {
				categoryNames[categoryID] = transaction.Category.Name
			}
		}
	}

	// Generate suggestions
	suggestions := make(map[string]interface{})
	var categorySuggestions []map[string]interface{}

	for categoryID, amounts := range categorySpending {
		if len(amounts) == 0 {
			continue
		}

		// Calculate average spending
		total := 0.0
		for _, amount := range amounts {
			total += amount
		}
		average := total / float64(len(amounts))

		// Calculate standard deviation for variability
		variance := 0.0
		for _, amount := range amounts {
			variance += (amount - average) * (amount - average)
		}
		stdDev := 0.0
		if len(amounts) > 1 {
			stdDev = variance / float64(len(amounts)-1)
		}

		// Suggest budget with buffer for variability
		buffer := stdDev * 0.5 // 50% of standard deviation as buffer
		suggestedAmount := average + buffer

		// Round to nearest 10
		suggestedAmount = float64(int(suggestedAmount/10+0.5)) * 10

		categorySuggestions = append(categorySuggestions, map[string]interface{}{
			"category_id":       categoryID,
			"category_name":     categoryNames[categoryID],
			"current_average":   average,
			"suggested_budget":  suggestedAmount,
			"transaction_count": len(amounts),
			"variability":       stdDev,
			"confidence":        calculateConfidence(len(amounts)),
		})
	}

	suggestions["categories"] = categorySuggestions
	suggestions["analysis_period"] = "3 months"
	suggestions["total_suggestions"] = len(categorySuggestions)

	return suggestions, nil
}

// calculateConfidence calculates confidence level based on sample size
func calculateConfidence(sampleSize int) string {
	switch {
	case sampleSize >= 30:
		return "high"
	case sampleSize >= 10:
		return "medium"
	default:
		return "low"
	}
}

// CalculateBudgetStatus determines budget status with sophisticated logic
func CalculateBudgetStatus(budgetAmount, spentAmount float64, daysRemaining int, categoryName string) string {
	utilizationRate := (spentAmount / budgetAmount) * 100
	remainingAmount := budgetAmount - spentAmount

	// Size-based threshold
	sizeThreshold := getWarningThreshold(budgetAmount)

	// Category multiplier
	categoryMultiplier := getCategoryMultiplier(categoryName)
	adjustedThreshold := sizeThreshold * categoryMultiplier

	// Check remaining amount logic
	remainingWarning := shouldWarnByRemaining(remainingAmount, budgetAmount)

	// Time-based adjustment
	timeAdjustment := getTimeBasedWarning(utilizationRate, daysRemaining, budgetAmount)

	// Final status determination
	if utilizationRate >= 100 {
		return "over_budget"
	}

	if utilizationRate >= adjustedThreshold || remainingWarning || timeAdjustment == "warning" {
		return "warning"
	}

	// Additional "caution" state for moderate utilization
	if utilizationRate >= 60 && utilizationRate < adjustedThreshold {
		return "caution"
	}

	return "on_track"
}

// getWarningThreshold returns warning threshold based on budget size
func getWarningThreshold(budgetAmount float64) float64 {
	switch {
	case budgetAmount < 100:
		return 70.0 // Small budgets: warn earlier
	case budgetAmount < 500:
		return 75.0 // Medium budgets
	case budgetAmount < 2000:
		return 80.0 // Large budgets
	default:
		return 85.0 // Very large budgets: more lenient
	}
}

// shouldWarnByRemaining checks if remaining amount warrants a warning
func shouldWarnByRemaining(remainingAmount, budgetAmount float64) bool {
	// Always warn if remaining amount is very small
	if remainingAmount < 10 {
		return true
	}

	// Warn if remaining percentage is too low relative to budget size
	remainingPercent := (remainingAmount / budgetAmount) * 100
	return remainingPercent < 5.0
}

// getTimeBasedWarning adjusts warning based on time remaining
func getTimeBasedWarning(utilizationRate float64, daysRemaining int, budgetAmount float64) string {
	// If budget period is ending soon and utilization is reasonable, be more lenient
	if daysRemaining <= 3 && utilizationRate < 95 {
		return "lenient"
	}

	// If lots of time left and spending is moderate, be more strict
	if daysRemaining > 14 && utilizationRate >= 70 {
		return "warning"
	}

	return "normal"
}

// getCategoryMultiplier returns warning threshold multiplier based on category
func getCategoryMultiplier(categoryName string) float64 {
	categoryMultipliers := map[string]float64{
		"essentials":    0.9,  // Housing, utilities: more lenient
		"savings":       0.95, // Emergency fund: very strict
		"emergency":     0.95, // Emergency fund: very strict
		"discretionary": 0.7,  // Entertainment: stricter
		"entertainment": 0.7,  // Entertainment: stricter
		"one-time":      0.85, // Planned purchases: more lenient
		"planned":       0.85, // Planned purchases: more lenient
		"vacation":      0.85, // Vacation: more lenient
		"travel":        0.85, // Travel: more lenient
	}

	// Case-insensitive matching
	categoryLower := strings.ToLower(categoryName)
	for category, multiplier := range categoryMultipliers {
		if strings.Contains(categoryLower, category) {
			return multiplier
		}
	}

	return 1.0 // Default multiplier
}

// BudgetAlert represents a budget alert notification
type BudgetAlert struct {
	BudgetID      uuid.UUID `json:"budget_id"`
	CategoryName  string    `json:"category_name"`
	AlertType     string    `json:"alert_type"` // "warning", "danger", "over_budget"
	Message       string    `json:"message"`
	CurrentUsage  float64   `json:"current_usage"`
	BudgetAmount  float64   `json:"budget_amount"`
	DaysRemaining int       `json:"days_remaining"`
}

// GetBudgetAlerts godoc
// @Summary      Get budget alerts and notifications
// @Description  Fetch budget alerts for budgets approaching or exceeding thresholds
// @Tags         budgets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]BudgetAlert
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budgets/alerts [get]
func GetBudgetAlerts(c *gin.Context) {
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

	alerts, err := generateBudgetAlerts(userIDUUID)
	if err != nil {
		config.Logger.Errorf("Error generating budget alerts for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate alerts"})
		return
	}

	config.Logger.Infof("Generated %d budget alerts for user %s", len(alerts), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// generateBudgetAlerts checks all budgets and generates alerts for thresholds
func generateBudgetAlerts(userID uuid.UUID) ([]BudgetAlert, error) {
	var budgets []models.Budget
	if err := config.GetDB().Where("user_id = ?", userID).Preload("Category").Find(&budgets).Error; err != nil {
		return nil, err
	}

	var alerts []BudgetAlert
	now := time.Now()

	for _, budget := range budgets {
		// Calculate current spending
		spentAmount, err := calculateBudgetSpending(budget, "current")
		if err != nil {
			config.Logger.Errorf("Error calculating spending for budget %s: %v", budget.ID, err)
			continue
		}

		utilizationRate := (spentAmount / budget.Amount) * 100
		daysRemaining := int(budget.EndDate.Sub(now).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Check for alerts
		if utilizationRate >= 100 {
			alerts = append(alerts, BudgetAlert{
				BudgetID:      budget.ID,
				CategoryName:  budget.Category.Name,
				AlertType:     "over_budget",
				Message:       fmt.Sprintf("You've exceeded your %s budget by $%.2f", budget.Category.Name, spentAmount-budget.Amount),
				CurrentUsage:  utilizationRate,
				BudgetAmount:  budget.Amount,
				DaysRemaining: daysRemaining,
			})
		} else if utilizationRate >= 90 {
			alerts = append(alerts, BudgetAlert{
				BudgetID:      budget.ID,
				CategoryName:  budget.Category.Name,
				AlertType:     "danger",
				Message:       fmt.Sprintf("You're at %.1f%% of your %s budget with %d days remaining", utilizationRate, budget.Category.Name, daysRemaining),
				CurrentUsage:  utilizationRate,
				BudgetAmount:  budget.Amount,
				DaysRemaining: daysRemaining,
			})
		} else if utilizationRate >= 75 {
			alerts = append(alerts, BudgetAlert{
				BudgetID:      budget.ID,
				CategoryName:  budget.Category.Name,
				AlertType:     "warning",
				Message:       fmt.Sprintf("You've used %.1f%% of your %s budget", utilizationRate, budget.Category.Name),
				CurrentUsage:  utilizationRate,
				BudgetAmount:  budget.Amount,
				DaysRemaining: daysRemaining,
			})
		}

		// Check for end-of-period alerts
		if daysRemaining <= 3 && daysRemaining > 0 && utilizationRate < 100 {
			alerts = append(alerts, BudgetAlert{
				BudgetID:      budget.ID,
				CategoryName:  budget.Category.Name,
				AlertType:     "warning",
				Message:       fmt.Sprintf("Your %s budget period ends in %d days. You've used %.1f%% so far", budget.Category.Name, daysRemaining, utilizationRate),
				CurrentUsage:  utilizationRate,
				BudgetAmount:  budget.Amount,
				DaysRemaining: daysRemaining,
			})
		}
	}

	return alerts, nil
}
