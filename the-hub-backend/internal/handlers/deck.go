package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// Deck handlers

// Get all decks
func GetDecks(c *gin.Context) {
	var decks []models.Deck
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist",
		})
		return
	}

	result := config.GetDB().Where("user_id = ?", userID).Find(&decks)

	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"decks": decks,
	})
}

// Get a specific deck
func GetDeck(c *gin.Context) {
	deckID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Fatal("GetDeck error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting deck",
		})
		return
	}

	var deck models.Deck
	result := config.GetDB().First(&deck, deckID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Deck does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deck": deck,
	})
}

// Create a deck
func CreateDeck(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create Deck"})
		return
	}

	deck := models.Deck{
		Name:   input.Name,
		UserID: c.MustGet("userID").(uint),
	}

	if err := config.GetDB().Create(&deck).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Deck"})
		return
	}
	c.JSON(http.StatusCreated, deck)
}

// Update a specific deck
func UpdateDeck(c *gin.Context) {
	var deck models.Deck

	deckID := c.Param("ID")
	if err := config.GetDB().First(&deck, deckID).Error; err != nil {
		log.Println("Error ID: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Deck not found",
		})
		return
	}

	var input struct {
		Name *string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedDeck := map[string]interface{}{}
	if input.Name != nil {
		updatedDeck["name"] = *input.Name
	}

	if err := config.GetDB().Model(&deck).Updates(updatedDeck).Error; err != nil {
		log.Println("Error updating deck:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// reload the deck to get updated data
	if err := config.GetDB().First(&deck, deck.ID).Error; err != nil {
		log.Println("Error retrieving updated deck:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving updated deck"})
		return
	}

	c.JSON(http.StatusOK, deck)
}

// Delete a specific deck
func DeleteDeck(c *gin.Context) {
	var deck models.Deck

	deckID := c.Param("ID")
	if err := config.GetDB().First(&deck, deckID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Deck not found",
		})
		return
	}

	if err := config.GetDB().Delete(&deck).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deck)
}
