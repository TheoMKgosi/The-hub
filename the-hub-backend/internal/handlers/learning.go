package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

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

// Card handlers

// Get all cards for a deck
func GetCards(c *gin.Context) {
	deckID, err := strconv.Atoi(c.Param("deckID"))
	if err != nil {
		log.Println("GetCards error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID",
		})
		return
	}

	var cards []models.Card
	result := config.GetDB().Where("deck_id = ?", deckID).Find(&cards)

	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
	})
}

// Get cards due for review
func GetDueCards(c *gin.Context) {
	deckID, err := strconv.Atoi(c.Param("deckID"))
	if err != nil {
		log.Println("GetDueCards error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid deck ID",
		})
		return
	}

	var cards []models.Card
	now := time.Now()
	result := config.GetDB().Where("deck_id = ? AND next_review <= ?", deckID, now).Find(&cards)

	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
		"count": len(cards),
	})
}

// Get a specific card
func GetCard(c *gin.Context) {
	cardID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Fatal("GetCard error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting card",
		})
		return
	}

	var card models.Card
	result := config.GetDB().First(&card, cardID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Card does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"card": card,
	})
}

// Create a card
func CreateCard(c *gin.Context) {
	var input struct {
		DeckID   uint   `json:"deck_id" binding:"required"`
		Question string `json:"question" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create Card"})
		return
	}

	card := models.Card{
		DeckID:       input.DeckID,
		Question:     input.Question,
		Answer:       input.Answer,
		Easiness:     2.5,
		Interval:     1,
		Repetitions:  0,
		LastReviewed: time.Time{},
		NextReview:   time.Now(),
	}

	if err := config.GetDB().Create(&card).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create Card"})
		return
	}
	c.JSON(http.StatusCreated, card)
}

// Update a specific card
func UpdateCard(c *gin.Context) {
	var card models.Card

	cardID := c.Param("ID")
	if err := config.GetDB().First(&card, cardID).Error; err != nil {
		log.Println("Error ID: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found",
		})
		return
	}

	var input struct {
		Question *string `json:"question"`
		Answer   *string `json:"answer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCard := map[string]interface{}{}
	if input.Question != nil {
		updatedCard["question"] = *input.Question
	}
	if input.Answer != nil {
		updatedCard["answer"] = *input.Answer
	}

	if err := config.GetDB().Model(&card).Updates(updatedCard).Error; err != nil {
		log.Println("Error updating card:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// reload the card to get updated data
	if err := config.GetDB().First(&card, card.ID).Error; err != nil {
		log.Println("Error retrieving updated card:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving updated card"})
		return
	}

	c.JSON(http.StatusOK, card)
}

// Review a card (SM-2 spaced repetition algorithm)
func ReviewCard(c *gin.Context) {
	var card models.Card

	cardID := c.Param("ID")
	if err := config.GetDB().First(&card, cardID).Error; err != nil {
		log.Println("Error ID: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found",
		})
		return
	}

	var input struct {
		Quality int `json:"quality" binding:"required,min=0,max=5"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quality must be between 0 and 5"})
		return
	}

	// SM-2 algorithm implementation
	quality := float64(input.Quality)
	
	if quality >= 3 {
		if card.Repetitions == 0 {
			card.Interval = 1
		} else if card.Repetitions == 1 {
			card.Interval = 6
		} else {
			card.Interval = int(float64(card.Interval) * card.Easiness)
		}
		card.Repetitions++
	} else {
		card.Repetitions = 0
		card.Interval = 1
	}

	card.Easiness = card.Easiness + (0.1 - (5-quality)*(0.08+(5-quality)*0.02))
	if card.Easiness < 1.3 {
		card.Easiness = 1.3
	}

	card.LastReviewed = time.Now()
	card.NextReview = time.Now().AddDate(0, 0, card.Interval)

	if err := config.GetDB().Save(&card).Error; err != nil {
		log.Println("Error updating card:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"card":          card,
		"next_interval": card.Interval,
		"next_review":   card.NextReview,
	})
}

// Delete a specific card
func DeleteCard(c *gin.Context) {
	var card models.Card

	cardID := c.Param("ID")
	if err := config.GetDB().First(&card, cardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Card not found",
		})
		return
	}

	if err := config.GetDB().Delete(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}
