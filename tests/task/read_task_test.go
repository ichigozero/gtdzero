package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/tests"
	"github.com/stretchr/testify/assert"
)

type tasksJSON struct {
	Tasks []*models.Task `json:"tasks"`
}

func TestGetTasks(t *testing.T) {
	router := tests.SetUp()

	w := httptest.NewRecorder()
	tests.Login(router, w)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todo/api/v1.0/tasks", nil)
	router.ServeHTTP(w, req)

	var data tasksJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

type taskJSON struct {
	Task *models.Task `json:"task"`
}

func TestGetTask(t *testing.T) {
	router := tests.SetUp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todo/api/v1.0/task/1", nil)
	router.ServeHTTP(w, req)

	var data taskJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFailToGetTask(t *testing.T) {
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
		req, _ := http.NewRequest("GET", st.uri, nil)
		router.ServeHTTP(w, req)

		var data tests.ErrorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code)
	}
}
