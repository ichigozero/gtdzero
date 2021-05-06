package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
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

type mockRedis struct{}

func (r *mockRedis) Del(keys ...string) *redis.IntCmd {
	return redis.NewIntCmd()
}

func (r *mockRedis) Set(
	key string,
	value interface{},
	expiration time.Duration,
) *redis.StatusCmd {
	return redis.NewStatusCmd()
}

func (r *mockRedis) Get(key string) *redis.StringCmd {
	return redis.NewStringCmd()
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
	client := auth.NewAuthClient(&mockRedis{})

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
