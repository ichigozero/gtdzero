package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/tests"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTask(t *testing.T) {
	router := tests.SetUp()
	w := httptest.NewRecorder()

	accessToken, _ := tests.Login(router, w)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/todo/api/v1.0/task/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	router.ServeHTTP(w, req)

	var data tests.ResultJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFailToDeleteTask(t *testing.T) {
	router := tests.SetUp()
	subtests := []struct {
		uri          string
		responseCode int
	}{
		{
			uri:          "/todo/api/v1.0/task/a",
			responseCode: http.StatusBadRequest,
		},
		{
			uri:          "/todo/api/v1.0/task/3",
			responseCode: http.StatusNotFound,
		},
	}

	w := httptest.NewRecorder()

	accessToken, _ := tests.Login(router, w)

	for _, st := range subtests {
		w = httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", st.uri, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		router.ServeHTTP(w, req)

		var data tests.ErrorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code)
	}
}
