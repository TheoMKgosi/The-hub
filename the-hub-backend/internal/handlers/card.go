package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	deckID, err := uuid.Parse(deckIDStr)
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
		"question":    true,
		"answer":      true,
		"easiness":    true,
		"interval":    true,
		"next_review": true,
		"created_at":  true,
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
	deckID, err := uuid.Parse(deckIDStr)
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
	cardID, err := uuid.Parse(cardIDStr)
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
	DeckID   uuid.UUID `json:"deck_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Question string    `json:"question" binding:"required" example:"What is the capital of France?"`
	Answer   string    `json:"answer" binding:"required" example:"Paris"`
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
	cardID, err := uuid.Parse(cardIDStr)
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
	cardID, err := uuid.Parse(cardIDStr)
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
	cardID, err := uuid.Parse(cardIDStr)
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

// ExportCardRequest represents the export request
type ExportCardRequest struct {
	Format string `form:"format" binding:"required,oneof=json csv"`
}

// ExportCard represents a card for export (without internal fields)
type ExportCard struct {
	Question     string     `json:"question"`
	Answer       string     `json:"answer"`
	Easiness     float64    `json:"easiness"`
	Interval     int        `json:"interval"`
	Repetitions  int        `json:"repetitions"`
	LastReviewed *time.Time `json:"last_reviewed,omitempty"`
	NextReview   time.Time  `json:"next_review"`
}

// ExportData represents the complete export structure
type ExportData struct {
	DeckName   string       `json:"deck_name"`
	ExportedAt time.Time    `json:"exported_at"`
	Cards      []ExportCard `json:"cards"`
}

// ExportCards godoc
// @Summary      Export cards from a deck
// @Description  Export all cards from a specific deck in JSON or CSV format
// @Tags         cards
// @Accept       json
// @Produce      json,csv
// @Security     BearerAuth
// @Param        deckID   path      string  true   "Deck ID"
// @Param        format   query     string  true   "Export format (json or csv)"  Enums(json,csv)
// @Success      200      {object}  ExportData  "JSON export"
// @Success      200      {string}  string      "CSV export"
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /decks/{deckID}/cards/export [get]
func ExportCards(c *gin.Context) {
	deckIDStr := c.Param("deckID")
	deckID, err := uuid.Parse(deckIDStr)
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
		config.Logger.Warnf("Deck ID %s not found for user %v: %v", deckID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	// Get format parameter
	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format must be 'json' or 'csv'"})
		return
	}

	// Fetch all cards for the deck
	var cards []models.Card
	if err := config.GetDB().Where("deck_id = ?", deckID).Order("created_at").Find(&cards).Error; err != nil {
		config.Logger.Errorf("Error fetching cards for deck %s: %v", deckID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch cards"})
		return
	}

	config.Logger.Infof("Exporting %d cards from deck %s in %s format", len(cards), deckID, format)

	if format == "json" {
		exportData := ExportData{
			DeckName:   deck.Name,
			ExportedAt: time.Now(),
			Cards:      make([]ExportCard, len(cards)),
		}

		for i, card := range cards {
			var lastReviewed *time.Time
			if !card.LastReviewed.IsZero() {
				lastReviewed = &card.LastReviewed
			}

			exportData.Cards[i] = ExportCard{
				Question:     card.Question,
				Answer:       card.Answer,
				Easiness:     card.Easiness,
				Interval:     card.Interval,
				Repetitions:  card.Repetitions,
				LastReviewed: lastReviewed,
				NextReview:   card.NextReview,
			}
		}

		// Set headers for file download
		filename := fmt.Sprintf("%s_cards_%s.json", strings.ReplaceAll(deck.Name, " ", "_"), time.Now().Format("2006-01-02"))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Type", "application/json")

		c.JSON(http.StatusOK, exportData)
	} else { // CSV format
		// Set headers for CSV download
		filename := fmt.Sprintf("%s_cards_%s.csv", strings.ReplaceAll(deck.Name, " ", "_"), time.Now().Format("2006-01-02"))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Header("Content-Type", "text/csv")

		// Create CSV writer
		c.Writer.WriteHeader(http.StatusOK)

		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		// Write header
		header := []string{"question", "answer", "easiness", "interval", "repetitions", "last_reviewed", "next_review"}
		if err := writer.Write(header); err != nil {
			config.Logger.Errorf("Error writing CSV header: %v", err)
			return
		}

		// Write cards
		for _, card := range cards {
			lastReviewed := ""
			if !card.LastReviewed.IsZero() {
				lastReviewed = card.LastReviewed.Format(time.RFC3339)
			}

			record := []string{
				strings.ReplaceAll(card.Question, "\n", "\\n"), // Escape newlines
				strings.ReplaceAll(card.Answer, "\n", "\\n"),   // Escape newlines
				strconv.FormatFloat(card.Easiness, 'f', 2, 64),
				strconv.Itoa(card.Interval),
				strconv.Itoa(card.Repetitions),
				lastReviewed,
				card.NextReview.Format(time.RFC3339),
			}

			if err := writer.Write(record); err != nil {
				config.Logger.Errorf("Error writing CSV record: %v", err)
				return
			}
		}
	}

	config.Logger.Infof("Successfully exported %d cards from deck %s", len(cards), deckID)
}

// ImportCard represents a card for import
type ImportCard struct {
	Question     string  `json:"question" csv:"question"`
	Answer       string  `json:"answer" csv:"answer"`
	Easiness     float64 `json:"easiness,omitempty" csv:"easiness"`
	Interval     int     `json:"interval,omitempty" csv:"interval"`
	Repetitions  int     `json:"repetitions,omitempty" csv:"repetitions"`
	LastReviewed *string `json:"last_reviewed,omitempty" csv:"last_reviewed"`
	NextReview   *string `json:"next_review,omitempty" csv:"next_review"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	SuccessCount int           `json:"success_count"`
	ErrorCount   int           `json:"error_count"`
	Errors       []ImportError `json:"errors,omitempty"`
}

// ImportError represents an error that occurred during import
type ImportError struct {
	Row   int         `json:"row,omitempty"`
	Field string      `json:"field,omitempty"`
	Error string      `json:"error"`
	Card  *ImportCard `json:"card,omitempty"`
}

// ImportCards godoc
// @Summary      Import cards to a deck
// @Description  Import cards from JSON or CSV file to a specific deck
// @Tags         cards
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        deckID   path      string  true   "Deck ID"
// @Param        format   query     string  true   "Import format (json or csv)"  Enums(json,csv)
// @Param        file     formData  file    true   "File to import"
// @Success      200      {object}  ImportResult
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /decks/{deckID}/cards/import [post]
func ImportCards(c *gin.Context) {
	deckIDStr := c.Param("deckID")
	deckID, err := uuid.Parse(deckIDStr)
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
		config.Logger.Warnf("Deck ID %s not found for user %v: %v", deckID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	// Get format parameter
	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format must be 'json' or 'csv'"})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		config.Logger.Warnf("Error getting uploaded file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file size (10MB limit)
	const maxFileSize = 10 << 20 // 10MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 10MB"})
		return
	}

	var importCards []ImportCard
	var errors []ImportError

	if format == "json" {
		importCards, errors = ParseJSONImport(file)
	} else { // CSV format
		importCards, errors = ParseCSVImport(file)
	}

	if len(errors) > 0 && len(importCards) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "No valid cards found in import file",
			"errors": errors,
		})
		return
	}

	// Validate and prepare cards for import
	var validCards []models.Card
	for i, importCard := range importCards {
		if err := ValidateImportCard(importCard); err != nil {
			errors = append(errors, ImportError{
				Row:   i + 1,
				Error: err.Error(),
				Card:  &importCard,
			})
			continue
		}

		// Convert to model
		card := models.Card{
			DeckID:      deckID,
			Question:    strings.ReplaceAll(importCard.Question, "\\n", "\n"), // Unescape newlines
			Answer:      strings.ReplaceAll(importCard.Answer, "\\n", "\n"),   // Unescape newlines
			Easiness:    importCard.Easiness,
			Interval:    importCard.Interval,
			Repetitions: importCard.Repetitions,
			NextReview:  time.Now(),
		}

		// Set default values if not provided
		if card.Easiness == 0 {
			card.Easiness = 2.5
		}
		if card.Interval == 0 {
			card.Interval = 1
		}

		// Parse dates if provided
		if importCard.LastReviewed != nil && *importCard.LastReviewed != "" {
			if parsedTime, err := time.Parse(time.RFC3339, *importCard.LastReviewed); err == nil {
				card.LastReviewed = parsedTime
			}
		}

		if importCard.NextReview != nil && *importCard.NextReview != "" {
			if parsedTime, err := time.Parse(time.RFC3339, *importCard.NextReview); err == nil {
				card.NextReview = parsedTime
			}
		}

		validCards = append(validCards, card)
	}

	// Limit import to 1000 cards
	if len(validCards) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many cards. Maximum 1000 cards per import"})
		return
	}

	// Import cards in a transaction
	tx := config.GetDB().Begin()
	successCount := 0

	for _, card := range validCards {
		if err := tx.Create(&card).Error; err != nil {
			config.Logger.Errorf("Error importing card: %v", err)
			errors = append(errors, ImportError{
				Error: fmt.Sprintf("Failed to import card '%s': %v", card.Question, err),
			})
			continue
		}
		successCount++
	}

	if successCount == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "No cards were successfully imported",
			"errors": errors,
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		config.Logger.Errorf("Error committing import transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete import"})
		return
	}

	result := ImportResult{
		SuccessCount: successCount,
		ErrorCount:   len(errors),
		Errors:       errors,
	}

	config.Logger.Infof("Successfully imported %d cards to deck %s, %d errors", successCount, deckID, len(errors))
	c.JSON(http.StatusOK, result)
}

// ParseJSONImport parses cards from JSON format
func ParseJSONImport(file multipart.File) ([]ImportCard, []ImportError) {
	var errors []ImportError
	var importData struct {
		Cards []ImportCard `json:"cards"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&importData); err != nil {
		errors = append(errors, ImportError{
			Error: fmt.Sprintf("Invalid JSON format: %v", err),
		})
		return nil, errors
	}

	return importData.Cards, nil
}

// ParseCSVImport parses cards from CSV format
func ParseCSVImport(file multipart.File) ([]ImportCard, []ImportError) {
	var errors []ImportError
	var cards []ImportCard

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		errors = append(errors, ImportError{
			Error: fmt.Sprintf("Error reading CSV: %v", err),
		})
		return nil, errors
	}

	if len(records) < 2 {
		errors = append(errors, ImportError{
			Error: "CSV file must contain at least a header row and one data row",
		})
		return nil, errors
	}

	// Parse header
	header := records[0]
	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// Check required columns
	questionIdx, hasQuestion := headerMap["question"]
	answerIdx, hasAnswer := headerMap["answer"]
	if !hasQuestion || !hasAnswer {
		errors = append(errors, ImportError{
			Error: "CSV must contain 'question' and 'answer' columns",
		})
		return nil, errors
	}

	// Parse data rows
	for _, record := range records[1:] {
		if len(record) == 0 {
			continue // Skip empty rows
		}

		var card ImportCard

		// Required fields
		if questionIdx < len(record) {
			card.Question = record[questionIdx]
		}
		if answerIdx < len(record) {
			card.Answer = record[answerIdx]
		}

		// Optional fields
		if easinessIdx, ok := headerMap["easiness"]; ok && easinessIdx < len(record) {
			if val, err := strconv.ParseFloat(record[easinessIdx], 64); err == nil {
				card.Easiness = val
			}
		}

		if intervalIdx, ok := headerMap["interval"]; ok && intervalIdx < len(record) {
			if val, err := strconv.Atoi(record[intervalIdx]); err == nil {
				card.Interval = val
			}
		}

		if repetitionsIdx, ok := headerMap["repetitions"]; ok && repetitionsIdx < len(record) {
			if val, err := strconv.Atoi(record[repetitionsIdx]); err == nil {
				card.Repetitions = val
			}
		}

		if lastReviewedIdx, ok := headerMap["last_reviewed"]; ok && lastReviewedIdx < len(record) && record[lastReviewedIdx] != "" {
			card.LastReviewed = &record[lastReviewedIdx]
		}

		if nextReviewIdx, ok := headerMap["next_review"]; ok && nextReviewIdx < len(record) && record[nextReviewIdx] != "" {
			card.NextReview = &record[nextReviewIdx]
		}

		cards = append(cards, card)
	}

	return cards, errors
}

// ValidateImportCard validates an import card
func ValidateImportCard(card ImportCard) error {
	if strings.TrimSpace(card.Question) == "" {
		return fmt.Errorf("question is required")
	}
	if strings.TrimSpace(card.Answer) == "" {
		return fmt.Errorf("answer is required")
	}
	if len(card.Question) > 1000 {
		return fmt.Errorf("question too long (max 1000 characters)")
	}
	if len(card.Answer) > 2000 {
		return fmt.Errorf("answer too long (max 2000 characters)")
	}
	if card.Easiness < 0 || card.Easiness > 5 {
		return fmt.Errorf("easiness must be between 0 and 5")
	}
	if card.Interval < 0 {
		return fmt.Errorf("interval must be non-negative")
	}
	if card.Repetitions < 0 {
		return fmt.Errorf("repetitions must be non-negative")
	}
	return nil
}
