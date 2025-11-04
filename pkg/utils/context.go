package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Get user_id safely from Gin context
func GetUserID(c *gin.Context) string {
	val, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return val.(string)
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
