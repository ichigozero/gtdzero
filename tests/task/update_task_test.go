package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/models"
	"github.com/stretchr/testify/assert"
)

type resultJSON struct {
	Result bool `json:"result"`
}

func TestUpdateTask(t *testing.T) {
	router := setUp()
	jsonStr, _ := json.Marshal(
		&models.UpdateTaskTemplate{
			Title:       "Title",
			Description: "Description",
			Done:        true,
		},
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"PUT",
		"/todo/api/v1.0/task/1",
		bytes.NewBuffer(jsonStr),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var data taskJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(
		t,
		&models.Task{
			ID:          1,
			Title:       "Title",
			Description: "Description",
			Done:        true,
		},
		data.Task,
	)
}

func TestFailToUpdateTask(t *testing.T) {
	router := setUp()
	task := &models.UpdateTaskTemplate{
		Title:       "Title",
		Description: "Description",
		Done:        true,
	}

	subtests := []struct {
		uri          string
		contentType  string
		task         interface{}
		responseCode int
		message      string
	}{
		{
			uri:          "/todo/api/v1.0/task/a",
			contentType:  "application/json",
			task:         task,
			responseCode: http.StatusBadRequest,
			message:      "Invalid ID",
		},
		{
			uri:          "/todo/api/v1.0/task/3",
			contentType:  "application/json",
			task:         task,
			responseCode: http.StatusNotFound,
			message:      "Task not found",
		},
		{
			uri:          "/todo/api/v1.0/task/1",
			contentType:  "text/html",
			task:         task,
			responseCode: http.StatusBadRequest,
			message:      "Invalid content type",
		},
		{
			uri:          "/todo/api/v1.0/task/1",
			contentType:  "application/json",
			task:         mockTemplate{},
			responseCode: http.StatusBadRequest,
			message:      "Invalid input",
		},
	}

	for _, st := range subtests {
		jsonStr, _ := json.Marshal(st.task)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", st.uri, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", st.contentType)
		router.ServeHTTP(w, req)

		var data errorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code, st.message)
	}
}
