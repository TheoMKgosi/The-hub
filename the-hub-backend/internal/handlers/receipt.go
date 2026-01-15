package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Get all receipts for user
func GetReceipts(c *gin.Context) {
	var receipts []models.Receipt
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
	}).Where("user_id = ?", userIDUUID).Order("created_at desc").Find(&receipts).Error; err != nil {
		log.Println("Error fetching receipts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch receipts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"receipts": receipts})
}

// Create new receipt
func CreateReceipt(c *gin.Context) {
	var input struct {
		Title      string   `json:"title" binding:"required"`
		ImageData  string   `json:"image_data" binding:"required"`
		Amount     *float64 `json:"amount"`
		Date       *string  `json:"date"`
		CategoryID *string  `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		log.Println("userID not found in context during receipt creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		log.Printf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Parse date if provided
	var receiptDate *time.Time
	if input.Date != nil && *input.Date != "" {
		layout := "2006-01-02"
		parsedDate, err := time.Parse(layout, *input.Date)
		if err != nil {
			log.Println("Error parsing date: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		receiptDate = &parsedDate
	}

	// Create directory structure based on month and year
	now := time.Now()
	yearMonth := now.Format("2006-01")
	dirPath := filepath.Join("uploads", "receipts", yearMonth)

	// Ensure directory exists
	if err := config.EnsureDir(dirPath); err != nil {
		log.Println("Error creating directory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
		return
	}

	// Generate filename
	filename := uuid.New().String() + ".jpg"
	filePath := filepath.Join(dirPath, filename)

	// Store relative path (without "uploads/" prefix) in database
	imagePath := filepath.Join("receipts", yearMonth, filename)

	// Save image data to file
	if err := config.SaveBase64Image(input.ImageData, filePath); err != nil {
		log.Println("Error saving image:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save receipt image"})
		return
	}

	receipt := models.Receipt{
		Title:     input.Title,
		ImagePath: imagePath,
		Amount:    input.Amount,
		Date:      receiptDate,
		UserID:    userIDUUID,
	}

	// Handle optional category
	if input.CategoryID != nil && *input.CategoryID != "" {
		categoryID, err := uuid.Parse(*input.CategoryID)
		if err != nil {
			log.Println("Error parsing category ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		receipt.CategoryID = &categoryID
	}

	if err := config.GetDB().Create(&receipt).Error; err != nil {
		log.Println("Error creating receipt: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create receipt"})
		return
	}

	// Load the category if it was set
	if receipt.CategoryID != nil {
		if err := config.GetDB().Preload("Category").First(&receipt, receipt.ID).Error; err != nil {
			log.Println("Error loading category:", err)
		}
	}

	c.JSON(http.StatusCreated, receipt)
}

// Update a receipt
func UpdateReceipt(c *gin.Context) {
	id := c.Param("ID")
	var receipt models.Receipt

	if err := config.GetDB().First(&receipt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	var input struct {
		Title      *string  `json:"title"`
		Amount     *float64 `json:"amount"`
		Date       *string  `json:"date"`
		CategoryID *string  `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error input: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Amount != nil {
		updates["amount"] = *input.Amount
	}
	if input.Date != nil {
		if *input.Date == "" {
			updates["date"] = nil
		} else {
			layout := "2006-01-02"
			date, err := time.Parse(layout, *input.Date)
			if err != nil {
				log.Println("Error timeParse: ", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
				return
			}
			updates["date"] = date
		}
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

	if err := config.GetDB().Model(&receipt).Updates(updates).Error; err != nil {
		log.Println("Error update: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receipt"})
		return
	}

	// Load the updated receipt with category
	if err := config.GetDB().Preload("Category").First(&receipt, id).Error; err != nil {
		log.Println("Error loading updated receipt:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated receipt"})
		return
	}

	c.JSON(http.StatusOK, receipt)
}

// Delete receipt
func DeleteReceipt(c *gin.Context) {
	id := c.Param("ID")
	var receipt models.Receipt

	if err := config.GetDB().First(&receipt, "id = ?", id).Error; err != nil {
		log.Println("Error fetch: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	// Delete the image file
	fullImagePath := filepath.Join("uploads", receipt.ImagePath)
	if err := config.DeleteFile(fullImagePath); err != nil {
		log.Println("Error deleting image file:", err)
		// Don't fail the request if file deletion fails, just log it
	}

	if err := config.GetDB().Delete(&receipt).Error; err != nil {
		log.Println("Error delete: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete receipt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Receipt deleted"})
}
