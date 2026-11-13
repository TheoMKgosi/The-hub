package unit

import (
	"testing"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/handlers"
)

func TestNormalizeTags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "trims whitespace and drops empties",
			input:    []string{"  work  ", "", "   ", "focus"},
			expected: []string{"work", "focus"},
		},
		{
			name:     "deduplicates case-insensitively keeping first casing",
			input:    []string{"Deep Work", "deep work", "DEEP", "admin"},
			expected: []string{"Deep Work", "DEEP", "admin"},
		},
		{
			name:     "empty input yields empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "nil input yields empty slice",
			input:    nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlers.NormalizeTags(tt.input)

			if len(got) != len(tt.expected) {
				t.Fatalf("NormalizeTags() length = %d, expected %d (%v)", len(got), len(tt.expected), got)
			}
			for i := range tt.expected {
				if got[i] != tt.expected[i] {
					t.Errorf("NormalizeTags()[%d] = %q, expected %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestComputeStreak(t *testing.T) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)
	date := func(daysAgo int) time.Time {
		return today.AddDate(0, 0, -daysAgo)
	}

	tests := []struct {
		name     string
		dates    []time.Time
		today    time.Time
		expected int
	}{
		{
			name:     "no entries means zero streak",
			dates:    []time.Time{},
			today:    today,
			expected: 0,
		},
		{
			name:     "logged today counts from today",
			dates:    []time.Time{date(0), date(1), date(2)},
			today:    today,
			expected: 3,
		},
		{
			name:     "streak ending yesterday still counts",
			dates:    []time.Time{date(1), date(2), date(3)},
			today:    today,
			expected: 3,
		},
		{
			name:     "gap breaks the streak",
			dates:    []time.Time{date(0), date(2), date(3)},
			today:    today,
			expected: 1,
		},
		{
			name:     "streak broken long ago is zero",
			dates:    []time.Time{date(5), date(6)},
			today:    today,
			expected: 0,
		},
		{
			name:     "duplicate dates do not inflate the streak",
			dates:    []time.Time{date(0), date(0), date(1)},
			today:    today,
			expected: 2,
		},
		{
			name:     "times during the day normalize to the same date",
			dates:    []time.Time{date(0).Add(3 * time.Hour), date(1).Add(6 * time.Hour)},
			today:    today,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlers.ComputeStreak(tt.dates, tt.today)
			if got != tt.expected {
				t.Errorf("ComputeStreak() = %d, expected %d", got, tt.expected)
			}
		})
	}
}
