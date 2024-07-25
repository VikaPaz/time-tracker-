package main

import (
	"github.com/VikaPaz/time_tracker/internal/app"
	"github.com/sirupsen/logrus"
)

// @title Time Tracker API
// @description This is time_tracker server.
// @host localhost:8000
func main() {
	logger := NewLogger(logrus.DebugLevel, &logrus.TextFormatter{
		FullTimestamp: true,
	})

	err := app.Run(logger)
	if err != nil {
		logger.Fatalln(err)
	}
}

func NewLogger(level logrus.Level, formatter logrus.Formatter) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	return logger
}
