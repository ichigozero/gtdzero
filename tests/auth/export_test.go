package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
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

	ac := controllers.NewAuthController(db)
	routers.SetAuthRoutes(r, ac)

	return r
}

type tokenJSON struct {
	Token string `json:"token"`
}
