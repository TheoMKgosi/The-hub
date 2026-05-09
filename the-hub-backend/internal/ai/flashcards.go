package ai

import (
	"fmt"
)

type PDFToFlashcardInput struct {
	PDFBase64 string `json:"pdf_base64"`
	NumCards  int    `json:"num_cards"`
}

type FlashcardFromPDF struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type GenerateFlashcardsRequest struct {
	PDF         string `json:"pdf" binding:"required"`
	NumCards    int    `json:"num_cards"`
	DeckID      string `json:"deck_id"`
	NewDeckName string `json:"new_deck_name"`
	Instruction string `json:"instruction"`
}

func (c *OpenRouterClient) GenerateFlashcardsFromPDF(pdfBase64 string, numCards int, instruction string) (string, error) {
	var systemPrompt string
	if instruction != "" {
		systemPrompt = fmt.Sprintf(`You are a learning assistant. Generate flashcards from provided PDF content.
%s
Focus on key concepts, definitions, and important points.`, instruction)
	} else {
		systemPrompt = `You are a learning assistant. Generate flashcards from provided PDF content (slides or documents).
Focus on:
- Key concepts and definitions
- Important formulas or relationships
- Critical steps or processes
- Key terminology
Generate clear, concise questions and answers that test understanding.`
	}

	userContent := fmt.Sprintf(`Generate %d flashcards from the provided PDF.

Respond with a JSON array of objects, each containing:
- "front": question, term, or concept
- "back": answer, definition, or explanation
- "category": topic category or section (optional)

Extract the content from the PDF document and create informative flashcards.`, numCards)

	return c.GenerateWithDocument(pdfBase64, "application/pdf", userContent, systemPrompt)
}
