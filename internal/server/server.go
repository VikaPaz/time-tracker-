package server

import (
	taskHandler "github.com/VikaPaz/time_tracker/internal/server/task"
	userHandler "github.com/VikaPaz/time_tracker/internal/server/user"
	"github.com/go-chi/chi/v5"
)

type ImplServer struct {
	user userHandler.User
	task taskHandler.Task
}

func NewServer(user userHandler.User, task taskHandler.Task) *ImplServer {
	return &ImplServer{
		user: user,
		task: task,
	}
}

func (i *ImplServer) Handlers() *chi.Mux {
	r := chi.NewRouter()

	u := userHandler.NewHandler(i.user)
	t := taskHandler.NewHandler(i.task)

	r.Mount("/user", u.Router())
	r.Mount("/task", t.Router())

	return r
}
