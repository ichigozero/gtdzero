package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
)

func SetTaskRoutes(router *gin.Engine, tc *controllers.TaskController) {
	v1 := router.Group("/todo/api/v1.0")
	{
		v1.GET("/tasks", tc.GetTasks)
		v1.GET("/task/:id", tc.GetTask)
		v1.POST("/tasks", tc.CreateTask)
	}
}
