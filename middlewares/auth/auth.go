package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero"
	"github.com/ichigozero/gtdzero/libs/auth"
)

func AccessTokenValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		strArr := strings.Split(bearerToken, " ")
		if len(strArr) != 2 {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid authorization request header"},
			)
			c.Abort()
			return
		}

		claims, err := auth.GetTokenClaims(strArr[1], gtdzero.AccessSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		details, err := auth.ExtractAccessToken(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("token_details", details)
		c.Next()
	}
}
