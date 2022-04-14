package gorm

import (
	"github.com/ichigozero/gtdzero/models"
	libgorm "gorm.io/gorm"
)

type taskDB struct {
	db *libgorm.DB
}

func NewTaskDB(db *libgorm.DB) models.TaskDB {
	return &taskDB{db}
}

func (td *taskDB) GetTasks(userID uint64) []*models.Task {
	var tasks []*models.Task
	td.db.Where("user_id = ?", userID).Find(&tasks)

	return tasks
}

func (td *taskDB) GetTask(userID uint64, taskID uint64) (*models.Task, error) {
	var task *models.Task
	result := td.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task)

	return task, result.Error
}

func (td *taskDB) CreateTask(userID uint64, t *models.NewTaskTemplate) *models.Task {
	task := models.Task{Title: t.Title, Description: t.Description, Done: false, UserID: userID}
	td.db.Create(&task)

	return &task
}

func (td *taskDB) UpdateTask(t *models.Task) error {
	task := models.Task{ID: t.ID}
	result := td.db.Model(&task).Updates(
		map[string]interface{}{
			"title":       t.Title,
			"description": t.Description,
			"done":        t.Done,
			"user_id":     t.UserID,
		})

	return result.Error
}

func (td *taskDB) DeleteTask(userID uint64, taskID uint64) error {
	task := models.Task{ID: taskID}
	result := td.db.Where("user_id", userID).Delete(&task)

	return result.Error
}
