package task

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/ichigozero/gtdzero/controllers"
	"github.com/ichigozero/gtdzero/models"
	"github.com/ichigozero/gtdzero/routers"
)

type mockDB struct {
	Tasks []*models.Task
}

func (d *mockDB) GetTasks() []*models.Task {
	return d.Tasks
}

func (d *mockDB) GetTask(id uint64) (*models.Task, error) {
	for _, task := range d.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, errors.New("not found")
}

func (d *mockDB) CreateTask(t *models.NewTaskTemplate) *models.Task {
	newTask := &models.Task{
		ID:          d.Tasks[len(d.Tasks)-1].ID + 1,
		Title:       t.Title,
		Description: t.Description,
		Done:        false,
	}

	d.Tasks = append(d.Tasks, newTask)

	return newTask
}

func (d *mockDB) UpdateTask(t *models.Task) error {
	return nil
}

func (d *mockDB) DeleteTask(id uint64) error {
	for index, task := range d.Tasks {
		if task.ID == id {
			d.Tasks = append(d.Tasks[:index], d.Tasks[index+1:]...)
			return nil
		}
	}

	return errors.New("not found")
}

func setUp() *gin.Engine {
	r := gin.Default()
	db := &mockDB{
		[]*models.Task{
			{
				ID:          1,
				Title:       "Buy groceries",
				Description: "Milk, Cheese, Pizza, Fruit, Tylenol",
				Done:        false,
			},
			{
				ID:          2,
				Title:       "Learn Go",
				Description: "Need to find a good Go tutorial on the web",
				Done:        false,
			},
		},
	}

	tc := controllers.NewTaskController(db)
	routers.SetTaskRoutes(r, tc)

	return r
}
