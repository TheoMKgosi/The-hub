package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TheoMKgosi/The-hub/internal/models"
)

// ThemeSuggestion is a proposed theme grouping produced by the AI, including
// the tasks it recommends assigning to it (as task IDs echoed from the input).
type ThemeSuggestion struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Color       string   `json:"color"`
	Days        []string `json:"days"`
	Tasks       []string `json:"tasks"`
	Reasoning   string   `json:"reasoning"`
}

// ThemeTaskInput is a lightweight representation of a task sent to the theme
// suggestion assistant.
type ThemeTaskInput struct {
	TaskID      string   `json:"task_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    *int     `json:"priority"`
	Category    string   `json:"category"`
	TaskType    string   `json:"task_type"`
	Tags        []string `json:"tags"`
}

// SuggestThemes asks the AI to group the user's tasks into themed days.
// The AI proposes theme names, a weekday assignment for each theme, and the
// tasks that belong to each theme. Weekdays are constrained so that each
// weekday is assigned to at most one theme.
func (c *OpenRouterClient) SuggestThemes(tasks []models.Task) (string, error) {
	inputs := make([]ThemeTaskInput, 0, len(tasks))
	for _, t := range tasks {
		inputs = append(inputs, ThemeTaskInput{
			TaskID:      t.ID.String(),
			Title:       t.Title,
			Description: t.Description,
			Priority:    t.Priority,
			Category:    t.Category,
			TaskType:    t.TaskType,
			Tags:        t.Tags,
		})
	}

	tasksJSON, _ := json.Marshal(inputs)

	systemPrompt := `You are a themed-day planning assistant. The user practices "theme days": every weekday (monday..sunday) is assigned to exactly one theme, and on a given day only the tasks belonging to that day's theme are shown.

Group the user's tasks into themes (e.g. "Deep Work", "Admin", "Health", "Learning"). Rules:
- Each task must be assigned to exactly one theme.
- Every theme must be assigned one or more weekdays.
- Each weekday may belong to AT MOST ONE theme across all proposed themes (one theme per day).
- Prefer a small number of themes (2-5) that group similar tasks together.
- Use the task's task_id values exactly as given when listing which tasks belong to each theme.

Respond with ONLY a valid JSON array and no other text. Each object must be exactly:
- "name": theme name (short, capitalized)
- "description": one short sentence describing the theme
- "color": a hex color string like "#3B82F6"
- "days": array of lowercase weekday names (monday..sunday) for this theme
- "tasks": array of task_id strings assigned to this theme
- "reasoning": one short sentence justifying the grouping and day assignment`

	userContent := fmt.Sprintf(`Group these tasks into themed days:
%s

Every weekday must belong to at most one theme, and every task must be assigned to exactly one theme.`, strings.TrimSpace(string(tasksJSON)))

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.5,
		MaxTokens:   4096,
	})
}
