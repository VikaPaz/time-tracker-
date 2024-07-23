package app

import (
	client "github.com/VikaPaz/time_tracker/internal/clients"
	"github.com/VikaPaz/time_tracker/internal/repository"
	"github.com/VikaPaz/time_tracker/internal/repository/task"
	"github.com/VikaPaz/time_tracker/internal/repository/user"
	"github.com/VikaPaz/time_tracker/internal/server"
	taskService "github.com/VikaPaz/time_tracker/internal/service/task"
	userService "github.com/VikaPaz/time_tracker/internal/service/user"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func Run() error {
	if err := godotenv.Overload(); err != nil {
		log.Print("No .env file found")
	}

	confPostgres := repository.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	dbConn, err := repository.Connection(confPostgres)
	if err != nil {
		return err
	}
	userRepo := user.NewRepository(dbConn)
	taskRepo := task.NewRepository(dbConn)

	userInf := client.NewClient("http://127.0.0.1:8080")

	userService := userService.NewService(userRepo, userInf)
	taskService := taskService.NewService(taskRepo)

	srv := server.NewServer(userService, taskService)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), srv.Handlers())
	//err = http.ListenAndServe(":8080", srv.Handlers())

	return err
}
