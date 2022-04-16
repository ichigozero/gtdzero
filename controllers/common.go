package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/libs/auth"
)

func getUserID(c *gin.Context, client auth.AuthClient) uint64 {
	details := getAccessTokenDetails(c, client)
	if details == nil {
		return 0
	}

	userID, err := client.Fetch(details.UUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return 0
	}

	if userID != details.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0
	}

	return userID
}

func getAccessTokenDetails(c *gin.Context, client auth.AuthClient) *auth.AccessTokenDetails {
	td, exists := c.Get("token_details")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "unable to get token details"},
		)
		return nil
	}

	details, ok := td.(*auth.AccessTokenDetails)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "unable to parse token details"},
		)
		return nil
	}

	return details
}
