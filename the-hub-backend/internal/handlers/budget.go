package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// Get all budgets for a user
func GetBudgets(c *gin.Context) {
	var budgets []models.Budget
	userID := c.MustGet("userID").(uint)

	if err := config.GetDB().Where("user_id = ?", userID).Preload("Category").Find(&budgets).Error; err != nil {
		log.Println("GetBudgets error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch budgets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budgets": budgets})
}

// Get a single budget
func GetBudget(c *gin.Context) {
	id := c.Param("id")
	var budget models.Budget

	if err := config.GetDB().Preload("Category").First(&budget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// Create a budget
func CreateBudget(c *gin.Context) {
	var input struct {
		CategoryID uint    `json:"category_id" binding:"required"`
		IncomeID   *uint   `json:"income_id"`
		Amount     float64 `json:"amount" binding:"required"`
		StartDate  string  `json:"start_date" binding:"required"`
		EndDate    string  `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, input.StartDate)
	if err != nil {
		log.Println("Error startDate: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
	}
	endDate, err := time.Parse(layout, input.EndDate)
	if err != nil {
		log.Println("Error endDate: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
	}

	budget := models.Budget{
		CategoryID: input.CategoryID,
		Amount:     input.Amount,
		IncomeID:   input.IncomeID,
		StartDate:  startDate,
		EndDate:    endDate,
		UserID:     c.MustGet("userID").(uint),
	}

	if err := config.GetDB().Create(&budget).Error; err != nil {
		log.Println("Error create: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create budget"})
		return
	}

	c.JSON(http.StatusCreated, budget)
}

// Update budget
func UpdateBudget(c *gin.Context) {
	id := c.Param("id")
	var budget models.Budget

	if err := config.GetDB().First(&budget, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	var input struct {
		Amount    *float64   `json:"amount"`
		StartDate *time.Time `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Amount != nil {
		updates["amount"] = *input.Amount
	}
	if input.StartDate != nil {
		updates["start_date"] = *input.StartDate
	}
	if input.EndDate != nil {
		updates["end_date"] = *input.EndDate
	}

	if err := config.GetDB().Model(&budget).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

// Delete budget
func DeleteBudget(c *gin.Context) {
	budgetID := c.Param("ID")
	incomeID := c.Param("incomeID")

	var budget models.Budget
	db := config.GetDB()

	// Use both income_id and id to find the budget
	if err := db.Where("id = ? AND income_id = ?", budgetID, incomeID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found for this income"})
		return
	}

	if err := db.Delete(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted"})
}
