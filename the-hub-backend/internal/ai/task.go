package ai

// // GenerateGoalTaskRecommendations generates AI-powered task recommendations for a goal
// func (c *OpenRouterClient) GenerateGoalTaskRecommendations(goalTitle, goalDescription string, existingTasks []string) (string, error) {
// 	prompt := fmt.Sprintf(`You are an AI assistant helping users break down their goals into actionable tasks.
//
// Goal Title: %s
// Goal Description: %s
//
// Existing Tasks:
// %s
//
// Please suggest 3-5 specific, actionable tasks that would help achieve this goal. For each task, provide:
// 1. A clear, concise title
// 2. A brief description of what the task involves
// 3. A priority level (1-5, where 5 is highest priority)
// 4. An estimated time in hours
// 5. A brief explanation of why this task is important for achieving the goal
//
// Return the response as a JSON array of objects with the following structure:
// [
//   {
//     "title": "Task Title",
//     "description": "Task description",
//     "priority": 3,
//     "estimated_hours": 2,
//     "reasoning": "Why this task is important"
//   }
// ]
//
// Make sure the tasks are:
// - Specific and actionable
// - Different from existing tasks
// - Realistic in scope
// - Ordered by priority (highest first)
// - Focused on achieving the goal`, goalTitle, goalDescription, formatTasks(existingTasks))
//
// 	messages := []Message{
// 		{
// 			Role:    "system",
// 			Content: "You are a goal achievement assistant that helps break down goals into actionable tasks. Always respond with valid JSON.",
// 		},
// 		{
// 			Role:    "user",
// 			Content: prompt,
// 		},
// 	}
//
// 	response, err := c.MakeRequest("anthropic/claude-3-haiku", messages)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	if len(response.Choices) == 0 {
// 		return "", fmt.Errorf("no response choices returned from OpenRouter")
// 	}
//
// 	return response.Choices[0].Message.Content, nil
// }
//
//
