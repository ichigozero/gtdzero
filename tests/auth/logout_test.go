package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/tests"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	router := tests.Setup()
	w := httptest.NewRecorder()

	tokens, _ := tests.Login(router, w)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens["access_token"]))
	router.ServeHTTP(w, req)

	var re tests.ResultJSON
	err := json.NewDecoder(w.Body).Decode(&re)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFailedLogout(t *testing.T) {
	router := tests.Setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var data tests.ErrorJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
