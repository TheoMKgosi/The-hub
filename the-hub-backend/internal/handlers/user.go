package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/TheoMKgosi/The-hub/internal/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Task handlers

// Get all tasks
func GetUsers(c *gin.Context) {
	// TODO: Decide if to implement
	var tasks []models.Task
	result := config.GetDB().Find(&tasks)

	if result.Error != nil {
		log.Fatal(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})

}

// Get a specific task
func Login(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := config.GetDB().Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not login user",
		})
		return
	}

	if !util.CheckPasswordHash(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// Create a task
func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}

	user := models.User{
		Email:    input.Email,
		Name:     input.Name,
		Password: hashedPassword,
	}

	if err := config.GetDB().Create(&user).Error; err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":"Registration successful",
		"token": token,
	})
}

// Update a specific task
func UpdateUser(c *gin.Context) {
	var user models.User

	userID := c.Param("ID")
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		log.Println("Error ID: ", err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error JSON: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser := map[string]interface{}{}
	if input.Email != nil {
		updatedUser["email"] = *input.Email
	}
	if input.Name != nil {
		updatedUser["name"] = *input.Name
	}
	if input.Password != nil {
		updatedUser["password"] = *input.Password
	}

	if err := config.GetDB().Model(&user).Updates(updatedUser).Error; err != nil {
		log.Println("Error updating task:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// reload the task to get updated due date
	if err := config.GetDB().First(&user, user.ID).Error; err != nil {
		log.Println("Error retrieving updated task:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving updated task"})
		return
	}

	c.JSON(http.StatusOK, user)

}

// Delete a specific user
func DeleteUser(c *gin.Context) {
	var user models.User

	userID := c.Param("ID")
	if err := config.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	if err := config.GetDB().Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}

