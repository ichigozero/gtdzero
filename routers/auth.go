package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
	"github.com/ichigozero/gtdzero/middlewares/auth"
)

func SetAuthRoutes(router *gin.Engine, ac *controllers.AuthController) {
	router.POST("/login", ac.Login)
	router.POST("/logout", auth.AccessTokenValidator(), ac.Logout)
	router.POST("/refresh", ac.Refresh)
}
