package models

type Task struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TaskDB interface {
	GetTasks() []*Task
}
