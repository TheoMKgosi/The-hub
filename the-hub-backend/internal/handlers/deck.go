package handlers

import (
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// GetDecks godoc
// @Summary      Get all decks
// @Description  Fetch decks for the logged-in user with optional ordering
// @Tags         decks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_by  query     string  false  "Order by field (name, created_at)"  default(created_at)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.Deck
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /decks [get]
func GetDecks(c *gin.Context) {
	var decks []models.Deck
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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

	config.Logger.Infof("Fetching decks for user ID: %v with order: %s", userID, orderClause)
	if err := config.GetDB().Where("user_id = ?", userID).Order(orderClause).Find(&decks).Error; err != nil {
		config.Logger.Errorf("Error fetching decks for user %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch decks"})
		return
	}

	config.Logger.Infof("Found %d decks for user ID %v", len(decks), userID)
	c.JSON(http.StatusOK, gin.H{"decks": decks})
}

// GetDeck godoc
// @Summary      Get a specific deck
// @Description  Fetch a specific deck by ID for the logged-in user
// @Tags         decks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Deck ID"
// @Success      200  {object}  map[string]models.Deck
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /decks/{ID} [get]
func GetDeck(c *gin.Context) {
	deckIDStr := c.Param("ID")
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

	config.Logger.Infof("Fetching deck ID: %d for user ID: %v", deckID, userID)
	var deck models.Deck
	// Ensure user can only access their own decks
	if err := config.GetDB().Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error; err != nil {
		config.Logger.Errorf("Deck ID %d not found for user %v: %v", deckID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved deck ID %d for user %v", deckID, userID)
	c.JSON(http.StatusOK, gin.H{"deck": deck})
}

// CreateDeckRequest represents the request body for creating a deck
type CreateDeckRequest struct {
	Name string `json:"name" binding:"required" example:"Spanish Vocabulary"`
}

// CreateDeck godoc
// @Summary      Create a new deck
// @Description  Create a new flashcard deck for the logged-in user
// @Tags         decks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        deck  body      CreateDeckRequest  true  "Deck creation data"
// @Success      201   {object}  models.Deck
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /decks [post]
func CreateDeck(c *gin.Context) {
	var input CreateDeckRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid deck input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for deck", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during deck creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Check for duplicate deck name for this user
	var existingDeck models.Deck
	if err := config.GetDB().Where("name = ? AND user_id = ?", input.Name, userIDUint).First(&existingDeck).Error; err == nil {
		config.Logger.Warnf("Duplicate deck name '%s' for user %d", input.Name, userIDUint)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deck name already exists"})
		return
	}

	deck := models.Deck{
		Name:   input.Name,
		UserID: userIDUint,
	}

	config.Logger.Infof("Creating deck for user %d: %s", userIDUint, input.Name)
	if err := config.GetDB().Create(&deck).Error; err != nil {
		config.Logger.Errorf("Error creating deck for user %d: %v", userIDUint, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create deck"})
		return
	}

	config.Logger.Infof("Successfully created deck ID %d for user %d", deck.ID, userIDUint)
	c.JSON(http.StatusCreated, deck)
}

// UpdateDeckRequest represents the request body for updating a deck
type UpdateDeckRequest struct {
	Name *string `json:"name" example:"Updated deck name"`
}

// UpdateDeck godoc
// @Summary      Update a deck
// @Description  Update a specific deck by ID for the logged-in user
// @Tags         decks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID    path      int                true  "Deck ID"
// @Param        deck  body      UpdateDeckRequest  true  "Deck update data"
// @Success      200   {object}  models.Deck
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /decks/{ID} [put]
func UpdateDeck(c *gin.Context) {
	deckIDStr := c.Param("ID")
	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid deck ID param for update: %s", deckIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during deck update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var deck models.Deck
	// Ensure user can only update their own decks
	if err := config.GetDB().Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error; err != nil {
		config.Logger.Warnf("Deck not found for update: ID %d, User %v", deckID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	var input UpdateDeckRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for deck ID %d: %v", deckID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != nil {
		// Check for duplicate deck name for this user (excluding current deck)
		var existingDeck models.Deck
		if err := config.GetDB().Where("name = ? AND user_id = ? AND id != ?", *input.Name, userID, deckID).First(&existingDeck).Error; err == nil {
			config.Logger.Warnf("Duplicate deck name '%s' for user %v during update", *input.Name, userID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Deck name already exists"})
			return
		}
		updates["name"] = *input.Name
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for deck update: ID %d", deckID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating deck ID %d for user %v with data: %+v", deckID, userID, updates)
	if err := config.GetDB().Model(&deck).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update deck ID %d: %v", deckID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update deck"})
		return
	}

	// Reload the updated deck
	if err := config.GetDB().First(&deck, deck.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated deck ID %d: %v", deck.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated deck"})
		return
	}

	config.Logger.Infof("Successfully updated deck ID %d for user %v", deck.ID, userID)
	c.JSON(http.StatusOK, deck)
}

// DeleteDeck godoc
// @Summary      Delete a deck
// @Description  Delete a specific deck by ID for the logged-in user
// @Tags         decks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Deck ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /decks/{ID} [delete]
func DeleteDeck(c *gin.Context) {
	deckIDStr := c.Param("ID")
	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid deck ID param for delete: %s", deckIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during deck deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var deck models.Deck
	// Ensure user can only delete their own decks
	if err := config.GetDB().Where("id = ? AND user_id = ?", deckID, userID).First(&deck).Error; err != nil {
		config.Logger.Warnf("Deck not found for delete: ID %d, User %v", deckID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	// Check if deck has cards
	var cardCount int64
	if err := config.GetDB().Model(&models.Card{}).Where("deck_id = ?", deckID).Count(&cardCount).Error; err != nil {
		config.Logger.Errorf("Error checking card count for deck ID %d: %v", deckID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check deck usage"})
		return
	}

	if cardCount > 0 {
		config.Logger.Warnf("Cannot delete deck ID %d: contains %d cards", deckID, cardCount)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete deck that contains cards. Please delete all cards first."})
		return
	}

	config.Logger.Infof("Deleting deck ID %d for user %v", deckID, userID)
	if err := config.GetDB().Delete(&deck).Error; err != nil {
		config.Logger.Errorf("Failed to delete deck ID %d: %v", deckID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deck"})
		return
	}

	config.Logger.Infof("Successfully deleted deck ID %d for user %v", deckID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Deck deleted successfully", "deck": deck})
}
