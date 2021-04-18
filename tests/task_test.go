package tests

import (
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
