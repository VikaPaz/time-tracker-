package task

import (
	"encoding/json"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Handler struct {
	service Task
	log     *logrus.Logger
}

type Task interface {
	CreateTask(task models.UserTask) (models.Task, error)
	StartTask(taskID uuid.UUID) error
	StopTask(taskID uuid.UUID) error
	GetTasks(request models.LaborTimeRequest) (models.GetTaskResponse, error)
}

func NewHandler(service Task, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     logger,
	}
}

func (rs *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/new", rs.new)
	r.Get("/get", rs.get)
	r.Patch("/start", rs.start)
	r.Patch("/stop", rs.stop)

	return r
}

// @Summary Create new task
// @Description Handles request to create a new task for a user.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body models.UserTask true "User ID and text"
// @Success 200 {object} models.Task "Details of the newly created task"
// @Failure 400
// @Failure 500
// @Router /task/new [post]
func (rs *Handler) new(w http.ResponseWriter, r *http.Request) {
	t := models.UserTask{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil || t.UserID == nil || t.Text == nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rs.log.Infof("Creating new task")
	newTask, err := rs.service.CreateTask(t)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(newTask)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// @Summary Get tasks
// @Description Handles request to get tasks for a user based on user ID.
// @Tags tasks
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param start_time query string false "Start Time"
// @Param end_time query string false "End Time"
// @Success 200 {array} models.GetTaskResponse "User ID, list of tasks and total"
// @Failure 400
// @Failure 500
// @Router /task/get [get]
func (rs *Handler) get(w http.ResponseWriter, r *http.Request) {
	p := models.LaborTimeRequest{}

	// Получаем параметры запроса из URL
	params := r.URL.Query()

	// Извлекаем и устанавливаем UserID
	if userIDStr := params.Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid UUID for user_id", http.StatusBadRequest)
			return
		}
		p.UserID = &userID
	} else {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	// Извлекаем и устанавливаем StartTime
	if startTimeStr := params.Get("start_time"); startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time parameter (RFC3339 format expected)", http.StatusBadRequest)
			return
		}
		p.StartTime = &startTime
	}

	// Извлекаем и устанавливаем EndTime
	if endTimeStr := params.Get("end_time"); endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time parameter (RFC3339 format expected)", http.StatusBadRequest)
			return
		}
		p.EndTime = &endTime
	}

	rs.log.Infof("Getting tasks")
	tasks, err := rs.service.GetTasks(p)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary Start timer for task
// @Description Handles request to start a timer for a task.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body models.Timer true "Task ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /task/start [patch]
func (rs *Handler) start(w http.ResponseWriter, r *http.Request) {
	timer := models.Timer{}
	err := json.NewDecoder(r.Body).Decode(&timer)
	if err != nil || timer.TaskID == uuid.Nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rs.log.Infof("Starting timer")
	err = rs.service.StartTask(timer.TaskID)
	if err != nil {
		rs.log.Error(err)
		if err == models.ErrTimerStarted {
			w.WriteHeader(http.StatusBadRequest)
		}
		if err == models.ErrStartTimer {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

// @Summary Stop timer for task
// @Description Handles request to stop a timer for a task.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body models.Timer true "Task ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /task/stop [patch]
func (rs *Handler) stop(w http.ResponseWriter, r *http.Request) {
	timer := models.Timer{}
	err := json.NewDecoder(r.Body).Decode(&timer)
	if err != nil || timer.TaskID == uuid.Nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rs.log.Infof("Stopping timer")
	err = rs.service.StopTask(timer.TaskID)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
