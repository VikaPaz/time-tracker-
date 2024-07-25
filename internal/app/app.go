package app

import (
	"database/sql"
	client "github.com/VikaPaz/time_tracker/internal/clients"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/VikaPaz/time_tracker/internal/repository"
	"github.com/VikaPaz/time_tracker/internal/repository/task"
	"github.com/VikaPaz/time_tracker/internal/repository/user"
	"github.com/VikaPaz/time_tracker/internal/server"
	taskService "github.com/VikaPaz/time_tracker/internal/service/task"
	userService "github.com/VikaPaz/time_tracker/internal/service/user"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
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

	err = runMigrations(logger, dbConn)
	if err != nil {
		logger.Errorf("can't run migrations")
		return err
	}

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

func runMigrations(logger *logrus.Logger, dbConn *sql.DB) error {
	upMigration, err := strconv.ParseBool(os.Getenv("RUN_MIGRATION"))
	if err != nil {
		return err
	}

	if !upMigration {
		return nil
	}

	migrationDir := os.Getenv("MIGRATION_DIR")
	if migrationDir == "" {
		logger.Infof("no migration dir provided; skipping migrations")
		return nil
	}
	err = goose.Up(dbConn, os.Getenv("MIGRATION_DIR"))
	if err != nil {
		return err
	}
	logger.Infof("migrations are applied successfully")

	return nil
}
