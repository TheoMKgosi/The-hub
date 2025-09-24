package unit

import (
	"bytes"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/TheoMKgosi/The-hub/internal/handlers"
	"github.com/stretchr/testify/assert"
)

// mockMultipartFile implements multipart.File for testing
type mockMultipartFile struct {
	*bytes.Reader
}

func (m *mockMultipartFile) Close() error {
	return nil
}

func newMockMultipartFile(data string) multipart.File {
	return &mockMultipartFile{
		Reader: bytes.NewReader([]byte(data)),
	}
}

func TestValidateImportCard(t *testing.T) {
	tests := []struct {
		name     string
		card     handlers.ImportCard
		expected bool
	}{
		{
			name: "valid card",
			card: handlers.ImportCard{
				Question: "What is 2+2?",
				Answer:   "4",
			},
			expected: true,
		},
		{
			name: "missing question",
			card: handlers.ImportCard{
				Answer: "4",
			},
			expected: false,
		},
		{
			name: "missing answer",
			card: handlers.ImportCard{
				Question: "What is 2+2?",
			},
			expected: false,
		},
		{
			name: "question too long",
			card: handlers.ImportCard{
				Question: strings.Repeat("a", 1001),
				Answer:   "4",
			},
			expected: false,
		},
		{
			name: "answer too long",
			card: handlers.ImportCard{
				Question: "What is 2+2?",
				Answer:   strings.Repeat("a", 2001),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handlers.ValidateImportCard(tt.card)
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestParseJSONImport(t *testing.T) {
	jsonData := `{
		"cards": [
			{
				"question": "What is 2+2?",
				"answer": "4",
				"easiness": 2.5
			},
			{
				"question": "What is the capital of France?",
				"answer": "Paris"
			}
		]
	}`

	file := newMockMultipartFile(jsonData)
	cards, errors := handlers.ParseJSONImport(file)

	assert.Len(t, cards, 2)
	assert.Len(t, errors, 0)
	assert.Equal(t, "What is 2+2?", cards[0].Question)
	assert.Equal(t, "4", cards[0].Answer)
	assert.Equal(t, 2.5, cards[0].Easiness)
	assert.Equal(t, "What is the capital of France?", cards[1].Question)
	assert.Equal(t, "Paris", cards[1].Answer)
}

func TestParseCSVImport(t *testing.T) {
	csvData := `question,answer,easiness,interval,repetitions
What is 2+2?,4,2.5,1,0
What is the capital of France?,Paris,2.5,1,0`

	file := newMockMultipartFile(csvData)
	cards, errors := handlers.ParseCSVImport(file)

	assert.Len(t, cards, 2)
	assert.Len(t, errors, 0)
	assert.Equal(t, "What is 2+2?", cards[0].Question)
	assert.Equal(t, "4", cards[0].Answer)
	assert.Equal(t, 2.5, cards[0].Easiness)
	assert.Equal(t, "What is the capital of France?", cards[1].Question)
	assert.Equal(t, "Paris", cards[1].Answer)
}

func TestParseCSVImport_InvalidHeader(t *testing.T) {
	csvData := `invalid,header
What is 2+2?,4`

	file := newMockMultipartFile(csvData)
	cards, errors := handlers.ParseCSVImport(file)

	assert.Len(t, cards, 0)
	assert.Len(t, errors, 1)
	assert.Contains(t, errors[0].Error, "CSV must contain 'question' and 'answer' columns")
}

// Note: Integration tests for actual HTTP endpoints would require database setup
// and are better suited for the integration test suite
