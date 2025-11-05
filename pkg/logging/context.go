package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey struct{}

var loggerKey = contextKey{}

// WithLogger attaches a *logrus.Entry to a context.
func WithContextLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the *logrus.Entry from context or returns a default logger.
func FromContext(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	if l, ok := ctx.Value(loggerKey).(*logrus.Entry); ok {
		return l
	}
	return logrus.NewEntry(logrus.StandardLogger())
}

// WithLayer returns a derived logger from context with an added "layer" field.
func WithLayer(ctx context.Context, layer, component string) *logrus.Entry {
	return FromContext(ctx).WithFields(logrus.Fields{
		"layer":     layer,
		"component": component,
	})
}
