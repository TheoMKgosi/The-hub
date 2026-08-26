package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/models"
)

const defaultSchedulingDays = 7

// ScheduleSlot represents the result of the AI scheduling pending tasks into free time slots.
type ScheduleSlot struct {
	Title     string    `json:"title"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Reasoning string    `json:"reasoning"`
}

// ScheduleTaskInput is a lightweight representation of a task sent to the scheduling assistant.
type ScheduleTaskInput struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Priority     *int       `json:"priority"`
	DueDate      *time.Time `json:"due_date"`
	TimeEstimate *int       `json:"time_estimate_minutes"`
	Category     string     `json:"category"`
	TaskType     string     `json:"task_type"`
}

// ScheduledEventInput is a lightweight representation of an existing calendar event.
type ScheduledEventInput struct {
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// buildScheduleInputs converts tasks and existing events into the lightweight prompt payloads.
func buildScheduleInputs(tasks []models.Task, existingEvents []models.ScheduledTask) ([]ScheduleTaskInput, []ScheduledEventInput) {
	taskInputs := make([]ScheduleTaskInput, 0, len(tasks))
	for _, t := range tasks {
		taskInputs = append(taskInputs, ScheduleTaskInput{
			Title:        t.Title,
			Description:  t.Description,
			Priority:     t.Priority,
			DueDate:      t.DueDate,
			TimeEstimate: t.TimeEstimate,
			Category:     t.Category,
			TaskType:     t.TaskType,
		})
	}

	eventInputs := make([]ScheduledEventInput, 0, len(existingEvents))
	for _, s := range existingEvents {
		eventInputs = append(eventInputs, ScheduledEventInput{
			Title: s.Title,
			Start: s.Start,
			End:   s.End,
		})
	}

	return taskInputs, eventInputs
}

// ScheduleTasksIntoSchedule asks the AI to slot the given tasks into free time slots,
// avoiding conflicts with the user's existing calendar events.
func (c *OpenRouterClient) ScheduleTasksIntoSchedule(tasks []models.Task, existingEvents []models.ScheduledTask, days int) (string, error) {
	if days <= 0 {
		days = defaultSchedulingDays
	}

	taskInputs, eventInputs := buildScheduleInputs(tasks, existingEvents)

	tasksJSON, _ := json.Marshal(taskInputs)
	eventsJSON, _ := json.Marshal(eventInputs)

	now := time.Now()
	windowEnd := now.AddDate(0, 0, days)

	systemPrompt := `You are an expert calendar scheduling assistant. You are given a list of tasks and the user's existing calendar events. Assign each task to a realistic time slot within the scheduling window.

Rules:
- Slots must NOT overlap each other or any existing calendar event.
- Use each task's time_estimate_minutes to set the slot duration (default to 60 minutes when missing). Dispatch hard limits: each slot cannot exceed 24 hours.
- Stay within normal working / productive hours (approx. 09:00-18:00) unless a task's context clearly requires otherwise.
- Prioritize tasks with earlier due dates, then higher priority values (1 = highest). Order output by suggested start time.
- Do not schedule tasks whose title matches an existing calendar event.
- Vary suggestions across the window rather than stacking everything on day one.

Respond with ONLY a valid JSON array and no other text. Each object must be exactly:
- "title": the task title
- "start": RFC 3339 start timestamp
- "end": RFC 3339 end timestamp
- "reasoning": one short sentence justifying the chosen slot`

	userContent := fmt.Sprintf(`Current time: %s
Scheduling window (inclusive): %s to %s

Tasks to schedule:
%s

Existing calendar events (must not overlap):
%s

Assign each task to a slot. Output as many tasks as can fit without conflicts.`,
		now.Format(time.RFC3339),
		now.Format(time.RFC3339),
		windowEnd.Format(time.RFC3339),
		strings.TrimSpace(string(tasksJSON)),
		strings.TrimSpace(string(eventsJSON)),
	)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.3,
		MaxTokens:   8192,
	})
}
