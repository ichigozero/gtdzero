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
	id, err := strconv.Atoi(c.Param("id"))
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
	if c.Request.Header.Get("Content-Type") != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	var json models.NewTaskTemplate

	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	newTask := t.db.CreateTask(&json)

	c.JSON(http.StatusCreated, gin.H{"task": newTask})
}
