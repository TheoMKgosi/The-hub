package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// OpenRouterClient handles communication with OpenRouter API
type OpenRouterClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// OpenRouterRequest represents a request to OpenRouter
type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response from OpenRouter
type OpenRouterResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice represents a completion choice
type Choice struct {
	Message Message `json:"message"`
}

// APIError represents an API error
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// NewOpenRouterClient creates a new OpenRouter client
func NewOpenRouterClient() (*OpenRouterClient, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}

	return &OpenRouterClient{
		apiKey:  apiKey,
		baseURL: "https://openrouter.ai/api/v1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// MakeRequest sends a request to OpenRouter with enhanced error handling
func (c *OpenRouterClient) MakeRequest(model string, messages []Message) (*OpenRouterResponse, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages cannot be empty")
	}

	req := OpenRouterRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://the-hub-app.com")
	httpReq.Header.Set("X-Title", "The Hub AI Assistant")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle different HTTP status codes
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("OpenRouter API authentication failed - check your API key")
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("OpenRouter API rate limit exceeded - please try again later")
		}
		if resp.StatusCode == http.StatusBadRequest {
			return nil, fmt.Errorf("OpenRouter API bad request: %s", string(body))
		}
		return nil, fmt.Errorf("OpenRouter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response OpenRouterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("OpenRouter API error: %s (Type: %s)", response.Error.Message, response.Error.Type)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("OpenRouter API returned no choices in response")
	}

	return &response, nil
}

// ParseNaturalLanguage uses OpenRouter to parse natural language task input
func (c *OpenRouterClient) ParseNaturalLanguage(input string) (string, string, *int, *time.Time, error) {
	prompt := fmt.Sprintf(`Parse the following natural language task description and extract:
1. Task title
2. Task description
3. Priority level (1-5, where 5 is highest)
4. Due date (if mentioned)

Input: "%s"

Return the result as a JSON object with keys: title, description, priority, due_date.
For due_date, use ISO format if a date is mentioned, otherwise null.
For priority, use 3 as default if not specified.`, input)

	messages := []Message{
		{
			Role:    "system",
			Content: "You are a task parsing assistant. Extract structured information from natural language task descriptions.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.MakeRequest("anthropic/claude-3-haiku", messages)
	if err != nil {
		return "", "", nil, nil, err
	}

	if len(response.Choices) == 0 {
		return "", "", nil, nil, fmt.Errorf("no response choices returned from OpenRouter")
	}

	// Parse the JSON response
	var result struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Priority    *int    `json:"priority"`
		DueDate     *string `json:"due_date"`
	}

	content := response.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to parse OpenRouter response: %w", err)
	}

	var dueDate *time.Time
	if result.DueDate != nil {
		if parsed, err := time.Parse(time.RFC3339, *result.DueDate); err == nil {
			dueDate = &parsed
		}
	}

	return result.Title, result.Description, result.Priority, dueDate, nil
}

// Helper functions
func formatTasks(tasks []string) string {
	if len(tasks) == 0 {
		return "No tasks provided"
	}
	result := ""
	for i, task := range tasks {
		result += fmt.Sprintf("%d. %s\n", i+1, task)
	}
	return result
}

// GenerateGoalTaskRecommendations generates AI-powered task recommendations for a goal
func (c *OpenRouterClient) GenerateGoalTaskRecommendations(goalTitle, goalDescription string, existingTasks []string) (string, error) {
	prompt := fmt.Sprintf(`You are an AI assistant helping users break down their goals into actionable tasks.

Goal Title: %s
Goal Description: %s

Existing Tasks:
%s

Please suggest 3-5 specific, actionable tasks that would help achieve this goal. For each task, provide:
1. A clear, concise title
2. A brief description of what the task involves
3. A priority level (1-5, where 5 is highest priority)
4. An estimated time in hours
5. A brief explanation of why this task is important for achieving the goal

Return the response as a JSON array of objects with the following structure:
[
  {
    "title": "Task Title",
    "description": "Task description",
    "priority": 3,
    "estimated_hours": 2,
    "reasoning": "Why this task is important"
  }
]

Make sure the tasks are:
- Specific and actionable
- Different from existing tasks
- Realistic in scope
- Ordered by priority (highest first)
- Focused on achieving the goal`, goalTitle, goalDescription, formatTasks(existingTasks))

	messages := []Message{
		{
			Role:    "system",
			Content: "You are a goal achievement assistant that helps break down goals into actionable tasks. Always respond with valid JSON.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.MakeRequest("anthropic/claude-3-haiku", messages)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from OpenRouter")
	}

	return response.Choices[0].Message.Content, nil
}

func formatEvents(events []string) string {
	if len(events) == 0 {
		return "No existing events"
	}
	result := ""
	for i, event := range events {
		result += fmt.Sprintf("%d. %s\n", i+1, event)
	}
	return result
}
