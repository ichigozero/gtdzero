package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/libs/auth"
	"github.com/ichigozero/gtdzero/models"
	"golang.org/x/crypto/bcrypt"
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

	user, err := a.db.GetUser(json.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login"})
		return
	}

	tokens, err := createStoreToken(a.tokenizer, a.client, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})
}

func (a *AuthController) Logout(c *gin.Context) {
	details := getAccessTokenDetails(c, a.client)
	if details == nil {
		return
	}

	err := a.client.Delete(details.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshUUID := auth.GenerateRefreshUUID(details.UUID)
	err = a.client.Delete(refreshUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func (a *AuthController) Refresh(c *gin.Context) {
	m := map[string]string{}
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	t := m["refresh_token"]

	cl, err := auth.GetTokenClaims(t, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	d, err := auth.ExtractRefreshToken(cl)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userID, err := a.client.Fetch(d.RefreshUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if d.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = a.client.Delete(d.AccessUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = a.client.Delete(d.RefreshUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokens, err := createStoreToken(a.tokenizer, a.client, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})
}

func createStoreToken(t auth.Tokenizer, c auth.AuthClient, userID uint64) (map[string]string, error) {
	tokenDetails, err := t.Create(userID)
	if err != nil {
		return nil, err
	}

	err = c.Store(userID, tokenDetails)
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}

	return tokens, nil
}
