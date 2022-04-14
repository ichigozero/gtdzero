package gorm

import (
	"github.com/ichigozero/gtdzero/models"
	libgorm "gorm.io/gorm"
)

type userDB struct {
	db *libgorm.DB
}

func NewUserDB(db *libgorm.DB) models.UserDB {
	return &userDB{db}
}

func (u *userDB) GetUser(username string, password string) (*models.User, error) {
	var user models.User
	result := u.db.Where("username = ? AND password = ?", username, password).First(&user)

	return &user, result.Error
}
