package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/libs/auth"
)

func TokenValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.ValidateToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
