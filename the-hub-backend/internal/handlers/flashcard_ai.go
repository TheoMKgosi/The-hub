package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/TheoMKgosi/The-hub/internal/ai"
	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func extractJSON(s string) string {
	s = strings.TrimSpace(s)

	start := strings.Index(s, "[")
	if start == -1 {
		start = strings.Index(s, "{")
	}
	if start == -1 {
		return ""
	}

	end := strings.LastIndex(s, "]")
	if end == -1 {
		end = strings.LastIndex(s, "}")
	}
	if end == -1 || end < start {
		return ""
	}

	return s[start : end+1]
}

type GenerateFlashcardsRequest struct {
	PDF          string `json:"pdf"`
	NumCards     int    `json:"num_cards"`
	DeckID       string `json:"deck_id"`
	NewDeckName string `json:"new_deck_name"`
	Instruction string `json:"instruction"`
}

type FlashcardFromPDF struct {
	Front     string `json:"front"`
	Back     string `json:"back"`
	Category string `json:"category"`
}

type FlashcardPreview struct {
	Front     string `json:"front"`
	Back     string `json:"back"`
	Category string `json:"category"`
}

type GenerateFlashcardsResponse struct {
	Cards   []FlashcardPreview `json:"cards"`
	DeckID string            `json:"deck_id,omitempty"`
	DeckName string          `json:"deck_name,omitempty"`
	Message string          `json:"message"`
}

func GenerateFlashcardsFromPDF(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var req GenerateFlashcardsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	config.Logger.Infof("PDF data length: %d, NumCards: %d, DeckID: %s, NewDeckName: %s",
		len(req.PDF), req.NumCards, req.DeckID, req.NewDeckName)

	if req.PDF == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF data is required"})
		return
	}

	if req.NumCards < 5 {
		req.NumCards = 10
	}
	if req.NumCards > 50 {
		req.NumCards = 50
	}

	pdfBase64 := req.PDF
	if pdfBase64 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF data is required"})
		return
	}

	config.Logger.Infof("PDF data first 100 chars: %s", pdfBase64[:min(100, len(pdfBase64))])

	_, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		_, err = base64.URLEncoding.DecodeString(pdfBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PDF encoding"})
			return
		}
	}

	client, err := ai.GetOpenRouterClient()
	if err != nil {
		config.Logger.Errorf("Failed to get AI client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable"})
		return
	}

	systemPrompt := `You are a learning assistant. Generate flashcards from provided PDF content (slides or documents).
Focus on:
- Key concepts and definitions
- Important formulas or relationships
- Critical steps or processes
- Key terminology
Generate clear, concise questions and answers that test understanding.

Respond with ONLY a JSON array of objects, no other text. Each object contains:
- "front": question, term, or concept
- "back": answer, definition, or explanation
- "category": topic category (optional)`

	userPrompt := fmt.Sprintf("Generate %d flashcards from the provided PDF. Extract key concepts, definitions, and important points. Format as JSON array.", req.NumCards)
	if req.Instruction != "" {
		userPrompt = fmt.Sprintf("Generate %d flashcards from the provided PDF. %s Format as JSON array.", req.NumCards, req.Instruction)
	}

	aiResponse, err := client.GenerateWithDocument(pdfBase64, "application/pdf", userPrompt, systemPrompt)
	if err != nil {
		config.Logger.Errorf("Failed to generate flashcards: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate flashcards"})
		return
	}

	config.Logger.Infof("AI Response: %s", aiResponse)

	var flashcards []FlashcardFromPDF
	if err := json.Unmarshal([]byte(aiResponse), &flashcards); err != nil {
		jsonStr := extractJSON(aiResponse)
		if jsonStr == "" {
			config.Logger.Errorf("Failed to extract JSON from AI response: %s", aiResponse)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
			return
		}

		if err := json.Unmarshal([]byte(jsonStr), &flashcards); err != nil {
			config.Logger.Errorf("Failed to parse AI response: %v, raw: %s", err, jsonStr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
			return
		}
	}

	if len(flashcards) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No flashcards generated"})
		return
	}

	var deckID uuid.UUID
	deckName := req.NewDeckName

	if req.DeckID != "" {
		deckID, err = uuid.Parse(req.DeckID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
			return
		}

		var deck models.Deck
		if err := config.GetDB().Where("id = ? AND user_id = ?", deckID, userIDUUID).First(&deck).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
			return
		}
		deckName = deck.Name
	} else {
		if req.NewDeckName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Deck name required when creating new deck"})
			return
		}

		newDeck := models.Deck{
			ID:     uuid.New(),
			Name:   req.NewDeckName,
			UserID: userIDUUID,
		}

		if err := config.GetDB().Create(&newDeck).Error; err != nil {
			config.Logger.Errorf("Failed to create deck: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create deck"})
			return
		}

		deckID = newDeck.ID
	}

	cards := make([]FlashcardPreview, 0, len(flashcards))
	for _, fc := range flashcards {
		card := models.Card{
			ID:       uuid.New(),
			DeckID:   deckID,
			Question: fc.Front,
			Answer:   fc.Back,
		}

		if err := config.GetDB().Create(&card).Error; err != nil {
			config.Logger.Errorf("Failed to create card: %v", err)
			continue
		}

		cards = append(cards, FlashcardPreview{
			Front:     fc.Front,
			Back:     fc.Back,
			Category: fc.Category,
		})
	}

	c.JSON(http.StatusOK, GenerateFlashcardsResponse{
		Cards:    cards,
		DeckID:   deckID.String(),
		DeckName: deckName,
		Message:  fmt.Sprintf("Created %d flashcards", len(cards)),
	})
}
