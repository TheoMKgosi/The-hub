package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/TheoMKgosi/The-hub/internal/config"
	"github.com/TheoMKgosi/The-hub/internal/models"
	"github.com/gin-gonic/gin"
)

// Goal handlers
func GetGoals(c *gin.Context) {
	var goals []models.Goal

	if err := config.GetDB().Preload("Tasks").Find(&goals).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error,
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"goals": goals,
	})

}

func GetGoal(c *gin.Context) {
	goalID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Fatal("GetGoal error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting goal",
		})
		return
	}

	var goal models.Goal
	if err := config.GetDB().Preload("Tasks").First(&goal, goalID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Goal does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"goal": goal,
	})
}

func CreateGoal(c *gin.Context) {
	var goal models.Goal

	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.GetDB().Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

func UpdateGoal(c *gin.Context) {
	var goal models.Goal

	goalID := c.Param("ID")
	if err := config.GetDB().First(&goal, goalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedGoal := models.Goal{
		Title:       input.Title,
		Description: input.Description,
	}

	if err := config.GetDB().Model(&goal).Updates(updatedGoal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)

}

func DeleteGoal(c *gin.Context) {
	var goal models.Goal

	goalID := c.Param("ID")
	if err := config.GetDB().First(&goal, goalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Goal not found",
		})
		return
	}

	if err := config.GetDB().Delete(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, goal)

}
