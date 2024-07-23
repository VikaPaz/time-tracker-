package task

import (
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/VikaPaz/time_tracker/internal/server/task"
	"github.com/google/uuid"
)

type TaskService struct {
	repo Repository
}

type Repository interface {
	Create(task models.Task) (models.Task, error)
	Get(request models.LaborTimeRequest) (models.LaborTimeResponse, error)
	Start(taskID uuid.UUID) error
	Stop(taskID uuid.UUID) error
	IsStarted(taskID uuid.UUID) (bool, error)
}

func NewService(repo Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (t *TaskService) CreateTask(userTask task.UserTask) (models.Task, error) {
	newTask := models.Task{
		Task:   *userTask.Text,
		UserID: *userTask.UserID,
	}

	result, err := t.repo.Create(newTask)
	if err != nil {
		return models.Task{}, err
	}

	return result, nil
}

func (t *TaskService) StartTask(taskID uuid.UUID) error {
	isStarted, err := t.repo.IsStarted(taskID)
	if err != nil {
		return err
	}
	if isStarted {
		return nil
	}

	err = t.repo.Start(taskID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) StopTask(taskID uuid.UUID) error {
	err := t.repo.Stop(taskID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) GetTasks(request models.LaborTimeRequest) (models.LaborTimeResponse, error) {
	result, err := t.repo.Get(request)
	if err != nil {
		return models.LaborTimeResponse{}, err
	}

	return result, nil
}
