package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// GetCards godoc
// @Summary      Get all cards for a deck
// @Description  Fetch cards for a specific deck with optional ordering
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        deckID    path      int     true   "Deck ID"
// @Param        order_by  query     string  false  "Order by field (question, answer, easiness, interval, next_review, created_at)"  default(created_at)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.Card
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /decks/{deckID}/cards [get]
func GetCards(c *gin.Context) {
	deckIDStr := c.Param("deckID")
	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid deck ID param: %s", deckIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify deck belongs to user
	var deck models.Deck
	if err := config.GetDB().Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error; err != nil {
		config.Logger.Warnf("Deck ID %d not found for user %v: %v", deckID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	// Get query parameters for ordering
	orderBy := c.DefaultQuery("order_by", "created_at")
	sortDir := c.DefaultQuery("sort", "asc")

	// Validate order_by parameter
	validOrderFields := map[string]bool{
		"question":     true,
		"answer":       true,
		"easiness":     true,
		"interval":     true,
		"next_review":  true,
		"created_at":   true,
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

	var cards []models.Card
	config.Logger.Infof("Fetching cards for deck ID: %d with order: %s", deckID, orderClause)
	if err := config.GetDB().Where("deck_id = ?", deckID).Order(orderClause).Find(&cards).Error; err != nil {
		config.Logger.Errorf("Error fetching cards for deck %d: %v", deckID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch cards"})
		return
	}

	config.Logger.Infof("Found %d cards for deck ID %d", len(cards), deckID)
	c.JSON(http.StatusOK, gin.H{"cards": cards})
}

// GetDueCards godoc
// @Summary      Get cards due for review
// @Description  Fetch cards that are due for review in a specific deck
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        deckID   path      int  true  "Deck ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /decks/{deckID}/cards/due [get]
func GetDueCards(c *gin.Context) {
	deckIDStr := c.Param("deckID")
	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid deck ID param: %s", deckIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify deck belongs to user
	var deck models.Deck
	if err := config.GetDB().Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error; err != nil {
		config.Logger.Warnf("Deck ID %d not found for user %v: %v", deckID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	var cards []models.Card
	now := time.Now()
	config.Logger.Infof("Fetching due cards for deck ID: %d", deckID)
	if err := config.GetDB().Where("deck_id = ? AND next_review <= ?", deckID, now).Find(&cards).Error; err != nil {
		config.Logger.Errorf("Error fetching due cards for deck %d: %v", deckID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch due cards"})
		return
	}

	config.Logger.Infof("Found %d due cards for deck ID %d", len(cards), deckID)
	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
		"count": len(cards),
	})
}

// GetCard godoc
// @Summary      Get a specific card
// @Description  Fetch a specific card by ID
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Card ID"
// @Success      200  {object}  map[string]models.Card
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /cards/{ID} [get]
func GetCard(c *gin.Context) {
	cardIDStr := c.Param("ID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid card ID param: %s", cardIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var card models.Card
	config.Logger.Infof("Fetching card ID: %d for user ID: %v", cardID, userID)
	
	// Join with decks table to ensure user owns the deck that contains this card
	if err := config.GetDB().Joins("JOIN decks ON cards.deck_id = decks.id").
		Where("cards.id = ? AND decks.user_id = ?", cardID, userID).
		First(&card).Error; err != nil {
		config.Logger.Errorf("Card ID %d not found for user %v: %v", cardID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved card ID %d for user %v", cardID, userID)
	c.JSON(http.StatusOK, gin.H{"card": card})
}

// CreateCardRequest represents the request body for creating a card
type CreateCardRequest struct {
	DeckID   uint   `json:"deck_id" binding:"required" example:"1"`
	Question string `json:"question" binding:"required" example:"What is the capital of France?"`
	Answer   string `json:"answer" binding:"required" example:"Paris"`
}

// CreateCard godoc
// @Summary      Create a new card
// @Description  Create a new flashcard in a deck
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        card  body      CreateCardRequest  true  "Card creation data"
// @Success      201   {object}  models.Card
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /cards [post]
func CreateCard(c *gin.Context) {
	var input CreateCardRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid card input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for card", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during card creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify deck belongs to user
	var deck models.Deck
	if err := config.GetDB().Where("id = ? AND user_id = ?", input.DeckID, userID).First(&deck).Error; err != nil {
		config.Logger.Warnf("Deck ID %d not found for user %v during card creation: %v", input.DeckID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
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

	config.Logger.Infof("Creating card for deck %d: %s", input.DeckID, input.Question)
	if err := config.GetDB().Create(&card).Error; err != nil {
		config.Logger.Errorf("Error creating card for deck %d: %v", input.DeckID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create card"})
		return
	}

	config.Logger.Infof("Successfully created card ID %d for deck %d", card.ID, input.DeckID)
	c.JSON(http.StatusCreated, card)
}

// UpdateCardRequest represents the request body for updating a card
type UpdateCardRequest struct {
	Question *string `json:"question" example:"Updated question"`
	Answer   *string `json:"answer" example:"Updated answer"`
}

// UpdateCard godoc
// @Summary      Update a card
// @Description  Update a specific card by ID
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID    path      int                true  "Card ID"
// @Param        card  body      UpdateCardRequest  true  "Card update data"
// @Success      200   {object}  models.Card
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /cards/{ID} [put]
func UpdateCard(c *gin.Context) {
	cardIDStr := c.Param("ID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid card ID param for update: %s", cardIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during card update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var card models.Card
	// Join with decks table to ensure user owns the deck that contains this card
	if err := config.GetDB().Joins("JOIN decks ON cards.deck_id = decks.id").
		Where("cards.id = ? AND decks.user_id = ?", cardID, userID).
		First(&card).Error; err != nil {
		config.Logger.Warnf("Card not found for update: ID %d, User %v", cardID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	var input UpdateCardRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for card ID %d: %v", cardID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Question != nil {
		updates["question"] = *input.Question
	}
	if input.Answer != nil {
		updates["answer"] = *input.Answer
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for card update: ID %d", cardID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating card ID %d for user %v with data: %+v", cardID, userID, updates)
	if err := config.GetDB().Model(&card).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update card ID %d: %v", cardID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card"})
		return
	}

	// Reload the updated card
	if err := config.GetDB().First(&card, card.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated card ID %d: %v", card.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated card"})
		return
	}

	config.Logger.Infof("Successfully updated card ID %d for user %v", card.ID, userID)
	c.JSON(http.StatusOK, card)
}

// ReviewCardRequest represents the request body for reviewing a card
type ReviewCardRequest struct {
	Quality int `json:"quality" binding:"required,min=0,max=5" example:"4"`
}

// ReviewCard godoc
// @Summary      Review a card
// @Description  Review a card using SM-2 spaced repetition algorithm
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID      path      int                true  "Card ID"
// @Param        review  body      ReviewCardRequest  true  "Review data with quality (0-5)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /cards/{ID}/review [post]
func ReviewCard(c *gin.Context) {
	cardIDStr := c.Param("ID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid card ID param for review: %s", cardIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during card review")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var card models.Card
	// Join with decks table to ensure user owns the deck that contains this card
	if err := config.GetDB().Joins("JOIN decks ON cards.deck_id = decks.id").
		Where("cards.id = ? AND decks.user_id = ?", cardID, userID).
		First(&card).Error; err != nil {
		config.Logger.Warnf("Card not found for review: ID %d, User %v", cardID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	var input ReviewCardRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid review input for card ID %d: %v", cardID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quality must be between 0 and 5", "details": err.Error()})
		return
	}

	// SM-2 algorithm implementation
	quality := float64(input.Quality)
	config.Logger.Infof("Reviewing card ID %d with quality %d", cardID, input.Quality)
	
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
		config.Logger.Errorf("Error updating card after review ID %d: %v", cardID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card after review"})
		return
	}

	config.Logger.Infof("Successfully reviewed card ID %d, next review: %v", cardID, card.NextReview)
	c.JSON(http.StatusOK, gin.H{
		"card":          card,
		"next_interval": card.Interval,
		"next_review":   card.NextReview,
	})
}

// DeleteCard godoc
// @Summary      Delete a card
// @Description  Delete a specific card by ID
// @Tags         cards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Card ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /cards/{ID} [delete]
func DeleteCard(c *gin.Context) {
	cardIDStr := c.Param("ID")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid card ID param for delete: %s", cardIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during card deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var card models.Card
	// Join with decks table to ensure user owns the deck that contains this card
	if err := config.GetDB().Joins("JOIN decks ON cards.deck_id = decks.id").
		Where("cards.id = ? AND decks.user_id = ?", cardID, userID).
		First(&card).Error; err != nil {
		config.Logger.Warnf("Card not found for delete: ID %d, User %v", cardID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	config.Logger.Infof("Deleting card ID %d for user %v", cardID, userID)
	if err := config.GetDB().Delete(&card).Error; err != nil {
		config.Logger.Errorf("Failed to delete card ID %d: %v", cardID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card"})
		return
	}

	config.Logger.Infof("Successfully deleted card ID %d for user %v", cardID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Card deleted successfully", "card": card})
}
