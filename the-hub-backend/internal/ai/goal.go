package ai

import (
	"encoding/json"
	"fmt"
)

func (c *OpenRouterClient) GenerateGoalTaskRecommendations(goalTitle, goalDesc string, existingTasks []string, goalPriority *int, dueDate, category string, blockedBy, neededBy []map[string]string) (string, error) {
	systemPrompt := `You are a task recommendation assistant. Generate relevant tasks for goals based on existing tasks to avoid duplicates.
Generate realistic and actionable task recommendations. Respond ONLY with a valid JSON array, no other text.`

	existingTasksJSON, _ := json.Marshal(existingTasks)

	var priorityInfo string
	if goalPriority != nil {
		priorityInfo = fmt.Sprintf("Goal priority: %d (1=highest, 5=lowest)", *goalPriority)
	}

	var dueDateInfo string
	if dueDate != "" {
		dueDateInfo = fmt.Sprintf("Goal due date: %s", dueDate)
	}

	var categoryInfo string
	if category != "" {
		categoryInfo = fmt.Sprintf("Goal category: %s", category)
	}

	var blockedInfo string
	if len(blockedBy) > 0 {
		blockedInfo = "\n⚠️ BLOCKED TASKS (waiting on external tasks):\n"
		for _, b := range blockedBy {
			blockedInfo += fmt.Sprintf("  - %s [%s]: %s\n", b["title"], b["status"], b["reason"])
		}
	}

	var neededInfo string
	if len(neededBy) > 0 {
		neededInfo = "\n📌 TASKS NEEDED BY OTHER GOALS:\n"
		for _, n := range neededBy {
			neededInfo += fmt.Sprintf("  - %s [%s]: %s\n", n["title"], n["status"], n["reason"])
		}
	}

	userContent := fmt.Sprintf(`Goal: %s
Description: %s
%s
%s
%s
%s
%s

Existing related tasks:
%s

Generate 5-10 new task recommendations for this goal. Response should be a JSON array with objects containing:
- "title": task title (action-specific)
- "description": task description (what needs to be done)
- "priority": 1-5 priority based on goal priority and task importance
- "estimated_hours": estimated hours (1-8, be realistic)
- "reasoning": brief explanation (1-2 sentences)

Ensure new tasks don't duplicate existing ones. Focus on actionable items that directly contribute to achieving the goal.`, goalTitle, goalDesc, priorityInfo, dueDateInfo, categoryInfo, blockedInfo, neededInfo, string(existingTasksJSON))

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.7,
		MaxTokens:   4096,
	})
}
