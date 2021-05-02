package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/models"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	router := setUp()
	jsonStr, _ := json.Marshal(
		&models.User{
			Username: "john",
			Password: "password",
		},
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var data tokenJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}
