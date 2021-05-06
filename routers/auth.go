package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
	"github.com/ichigozero/gtdzero/middlewares/contenttype"
)

func SetAuthRoutes(router *gin.Engine, ac *controllers.AuthController) {
	router.POST("/login", contenttype.AllowOnlyJSON(), ac.Login)
	router.POST("/logout", ac.Logout)
}
