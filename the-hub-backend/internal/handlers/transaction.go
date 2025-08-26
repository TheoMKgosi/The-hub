package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Get all transactions for user
func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction
	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := config.GetDB().Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Where("user_id = ?", userIDUUID).Order("date desc").Find(&transactions).Error; err != nil {
		log.Println("Error fetching transactions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// Create new transaction
func CreateTransaction(c *gin.Context) {
	var input struct {
		Description string  `json:"description" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=income expense"`
		Date        string  `json:"date" binding:"required"`
		CategoryID  *string `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, input.Date)
	if err != nil {
		log.Println("Error timeParse: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context during transaction creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	transaction := models.Transaction{
		Description: input.Description,
		Amount:      input.Amount,
		Type:        input.Type,
		Date:        date,
		UserID:      userIDUUID,
	}

	// Handle optional category
	if input.CategoryID != nil && *input.CategoryID != "" {
		categoryID, err := uuid.Parse(*input.CategoryID)
		if err != nil {
			log.Println("Error parsing category ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		transaction.CategoryID = &categoryID
	}

	if err := config.GetDB().Create(&transaction).Error; err != nil {
		log.Println("Error create: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Load the category if it was set
	if transaction.CategoryID != nil {
		if err := config.GetDB().Preload("Category").First(&transaction, transaction.ID).Error; err != nil {
			log.Println("Error loading category:", err)
		}
	}

	c.JSON(http.StatusCreated, transaction)
}

// Update a transaction
func UpdateTransaction(c *gin.Context) {
	id := c.Param("ID")
	var transaction models.Transaction

	if err := config.GetDB().First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	var input struct {
		Description *string  `json:"description"`
		Amount      *float64 `json:"amount"`
		Type        *string  `json:"type"`
		Date        *string  `json:"date"`
		CategoryID  *string  `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Amount != nil {
		updates["amount"] = *input.Amount
	}
	if input.Type != nil {
		if *input.Type != "income" && *input.Type != "expense" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'income' or 'expense'"})
			return
		}
		updates["type"] = *input.Type
	}
	if input.Date != nil {
		layout := "2006-01-02"
		date, err := time.Parse(layout, *input.Date)
		if err != nil {
			log.Println("Error timeParse: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		updates["date"] = date
	}
	if input.CategoryID != nil {
		if *input.CategoryID == "" {
			updates["category_id"] = nil
		} else {
			categoryID, err := uuid.Parse(*input.CategoryID)
			if err != nil {
				log.Println("Error parsing category ID:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
				return
			}
			updates["category_id"] = categoryID
		}
	}

	if err := config.GetDB().Model(&transaction).Updates(updates).Error; err != nil {
		log.Println("Error update: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	// Load the updated transaction with category
	if err := config.GetDB().Preload("Category").First(&transaction, id).Error; err != nil {
		log.Println("Error loading updated transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Delete transaction
func DeleteTransaction(c *gin.Context) {
	id := c.Param("ID")
	var transaction models.Transaction

	if err := config.GetDB().First(&transaction, id).Error; err != nil {
		log.Println("Error fetch: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := config.GetDB().Delete(&transaction).Error; err != nil {
		log.Println("Error delete: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}
