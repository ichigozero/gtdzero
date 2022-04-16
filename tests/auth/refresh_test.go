package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/tests"
	"github.com/stretchr/testify/assert"
)

type authTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type tokenTemplate struct {
	RefreshToken string `json:"refresh_token"`
}

func TestRefresh(t *testing.T) {
	router := tests.Setup()
	w := httptest.NewRecorder()

	tokens, _ := tests.Login(router, w)
	jsonStr, _ := json.Marshal(
		&tokenTemplate{
			RefreshToken: tokens["refresh_token"],
		},
	)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var at authTokens
	err := json.NewDecoder(w.Body).Decode(&at)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotEqual(t, tokens["refresh_token"], at.RefreshToken)
}

func TestFailedRefresh(t *testing.T) {
	router := tests.Setup()
	jsonStr, _ := json.Marshal(
		&tokenTemplate{
			RefreshToken: "dummy-token",
		},
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var data tests.ErrorJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
