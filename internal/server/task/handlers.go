package task

import (
	"encoding/json"
	"fmt"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	service Task
}

type UserTask struct {
	UserID *uuid.UUID `json:"user_id"`
	Text   *string    `json:"text"`
}

type Timer struct {
	TaskID uuid.UUID `json:"task_id"`
}

type Task interface {
	CreateTask(task UserTask) (models.Task, error)
	StartTask(taskID uuid.UUID) error
	StopTask(taskID uuid.UUID) error
	GetTasks(request models.LaborTimeRequest) (models.LaborTimeResponse, error)
}

func NewHandler(service Task) *Handler {
	return &Handler{
		service: service,
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

func (rs *Handler) new(w http.ResponseWriter, r *http.Request) {
	t := UserTask{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil || t.UserID == nil || t.Text == nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newTask, err := rs.service.CreateTask(t)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(newTask)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (rs *Handler) get(w http.ResponseWriter, r *http.Request) {
	p := models.LaborTimeRequest{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil || p.UserID == nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks, err := rs.service.GetTasks(p)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (rs *Handler) start(w http.ResponseWriter, r *http.Request) {
	timer := Timer{}
	err := json.NewDecoder(r.Body).Decode(&timer)
	if err != nil || timer.TaskID == uuid.Nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rs.service.StartTask(timer.TaskID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rs *Handler) stop(w http.ResponseWriter, r *http.Request) {
	timer := Timer{}
	err := json.NewDecoder(r.Body).Decode(&timer)
	if err != nil || timer.TaskID == uuid.Nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rs.service.StopTask(timer.TaskID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
