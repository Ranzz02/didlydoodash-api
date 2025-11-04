package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const ginLoggerKey = "logger"

// WithLogger attaches a logrus.Entry to the context
func WithLogger(c *gin.Context, logger *logrus.Entry) {
	c.Set(ginLoggerKey, logger)
}

// GetLogger extracts a *logrus.Entry from context, or returns a default one
func GetLogger(c *gin.Context) *logrus.Entry {
	if l, exists := c.Get(ginLoggerKey); exists {
		if logger, ok := l.(*logrus.Entry); ok {
			return logger
		}
	}
	// Fallback: return a new generic logger if missing
	return logrus.NewEntry(logrus.StandardLogger())
}
