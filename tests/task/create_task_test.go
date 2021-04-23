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

func TestCreateTask(t *testing.T) {
	router := setUp()
	jsonStr, _ := json.Marshal(
		&models.NewTaskTemplate{
			Title:       "Title",
			Description: "Description",
		},
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/todo/api/v1.0/tasks",
		bytes.NewBuffer(jsonStr),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var data taskJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}

type mockTemplate struct{}

func TestFailToCreateTask(t *testing.T) {
	router := setUp()
	subtests := []struct {
		task        interface{}
		contentType string
		message     string
	}{
		{
			task:        &mockTemplate{},
			contentType: "application/json",
			message:     "Invalid input",
		},
		{
			task: &models.NewTaskTemplate{
				Title:       "",
				Description: "Description",
			},
			contentType: "application/json",
			message:     "Missing title from input",
		},
		{
			task: &models.NewTaskTemplate{
				Title:       "Title",
				Description: "Description",
			},
			contentType: "text/html",
			message:     "Invalid content type",
		},
	}

	for _, st := range subtests {
		jsonStr, _ := json.Marshal(st.task)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"POST",
			"/todo/api/v1.0/tasks",
			bytes.NewBuffer(jsonStr),
		)
		req.Header.Set("Content-Type", st.contentType)
		router.ServeHTTP(w, req)

		var data errorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code, st.message)
	}
}
