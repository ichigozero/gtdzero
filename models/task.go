package models

type Task struct {
	ID          uint64 `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type NewTaskTemplate struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskTemplate struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TaskDB interface {
	GetTasks() []*Task
	GetTask(id uint64) (*Task, error)
	CreateTask(t *NewTaskTemplate) *Task
	UpdateTask(t *Task) error
	DeleteTask(id uint64) error
}
