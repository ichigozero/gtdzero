package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/libs/auth"
	"github.com/ichigozero/gtdzero/models"
)

type TaskController struct {
	db     models.TaskDB
	client auth.AuthClient
}

func NewTaskController(
	db models.TaskDB,
	client auth.AuthClient,
) *TaskController {
	return &TaskController{db, client}
}

func (t *TaskController) GetTasks(c *gin.Context) {
	userID := getUserID(c, t.client)
	if userID == 0 {
		return
	}

	tasks := t.db.GetTasks(userID)

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (t *TaskController) GetTask(c *gin.Context) {
	userID := getUserID(c, t.client)
	if userID == 0 {
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task, err := t.db.GetTask(userID, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (t *TaskController) CreateTask(c *gin.Context) {
	userID := getUserID(c, t.client)
	if userID == 0 {
		return
	}

	var json models.NewTaskTemplate

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	newTask := t.db.CreateTask(userID, &json)

	c.JSON(http.StatusCreated, gin.H{"task": newTask})
}

func (t *TaskController) UpdateTask(c *gin.Context) {
	userID := getUserID(c, t.client)
	if userID == 0 {
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	var json models.UpdateTaskTemplate

	err = c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	_, err = t.db.GetTask(userID, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	updatedTask := &models.Task{
		ID:          taskID,
		Title:       json.Title,
		Description: json.Description,
		Done:        json.Done,
		UserID:      userID,
	}

	t.db.UpdateTask(updatedTask)

	c.JSON(http.StatusOK, gin.H{"task": updatedTask})
}

func (t *TaskController) DeleteTask(c *gin.Context) {
	userID := getUserID(c, t.client)
	if userID == 0 {
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	err = t.db.DeleteTask(userID, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}
