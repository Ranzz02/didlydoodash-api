package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New creates and configures a base *logrus.Logger instance.
func New(env string) *logrus.Logger {
	logger := logrus.New()
	logger.Out = os.Stdout

	// Choose formatter based on environment
	if env == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			PadLevelText:    true,
		})
		logger.SetReportCaller(true) // shows file:line in dev mode
		logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
