package models

type Task struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type NewTaskTemplate struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type TaskDB interface {
	GetTasks() []*Task
	GetTask(id int) (*Task, error)
	CreateTask(t *NewTaskTemplate) *Task
}
