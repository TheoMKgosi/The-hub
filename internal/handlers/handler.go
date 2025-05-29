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

// Task handlers
func GetTasks(c *gin.Context) {
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

func GetTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("ID"))
	if err != nil {
		log.Fatal("GetTask error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting task",
		})
		return
	}

	var task models.Task
	result := config.GetDB().First(&task, taskID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Task does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task

	taskID := c.Param("ID")
	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
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

	updatedTask := models.Task{
		Title:       input.Title,
		Description: input.Description,
	}

	if err := config.GetDB().Model(&task).Updates(updatedTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)

}

func DeleteTask(c *gin.Context) {
	var task models.Task

	taskID := c.Param("ID")
	if err := config.GetDB().First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
		})
		return
	}

	if err := config.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)

}
