package handlers

import (
	"net/http"
	"time"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminUserResponse represents user data for admin view
type AdminUserResponse struct {
	ID        uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdminStatsResponse represents system statistics
type AdminStatsResponse struct {
	TotalUsers         int64 `json:"total_users"`
	ActiveUsers        int64 `json:"active_users"`
	TotalTasks         int64 `json:"total_tasks"`
	TotalGoals         int64 `json:"total_goals"`
	TotalDecks         int64 `json:"total_decks"`
	TotalStudySessions int64 `json:"total_study_sessions"`
}

// UpdateUserRoleRequest represents request to update user role
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin"`
}

// GetAllUsers godoc
// @Summary      Get all users (admin only)
// @Description  Fetch all users with their details
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string][]AdminUserResponse
// @Failure      500  {object}  map[string]string
// @Router       /admin/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.GetDB().Find(&users).Error; err != nil {
		config.Logger.Errorf("Error fetching all users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}

	// Convert to admin response format
	adminUsers := make([]AdminUserResponse, len(users))
	for i, user := range users {
		adminUsers[i] = AdminUserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	config.Logger.Infof("Admin fetched %d users", len(users))
	c.JSON(http.StatusOK, gin.H{"users": adminUsers})
}

// UpdateUserRole godoc
// @Summary      Update user role (admin only)
// @Description  Change a user's role (user/admin)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userID   path      string  true  "User ID"
// @Param        role     body      UpdateUserRoleRequest  true  "New role"
// @Success      200      {object}  AdminUserResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /admin/users/{userID}/role [put]
func UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid role update input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found for role update: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update role
	if err := config.GetDB().Model(&user).Update("role", input.Role).Error; err != nil {
		config.Logger.Errorf("Error updating user role ID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	// Reload user to get updated data
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Errorf("Error reloading user after role update ID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload updated user"})
		return
	}

	adminUser := AdminUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	config.Logger.Infof("User role updated successfully: ID %s, new role: %s", userID, input.Role)
	c.JSON(http.StatusOK, adminUser)
}

// DeleteUserAdmin godoc
// @Summary      Delete user (admin only)
// @Description  Delete any user account
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userID   path      string  true  "User ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /admin/users/{userID} [delete]
func DeleteUserAdmin(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param for admin delete: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found for admin deletion: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	config.Logger.Infof("Admin deleting user account: ID %s, Email: %s", user.ID, user.Email)

	if err := config.GetDB().Delete(&user).Error; err != nil {
		config.Logger.Errorf("Failed to delete user ID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
		return
	}

	config.Logger.Infof("User account deleted by admin successfully: ID %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "User account deleted successfully",
		"user_id": userID,
	})
}

// GetSystemStats godoc
// @Summary      Get system statistics (admin only)
// @Description  Fetch overall system statistics
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  AdminStatsResponse
// @Failure      500  {object}  map[string]string
// @Router       /admin/stats [get]
func GetSystemStats(c *gin.Context) {
	var stats AdminStatsResponse

	// Count total users
	if err := config.GetDB().Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		config.Logger.Errorf("Error counting users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user statistics"})
		return
	}

	// Count active users (users created in last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := config.GetDB().Model(&models.User{}).Where("created_at >= ?", thirtyDaysAgo).Count(&stats.ActiveUsers).Error; err != nil {
		config.Logger.Errorf("Error counting active users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch active user statistics"})
		return
	}

	// Count total tasks
	if err := config.GetDB().Model(&models.Task{}).Count(&stats.TotalTasks).Error; err != nil {
		config.Logger.Errorf("Error counting tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch task statistics"})
		return
	}

	// Count total goals
	if err := config.GetDB().Model(&models.Goal{}).Count(&stats.TotalGoals).Error; err != nil {
		config.Logger.Errorf("Error counting goals: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch goal statistics"})
		return
	}

	// Count total decks
	if err := config.GetDB().Model(&models.Deck{}).Count(&stats.TotalDecks).Error; err != nil {
		config.Logger.Errorf("Error counting decks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch deck statistics"})
		return
	}

	// Count total study sessions
	if err := config.GetDB().Model(&models.StudySession{}).Count(&stats.TotalStudySessions).Error; err != nil {
		config.Logger.Errorf("Error counting study sessions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch study session statistics"})
		return
	}

	config.Logger.Info("Admin fetched system statistics")
	c.JSON(http.StatusOK, stats)
}

// PromoteToAdmin godoc
// @Summary      Promote user to admin (admin only)
// @Description  Grant admin privileges to a user
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userID   path      string  true  "User ID to promote"
// @Success      200      {object}  AdminUserResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /admin/users/{userID}/promote [post]
func PromoteToAdmin(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		config.Logger.Warnf("Invalid user ID param for promotion: %s", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Warnf("User not found for promotion: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role == "admin" {
		config.Logger.Warnf("User already admin: ID %s", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already an admin"})
		return
	}

	// Update role to admin
	if err := config.GetDB().Model(&user).Update("role", "admin").Error; err != nil {
		config.Logger.Errorf("Error promoting user to admin ID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user to admin"})
		return
	}

	// Reload user to get updated data
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		config.Logger.Errorf("Error reloading user after promotion ID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not reload promoted user"})
		return
	}

	adminUser := AdminUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	config.Logger.Infof("User promoted to admin successfully: ID %s", userID)
	c.JSON(http.StatusOK, adminUser)
}
