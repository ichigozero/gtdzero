package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/tests"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTask(t *testing.T) {
	router := tests.SetUp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/todo/api/v1.0/task/1", nil)
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

	for _, st := range subtests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", st.uri, nil)
		router.ServeHTTP(w, req)

		var data tests.ErrorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code)
	}
}
