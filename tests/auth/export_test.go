package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/ichigozero/gtdzero/controllers"
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

func (r *mockRedis) Set(
	key string,
	value interface{},
	expiration time.Duration,
) *redis.StatusCmd {
	return redis.NewStatusCmd()
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
	rc := &mockRedis{}

	ac := controllers.NewAuthController(db, rc)
	routers.SetAuthRoutes(r, ac)

	return r
}

type tokenJSON struct {
	Tokens map[string]string `json:"tokens"`
}
