package contenttype

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AllowOnlyJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		c.Next()
	}
}
