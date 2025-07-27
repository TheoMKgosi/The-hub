package handlers

import (
	"log"
	"net/http"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
)

// Get all users
func GetUsers(c *gin.Context) {
	var users []models.User
	result := config.GetDB().Find(&users)

	if result.Error != nil {
		log.Printf("Error fetching users: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}

	log.Printf("Fetched %d tasks", len(users))
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Login user
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid login input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("User not found with email: %s", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !util.CheckPasswordHash(user.Password, input.Password) {
		log.Printf("Password mismatch for user: %s", input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Token generation failed for user ID %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("User login successful: ID %d", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// Register user
func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid registration input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		log.Printf("Password hashing failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user := models.User{
		Email:    input.Email,
		Name:     input.Name,
		Password: hashedPassword,
	}

	if err := config.GetDB().Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("JWT generation failed for new user ID %d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	log.Printf("User registered successfully: ID %d", user.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"token":   token,
	})
}

// Update a specific user
func UpdateUser(c *gin.Context) {
	var user models.User
	userID := c.Param("ID")

	if err := config.GetDB().First(&user, userID).Error; err != nil {
		log.Printf("User not found for update: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid update input for user ID %s: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Password != nil {
		hashedPassword, err := util.HashPassword(*input.Password)
		if err != nil {
			log.Printf("Password hashing failed for user ID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		updates["password"] = hashedPassword
	}

	if err := config.GetDB().Model(&user).Updates(updates).Error; err != nil {
		log.Printf("Error updating user ID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	log.Printf("User updated successfully: ID %s", userID)
	c.JSON(http.StatusOK, user)
}

// Delete a specific user
func DeleteUser(c *gin.Context) {
	var user models.User
	userID := c.Param("ID")

	if err := config.GetDB().First(&user, userID).Error; err != nil {
		log.Printf("User not found for deletion: ID %s", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.GetDB().Delete(&user).Error; err != nil {
		log.Printf("Failed to delete user ID %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	log.Printf("User deleted successfully: ID %s", userID)
	c.JSON(http.StatusOK, user)
}

