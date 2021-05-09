package auth

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

type mockDB struct {
	Users []*models.User
}

func (d *mockDB) GetUser(
	username string,
	password string,
) (*models.User, error) {
	for _, user := range d.Users {
		if user.Username == username && user.Password == password {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}

type authClientMock struct {
	userID uint64
}

func (a *authClientMock) Store(userID uint64, td *auth.TokenDetails) error {
	a.userID = userID
	return nil
}

func (a *authClientMock) Fetch(r *http.Request) (uint64, error) {
	return a.userID, nil
}

func (a *authClientMock) Delete(r *http.Request) (uint64, error) {
	if a.userID == 0 {
		return a.userID, errors.New("")
	}

	deleted := a.userID
	a.userID = 0

	return deleted, nil
}

func setUp() *gin.Engine {
	r := gin.Default()

	db := &mockDB{
		[]*models.User{
			{
				ID:       1,
				Username: "john",
				Password: "password",
			},
		},
	}
	tokenizer := auth.NewTokenizer()
	client := &authClientMock{}

	ac := controllers.NewAuthController(db, tokenizer, client)
	routers.SetAuthRoutes(r, ac)

	return r
}

type tokenJSON struct {
	Tokens map[string]string `json:"tokens"`
}

type resultJSON struct {
	Result bool `json:"result"`
}

type errorJSON struct {
	Error string `json:"error"`
}

func login(router *gin.Engine, w *httptest.ResponseRecorder) {
	jsonStr, _ := json.Marshal(
		&models.UserLoginTemplate{
			Username: "john",
			Password: "password",
		},
	)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
}
