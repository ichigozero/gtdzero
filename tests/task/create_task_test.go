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

func TestCreateTask(t *testing.T) {
	router := tests.Setup()
	jsonStr, _ := json.Marshal(
		&models.NewTaskTemplate{
			Title:       "Title",
			Description: "Description",
		},
	)
	w := httptest.NewRecorder()

	accessToken, _ := tests.Login(router, w)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/todo/api/v1.0/tasks",
		bytes.NewBuffer(jsonStr),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	router.ServeHTTP(w, req)

	var data taskJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
}

type mockTemplate struct{}

func TestFailToCreateTask(t *testing.T) {
	router := tests.Setup()
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
			contentType: "text/xml",
			message:     "Invalid content type",
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

		req, _ := http.NewRequest(
			"POST",
			"/todo/api/v1.0/tasks",
			bytes.NewBuffer(buf),
		)
		req.Header.Set("Content-Type", st.contentType)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		router.ServeHTTP(w, req)

		var data tests.ErrorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code, st.message)
	}
}
