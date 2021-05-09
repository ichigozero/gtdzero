package auth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/tests"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	router := tests.SetUp()
	w := httptest.NewRecorder()

	tests.Login(router, w)

	var data tests.TokenJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestFailedLogin(t *testing.T) {
	router := tests.SetUp()
	subtests := []struct {
		user         interface{}
		contentType  string
		responseCode int
		message      string
	}{
		{
			user: &models.UserLoginTemplate{
				Username: "",
				Password: "",
			},
			contentType:  "application/json",
			responseCode: http.StatusBadRequest,
			message:      "Invalid input",
		},
		{
			user: &models.UserLoginTemplate{
				Username: "john",
				Password: "password",
			},
			contentType:  "text/xml",
			responseCode: http.StatusBadRequest,
			message:      "Invalid content type",
		},
		{
			user: &models.UserLoginTemplate{
				Username: "jean",
				Password: "password",
			},
			contentType:  "application/json",
			responseCode: http.StatusUnauthorized,
			message:      "User not found",
		},
	}

	for _, st := range subtests {
		w := httptest.NewRecorder()

		var buf []byte

		if st.contentType == "application/json" {
			buf, _ = json.Marshal(st.user)
		} else {
			buf, _ = xml.Marshal(st.user)
		}

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(buf))
		req.Header.Set("Content-Type", st.contentType)
		router.ServeHTTP(w, req)

		var data tests.ErrorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code, st.message)
	}
}
