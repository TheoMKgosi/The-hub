package ai

import (
	"encoding/json"
	"fmt"
)

type TaskEnhancementInput struct {
	TaskID      string `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SubtaskSuggestion struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	EstimatedHours int    `json:"estimated_hours"`
}

type TaskEnhancement struct {
	TaskID         string              `json:"task_id"`
	OriginalTitle  string              `json:"original_title"`
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	Priority       int                 `json:"priority"`
	EstimatedHours int                 `json:"estimated_hours"`
	Subtasks       []SubtaskSuggestion `json:"subtasks"`
}

func (c *OpenRouterClient) EnhanceTasks(tasks []TaskEnhancementInput) (string, error) {
	tasksJSON, _ := json.Marshal(tasks)

	systemPrompt := `You are a task enhancement assistant. Analyze each task and provide improvements including:
- Better, more specific title
- Enhanced description with actionable details
- Priority level (1-5, where 1 is highest)
- Time estimate in hours
- For large tasks (>4 hours), break into smaller subtasks

Respond with a JSON array of objects, each containing:
- "task_id": the original task ID
- "original_title": the original task title
- "title": improved task title
- "description": enhanced description with actionable details
- "priority": 1-5 priority level
- "estimated_hours": estimated time in hours (1-40) integer not float
- "subtasks": array of subtask objects with title, description, estimated_hours (only for tasks that will take 4+ hours)

For tasks under 4 hours, return empty subtasks array.`

	userContent := fmt.Sprintf(`Analyze and enhance these tasks:
%s

Provide improvements for each task. Focus on making titles specific and descriptions actionable.`, tasksJSON)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.7,
		MaxTokens:   8192,
	})
}

func (c *OpenRouterClient) SuggestGoalSteps(title, description string) (string, error) {
	systemPrompt := `You are a goal planning assistant. Break down goals into actionable steps.`
	userContent := fmt.Sprintf(`Goal: %s
Description: %s

Break this goal down into specific, actionable steps. Respond with a JSON array of objects, each containing:
- "title": step title
- "description": step description
- "priority": 1-5 priority level
- "estimated_hours": estimated time in hours integer not float
- "reasoning": why this step is important`, title, description)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.7,
		MaxTokens:   4096,
	})
}

func (c *OpenRouterClient) ProcessTasks(tasks []string, instruction string) (string, error) {
	tasksJSON, _ := json.Marshal(tasks)

	systemPrompt := `You are a task organization assistant. Analyze the provided list of tasks and organize them effectively.`
	userContent := fmt.Sprintf(`%s

Here is my task list:
%s

Please analyze and organize these tasks. Respond with a JSON array of objects, each containing:
- "title": task title
- "description": brief description
- "priority": 1-5 priority level
- "estimated_hours": estimated time in hours
- "reasoning": why you organized them this way`, instruction, tasksJSON)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.7,
		MaxTokens:   4096,
	})
}
