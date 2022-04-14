package models

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type UserLoginTemplate struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDB interface {
	GetUser(username string) (*User, error)
}
