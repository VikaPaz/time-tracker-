package app

import (
	client "github.com/VikaPaz/time_tracker/internal/clients"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/VikaPaz/time_tracker/internal/repository"
	"github.com/VikaPaz/time_tracker/internal/repository/task"
	"github.com/VikaPaz/time_tracker/internal/repository/user"
	"github.com/VikaPaz/time_tracker/internal/server"
	taskService "github.com/VikaPaz/time_tracker/internal/service/task"
	userService "github.com/VikaPaz/time_tracker/internal/service/user"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Run(logger *logrus.Logger) error {
	if err := godotenv.Overload(); err != nil {
		logger.Errorf("Error loading .env file")
		return models.ErrLoadEnvFailed
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
		logger.Errorf("Error connecting to database")
		return err
	}
	logger.Infof("Connected to PostgreSQL")

	userRepo := user.NewRepository(dbConn, logger)
	taskRepo := task.NewRepository(dbConn, logger)

	userInf, err := client.NewClient(os.Getenv("INFO_SERVER"), logger)
	if err != nil {
		return err
	}

	userService := userService.NewService(userRepo, userInf, logger)
	taskService := taskService.NewService(taskRepo, logger)

	srv := server.NewServer(userService, taskService, logger)

	logger.Infof("Running server on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), srv.Handlers())
	if err != nil {
		logger.Errorf("Error starting server")
		return models.ErrServerFailed
	}

	return err
}
