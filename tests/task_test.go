package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/routers"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	Tasks []*models.Task
}

func (d *mockDB) GetTasks() []*models.Task {
	return d.Tasks
}

func (d *mockDB) GetTask(id int) (*models.Task, error) {
	for _, task := range d.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, errors.New("not found")
}

func (d *mockDB) CreateTask(t *models.NewTaskTemplate) *models.Task {
	newTask := &models.Task{
		ID:          d.Tasks[len(d.Tasks)-1].ID + 1,
		Title:       t.Title,
		Description: t.Description,
		Done:        false,
	}

	d.Tasks = append(d.Tasks, newTask)

	return newTask
}

func (d *mockDB) UpdateTask(t *models.Task) error {
	return nil
}

func (d *mockDB) DeleteTask(id int) error {
	for index, task := range d.Tasks {
		if task.ID == id {
			d.Tasks = append(d.Tasks[:index], d.Tasks[index+1:]...)
			return nil
		}
	}

	return errors.New("not found")
}

func setUp() *gin.Engine {
	r := gin.Default()
	db := &mockDB{
		[]*models.Task{
			{
				ID:          1,
				Title:       "Buy groceries",
				Description: "Milk, Cheese, Pizza, Fruit, Tylenol",
				Done:        false,
			},
			{
				ID:          2,
				Title:       "Learn Go",
				Description: "Need to find a good Go tutorial on the web",
				Done:        false,
			},
		},
	}

	tc := controllers.NewTaskController(db)
	routers.SetTaskRoutes(r, tc)

	return r
}

type tasksJSON struct {
	Tasks []*models.Task `json:"tasks"`
}

func TestGetTasks(t *testing.T) {
	router := setUp()

	w := httptest.NewRecorder()
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
	router := setUp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/todo/api/v1.0/task/1", nil)
	router.ServeHTTP(w, req)

	var data taskJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

type errorJSON struct {
	Error string `json:"error"`
}

func TestFailToGetTask(t *testing.T) {
	router := setUp()
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

		var data errorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code)
	}
}

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
			task:         &models.UpdateTaskTemplate{},
			responseCode: http.StatusBadRequest,
			message:      "Invalid ID",
		},
		{
			uri:          "/todo/api/v1.0/task/3",
			contentType:  "application/json",
			task:         &models.UpdateTaskTemplate{},
			responseCode: http.StatusNotFound,
			message:      "Task not found",
		},
		{
			uri:          "/todo/api/v1.0/task/1",
			contentType:  "text/html",
			task:         &models.UpdateTaskTemplate{},
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

func TestDeleteTask(t *testing.T) {
	router := setUp()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/todo/api/v1.0/task/1", nil)
	router.ServeHTTP(w, req)

	var data resultJSON
	err := json.NewDecoder(w.Body).Decode(&data)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFailToDeleteTask(t *testing.T) {
	router := setUp()
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

		var data errorJSON
		err := json.NewDecoder(w.Body).Decode(&data)

		assert.Nil(t, err)
		assert.Equal(t, st.responseCode, w.Code)
	}
}
