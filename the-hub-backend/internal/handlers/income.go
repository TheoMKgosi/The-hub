package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get all incomes for user
func GetIncomes(c *gin.Context) {
	var incomes []models.Income
	userID := c.MustGet("userID").(uint)

	if err := config.GetDB().Preload("Budgets.Category", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Where("user_id = ?", userID).Order("created_at desc").Find(&incomes).Error; err != nil {
		log.Println("Error fetching incomes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch incomes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"incomes": incomes})
}

// Create new income
func CreateIncome(c *gin.Context) {
	var input struct {
		Source     string  `json:"source" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		ReceivedAt string  `json:"received_at" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	layout := "2006-01-02"
	receivedAt, err := time.Parse(layout, input.ReceivedAt)
	if err != nil {
		log.Println("Error timeParse: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create income"})
		return
	}

	income := models.Income{
		Source:     input.Source,
		Amount:     input.Amount,
		ReceivedAt: receivedAt,
		UserID:     c.MustGet("userID").(uint),
	}

	if err := config.GetDB().Create(&income).Error; err != nil {
		log.Println("Error create: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create income"})
		return
	}

	c.JSON(http.StatusCreated, income)
}

// Update an income
func UpdateIncome(c *gin.Context) {
	id := c.Param("id")
	var income models.Income

	if err := config.GetDB().First(&income, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Income not found"})
		return
	}

	var input struct {
		Source     *string    `json:"source"`
		Amount     *float64   `json:"amount"`
		ReceivedAt *time.Time `json:"received_at"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Source != nil {
		updates["source"] = *input.Source
	}
	if input.Amount != nil {
		updates["amount"] = *input.Amount
	}
	if input.ReceivedAt != nil {
		updates["received_at"] = *input.ReceivedAt
	}

	if err := config.GetDB().Model(&income).Updates(updates).Error; err != nil {
		log.Println("Error create: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update income"})
		return
	}

	c.JSON(http.StatusOK, income)
}

// Delete income
func DeleteIncome(c *gin.Context) {
	id := c.Param("id")
	var income models.Income

	if err := config.GetDB().First(&income, id).Error; err != nil {
		log.Println("Error fetch: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Income not found"})
		return
	}

	if err := config.GetDB().Delete(&income).Error; err != nil {
		log.Println("Error delete: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete income"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income deleted"})
}
