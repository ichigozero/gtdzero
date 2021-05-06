package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/libs/auth"
	"github.com/ichigozero/gtdzero/models"
)

type AuthController struct {
	db        models.UserDB
	tokenizer auth.Tokenizer
	client    auth.AuthClient
}

func NewAuthController(
	db models.UserDB,
	tokenizer auth.Tokenizer,
	client auth.AuthClient,
) *AuthController {
	return &AuthController{
		db:        db,
		tokenizer: tokenizer,
		client:    client,
	}
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

	tokenDetails, err := a.tokenizer.Create(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err = a.client.Store(user.ID, tokenDetails)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	tokens := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}

	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})
}

func (a *AuthController) Logout(c *gin.Context) {
	_, err := a.client.Delete(c.Request)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}
