package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type User interface {
	CreateUser(person Person, ctx context.Context) (models.User, error)
	DeleteUser(request models.DeleteUserRequest) error
	ChangeUser(user models.User) error
	GetUsers(request models.FilterRequest) (models.FilterResponse, error)
}

type Handler struct {
	service User
}

type Person struct {
	PassportNumber *string `json:"passportNumber,omitempty"`
}

func NewHandler(service User) *Handler {
	return &Handler{
		service: service,
	}
}

func (rs *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/new", rs.new)
	r.Delete("/delete", rs.del)
	r.Get("/get", rs.get)
	r.Patch("/change", rs.change)

	return r
}

func (rs *Handler) new(w http.ResponseWriter, r *http.Request) {
	p := Person{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var newUser models.User
	newUser, err = rs.service.CreateUser(p, ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(newUser)
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

func (rs *Handler) del(w http.ResponseWriter, r *http.Request) {
	request := models.DeleteUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = rs.service.DeleteUser(request)
	if err != nil {
		fmt.Println(err)
	}
}

func (rs *Handler) get(w http.ResponseWriter, r *http.Request) {
	filter := models.FilterRequest{}
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	users, err := rs.service.GetUsers(filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(users)
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

func (rs *Handler) change(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = rs.service.ChangeUser(user)
	if err != nil {
		fmt.Println(err)
		return
	}
}
