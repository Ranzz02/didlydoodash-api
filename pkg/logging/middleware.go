package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
)

// Middleware returns a Gin Middleware that attaches a contextual logger per request
func Middleware(base *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate or extract request ID
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID, _ = gonanoid.New()
		}

		// Create a request-scoped logger
		reqLogger := base.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
		})

		// Attach both contexts
		WithLogger(c, reqLogger)
		c.Request = c.Request.WithContext(WithContextLogger(c.Request.Context(), reqLogger))

		// Process the request
		c.Next()

		// After response - log status + latency
		status := c.Writer.Status()
		latency := time.Since(start)

		entry := reqLogger.WithFields(logrus.Fields{
			"status":  status,
			"latency": latency.String(),
		})

		switch {
		case status >= 500:
			entry.Error("http server error")
		case status >= 400:
			entry.Warn("http user error")
		default:
			entry.Info("http success")
		}
	}
}
