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
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

// GetTask godoc
// @Summary      Get a specific task
// @Description  Fetch a specific task by ID for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ID   path      string  true  "Task ID"
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
	if err := config.GetDB().Preload("Subtasks").Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Errorf("Task ID %s not found for user %s: %v", taskID, userIDUUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	config.Logger.Infof("Successfully retrieved task ID %s for user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"task": task})
}

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

	// Get filtering parameters
	status := c.Query("status")
	priority := c.Query("priority")
	goalID := c.Query("goal_id")
	search := c.Query("search")
	dueBefore := c.Query("due_before")
	dueAfter := c.Query("due_after")

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

	// Build query
	query := config.GetDB().Where("user_id = ? AND parent_task_id IS NULL", userIDUUID)

	// Apply filters
	if status != "" {
		validStatuses := map[string]bool{
			"pending":     true,
			"completed":   true,
			"in_progress": true,
		}
		if !validStatuses[status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status filter"})
			return
		}
		query = query.Where("status = ?", status)
	}

	if priority != "" {
		pri, err := strconv.Atoi(priority)
		if err != nil || pri < 1 || pri > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority filter"})
			return
		}
		query = query.Where("priority = ?", pri)
	}

	if goalID != "" {
		goalUUID, err := uuid.Parse(goalID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal_id filter"})
			return
		}
		query = query.Where("goal_id = ?", goalUUID)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if dueBefore != "" {
		beforeDate, err := time.Parse("2006-01-02", dueBefore)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_before format. Use YYYY-MM-DD"})
			return
		}
		query = query.Where("due_date <= ?", beforeDate)
	}

	if dueAfter != "" {
		afterDate, err := time.Parse("2006-01-02", dueAfter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_after format. Use YYYY-MM-DD"})
			return
		}
		query = query.Where("due_date >= ?", afterDate)
	}

	orderClause := orderBy + " " + sortDir

	config.Logger.Infof("Fetching tasks for user ID: %s with filters - status: %s, priority: %s, goal: %s, search: %s, order: %s",
		userIDUUID, status, priority, goalID, search, orderClause)

	if err := query.Order(orderClause).Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Error fetching tasks for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch tasks"})
		return
	}

	config.Logger.Infof("Found %d tasks for user ID %s", len(tasks), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title                string     `json:"title" example:"Complete project proposal"`
	Description          string     `json:"description" example:"Finish the quarterly project proposal document"`
	Priority             *int       `json:"priority" example:"3"`
	DueDate              *time.Time `json:"due_date" example:"2024-12-31T23:59:59Z"`
	OrderIndex           *int       `json:"order" example:"1"`
	GoalID               *uuid.UUID `json:"goal_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ParentTaskID         *uuid.UUID `json:"parent_task_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	TimeEstimate         *int       `json:"time_estimate_minutes" example:"60"`
	TemplateID           *uuid.UUID `json:"template_id" example:"550e8400-e29b-41d4-a716-446655440002"`
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
		Title:        input.Title,
		Description:  input.Description,
		Priority:     input.Priority,
		DueDate:      input.DueDate,
		OrderIndex:   order,
		GoalID:       input.GoalID,
		ParentTaskID: input.ParentTaskID,
		TimeEstimate: input.TimeEstimate,
		TemplateID:   input.TemplateID,
		UserID:       userIDUUID,
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

	// Recalculate goal progress if task is linked to a goal
	if task.GoalID != nil {
		var goal models.Goal
		if err := config.GetDB().Where("id = ? AND user_id = ?", *task.GoalID, userIDUUID).First(&goal).Error; err == nil {
			if err := goal.CalculateProgress(config.GetDB()); err != nil {
				config.Logger.Warnf("Failed to recalculate progress for goal %s: %v", goal.ID, err)
			} else {
				// Update goal in database
				config.GetDB().Model(&goal).Updates(map[string]interface{}{
					"progress":        goal.Progress,
					"total_tasks":     goal.TotalTasks,
					"completed_tasks": goal.CompletedTasks,
					"status":          goal.Status,
				})
			}
		}
	}

	// Send reminder notification if task has a due date
	if task.DueDate != nil && time.Until(*task.DueDate) > 0 {
		pushService := util.NewPushNotificationService(config.GetDB())
		go func() {
			// Send reminder 1 day before due date if it's more than 1 day away
			if time.Until(*task.DueDate) > 24*time.Hour {
				time.Sleep(time.Until(task.DueDate.Add(-24 * time.Hour)))
				pushService.SendTaskReminder(task.ID, userIDUUID, task.Title, task.DueDate)
			}
		}()
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
		config.Logger.Errorf("Error retrieving updated task ID %s: %v", task.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated task"})
		return
	}

	// Update parent task status if this is a subtask
	if task.ParentTaskID != nil {
		if err := task.UpdateParentStatus(config.GetDB()); err != nil {
			config.Logger.Warnf("Failed to update parent status for task %s: %v", task.ID, err)
		}
	}

	// Recalculate goal progress if task is linked to a goal
	if task.GoalID != nil {
		var goal models.Goal
		if err := config.GetDB().Where("id = ? AND user_id = ?", *task.GoalID, userIDUUID).First(&goal).Error; err == nil {
			if err := goal.CalculateProgress(config.GetDB()); err != nil {
				config.Logger.Warnf("Failed to recalculate progress for goal %s: %v", goal.ID, err)
			} else {
				// Update goal in database
				config.GetDB().Model(&goal).Updates(map[string]interface{}{
					"progress":        goal.Progress,
					"total_tasks":     goal.TotalTasks,
					"completed_tasks": goal.CompletedTasks,
					"status":          goal.Status,
				})
			}
		}
	}

	config.Logger.Infof("Successfully updated task ID %s for user %s", task.ID, userIDUUID)
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

	// Recalculate goal progress if task was linked to a goal
	if task.GoalID != nil {
		var goal models.Goal
		if err := config.GetDB().Where("id = ? AND user_id = ?", *task.GoalID, userIDUUID).First(&goal).Error; err == nil {
			if err := goal.CalculateProgress(config.GetDB()); err != nil {
				config.Logger.Warnf("Failed to recalculate progress for goal %s: %v", goal.ID, err)
			} else {
				// Update goal in database
				config.GetDB().Model(&goal).Updates(map[string]interface{}{
					"progress":        goal.Progress,
					"total_tasks":     goal.TotalTasks,
					"completed_tasks": goal.CompletedTasks,
					"status":          goal.Status,
				})
			}
		}
	}

	config.Logger.Infof("Successfully deleted task ID %s for user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully", "task": task})
}

// UndoDeleteTask godoc
// @Summary      Undo deletion of a task
// @Description  Restore a previously deleted task for the logged-in user
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
// @Router       /tasks/{ID}/undo-delete [patch]
func UndoDeleteTask(c *gin.Context) {
	taskIDStr := c.Param("ID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param for undo delete: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during task undo delete")
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
	// Find the soft deleted task
	if err := config.GetDB().Unscoped().Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Soft deleted task not found for undo: ID %s, User %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or not deleted"})
		return
	}

	// Check if task was deleted within the last 30 days (configurable undo window)
	if task.DeletedAt.Time.Before(time.Now().AddDate(0, 0, -30)) {
		config.Logger.Warnf("Task %s was deleted more than 30 days ago, cannot undo", taskID)
		c.JSON(http.StatusGone, gin.H{"error": "Task was deleted too long ago to undo"})
		return
	}

	config.Logger.Infof("Undoing deletion of task ID %s for user %s", taskID, userIDUUID)

	// Restore the task by clearing the deleted_at field
	if err := config.GetDB().Model(&task).Update("deleted_at", nil).Error; err != nil {
		config.Logger.Errorf("Failed to undo delete task ID %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore task"})
		return
	}

	// Recalculate goal progress if task was linked to a goal
	if task.GoalID != nil {
		var goal models.Goal
		if err := config.GetDB().Where("id = ? AND user_id = ?", *task.GoalID, userIDUUID).First(&goal).Error; err == nil {
			if err := goal.CalculateProgress(config.GetDB()); err != nil {
				config.Logger.Warnf("Failed to recalculate progress for goal %s: %v", goal.ID, err)
			} else {
				// Update goal in database
				config.GetDB().Model(&goal).Updates(map[string]interface{}{
					"progress":        goal.Progress,
					"total_tasks":     goal.TotalTasks,
					"completed_tasks": goal.CompletedTasks,
					"status":          goal.Status,
				})
			}
		}
	}

	config.Logger.Infof("Successfully restored task ID %s for user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Task restored successfully", "task": task})
}

// GetRecentlyDeletedTasks godoc
// @Summary      Get recently deleted tasks
// @Description  Fetch tasks that were deleted within the last 30 days for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.Task
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/recently-deleted [get]
func GetRecentlyDeletedTasks(c *gin.Context) {
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

	var deletedTasks []models.Task
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	if err := config.GetDB().Unscoped().Where("user_id = ? AND deleted_at IS NOT NULL AND deleted_at > ?", userIDUUID, thirtyDaysAgo).Order("deleted_at DESC").Find(&deletedTasks).Error; err != nil {
		config.Logger.Errorf("Failed to fetch recently deleted tasks for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deleted tasks"})
		return
	}

	config.Logger.Infof("Fetched %d recently deleted tasks for user %s", len(deletedTasks), userIDUUID)
	c.JSON(http.StatusOK, gin.H{"tasks": deletedTasks})
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

	config.Logger.Infof("Successfully updated scheduled task for task ID %s", task.ID)
	return nil
}

// CreateTaskDependencyRequest represents the request body for creating a task dependency
type CreateTaskDependencyRequest struct {
	DependsOnID uuid.UUID `json:"depends_on_id" binding:"required"` // The task that must be completed first
}

// CreateTaskDependency godoc
// @Summary      Create a task dependency
// @Description  Create a dependency relationship between two tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Param        dependency  body      CreateTaskDependencyRequest  true  "Dependency data"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/dependencies [post]
func CreateTaskDependency(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	var input CreateTaskDependencyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid dependency input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Verify both tasks exist and belong to the user
	var task, dependsOnTask models.Task
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := config.GetDB().Where("id = ? AND user_id = ?", input.DependsOnID, userIDUUID).First(&dependsOnTask).Error; err != nil {
		config.Logger.Warnf("Depends on task ID %s not found for user %s", input.DependsOnID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Depends on task not found"})
		return
	}

	// Prevent circular dependencies
	if taskID == input.DependsOnID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task cannot depend on itself"})
		return
	}

	// Check for existing dependency
	var existingDep models.TaskDependency
	if err := config.GetDB().Where("task_id = ? AND depends_on_id = ?", taskID, input.DependsOnID).First(&existingDep).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dependency already exists"})
		return
	}

	dependency := models.TaskDependency{
		TaskID:      taskID,
		DependsOnID: input.DependsOnID,
		UserID:      userIDUUID,
	}

	if err := config.GetDB().Create(&dependency).Error; err != nil {
		config.Logger.Errorf("Error creating task dependency: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create dependency"})
		return
	}

	config.Logger.Infof("Successfully created dependency: task %s depends on %s", taskID, input.DependsOnID)
	c.JSON(http.StatusCreated, gin.H{"message": "Dependency created successfully"})
}

// GetTaskDependencies godoc
// @Summary      Get task dependencies
// @Description  Get all dependencies for a specific task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Success      200  {object}  map[string][]models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/dependencies [get]
func GetTaskDependencies(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	// Verify task exists and belongs to user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	dependencies, err := task.GetDependencies(config.GetDB())
	if err != nil {
		config.Logger.Errorf("Error fetching dependencies for task %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch dependencies"})
		return
	}

	dependents, err := task.GetDependents(config.GetDB())
	if err != nil {
		config.Logger.Errorf("Error fetching dependents for task %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch dependents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dependencies": dependencies, // Tasks this task depends on
		"dependents":   dependents,   // Tasks that depend on this task
	})
}

// DeleteTaskDependency godoc
// @Summary      Delete a task dependency
// @Description  Delete a dependency relationship between two tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Param        dependsOnID   path      string  true  "Depends On Task ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/dependencies/{dependsOnID} [delete]
func DeleteTaskDependency(c *gin.Context) {
	taskIDStr := c.Param("taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid task ID param: %s", taskIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	dependsOnIDStr := c.Param("dependsOnID")
	dependsOnID, err := uuid.Parse(dependsOnIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid depends on task ID param: %s", dependsOnIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid depends on task ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Delete the dependency
	if err := config.GetDB().Where("task_id = ? AND depends_on_id = ? AND user_id = ?", taskID, dependsOnID, userIDUUID).Delete(&models.TaskDependency{}).Error; err != nil {
		config.Logger.Errorf("Error deleting task dependency: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete dependency"})
		return
	}

	config.Logger.Infof("Successfully deleted dependency: task %s no longer depends on %s", taskID, dependsOnID)
	c.JSON(http.StatusOK, gin.H{"message": "Dependency deleted successfully"})
}

// GetTaskSubtasks godoc
// @Summary      Get task subtasks
// @Description  Get all subtasks for a specific task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Success      200  {object}  map[string][]models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/subtasks [get]
func GetTaskSubtasks(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	// Verify task exists and belongs to user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	subtasks, err := task.GetSubtasks(config.GetDB())
	if err != nil {
		config.Logger.Errorf("Error fetching subtasks for task %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch subtasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subtasks": subtasks})
}

// StartTimeTrackingRequest represents the request body for starting time tracking
type StartTimeTrackingRequest struct {
	Description string `json:"description"`
}

// StartTimeTracking godoc
// @Summary      Start time tracking for a task
// @Description  Start tracking time for a specific task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Param        timeEntry  body      StartTimeTrackingRequest  true  "Time tracking data"
// @Success      200  {object}  models.TimeEntry
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/time/start [post]
func StartTimeTracking(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	// Verify task exists and belongs to user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input StartTimeTrackingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid time tracking input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	timeEntry, err := task.StartTimeTracking(config.GetDB(), userIDUUID, input.Description)
	if err != nil {
		config.Logger.Errorf("Error starting time tracking for task %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start time tracking"})
		return
	}

	config.Logger.Infof("Started time tracking for task %s by user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, timeEntry)
}

// StopTimeTracking godoc
// @Summary      Stop time tracking for a task
// @Description  Stop the currently running time tracking for a specific task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Success      200  {object}  models.TimeEntry
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/time/stop [post]
func StopTimeTracking(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	// Verify task exists and belongs to user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	timeEntry, err := task.StopTimeTracking(config.GetDB(), userIDUUID)
	if err != nil {
		config.Logger.Errorf("Error stopping time tracking for task %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not stop time tracking"})
		return
	}

	config.Logger.Infof("Stopped time tracking for task %s by user %s", taskID, userIDUUID)
	c.JSON(http.StatusOK, timeEntry)
}

// GetTaskTimeEntries godoc
// @Summary      Get time entries for a task
// @Description  Get all time tracking entries for a specific task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Success      200  {object}  map[string][]models.TimeEntry
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/time [get]
func GetTaskTimeEntries(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	// Verify task exists and belongs to user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var timeEntries []models.TimeEntry
	if err := config.GetDB().Where("task_id = ? AND user_id = ?", taskID, userIDUUID).Order("start_time DESC").Find(&timeEntries).Error; err != nil {
		config.Logger.Errorf("Error fetching time entries for task %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch time entries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"time_entries": timeEntries})
}

// CreateTaskTemplateRequest represents the request body for creating a task template
type CreateTaskTemplateRequest struct {
	Name                string `json:"name" binding:"required"`
	Description         string `json:"description"`
	Category            string `json:"category"`
	TitleTemplate       string `json:"title_template" binding:"required"`
	DescriptionTemplate string `json:"description_template"`
	Priority            *int   `json:"priority"`
	TimeEstimate        *int   `json:"time_estimate_minutes"`
	Tags                string `json:"tags"`
	IsPublic            *bool  `json:"is_public"`
}

// CreateTaskTemplate godoc
// @Summary      Create a task template
// @Description  Create a new task template for reuse
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        template  body      CreateTaskTemplateRequest  true  "Template data"
// @Success      201  {object}  models.TaskTemplate
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task-templates [post]
func CreateTaskTemplate(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var input CreateTaskTemplateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid template input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	isPublic := false
	if input.IsPublic != nil {
		isPublic = *input.IsPublic
	}

	template := models.TaskTemplate{
		UserID:              userIDUUID,
		Name:                input.Name,
		Description:         input.Description,
		Category:            input.Category,
		TitleTemplate:       input.TitleTemplate,
		DescriptionTemplate: input.DescriptionTemplate,
		Priority:            input.Priority,
		TimeEstimate:        input.TimeEstimate,
		Tags:                input.Tags,
		IsPublic:            isPublic,
	}

	if err := config.GetDB().Create(&template).Error; err != nil {
		config.Logger.Errorf("Error creating task template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create template"})
		return
	}

	config.Logger.Infof("Created task template %s for user %s", template.ID, userIDUUID)
	c.JSON(http.StatusCreated, template)
}

// GetTaskTemplates godoc
// @Summary      Get task templates
// @Description  Get all task templates for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.TaskTemplate
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task-templates [get]
func GetTaskTemplates(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var templates []models.TaskTemplate
	if err := config.GetDB().Where("user_id = ? OR is_public = ?", userIDUUID, true).Order("usage_count DESC, created_at DESC").Find(&templates).Error; err != nil {
		config.Logger.Errorf("Error fetching task templates for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// CreateTaskFromTemplate godoc
// @Summary      Create task from template
// @Description  Create a new task using a template
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        templateID   path      string  true  "Template ID"
// @Param        task  body      CreateTaskRequest  true  "Task data"
// @Success      201  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task-templates/{templateID}/create-task [post]
func CreateTaskFromTemplate(c *gin.Context) {
	templateIDStr := c.Param("templateID")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid template ID param: %s", templateIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify template exists and is accessible
	var template models.TaskTemplate
	if err := config.GetDB().Where("(id = ? AND user_id = ?) OR (id = ? AND is_public = ?)", templateID, userIDUUID, templateID, true).First(&template).Error; err != nil {
		config.Logger.Warnf("Template ID %s not found for user %s", templateID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	var input CreateTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid task input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Create task from template
	task := template.CreateFromTemplate(userIDUUID)

	// Override template values with input values if provided
	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Priority != nil {
		task.Priority = input.Priority
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.TimeEstimate != nil {
		task.TimeEstimate = input.TimeEstimate
	}
	if input.GoalID != nil {
		task.GoalID = input.GoalID
	}
	if input.ParentTaskID != nil {
		task.ParentTaskID = input.ParentTaskID
	}

	// Set order index
	order := 0
	if input.OrderIndex != nil {
		order = *input.OrderIndex
	} else {
		var maxOrder int
		if err := config.GetDB().Model(&models.Task{}).Where("user_id = ?", userIDUUID).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder).Error; err != nil {
			config.Logger.Warnf("Failed to get max order for user %s: %v", userIDUUID, err)
		}
		order = maxOrder + 1
	}
	task.OrderIndex = order

	if err := config.GetDB().Create(&task).Error; err != nil {
		config.Logger.Errorf("Error creating task from template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create task"})
		return
	}

	// Increment template usage count
	if err := config.GetDB().Model(&template).Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error; err != nil {
		config.Logger.Warnf("Failed to increment template usage count: %v", err)
	}

	config.Logger.Infof("Created task %s from template %s for user %s", task.ID, templateID, userIDUUID)
	c.JSON(http.StatusCreated, task)
}

// CreateRecurrenceRuleRequest represents the request body for creating a recurrence rule
type CreateRecurrenceRuleRequest struct {
	Name                string     `json:"name" binding:"required"`
	Description         string     `json:"description"`
	Frequency           string     `json:"frequency" binding:"required"` // daily, weekly, monthly, yearly
	Interval            *int       `json:"interval"`
	ByDay               string     `json:"by_day"`
	ByMonthDay          *int       `json:"by_month_day"`
	ByMonth             *int       `json:"by_month"`
	StartDate           *time.Time `json:"start_date"`
	EndDate             *time.Time `json:"end_date"`
	Count               *int       `json:"count"`
	TitleTemplate       string     `json:"title_template" binding:"required"`
	DescriptionTemplate string     `json:"description_template"`
	Priority            *int       `json:"priority"`
	TimeEstimate        *int       `json:"time_estimate_minutes"`
	DueDateOffset       *int       `json:"due_date_offset_days"`
}

// CreateTaskRecurrenceRule godoc
// @Summary      Create a recurrence rule
// @Description  Create a new recurrence rule for recurring tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        rule  body      CreateRecurrenceRuleRequest  true  "Recurrence rule data"
// @Success      201  {object}  models.RecurrenceRule
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /recurrence-rules [post]
func CreateTaskRecurrenceRule(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var input CreateRecurrenceRuleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid recurrence rule input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate frequency
	validFrequencies := map[string]bool{
		"daily":   true,
		"weekly":  true,
		"monthly": true,
		"yearly":  true,
	}
	if !validFrequencies[input.Frequency] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid frequency. Must be daily, weekly, monthly, or yearly"})
		return
	}

	interval := 1
	if input.Interval != nil && *input.Interval > 0 {
		interval = *input.Interval
	}

	rule := models.RecurrenceRule{
		UserID:              userIDUUID,
		Name:                input.Name,
		Description:         input.Description,
		Frequency:           input.Frequency,
		Interval:            interval,
		ByDay:               input.ByDay,
		ByMonthDay:          input.ByMonthDay,
		ByMonth:             input.ByMonth,
		StartDate:           input.StartDate,
		EndDate:             input.EndDate,
		Count:               input.Count,
		TitleTemplate:       input.TitleTemplate,
		DescriptionTemplate: input.DescriptionTemplate,
		Priority:            input.Priority,
		TimeEstimate:        input.TimeEstimate,
		DueDateOffset:       input.DueDateOffset,
	}

	if err := config.GetDB().Create(&rule).Error; err != nil {
		config.Logger.Errorf("Error creating recurrence rule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create recurrence rule"})
		return
	}

	config.Logger.Infof("Created recurrence rule %s for user %s", rule.ID, userIDUUID)
	c.JSON(http.StatusCreated, rule)
}

// GetRecurrenceRules godoc
// @Summary      Get recurrence rules
// @Description  Get all recurrence rules for the logged-in user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.RecurrenceRule
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /recurrence-rules [get]
func GetRecurrenceRules(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var rules []models.RecurrenceRule
	if err := config.GetDB().Where("user_id = ?", userIDUUID).Order("created_at DESC").Find(&rules).Error; err != nil {
		config.Logger.Errorf("Error fetching recurrence rules for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch recurrence rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recurrence_rules": rules})
}

// GenerateRecurringTasks godoc
// @Summary      Generate recurring tasks
// @Description  Generate task instances for a recurrence rule
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        ruleID   path      string  true  "Recurrence rule ID"
// @Param        count    query     int     false  "Number of tasks to generate (default: 1, max: 10)"
// @Success      200  {object}  map[string][]models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /recurrence-rules/{ruleID}/generate-tasks [post]
func GenerateRecurringTasks(c *gin.Context) {
	ruleIDStr := c.Param("ruleID")
	ruleID, err := uuid.Parse(ruleIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid rule ID param: %s", ruleIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	// Verify rule exists and belongs to user
	var rule models.RecurrenceRule
	if err := config.GetDB().Where("id = ? AND user_id = ?", ruleID, userIDUUID).First(&rule).Error; err != nil {
		config.Logger.Warnf("Recurrence rule ID %s not found for user %s", ruleID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Recurrence rule not found"})
		return
	}

	// Get count parameter
	countStr := c.DefaultQuery("count", "1")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		count = 1
	}
	if count > 10 {
		count = 10 // Limit to prevent abuse
	}

	// Generate occurrences
	fromDate := time.Now()
	occurrences := rule.GenerateOccurrences(fromDate, count)

	if len(occurrences) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid occurrences found for the recurrence rule"})
		return
	}

	var createdTasks []models.Task

	// Create tasks for each occurrence
	for _, occurrence := range occurrences {
		task := rule.CreateTaskFromRule(userIDUUID, occurrence)

		// Set order index
		var maxOrder int
		if err := config.GetDB().Model(&models.Task{}).Where("user_id = ?", userIDUUID).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder).Error; err != nil {
			config.Logger.Warnf("Failed to get max order for user %s: %v", userIDUUID, err)
		}
		task.OrderIndex = maxOrder + 1

		if err := config.GetDB().Create(&task).Error; err != nil {
			config.Logger.Errorf("Error creating recurring task: %v", err)
			continue // Continue with other tasks
		}

		createdTasks = append(createdTasks, *task)
	}

	config.Logger.Infof("Generated %d recurring tasks for rule %s", len(createdTasks), ruleID)
	c.JSON(http.StatusOK, gin.H{"tasks": createdTasks})
}

// GetTaskAnalytics godoc
// @Summary      Get task analytics for a user
// @Description  Get productivity analytics and insights for the logged-in user
// @Tags         analytics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period   query     string  false  "Period (daily, weekly, monthly)"  default(weekly)
// @Param        days     query     int     false  "Number of days to analyze"  default(30)
// @Success      200  {object}  models.UserProductivityInsights
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /analytics/tasks [get]
func GetTaskAnalytics(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	period := c.DefaultQuery("period", "weekly")
	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	insights, err := calculateProductivityInsights(config.GetDB(), userIDUUID, startDate, endDate, period)
	if err != nil {
		config.Logger.Errorf("Error calculating productivity insights for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not calculate analytics"})
		return
	}

	c.JSON(http.StatusOK, insights)
}

// GetTaskAnalyticsChart godoc
// @Summary      Get task analytics chart data
// @Description  Get chart data for productivity visualization
// @Tags         analytics
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        metric   query     string  true   "Metric to chart (completion_rate, productivity_score, time_spent)"
// @Param        days     query     int     false  "Number of days"  default(30)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /analytics/tasks/chart [get]
func GetTaskAnalyticsChart(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	metric := c.Query("metric")
	if metric == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric parameter is required"})
		return
	}

	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var analytics []models.TaskAnalytics
	if err := config.GetDB().Where("user_id = ? AND date BETWEEN ? AND ?", userIDUUID, startDate, endDate).
		Order("date").Find(&analytics).Error; err != nil {
		config.Logger.Errorf("Error fetching analytics data for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch analytics data"})
		return
	}

	// Format data for charts
	labels := make([]string, len(analytics))
	data := make([]float64, len(analytics))

	for i, analytic := range analytics {
		labels[i] = analytic.Date.Format("2006-01-02")

		switch metric {
		case "completion_rate":
			data[i] = analytic.CompletionRate
		case "productivity_score":
			data[i] = analytic.ProductivityScore
		case "time_spent":
			data[i] = float64(analytic.TotalTimeSpent)
		case "tasks_completed":
			data[i] = float64(analytic.CompletedTasks)
		default:
			data[i] = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": labels,
		"data":   data,
		"metric": metric,
	})
}

// calculateProductivityInsights calculates comprehensive productivity insights
func calculateProductivityInsights(db *gorm.DB, userID uuid.UUID, startDate, endDate time.Time, period string) (*models.UserProductivityInsights, error) {
	insights := &models.UserProductivityInsights{
		UserID:    userID,
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Get analytics data for the period
	var analytics []models.TaskAnalytics
	if err := db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&analytics).Error; err != nil {
		return nil, err
	}

	if len(analytics) == 0 {
		return insights, nil
	}

	// Calculate overall metrics
	totalTasks := 0
	totalCompleted := 0
	totalTimeSpent := 0
	totalTimeEstimated := 0
	completionRates := []float64{}
	productivityScores := []float64{}

	for _, analytic := range analytics {
		totalTasks += analytic.TotalTasks
		totalCompleted += analytic.CompletedTasks
		totalTimeSpent += analytic.TotalTimeSpent
		totalTimeEstimated += analytic.TotalTimeEstimated
		completionRates = append(completionRates, analytic.CompletionRate)
		productivityScores = append(productivityScores, analytic.ProductivityScore)
	}

	insights.TotalTasks = totalTasks
	insights.TotalCompleted = totalCompleted
	insights.TotalTimeSpent = totalTimeSpent

	if totalTasks > 0 {
		insights.AverageCompletionRate = average(completionRates)
		insights.AverageProductivityScore = average(productivityScores)
		insights.AverageTimePerTask = float64(totalTimeSpent) / float64(totalTasks)
	}

	// Calculate trends (compare first half vs second half)
	midPoint := len(analytics) / 2
	if midPoint > 0 {
		firstHalf := analytics[:midPoint]
		secondHalf := analytics[midPoint:]

		firstHalfAvg := average(extractField(firstHalf, "completion_rate"))
		secondHalfAvg := average(extractField(secondHalf, "completion_rate"))
		insights.CompletionRateTrend = secondHalfAvg - firstHalfAvg

		firstHalfProd := average(extractField(firstHalf, "productivity_score"))
		secondHalfProd := average(extractField(secondHalf, "productivity_score"))
		insights.ProductivityTrend = secondHalfProd - firstHalfProd
	}

	// Find best day
	bestScore := 0.0
	var bestDay *time.Time
	for _, analytic := range analytics {
		if analytic.ProductivityScore > bestScore {
			bestScore = analytic.ProductivityScore
			bestDay = &analytic.Date
		}
	}
	insights.BestDay = bestDay
	insights.BestDayScore = bestScore

	// Generate recommendations
	insights.Recommendations = generateRecommendations(insights)

	return insights, nil
}

// Helper functions
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func extractField(analytics []models.TaskAnalytics, field string) []float64 {
	values := make([]float64, len(analytics))
	for i, analytic := range analytics {
		switch field {
		case "completion_rate":
			values[i] = analytic.CompletionRate
		case "productivity_score":
			values[i] = analytic.ProductivityScore
		default:
			values[i] = 0
		}
	}
	return values
}

func generateRecommendations(insights *models.UserProductivityInsights) []string {
	recommendations := []string{}

	if insights.AverageCompletionRate < 0.7 {
		recommendations = append(recommendations, "Try breaking down large tasks into smaller, manageable subtasks")
	}

	if insights.ProductivityTrend < 0 {
		recommendations = append(recommendations, "Your productivity has been declining. Consider reviewing your task prioritization strategy")
	}

	if insights.AverageTimePerTask > 120 { // More than 2 hours per task
		recommendations = append(recommendations, "Tasks are taking longer than expected. Consider setting more realistic time estimates")
	}

	if insights.TotalTasks > 0 && insights.TotalCompleted == 0 {
		recommendations = append(recommendations, "You haven't completed any tasks recently. Focus on finishing smaller tasks first")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Great job! Keep up the good work with your current productivity habits")
	}

	return recommendations
}

// ShareTaskRequest represents the request body for sharing a task
type ShareTaskRequest struct {
	SharedWithID uuid.UUID `json:"shared_with_id" binding:"required"`
	Permission   string    `json:"permission" binding:"required"` // view, edit, admin
}

// ShareTask godoc
// @Summary      Share a task with another user
// @Description  Share a task with another user with specified permissions
// @Tags         collaboration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Param        share  body      ShareTaskRequest  true  "Share data"
// @Success      201  {object}  models.TaskShare
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/share [post]
func ShareTask(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	var input ShareTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid share input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate permission
	validPermissions := map[string]bool{
		"view":  true,
		"edit":  true,
		"admin": true,
	}
	if !validPermissions[input.Permission] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission. Must be view, edit, or admin"})
		return
	}

	// Verify task exists and belongs to user
	var task models.Task
	if err := config.GetDB().Where("id = ? AND user_id = ?", taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not found for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check if already shared
	var existingShare models.TaskShare
	if err := config.GetDB().Where("task_id = ? AND shared_with_id = ?", taskID, input.SharedWithID).First(&existingShare).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task already shared with this user"})
		return
	}

	// Create share
	share := models.TaskShare{
		TaskID:       taskID,
		OwnerID:      userIDUUID,
		SharedWithID: input.SharedWithID,
		Permission:   input.Permission,
	}

	if err := config.GetDB().Create(&share).Error; err != nil {
		config.Logger.Errorf("Error creating task share: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not share task"})
		return
	}

	config.Logger.Infof("Task %s shared by user %s with user %s", taskID, userIDUUID, input.SharedWithID)
	c.JSON(http.StatusCreated, share)
}

// GetSharedTasks godoc
// @Summary      Get tasks shared with the user
// @Description  Get all tasks that have been shared with the logged-in user
// @Tags         collaboration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]models.Task
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /shared-tasks [get]
func GetSharedTasks(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDUUID := userID.(uuid.UUID)

	var sharedTasks []models.Task
	if err := config.GetDB().Joins("JOIN task_shares ON tasks.id = task_shares.task_id").
		Where("task_shares.shared_with_id = ?", userIDUUID).
		Preload("User"). // Load the owner info
		Find(&sharedTasks).Error; err != nil {
		config.Logger.Errorf("Error fetching shared tasks for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch shared tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shared_tasks": sharedTasks})
}

// AddTaskCommentRequest represents the request body for adding a comment
type AddTaskCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// AddTaskComment godoc
// @Summary      Add a comment to a task
// @Description  Add a comment to a shared task
// @Tags         collaboration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Param        comment  body      AddTaskCommentRequest  true  "Comment data"
// @Success      201  {object}  models.TaskComment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/comments [post]
func AddTaskComment(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	var input AddTaskCommentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid comment input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Verify user has access to the task (either owns it or it's shared with them)
	var task models.Task
	if err := config.GetDB().Where("(id = ? AND user_id = ?) OR EXISTS(SELECT 1 FROM task_shares WHERE task_id = ? AND shared_with_id = ?)",
		taskID, userIDUUID, taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not accessible for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or not shared with you"})
		return
	}

	comment := models.TaskComment{
		TaskID:  taskID,
		UserID:  userIDUUID,
		Content: input.Content,
	}

	if err := config.GetDB().Create(&comment).Error; err != nil {
		config.Logger.Errorf("Error creating task comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add comment"})
		return
	}

	config.Logger.Infof("Comment added to task %s by user %s", taskID, userIDUUID)
	c.JSON(http.StatusCreated, comment)
}

// GetTaskComments godoc
// @Summary      Get comments for a task
// @Description  Get all comments for a shared task
// @Tags         collaboration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        taskID   path      string  true  "Task ID"
// @Success      200  {object}  map[string][]models.TaskComment
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{taskID}/comments [get]
func GetTaskComments(c *gin.Context) {
	taskIDStr := c.Param("taskID")
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
	userIDUUID := userID.(uuid.UUID)

	// Verify user has access to the task
	var task models.Task
	if err := config.GetDB().Where("(id = ? AND user_id = ?) OR EXISTS(SELECT 1 FROM task_shares WHERE task_id = ? AND shared_with_id = ?)",
		taskID, userIDUUID, taskID, userIDUUID).First(&task).Error; err != nil {
		config.Logger.Warnf("Task ID %s not accessible for user %s", taskID, userIDUUID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or not shared with you"})
		return
	}

	var comments []models.TaskComment
	if err := config.GetDB().Where("task_id = ?", taskID).
		Preload("User").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		config.Logger.Errorf("Error fetching comments for task %s: %v", taskID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
