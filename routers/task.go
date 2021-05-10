package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
	"github.com/ichigozero/gtdzero/middlewares/auth"
)

func SetTaskRoutes(router *gin.Engine, tc *controllers.TaskController) {
	v1 := router.Group("/todo/api/v1.0")
	{
		v1.GET("/tasks", auth.TokenValidator(), tc.GetTasks)
		v1.GET("/task/:id", auth.TokenValidator(), tc.GetTask)
		v1.POST("/tasks", tc.CreateTask)
		v1.PUT("/task/:id", tc.UpdateTask)
		v1.DELETE("/task/:id", tc.DeleteTask)
	}
}
