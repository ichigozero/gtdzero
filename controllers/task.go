package controllers

import (
	"net/http"

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
