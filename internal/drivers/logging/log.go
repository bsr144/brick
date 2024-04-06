package logging

import (
	"brick/internal/config"

	"github.com/sirupsen/logrus"
)

func NewLogger(loggingConfig config.Logging) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(loggingConfig.Level))

	switch loggingConfig.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{})
	}

	return log
}
