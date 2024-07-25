package server

import (
	_ "github.com/VikaPaz/time_tracker/docs"
	taskHandler "github.com/VikaPaz/time_tracker/internal/server/task"
	userHandler "github.com/VikaPaz/time_tracker/internal/server/user"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type ImplServer struct {
	user userHandler.User
	task taskHandler.Task
	log  *logrus.Logger
}

func NewServer(user userHandler.User, task taskHandler.Task, logger *logrus.Logger) *ImplServer {
	return &ImplServer{
		user: user,
		task: task,
		log:  logger,
	}
}

func (i *ImplServer) Handlers() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	u := userHandler.NewHandler(i.user, i.log)
	t := taskHandler.NewHandler(i.task, i.log)

	r.Mount("/user", u.Router())
	r.Mount("/task", t.Router())

	return r
}
