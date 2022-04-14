package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
	"github.com/ichigozero/gtdzero/libs/auth"
	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/routers"
)

type mockUserDB struct {
	Users []*models.User
}

func (d *mockUserDB) GetUser(
	username string,
) (*models.User, error) {
	for _, user := range d.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}

type mockTaskDB struct {
	Tasks []*models.Task
}

func (d *mockTaskDB) GetTasks(userID uint64) []*models.Task {
	var tasks []*models.Task

	for _, task := range d.Tasks {
		if task.UserID == userID {
			tasks = append(tasks, task)
		}
	}

	return d.Tasks
}

func (d *mockTaskDB) GetTask(
	userID uint64,
	taskID uint64,
) (*models.Task, error) {
	for _, task := range d.Tasks {
		if task.UserID == userID && task.ID == taskID {
			return task, nil
		}
	}
	return nil, errors.New("not found")
}

func (d *mockTaskDB) CreateTask(
	userID uint64,
	t *models.NewTaskTemplate,
) *models.Task {
	newTask := &models.Task{
		ID:          d.Tasks[len(d.Tasks)-1].ID + 1,
		Title:       t.Title,
		Description: t.Description,
		Done:        false,
		UserID:      userID,
	}

	d.Tasks = append(d.Tasks, newTask)

	return newTask
}

func (d *mockTaskDB) UpdateTask(t *models.Task) error {
	return nil
}

func (d *mockTaskDB) DeleteTask(userID uint64, taskID uint64) error {
	for index, task := range d.Tasks {
		if task.UserID == userID && task.ID == taskID {
			d.Tasks = append(d.Tasks[:index], d.Tasks[index+1:]...)
			return nil
		}
	}

	return errors.New("not found")
}

func Setup() *gin.Engine {
	user := models.User{
		ID:       1,
		Username: "john",
		Password: "$2a$12$udogIRFurk7EMHfALwSZZexk4K8salP3n7/bEV8pr8PAJ1Fztxcdq",
	}

	userDB := &mockUserDB{[]*models.User{&user}}
	taskDB := &mockTaskDB{
		[]*models.Task{
			{
				ID:          1,
				Title:       "Buy groceries",
				Description: "Milk, Cheese, Pizza, Fruit, Tylenol",
				Done:        false,
				UserID:      user.ID,
			},
			{
				ID:          2,
				Title:       "Learn Go",
				Description: "Need to find a good Go tutorial on the web",
				Done:        false,
				UserID:      user.ID,
			},
		},
	}
	tokenizer := auth.NewTokenizer()
	authClient := &AuthClientMock{}

	r := gin.Default()
	ac := controllers.NewAuthController(userDB, tokenizer, authClient)
	tc := controllers.NewTaskController(taskDB, authClient)

	routers.SetAuthRoutes(r, ac)
	routers.SetTaskRoutes(r, tc)

	return r
}

func Login(router *gin.Engine, w *httptest.ResponseRecorder) (string, error) {
	jsonStr, _ := json.Marshal(
		&models.UserLoginTemplate{
			Username: "john",
			Password: "password",
		},
	)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var data TokenJSON
	err := json.NewDecoder(w.Body).Decode(&data)
	if err != nil {
		return "", errors.New("unauthorized")
	}

	return data.Tokens["access_token"], nil
}

type AuthClientMock struct {
	userID uint64
}

func (a *AuthClientMock) Store(userID uint64, td *auth.TokenDetails) error {
	a.userID = userID
	return nil
}

func (a *AuthClientMock) Fetch(r *http.Request) (uint64, error) {
	if a.userID == 0 {
		return 0, errors.New("unauthorized")
	}

	return a.userID, nil
}

func (a *AuthClientMock) Delete(r *http.Request) (uint64, error) {
	if a.userID == 0 {
		return 0, errors.New("unauthorized")
	}

	deleted := a.userID
	a.userID = 0

	return deleted, nil
}

type TokenJSON struct {
	Tokens map[string]string `json:"tokens"`
}

type ResultJSON struct {
	Result bool `json:"result"`
}

type ErrorJSON struct {
	Error string `json:"error"`
}
