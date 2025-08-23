package handlers

import (
	"net/http"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetBudgetCategories godoc
// @Summary      Get all budget categories
// @Description  Fetch budget categories for the logged-in user with optional ordering
// @Tags         budget-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_by  query     string  false  "Order by field (name, created_at)"  default(created_at)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.BudgetCategory
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budget-categories [get]
func GetBudgetCategories(c *gin.Context) {
	var categories []models.BudgetCategory
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
		"name":       true,
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

	config.Logger.Infof("Fetching budget categories for user ID: %s with order: %s", userIDUUID, orderClause)
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Order(orderClause).Find(&categories).Error; err != nil {
		config.Logger.Errorf("Error fetching budget categories for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch categories"})
		return
	}

	config.Logger.Infof("Found %d budget categories for user ID %s", len(categories), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetBudgetCategory godoc
// @Summary      Get a specific budget category
// @Description  Fetch a specific budget category by ID for the logged-in user
// @Tags         budget-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Budget Category ID"
// @Success      200  {object}  map[string]models.BudgetCategory
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budget-categories/{ID} [get]
func GetBudgetCategory(c *gin.Context) {
	categoryIDStr := c.Param("ID")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid budget category ID param: %s", categoryIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget category ID"})
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

	config.Logger.Infof("Fetching budget category ID: %s for user ID: %s", categoryID, userIDUUID)
	var category models.BudgetCategory
	// Ensure user can only access their own categories
	if err := config.GetDB().Where("id = ? AND user_id = ?", categoryID, userIDUUID).First(&category).Error; err != nil {
		config.Logger.Errorf("Budget category ID %s not found for user %s: %v", categoryID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget category not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved budget category ID %s for user %s", categoryID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"category": category})
}

// CreateBudgetCategoryRequest represents the request body for creating a budget category
type CreateBudgetCategoryRequest struct {
	Name string `json:"name" binding:"required" example:"Groceries"`
}

// CreateBudgetCategory godoc
// @Summary      Create a new budget category
// @Description  Create a new budget category for the logged-in user
// @Tags         budget-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        category  body      CreateBudgetCategoryRequest  true  "Budget category creation data"
// @Success      201       {object}  models.BudgetCategory
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /budget-categories [post]
func CreateBudgetCategory(c *gin.Context) {
	var input CreateBudgetCategoryRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid budget category input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for budget category", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during budget category creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Check for duplicate category name for this user
	var existingCategory models.BudgetCategory
	if err := config.GetDB().Where("name = ? AND user_id = ?", input.Name, userIDUUID).First(&existingCategory).Error; err == nil {
		config.Logger.Warnf("Duplicate budget category name '%s' for user %s", input.Name, userIDUUID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name already exists"})
		return
	}

	category := models.BudgetCategory{
		Name:   input.Name,
		UserID: userIDUUID,
	}

	config.Logger.Infof("Creating budget category for user %s: %s", userIDUUID, input.Name)
	if err := config.GetDB().Create(&category).Error; err != nil {
		config.Logger.Errorf("Error creating budget category for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create budget category"})
		return
	}

	config.Logger.Infof("Successfully created budget category ID %s for user %s", category.ID, userIDUUID)
	c.JSON(http.StatusCreated, category)
}

// UpdateBudgetCategoryRequest represents the request body for updating a budget category
type UpdateBudgetCategoryRequest struct {
	Name *string `json:"name" example:"Updated category name"`
}

// UpdateBudgetCategory godoc
// @Summary      Update a budget category
// @Description  Update a specific budget category by ID for the logged-in user
// @Tags         budget-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID        path      int                          true  "Budget Category ID"
// @Param        category  body      UpdateBudgetCategoryRequest  true  "Budget category update data"
// @Success      200       {object}  models.BudgetCategory
// @Failure      400       {object}  map[string]string
// @Failure      401       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /budget-categories/{ID} [put]
func UpdateBudgetCategory(c *gin.Context) {
	categoryIDStr := c.Param("ID")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid budget category ID param for update: %s", categoryIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget category ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during budget category update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var category models.BudgetCategory
	// Ensure user can only update their own categories
	if err := config.GetDB().Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		config.Logger.Warnf("Budget category not found for update: ID %d, User %v", categoryID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget category not found"})
		return
	}

	var input UpdateBudgetCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for budget category ID %d: %v", categoryID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		// Check for duplicate category name for this user (excluding current category)
		var existingCategory models.BudgetCategory
		if err := config.GetDB().Where("name = ? AND user_id = ? AND id != ?", *input.Name, userID, categoryID).First(&existingCategory).Error; err == nil {
			config.Logger.Warnf("Duplicate budget category name '%s' for user %v during update", *input.Name, userID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category name already exists"})
			return
		}
		updates["name"] = *input.Name
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for budget category update: ID %d", categoryID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating budget category ID %d for user %v with data: %+v", categoryID, userID, updates)
	if err := config.GetDB().Model(&category).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update budget category ID %d: %v", categoryID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget category"})
		return
	}

	// Reload the updated category
	if err := config.GetDB().First(&category, category.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated budget category ID %d: %v", category.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated budget category"})
		return
	}

	config.Logger.Infof("Successfully updated budget category ID %d for user %v", category.ID, userID)
	c.JSON(http.StatusOK, category)
}

// DeleteBudgetCategory godoc
// @Summary      Delete a budget category
// @Description  Delete a specific budget category by ID for the logged-in user
// @Tags         budget-categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Budget Category ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /budget-categories/{ID} [delete]
func DeleteBudgetCategory(c *gin.Context) {
	categoryIDStr := c.Param("ID")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid budget category ID param for delete: %s", categoryIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget category ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during budget category deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var category models.BudgetCategory
	// Ensure user can only delete their own categories
	if err := config.GetDB().Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		config.Logger.Warnf("Budget category not found for delete: ID %d, User %v", categoryID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget category not found"})
		return
	}

	config.Logger.Infof("Deleting budget category ID %d for user %v", categoryID, userID)
	if err := config.GetDB().Delete(&category).Error; err != nil {
		config.Logger.Errorf("Failed to delete budget category ID %d: %v", categoryID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget category"})
		return
	}

	config.Logger.Infof("Successfully deleted budget category ID %d for user %v", categoryID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Budget category deleted successfully", "category": category})
}
