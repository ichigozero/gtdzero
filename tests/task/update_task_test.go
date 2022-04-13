package task

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/tests"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTask(t *testing.T) {
	router := tests.SetUp()
	jsonStr, _ := json.Marshal(
		&models.UpdateTaskTemplate{
			Title:       "Title",
			Description: "Description",
			Done:        true,
		},
	)
	w := httptest.NewRecorder()

	accessToken, _ := tests.Login(router, w)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest(
		"PUT",
		"/todo/api/v1.0/task/1",
		bytes.NewBuffer(jsonStr),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
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
	router := tests.SetUp()
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
			uri:         "/todo/api/v1.0/task/1",
			contentType: "text/xml",
			task: &models.UpdateTaskTemplate{
				Title:       "Title",
				Description: "Description",
				Done:        true,
			},
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

	w := httptest.NewRecorder()

	accessToken, _ := tests.Login(router, w)

	for _, st := range subtests {
		w = httptest.NewRecorder()

		var buf []byte

		if st.contentType == "application/json" {
			buf, _ = json.Marshal(st.task)
		} else {
			buf, _ = xml.Marshal(st.task)
		}

		req, _ := http.NewRequest("PUT", st.uri, bytes.NewBuffer(buf))
		req.Header.Set("Content-Type", st.contentType)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		router.ServeHTTP(w, req)

		var data tests.ErrorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code, st.message)
	}
}
