package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/libs/authtoken"
	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/services/redis"
)

type AuthController struct {
	db models.UserDB
	rc redis.Client
}

func NewAuthController(db models.UserDB, rc redis.Client) *AuthController {
	return &AuthController{db, rc}
}

func (a *AuthController) Login(c *gin.Context) {
	var json models.UserLoginTemplate
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	user, err := a.db.GetUser(json.Username, json.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login"})
		return
	}

	at, err := authtoken.Create(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = authtoken.StoreAuth(user.ID, at, a.rc)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  at.AccessToken,
		"refresh_token": at.RefreshToken,
	}

	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})
}
