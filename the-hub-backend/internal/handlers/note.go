package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateNoteRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdateNoteRequest struct {
	Title   *string  `json:"title"`
	Content *string  `json:"content"`
	Tags    []string `json:"tags"`
}

func GetNotes(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	search := c.Query("search")
	tag := c.Query("tag")

	query := config.GetDB().Where("user_id = ?", userID)

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchPattern, searchPattern)
	}

	if tag != "" {
		query = query.Where("tags::text LIKE ?", "%"+tag+"%")
	}

	var notes []models.Note
	if err := query.Order("updated_at DESC").Find(&notes).Error; err != nil {
		config.Logger.Errorf("Error fetching notes for user %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch notes"})
		return
	}

	var response []models.NoteResponse
	for _, note := range notes {
		response = append(response, note.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"notes": response})
}

func GetNote(c *gin.Context) {
	noteIDStr := c.Param("ID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid note ID param: %s", noteIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var note models.Note
	if err := config.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		config.Logger.Errorf("Note ID %s not found for user %v: %v", noteID, userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note.ToResponse()})
}

func CreateNote(c *gin.Context) {
	var input CreateNoteRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid note input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for note", "details": err.Error()})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during note creation")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		config.Logger.Errorf("Invalid userID type in context: %T", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	tagsJSON, _ := json.Marshal(input.Tags)

	note := models.Note{
		Title:   input.Title,
		Content: input.Content,
		Tags:    string(tagsJSON),
		UserID:  userIDUUID,
	}

	config.Logger.Infof("Creating note for user %s: %s", userIDUUID, input.Title)
	if err := config.GetDB().Create(&note).Error; err != nil {
		config.Logger.Errorf("Error creating note for user %s: %v", userIDUUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create note"})
		return
	}

	config.Logger.Infof("Successfully created note ID %s for user %s", note.ID, userIDUUID)
	c.JSON(http.StatusCreated, note.ToResponse())
}

func UpdateNote(c *gin.Context) {
	noteIDStr := c.Param("ID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid note ID param for update: %s", noteIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during note update")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var note models.Note
	if err := config.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		config.Logger.Warnf("Note not found for update: ID %s, User %v", noteID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	var input UpdateNoteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid update input for note ID %s: %v", noteID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Content != nil {
		updates["content"] = *input.Content
	}
	if input.Tags != nil {
		tagsJSON, _ := json.Marshal(input.Tags)
		updates["tags"] = string(tagsJSON)
	}

	if len(updates) == 0 {
		config.Logger.Warnf("No valid fields provided for note update: ID %s", noteID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	config.Logger.Infof("Updating note ID %s for user %v with data: %+v", noteID, userID, updates)
	if err := config.GetDB().Model(&note).Updates(updates).Error; err != nil {
		config.Logger.Errorf("Failed to update note ID %s: %v", noteID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	if err := config.GetDB().First(&note, note.ID).Error; err != nil {
		config.Logger.Errorf("Error retrieving updated note ID %s: %v", note.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated note"})
		return
	}

	config.Logger.Infof("Successfully updated note ID %s for user %v", note.ID, userID)
	c.JSON(http.StatusOK, note.ToResponse())
}

func DeleteNote(c *gin.Context) {
	noteIDStr := c.Param("ID")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid note ID param for delete: %s", noteIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		config.Logger.Warn("userID not found in context during note deletion")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var note models.Note
	if err := config.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		config.Logger.Warnf("Note not found for delete: ID %s, User %v", noteID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	config.Logger.Infof("Deleting note ID %s for user %v", noteID, userID)
	if err := config.GetDB().Delete(&note).Error; err != nil {
		config.Logger.Errorf("Failed to delete note ID %s: %v", noteID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	config.Logger.Infof("Successfully deleted note ID %s for user %v", noteID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
