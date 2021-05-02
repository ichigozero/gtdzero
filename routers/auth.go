package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
)

func SetAuthRoutes(router *gin.Engine, ac *controllers.AuthController) {
	router.POST("/login", ac.Login)
}
