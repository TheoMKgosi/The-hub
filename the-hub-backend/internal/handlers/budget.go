package handlers

import (
	"net/http"
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
