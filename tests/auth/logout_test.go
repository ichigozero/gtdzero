package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	router := setUp()
	w := httptest.NewRecorder()

	login(router, w)

	var tk tokenJSON
	json.NewDecoder(w.Body).Decode(&tk)

	//TODO mock RedisClient
	w = httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", tk.Tokens["access_token"]),
	)
	router.ServeHTTP(w, req)

	var re resultJSON
	err := json.NewDecoder(w.Body).Decode(&re)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFailedLogout(t *testing.T) {
	router := setUp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var data errorJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
