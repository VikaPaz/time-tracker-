package task

import (
	"fmt"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TaskService struct {
	repo Repository
	log  *logrus.Logger
}

type Repository interface {
	Create(task models.Task) (models.Task, error)
	Get(request models.LaborTimeRequest) (models.LaborTimeResponse, error)
	Start(taskID uuid.UUID) error
	Stop(taskID uuid.UUID) error
	IsStarted(taskID uuid.UUID) (bool, error)
}

func NewService(repo Repository, logger *logrus.Logger) *TaskService {
	return &TaskService{
		repo: repo,
		log:  logger,
	}
}

func (t *TaskService) CreateTask(userTask models.UserTask) (models.Task, error) {
	newTask := models.Task{
		Task:   *userTask.Text,
		UserID: *userTask.UserID,
	}

	t.log.Debugf("Creating task for user: %s", *userTask.UserID)
	result, err := t.repo.Create(newTask)
	if err != nil {
		return models.Task{}, err
	}

	return result, nil
}

func (t *TaskService) StartTask(taskID uuid.UUID) error {
	t.log.Debugf("Checking status timer with task ID %v", taskID)
	isStarted, err := t.repo.IsStarted(taskID)
	if err != nil {
		return err
	}
	if isStarted {
		return nil
	}

	t.log.Debugf("Starting timer with task ID %s", taskID)
	err = t.repo.Start(taskID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) StopTask(taskID uuid.UUID) error {
	t.log.Debugf("Stopping timer with task ID %v", taskID)
	err := t.repo.Stop(taskID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) GetTasks(request models.LaborTimeRequest) (models.GetTaskResponse, error) {
	t.log.Debugf("Getting tasks with user ID: %v", request.UserID)
	result, err := t.repo.Get(request)
	if err != nil {
		return models.GetTaskResponse{}, err
	}

	response := models.GetTaskResponse{
		UserID: result.UserID,
		Tasks:  []models.GetTaskInfo{},
	}
	for _, v := range result.Tasks {
		labor := models.GetTaskInfo{
			ID:   v.ID,
			Task: *v.Task,
		}
		h := *v.LaborTime / 360
		m := *v.LaborTime/60 - h
		labor.LaborTime = fmt.Sprintf("hours: %v minutes: %v", h, m)
		response.Tasks = append(response.Tasks, labor)
	}

	return response, nil
}
