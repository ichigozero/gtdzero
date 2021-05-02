package models

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type UserDB interface {
	GetUser(username string, password string) (*User, error)
}
