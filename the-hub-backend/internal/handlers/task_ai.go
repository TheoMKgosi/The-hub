package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/ai"
	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxTasksForAI = 5

type AITaskInput struct {
	TaskID      string `json:"task_id" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AISubtaskSuggestion struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	EstimatedHours int    `json:"estimated_hours"`
}

type AITaskEnhancement struct {
	TaskID          string               `json:"task_id"`
	OriginalTitle   string               `json:"original_title"`
	Title          string               `json:"title"`
	Description   string               `json:"description"`
	Priority       int                  `json:"priority"`
	EstimatedHours int                  `json:"estimated_hours"`
	Subtasks      []AISubtaskSuggestion `json:"subtasks"`
}

type AITaskPreviewResponse struct {
	Preview []AITaskEnhancement `json:"preview"`
	Message string          `json:"message"`
}

type AppliedTask struct {
	TaskID    string `json:"task_id" binding:"required"`
	Selected  bool   `json:"selected"`
}

type ApplyAITasksRequest struct {
	AppliedTasks []AppliedTask `json:"applied_tasks" binding:"required"`
}

func GetAITaskPreview(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var tasks []models.Task
	if err := config.GetDB().
		Where("user_id = ? AND ai_checked = ? AND status != ?", userIDUUID, false, "completed").
		Order("created_at ASC").
		Limit(maxTasksForAI).
		Find(&tasks).Error; err != nil {
		config.Logger.Errorf("Failed to fetch tasks for AI check: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, AITaskPreviewResponse{
			Preview: []AITaskEnhancement{},
			Message: "All tasks already optimized",
		})
		return
	}

	client, err := ai.GetOpenRouterClient()
	if err != nil {
		config.Logger.Errorf("Failed to get AI client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable"})
		return
	}

	taskInputs := make([]ai.TaskEnhancementInput, len(tasks))
	for i, task := range tasks {
		taskInputs[i] = ai.TaskEnhancementInput{
			TaskID:      task.ID.String(),
			Title:       task.Title,
			Description: task.Description,
		}
	}

	aiResponse, err := client.EnhanceTasks(taskInputs)
	if err != nil {
		config.Logger.Errorf("Failed to enhance tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI recommendations"})
		return
	}

	var enhancements []AITaskEnhancement
	if err := json.Unmarshal([]byte(aiResponse), &enhancements); err != nil {
		start := strings.Index(aiResponse, "[")
		end := strings.LastIndex(aiResponse, "]")
		if start == -1 || end == -1 || start >= end {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
			return
		}

		jsonStr := aiResponse[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &enhancements); err != nil {
			config.Logger.Errorf("Failed to parse AI response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
			return
		}
	}

	preview := make([]AITaskEnhancement, 0, len(enhancements))
	for _, e := range enhancements {
		for _, task := range tasks {
			if task.ID.String() == e.TaskID {
				preview = append(preview, AITaskEnhancement{
					TaskID:          e.TaskID,
					OriginalTitle:  task.Title,
					Title:          e.Title,
					Description:    e.Description,
					Priority:       e.Priority,
					EstimatedHours: e.EstimatedHours,
					Subtasks:       e.Subtasks,
				})
				break
			}
		}
	}

	message := fmt.Sprintf("Found %d tasks to optimize", len(preview))
	c.JSON(http.StatusOK, AITaskPreviewResponse{
		Preview: preview,
		Message: message,
	})
}

func ApplyAITasks(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var req ApplyAITasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.AppliedTasks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No tasks selected"})
		return
	}

	client, err := ai.GetOpenRouterClient()
	if err != nil {
		config.Logger.Errorf("Failed to get AI client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable"})
		return
	}

	taskInputs := make([]ai.TaskEnhancementInput, 0, len(req.AppliedTasks))
	taskMap := make(map[string]models.Task)

	for _, at := range req.AppliedTasks {
		if !at.Selected {
			continue
		}

		taskID, err := uuid.Parse(at.TaskID)
		if err != nil {
			continue
		}

		var task models.Task
		if err := config.GetDB().
			Where("id = ? AND user_id = ?", taskID, userIDUUID).
			First(&task).Error; err != nil {
			continue
		}

		taskMap[task.ID.String()] = task
		taskInputs = append(taskInputs, ai.TaskEnhancementInput{
			TaskID:      task.ID.String(),
			Title:       task.Title,
			Description: task.Description,
		})
	}

	if len(taskInputs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid tasks to apply"})
		return
	}

	aiResponse, err := client.EnhanceTasks(taskInputs)
	if err != nil {
		config.Logger.Errorf("Failed to enhance tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI recommendations"})
		return
	}

	var enhancements []AITaskEnhancement
	if err := json.Unmarshal([]byte(aiResponse), &enhancements); err != nil {
		start := strings.Index(aiResponse, "[")
		end := strings.LastIndex(aiResponse, "]")
		if start == -1 || end == -1 || start >= end {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
			return
		}

		jsonStr := aiResponse[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &enhancements); err != nil {
			config.Logger.Errorf("Failed to parse AI response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
			return
		}
	}

	now := time.Now()
	for _, e := range enhancements {
		task, exists := taskMap[e.TaskID]
		if !exists {
			continue
		}

		timeEstimate := e.EstimatedHours * 60

		if err := config.GetDB().Model(&task).Updates(map[string]interface{}{
			"title":              e.Title,
			"description":       e.Description,
			"priority":          e.Priority,
			"time_estimate": timeEstimate,
			"ai_checked":        true,
		}).Error; err != nil {
			config.Logger.Errorf("Failed to update task %s: %v", task.ID, err)
			continue
		}

		if len(e.Subtasks) > 0 {
			var maxOrderIndex int
			config.GetDB().Model(&models.Task{}).
				Where("user_id = ? AND parent_task_id = ?", userIDUUID, task.ID).
				Select("COALESCE(MAX(order_index), 0)").
				Scan(&maxOrderIndex)

			for i, sub := range e.Subtasks {
				subTimeEstimate := sub.EstimatedHours * 60
				subtask := models.Task{
					ID:            uuid.New(),
					Title:         sub.Title,
					Description:  sub.Description,
					UserID:        userIDUUID,
					ParentTaskID: &task.ID,
					Priority:     &e.Priority,
					Status:       "pending",
					OrderIndex:   maxOrderIndex + i + 1,
					TimeEstimate: &subTimeEstimate,
				}

				if err := config.GetDB().Create(&subtask).Error; err != nil {
					config.Logger.Errorf("Failed to create subtask: %v", err)
				}
			}
		}

		if task.ParentTaskID != nil {
			var parentTask models.Task
			if err := config.GetDB().Where("id = ?", *task.ParentTaskID).First(&parentTask).Error; err == nil {
				subtasks, _ := parentTask.GetSubtasks(config.GetDB())
				if len(subtasks) > 0 && parentTask.Status == "completed" {
					completed := 0
					for _, st := range subtasks {
						if st.Status == "completed" {
							completed++
						}
					}
					if completed != len(subtasks) {
						config.GetDB().Model(&parentTask).Updates(map[string]interface{}{
							"status":       "pending",
							"completed_at": nil,
						})
					}
				}
			}
		}
		_ = now
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Tasks updated successfully",
		"updated_count":  len(enhancements),
	})
}
