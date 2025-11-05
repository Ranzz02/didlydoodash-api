package utils

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

// Get user_id safely from Gin context
func GetUserID(c *gin.Context) string {
	val, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return val.(string)
}

func GetUserIDFromContext(ctx context.Context) string {
	val := ctx.Value(userIDKey)
	if id, ok := val.(string); ok {
		return id
	}
	return ""
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// ParseIntDefault safely parses an int with fallback
func ParseIntDefault(s string, fallback int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return fallback
}

func ParseBoolDefault(s string, fallback bool) bool {
	if n, err := strconv.ParseBool(s); err == nil {
		return n
	}
	return fallback
}
