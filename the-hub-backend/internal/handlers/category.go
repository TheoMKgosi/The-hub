package handlers

import (
	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Get all budget categories for user
func GetBudgetCategories(c *gin.Context) {
	var categories []models.BudgetCategory
	userID := c.MustGet("userID").(uint)

	if err := config.GetDB().Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		log.Println("Error fetching categories:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// Create new budget category
func CreateBudgetCategory(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	category := models.BudgetCategory{
		Name:   input.Name,
		UserID: c.MustGet("userID").(uint),
	}

	if err := config.GetDB().Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// Update a category
func UpdateBudgetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.BudgetCategory

	if err := config.GetDB().First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input struct {
		Name *string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}

	if err := config.GetDB().Model(&category).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete category
func DeleteBudgetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.BudgetCategory

	if err := config.GetDB().First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := config.GetDB().Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
