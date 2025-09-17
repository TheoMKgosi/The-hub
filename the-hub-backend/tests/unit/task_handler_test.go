package unit

import (
	"testing"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/handlers"
)

func TestParseNaturalLanguage(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedTitle string
		expectDueDate bool
	}{
		{
			name:          "next friday",
			input:         "test next friday",
			expectedTitle: "test",
			expectDueDate: true,
		},
		{
			name:          "next monday",
			input:         "buy groceries next monday",
			expectedTitle: "buy groceries",
			expectDueDate: true,
		},
		{
			name:          "this friday",
			input:         "meeting this friday",
			expectedTitle: "meeting",
			expectDueDate: true,
		},
		{
			name:          "tomorrow",
			input:         "call mom tomorrow",
			expectedTitle: "call mom",
			expectDueDate: true,
		},
		{
			name:          "no date",
			input:         "simple task",
			expectedTitle: "simple task",
			expectDueDate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, _, _, dueDate, err := handlers.ParseNaturalLanguage(tt.input)

			if err != nil {
				t.Errorf("ParseNaturalLanguage() error = %v", err)
				return
			}

			if title != tt.expectedTitle {
				t.Errorf("ParseNaturalLanguage() title = %v, expected %v", title, tt.expectedTitle)
			}

			if tt.expectDueDate && dueDate == nil {
				t.Errorf("ParseNaturalLanguage() expected due date but got nil")
			}

			if !tt.expectDueDate && dueDate != nil {
				t.Errorf("ParseNaturalLanguage() expected no due date but got %v", dueDate)
			}

			// For date tests, verify the date is reasonable
			if tt.expectDueDate && dueDate != nil {
				now := time.Now()
				if dueDate.Before(now) {
					t.Errorf("ParseNaturalLanguage() due date %v is in the past", dueDate)
				}
			}
		})
	}
}
