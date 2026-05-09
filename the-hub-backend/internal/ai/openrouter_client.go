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
	apiKey     string
	baseURL    string
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ContentBlock struct {
	Type     string    `json:"type,omitempty"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
	Document *Document `json:"document,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Document struct {
	URL        string `json:"url,omitempty"`
	Base64Data string `json:"base64_data,omitempty"`
	MimeType   string `json:"mime_type"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Options struct {
	Model       string
	Temperature float32
	MaxTokens   int
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

// TODO: Use a proper logging library instead of fmt.Println
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
		Stream:      false,
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

func (c *OpenRouterClient) GenerateWithDocument(pdfBase64, mimeType, prompt string, systemPrompt string) (string, error) {
	doc := &Document{
		Base64Data: pdfBase64,
		MimeType:   mimeType,
	}

	content := []ContentBlock{
		{Type: "document", Document: doc},
		{Type: "text", Text: prompt},
	}

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: content},
	}

	return c.SendMessage(messages, Options{
		Temperature: 0.3,
		MaxTokens:   8192,
		Model:       "google/gemini-2.5-pro",
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
