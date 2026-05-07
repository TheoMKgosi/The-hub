package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
)

var logger = config.Logger

type OpenRouterClient struct {
	httpClient *http.Client
	apiKey    string
	baseURL   string
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ContentBlock struct {
	Type     string `json:"type,omitempty"`
	Text     string `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
	Document *Document `json:"document,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Document struct {
	URL         string `json:"url,omitempty"`
	Base64Data string `json:"base64_data,omitempty"`
	MimeType   string `json:"mime_type"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []Message    `json:"messages"`
	Temperature float32      `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream     bool         `json:"stream,omitempty"`
}

type ChatResponse struct {
	ID        string   `json:"id"`
	Choices  []Choice `json:"choices"`
	Usage    Usage    `json:"usage"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message     Message  `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens     int `json:"total_tokens"`
}

type Options struct {
	Model        string
	Temperature  float32
	MaxTokens    int
}

const defaultModel = "qwen/qwen3.6-flash"
const defaultTemperature = 0.7
const defaultMaxTokens = 4096
const defaultBaseURL = "https://openrouter.ai/api/v1"

func NewOpenRouterClient() (*OpenRouterClient, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &OpenRouterClient{
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: baseURL,
	}, nil
}

func (c *OpenRouterClient) SendMessage(messages []Message, opts Options) (string, error) {
	model := opts.Model
	if model == "" {
		model = defaultModel
	}

	temperature := opts.Temperature
	if temperature == 0 {
		temperature = defaultTemperature
	}

	maxTokens := opts.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	reqBody := ChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Stream:     false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://projectlifeledger.com/")
	req.Header.Set("X-Title", "Project Life Ledger")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("API error: status %d, response: %v", resp.StatusCode, errResp)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return chatResp.Choices[0].Message.Content.(string), nil
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
	TaskID          string              `json:"task_id"`
	OriginalTitle   string              `json:"original_title"`
	Title          string              `json:"title"`
	Description   string              `json:"description"`
	Priority       int                 `json:"priority"`
	EstimatedHours int                 `json:"estimated_hours"`
	Subtasks      []SubtaskSuggestion   `json:"subtasks"`
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

func (c *OpenRouterClient) GenerateFlashcards(pdfBase64 string, numCards int) (string, error) {
	systemPrompt := `You are a learning assistant. Generate flashcards from provided content.`
	userContent := fmt.Sprintf(`Generate %d flashcards from the following PDF content (base64 encoded). 

Respond with a JSON array of objects, each containing:
- "front": question or term
- "back": answer or definition
- "category": topic category (optional)

PDF Content (base64):
%s`, numCards, pdfBase64)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.5,
		MaxTokens:   8192,
	})
}

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

func (c *OpenRouterClient) GenerateWithDocument(pdfBase64, mimeType, prompt string) (string, error) {
	doc := &Document{
		Base64Data: pdfBase64,
		MimeType:   mimeType,
	}

	content := []ContentBlock{
		{Type: "document", Document: doc},
		{Type: "text", Text: prompt},
	}

	messages := []Message{
		{Role: "user", Content: content},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.5,
		MaxTokens:   8192,
	})
}

func (c *OpenRouterClient) GenerateWithImage(imageBase64, prompt string) (string, error) {
	img := &ImageURL{
		URL: "data:image/jpeg;base64," + imageBase64,
	}

	content := []ContentBlock{
		{Type: "image_url", ImageURL: img},
		{Type: "text", Text: prompt},
	}

	messages := []Message{
		{Role: "user", Content: content},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.7,
		MaxTokens:   4096,
	})
}

func EncodeFileToBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

var aiClient *OpenRouterClient

func GetOpenRouterClient() (*OpenRouterClient, error) {
	if aiClient != nil {
		return aiClient, nil
	}

	var err error
	aiClient, err = NewOpenRouterClient()
	if err != nil {
		config.Logger.Warnw("Failed to initialize OpenRouter client", "error", err.Error())
		return nil, err
	}

	return aiClient, nil
}

func InitAI() {
	client, err := NewOpenRouterClient()
	if err != nil {
		logger.Warnw("OpenRouter client not initialized on startup", "error", err.Error())
		return
	}
	aiClient = client
}
