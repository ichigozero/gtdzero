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
