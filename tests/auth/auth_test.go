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
		&models.UserLoginTemplate{
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

func TestFailedLogin(t *testing.T) {
	router := setUp()
	subtests := []struct {
		user         interface{}
		responseCode int
		message      string
	}{
		{
			user: &models.UserLoginTemplate{
				Username: "",
				Password: "",
			},
			responseCode: http.StatusBadRequest,
			message:      "Invalid input",
		},
		{
			user: &models.UserLoginTemplate{
				Username: "jean",
				Password: "password",
			},
			responseCode: http.StatusUnauthorized,
			message:      "User not found",
		},
	}

	for _, st := range subtests {
		jsonStr, _ := json.Marshal(st.user)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var data tokenJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code, st.message)
	}
}
