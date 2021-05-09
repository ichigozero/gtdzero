package models

type Task struct {
	ID          uint64 `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	UserID      uint64 `json:"userid"`
	User        User   `json:"user"`
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
	GetTasks(userID uint64) []*Task
	GetTask(userID uint64, taskID uint64) (*Task, error)
	CreateTask(t *NewTaskTemplate) *Task
	UpdateTask(t *Task) error
	DeleteTask(id uint64) error
}
