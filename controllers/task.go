package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/models"
)

type TaskController struct {
	db models.TaskDB
}

func NewTaskController(db models.TaskDB) *TaskController {
	return &TaskController{db}
}

func (t *TaskController) GetTasks(c *gin.Context) {
	tasks := t.db.GetTasks()

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (t *TaskController) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task, err := t.db.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (t *TaskController) CreateTask(c *gin.Context) {
	var json models.NewTaskTemplate

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	newTask := t.db.CreateTask(&json)

	c.JSON(http.StatusCreated, gin.H{"task": newTask})
}

func (t *TaskController) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	_, err = t.db.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var json models.UpdateTaskTemplate

	err = c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	updatedTask := &models.Task{
		ID:          id,
		Title:       json.Title,
		Description: json.Description,
		Done:        json.Done,
	}

	t.db.UpdateTask(updatedTask)

	c.JSON(http.StatusOK, gin.H{"task": updatedTask})
}

func (t *TaskController) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	err = t.db.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}
