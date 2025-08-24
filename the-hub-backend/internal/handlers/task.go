package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

// GetTasks godoc
// @Summary      Get all tasks
// @Description  Fetch tasks for the logged-in user with optional ordering
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_by  query     string  false  "Order by field (order, priority, due_date, created_at)"  default(order)
// @Param        sort      query     string  false  "Sort direction (asc, desc)"  default(asc)
// @Success      200  {object}  map[string][]models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks [get]
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get query parameters for ordering
	orderBy := c.DefaultQuery("order_by", "order_index")
	sortDir := c.DefaultQuery("sort", "asc")

	// Validate order_by parameter
	validOrderFields := map[string]bool{
		"order_index": true,
		"priority":    true,
		"due_date":    true,
		"created_at":  true,
		"title":       true,
		"status":      true,
	}

	if !validOrderFields[orderBy] {
		config.Logger.Warnf("Invalid order_by parameter: %s", orderBy)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order_by parameter"})
		return
	}

	// Validate sort direction
	if sortDir != "asc" && sortDir != "desc" {
		config.Logger.Warnf("Invalid sort direction: %s", sortDir)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort direction. Use 'asc' or 'desc'"})
		return
	}

	orderClause := orderBy + " " + sortDir

	config.Logger.Infof("Fetching tasks for user ID: %s with order: %s", userIDUUID, orderClause)
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Order(orderClause).Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Error fetching tasks for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	config.Logger.Infof("Found %d tasks for user ID %s", len(tasks), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTask godoc
// @Summary      Get a specific task
// @Description  Fetch a specific task by ID for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Task ID"
// @Success      200  {object}  map[string]models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{ID} [get]
func GetTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	config.Logger.Infof("Fetching task ID: %s for user ID: %s", taskID, userIDUUID)
	var task models.Task
	// Ensure user can only access their own tasks
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Errorf("Task ID %s not found for user %s: %v", taskID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved task ID %s for user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title                string     `json:"title" example:"Complete project proposal"`
	Description          string     `json:"description" example:"Finish the quarterly project proposal document"`
	Priority             *int       `json:"priority" example:"3"`
	DueDate              *time.Time `json:"due_date" example:"2024-12-31T23:59:59Z"`
	OrderIndex           *int       `json:"order" example:"1"`
	GoalID               *uuid.UUID `json:"goal_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	NaturalLanguageInput *string    `json:"natural_language_input" example:"Buy groceries tomorrow at 5pm, high priority"`
	UseNaturalLanguage   *bool      `json:"use_natural_language" example:"true"`
}

// AIResponse represents the expected JSON response from OpenAI
type AIResponse struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    *int    `json:"priority"`
	DueDate     *string `json:"due_date"`
}

// parseNaturalLanguage parses natural language input to extract task details
func parseNaturalLanguage(input string) (string, string, *int, *time.Time, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		config.Logger.Warn("OPENAI_API_KEY not set, using default parsing")
		return parseNaturalLanguageFallback(input)
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	prompt := `Parse the following natural language task description and extract structured information.

Input: "` + input + `"

Return a JSON object with these exact fields:
{
  "title": "A concise, actionable title for the task (required)",
  "description": "Additional details or context (can be empty string)",
  "priority": 3,
  "due_date": null
}

Priority levels:
- 1 = Low/Trivial/Minor
- 2 = Low
- 3 = Medium/Normal (default if not specified)
- 4 = High/Important
- 5 = Urgent/Critical/ASAP

Date formats: Use ISO 8601 format (e.g., "2024-12-25T14:30:00Z") or null if no date mentioned.

Examples:
Input: "Buy groceries tomorrow urgent"
Output: {"title": "Buy groceries", "description": "", "priority": 5, "due_date": "2024-12-26T12:00:00Z"}

Input: "Finish report by Friday"
Output: {"title": "Finish report", "description": "", "priority": 3, "due_date": "2024-12-27T17:00:00Z"}

Input: "Call mom low priority"
Output: {"title": "Call mom", "description": "", "priority": 2, "due_date": null}`

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a task parsing assistant. Extract structured information from natural language task descriptions.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens: 300,
	})

	if err != nil {
		// Check for specific error types that indicate credit/usage issues
		errStr := err.Error()
		if strings.Contains(errStr, "insufficient_quota") ||
			strings.Contains(errStr, "billing") ||
			strings.Contains(errStr, "credit") ||
			strings.Contains(errStr, "quota") {
			config.Logger.Warnf("OpenAI API credit/usage error: %v - falling back to enhanced parsing", err)
		} else {
			config.Logger.Errorf("OpenAI API error: %v - falling back to enhanced parsing", err)
		}
		return parseNaturalLanguageFallback(input)
	}

	// Check if we got a valid response
	if len(resp.Choices) == 0 {
		config.Logger.Warn("OpenAI API returned no choices, using fallback parsing")
		return parseNaturalLanguageFallback(input)
	}

	// Check if the response has the expected structure
	if resp.Choices[0].FinishReason == "length" {
		config.Logger.Warn("OpenAI API response was truncated, using fallback parsing")
		return parseNaturalLanguageFallback(input)
	}

	// Parse the JSON response
	responseContent := strings.TrimSpace(resp.Choices[0].Message.Content)
	if responseContent == "" {
		config.Logger.Warn("OpenAI API returned empty response, using fallback parsing")
		return parseNaturalLanguageFallback(input)
	}

	// Try to parse the JSON response from AI
	var aiResponse AIResponse
	if err := json.Unmarshal([]byte(responseContent), &aiResponse); err != nil {
		config.Logger.Warnf("Failed to parse AI JSON response: %v, using fallback parsing", err)
		return parseNaturalLanguageFallback(input)
	}

	// Validate and convert AI response
	title := strings.TrimSpace(aiResponse.Title)
	if title == "" {
		config.Logger.Warn("AI returned empty title, using fallback parsing")
		return parseNaturalLanguageFallback(input)
	}

	// Convert priority if provided
	var priority *int
	if aiResponse.Priority != nil {
		if *aiResponse.Priority >= 1 && *aiResponse.Priority <= 5 {
			priority = aiResponse.Priority
		} else {
			config.Logger.Warnf("AI returned invalid priority %d, using fallback", *aiResponse.Priority)
			return parseNaturalLanguageFallback(input)
		}
	}

	// Parse due date if provided
	var dueDate *time.Time
	if aiResponse.DueDate != nil && *aiResponse.DueDate != "" && *aiResponse.DueDate != "null" {
		if parsedTime, err := time.Parse(time.RFC3339, *aiResponse.DueDate); err == nil {
			dueDate = &parsedTime
		} else if parsedTime, err := time.Parse("2006-01-02", *aiResponse.DueDate); err == nil {
			dueDate = &parsedTime
		} else {
			config.Logger.Warnf("AI returned unparseable date %s, ignoring", *aiResponse.DueDate)
		}
	}

	description := strings.TrimSpace(aiResponse.Description)

	config.Logger.Infof("Successfully parsed AI response - Title: '%s', Priority: %v, DueDate: %v", title, priority, dueDate)
	return title, description, priority, dueDate, nil
}

// parseNaturalLanguageFallback provides enhanced parsing when AI is not available
func parseNaturalLanguageFallback(input string) (string, string, *int, *time.Time, error) {
	// Normalize input: trim spaces and convert to lowercase for processing
	normalizedInput := strings.TrimSpace(strings.ToLower(input))
	originalInput := strings.TrimSpace(input)

	// Initialize default values
	priority := 3
	var dueDate *time.Time
	now := time.Now()

	// Priority keywords and their values
	priorityMap := map[string]int{
		"urgent":    5,
		"asap":      5,
		"critical":  5,
		"high":      4,
		"important": 4,
		"medium":    3,
		"normal":    3,
		"low":       2,
		"minor":     1,
		"trivial":   1,
	}

	// Date patterns with actual parsing logic
	datePatterns := map[string]func(time.Time) time.Time{
		"tomorrow":   func(t time.Time) time.Time { return t.AddDate(0, 0, 1) },
		"today":      func(t time.Time) time.Time { return t },
		"next week":  func(t time.Time) time.Time { return t.AddDate(0, 0, 7) },
		"next month": func(t time.Time) time.Time { return t.AddDate(0, 1, 0) },
		"end of week": func(t time.Time) time.Time {
			daysUntilSaturday := (6 - int(t.Weekday())) % 7
			return t.AddDate(0, 0, daysUntilSaturday)
		},
		"end of month": func(t time.Time) time.Time {
			return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 0, t.Location())
		},
	}

	// Day-specific patterns
	dayPatterns := map[string]int{
		"monday": 1, "tuesday": 2, "wednesday": 3, "thursday": 4,
		"friday": 5, "saturday": 6, "sunday": 0,
	}

	// Extract priority from input
	workingText := normalizedInput
	for keyword, pri := range priorityMap {
		if strings.Contains(workingText, keyword) {
			priority = pri
			break
		}
	}

	// Extract dates from input
	for pattern, dateFunc := range datePatterns {
		if strings.Contains(normalizedInput, pattern) {
			parsedDate := dateFunc(now)
			dueDate = &parsedDate
			break
		}
	}

	// Handle day-specific patterns (next Monday, this Friday)
	if dueDate == nil {
		for dayName, dayNum := range dayPatterns {
			if strings.Contains(normalizedInput, "next "+dayName) {
				daysUntil := (dayNum - int(now.Weekday()) + 7) % 7
				if daysUntil == 0 {
					daysUntil = 7
				} // Next week if today
				targetDate := now.AddDate(0, 0, daysUntil)
				dueDate = &targetDate
				break
			}
			if strings.Contains(normalizedInput, "this "+dayName) {
				daysUntil := (dayNum - int(now.Weekday()) + 7) % 7
				targetDate := now.AddDate(0, 0, daysUntil)
				dueDate = &targetDate
				break
			}
		}
	}

	// Handle time-specific patterns (at 3pm, at 14:30)
	if strings.Contains(normalizedInput, " at ") {
		timeRegex := regexp.MustCompile(`at (\d{1,2})(?::(\d{2}))?(am|pm)?`)
		matches := timeRegex.FindStringSubmatch(normalizedInput)
		if len(matches) > 0 {
			hour, _ := strconv.Atoi(matches[1])
			minute := 0
			if matches[2] != "" {
				minute, _ = strconv.Atoi(matches[2])
			}

			// Handle AM/PM
			if strings.ToLower(matches[3]) == "pm" && hour != 12 {
				hour += 12
			} else if strings.ToLower(matches[3]) == "am" && hour == 12 {
				hour = 0
			}

			// Set time on the due date
			if dueDate != nil {
				*dueDate = time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(),
					hour, minute, 0, 0, dueDate.Location())
			} else {
				// If no date but time specified, use today
				todayWithTime := time.Date(now.Year(), now.Month(), now.Day(),
					hour, minute, 0, 0, now.Location())
				dueDate = &todayWithTime
			}
		}
	}

	// Build title by removing date and priority keywords
	title := originalInput

	// Priority patterns to remove from title
	priorityPatterns := []string{
		"urgent", "asap", "critical", "high priority", "high",
		"important", "medium priority", "medium", "normal",
		"low priority", "low", "minor", "trivial",
	}

	// Date patterns to remove from title
	datePatternsToRemove := []string{
		"tomorrow", "today", "next week", "next month",
		"end of week", "end of month", "by \\w+", "at \\d{1,2}(?::\\d{2})?(am|pm)?",
		"on \\w+", "this \\w+", "in \\d+ days?", "in \\d+ weeks?", "in \\d+ months?",
	}

	// Remove priority keywords
	for _, pattern := range priorityPatterns {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(pattern) + `\b`)
		title = re.ReplaceAllString(title, "")
	}

	// Remove date keywords
	for _, pattern := range datePatternsToRemove {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(pattern) + `\b`)
		title = re.ReplaceAllString(title, "")
	}

	// Clean up the title: remove extra spaces, punctuation at the end
	title = strings.TrimSpace(title)
	title = strings.TrimRight(title, ".,!?;:")

	// If title is empty after cleaning, use original input
	if title == "" {
		title = originalInput
	}

	// Extract description if there's additional context after the main task
	description := ""

	config.Logger.Infof("Parsed task - Title: '%s', Priority: %d, DueDate: %v", title, priority, dueDate)
	return title, description, &priority, dueDate, nil
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        task  body      CreateTaskRequest  true  "Task creation data"
// @Success      201   {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks [post]
func CreateTask(c *gin.Context) {
	var input CreateTaskRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid task input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for task", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Handle natural language input
	if input.UseNaturalLanguage != nil && *input.UseNaturalLanguage && input.NaturalLanguageInput != nil {
		parsedTitle, parsedDescription, parsedPriority, parsedDueDate, err := parseNaturalLanguage(*input.NaturalLanguageInput)
		if err != nil {
			config.Logger.Errorf("Failed to parse natural language input: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse natural language input"})
			return
		}

		// Override the input fields with parsed values
		input.Title = parsedTitle
		input.Description = parsedDescription
		input.Priority = parsedPriority
		input.DueDate = parsedDueDate
	}

	// If no order is specified, set it to the next available position
	order := 0
	if input.OrderIndex != nil {
		order = *input.OrderIndex
	} else {
		// Get the highest order number and add 1
		var maxOrder int
		if err := config.GetDB().Model(&models.Task{}).Where("user_id = ?", userIDUUID).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder).Error; err != nil {
			config.Logger.Warnf("Failed to get max order for user %s: %v", userIDUUID, err)
		}
		order = maxOrder + 1
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		OrderIndex:  order,
		GoalID:      input.GoalID,
		UserID:      userIDUUID,
	}

	config.Logger.Infof("Creating task for user %s: %s with order %d", userIDUUID, input.Title, order)
	if err := config.GetDB().Create(&task).Error; err != nil {
		config.Logger.Errorf("Error creating task for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	// Handle scheduled task if due date is provided
	if task.DueDate != nil {
		if err := UpsertScheduledTask(task); err != nil {
			config.Logger.Warnf("Failed to create scheduled task for task ID %s: %v", task.ID, err)
			// Don't return error as the main task was created successfully
		}
	}

	config.Logger.Infof("Successfully created task ID %s for user %s", task.ID, userIDUUID)
	c.JSON(http.StatusCreated, task)
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title       *string    `json:"title" example:"Updated task title"`
	Description *string    `json:"description" example:"Updated task description"`
	Priority    *int       `json:"priority" example:"2"`
	Status      *string    `json:"status" example:"completed"`
	DueDate     *time.Time `json:"due_date" example:"2024-12-31T23:59:59Z"`
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Update a specific task by ID for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID    path      int                true  "Task ID"
// @Param        task  body      UpdateTaskRequest  true  "Task update data"
// @Success      200   {object}  models.Task
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tasks/{ID} [patch]
func UpdateTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param for update: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var task models.Task
	// Ensure user can only update their own tasks
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task not found for update: ID %s, User %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input UpdateTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for task ID %d: %v", taskID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.DueDate != nil {
		updates["due_date"] = *input.DueDate
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for task update: ID %d", taskID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating task ID %d for user %v with data: %+v", taskID, userID, updates)
	if err := config.GetDB().Model(&task).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update task ID %d: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Reload the updated task
	if err := config.GetDB().First(&task, task.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated task ID %d: %v", task.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	config.Logger.Infof("Successfully updated task ID %d for user %v", task.ID, userID)
	c.JSON(http.StatusOK, task)
}

// ReorderTasksRequest represents the request body for reordering tasks
type ReorderTasksRequest struct {
	TaskOrders []TaskOrderItem `json:"task_orders" binding:"required"`
}

type TaskOrderItem struct {
	TaskID uuid.UUID `json:"task_id" binding:"required"`
	Order  int       `json:"order" binding:"required"`
}

// ReorderTasks godoc
// @Summary      Reorder multiple tasks
// @Description  Update the order of multiple tasks at once
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        reorder  body      ReorderTasksRequest  true  "Task reorder data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /tasks/reorder [put]
func ReorderTasks(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task reorder")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input ReorderTasksRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid reorder input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if len(input.TaskOrders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tasks to reorder"})
		return
	}

	// Start transaction
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	config.Logger.Infof("Reordering %d tasks for user %v", len(input.TaskOrders), userID)

	// Update each task's order
	for _, item := range input.TaskOrders {
		var task models.Task
		// Ensure user owns the task
		if err := tx.Where("id = ? AND user_id = ?", item.TaskID, userID).First(&task).Error; err != nil {
			tx.Rollback()
			config.Logger.Warnf("Task ID %d not found for user %v during reorder", item.TaskID, userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "One or more tasks not found"})
			return
		}

		if err := tx.Model(&task).Update("order_index", item.Order).Error; err != nil {
			tx.Rollback()
			config.Logger.Errorf("Failed to update order for task ID %d: %v", item.TaskID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reorder tasks"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		config.Logger.Errorf("Failed to commit task reorder transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task order"})
		return
	}

	config.Logger.Infof("Successfully reordered tasks for user %v", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Tasks reordered successfully"})
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a specific task by ID for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      int  true  "Task ID"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{ID} [delete]
func DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param for delete: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var task models.Task
	// Ensure user can only delete their own tasks
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task not found for delete: ID %s, User %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	config.Logger.Infof("Deleting task ID %s for user %s", taskID, userIDUUID)
	if err := config.GetDB().Delete(&task).Error; err != nil {
		config.Logger.Errorf("Failed to delete task ID %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// Clean up scheduled task if it exists
	if err := config.GetDB().Where("task_id = ?", task.ID).Delete(&models.ScheduledTask{}).Error; err != nil {
		config.Logger.Warnf("Failed to delete scheduled task for task ID %s: %v", task.ID, err)
		// Don't return error as the main task was deleted successfully
	}

	config.Logger.Infof("Successfully deleted task ID %s for user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task": task})
}

// GetRecommendedTasks godoc
// @Summary      Get recommended tasks for the logged-in user
// @Description  Fetch tasks recommended for the logged-in user based on their goals, schedule, and learning topics
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.Task
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/recommendations [get]
func GetRecommendedTasks(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Placeholder implementation - replace with AI logic
	var recommendedTasks []models.Task

	config.Logger.Infof("Fetching recommended tasks for user ID: %s", userIDUUID)

	c.JSON(http.StatusOK, gin.H{"tasks": recommendedTasks})
}

// UpsertScheduledTask creates or updates a scheduled task based on the task's due date
func UpsertScheduledTask(task models.Task) error {
	db := config.GetDB()

	if task.DueDate == nil {
		config.Logger.Infof("Removing scheduled task for task ID %d (no due date)", task.ID)
		return db.Where("task_id = ?", task.ID).Delete(&models.ScheduledTask{}).Error
	}

	start := *task.DueDate
	end := start.Add(time.Hour)

	var scheduled models.ScheduledTask
	err := db.Where("task_id = ?", task.ID).First(&scheduled).Error

	if err != nil {
		// Create new scheduled task
		config.Logger.Infof("Creating new scheduled task for task ID %d", task.ID)
		scheduled = models.ScheduledTask{
			Title:  task.Title,
			Start:  start,
			End:    end,
			UserID: task.UserID,
		}
		if createErr := db.Create(&scheduled).Error; createErr != nil {
			config.Logger.Errorf("Failed to create scheduled task for task ID %d: %v", task.ID, createErr)
			return createErr
		}
		config.Logger.Infof("Successfully created scheduled task for task ID %d", task.ID)
		return nil
	}

	// Update existing scheduled task
	config.Logger.Infof("Updating existing scheduled task for task ID %d", task.ID)
	updateData := models.ScheduledTask{
		Title: task.Title,
		Start: start,
		End:   end,
	}
	if updateErr := db.Model(&scheduled).Updates(updateData).Error; updateErr != nil {
		config.Logger.Errorf("Failed to update scheduled task for task ID %d: %v", task.ID, updateErr)
		return updateErr
	}

	config.Logger.Infof("Successfully updated scheduled task for task ID %d", task.ID)
	return nil
}
