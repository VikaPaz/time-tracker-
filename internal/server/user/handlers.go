package user

import (
	"context"
	"encoding/json"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type User interface {
	CreateUser(person models.CreateUserRequest, ctx context.Context) (models.User, error)
	DeleteUser(request models.DeleteUserRequest) error
	ChangeUser(user models.User) error
	GetUsers(request models.FilterRequest) (models.FilterResponse, error)
}

type Handler struct {
	service User
	log     *logrus.Logger
}

func NewHandler(service User, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     logger,
	}
}

func (rs *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/new", rs.new)
	r.Delete("/delete", rs.del)
	r.Get("/get", rs.get)
	r.Patch("/set", rs.change)

	return r
}

// @Summary Creating a new user
// @Description Handles request to create a new user by passportNumber and returns the user information in JSON.
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "Passport"
// @Success 200 {object} models.User "Created user"
// @Failure 400
// @Failure 500
// @Router /user/new [post]
func (rs *Handler) new(w http.ResponseWriter, r *http.Request) {
	p := models.CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var newUser models.User

	rs.log.Infof("Creating new user")
	newUser, err = rs.service.CreateUser(p, ctx)
	if err != nil {
		rs.log.Error(err)
		if err == models.ErrGetUserResponse {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err == models.ErrUserExists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err == models.ErrInvalidPassword {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err == models.ErrCreateUserResponse {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(newUser)
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

// @Summary Delete user
// @Description Handles request to delete a user by ID.
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.DeleteUserRequest true "User ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /user/delete [delete]
func (rs *Handler) del(w http.ResponseWriter, r *http.Request) {
	request := models.DeleteUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rs.log.Infof("Deleting user: %v", request.ID)
	err = rs.service.DeleteUser(request)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// @Summary Get users
// @Description Handles request to get users by filter.
// @Tags users
// @Accept json
// @Produce json
// @Param fields.id query string false "User ID"
// @Param fields.passport query string false "User Passport"
// @Param fields.name query string false "Username"
// @Param fields.surname query string false "User Surname"
// @Param fields.patronymic query string false "User Patronymic"
// @Param fields.address query string false "User Address"
// @Param limit query int false "Maximum number of results"
// @Param offset query int false "Offset from the beginning of results"
// @Success 200 {object}  models.FilterResponse "List of users and total results"
// @Failure 400
// @Failure 500
// @Router /user/get [get]
func (rs *Handler) get(w http.ResponseWriter, r *http.Request) {
	filter := models.FilterRequest{}
	params := r.URL.Query()

	if idStr := params.Get("fields.id"); idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid UUID for fields.id", http.StatusBadRequest)
			return
		}
		filter.Fields.ID = &id
	}
	if passport := params.Get("fields.passport"); passport != "" {
		filter.Fields.Passport = &passport
	}
	if name := params.Get("fields.name"); name != "" {
		filter.Fields.Name = &name
	}
	if surname := params.Get("fields.surname"); surname != "" {
		filter.Fields.Surname = &surname
	}
	if patronymic := params.Get("fields.patronymic"); patronymic != "" {
		filter.Fields.Patronymic = &patronymic
	}
	if address := params.Get("fields.address"); address != "" {
		filter.Fields.Address = &address
	}

	if limitStr := params.Get("limit"); limitStr != "" {
		limit, err := strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		filter.Limit = limit
	}
	if offsetStr := params.Get("offset"); offsetStr != "" {
		offset, err := strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
		filter.Offset = offset
	}

	rs.log.Infof("Getting users")
	users, err := rs.service.GetUsers(filter)
	if err != nil {
		rs.log.Error(err)
		return
	}

	data, err := json.Marshal(users)
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

// @Summary Update user
// @Description Handles request to update user information
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.User true "User Information"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /user/set [patch]
func (rs *Handler) change(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rs.log.Infof("Changing user information")
	err = rs.service.ChangeUser(user)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
