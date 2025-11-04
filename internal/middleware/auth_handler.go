package middleware

import (
	"errors"
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.EnvConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract access token
		tokenString := utils.ExtractToken(c)
		if tokenString == "" {
			c.Error(utils.NewError(http.StatusUnauthorized, "no token provided", errors.New("no token provided in request")))
			c.Abort()
			return
		}

		token, err := utils.ValidateToken(cfg, tokenString, utils.AccessToken)
		if err != nil {
			c.Error(utils.NewError(http.StatusUnauthorized, "invalid token provided", err))
			c.Abort()
			return
		}

		sub, err := token.GetSubject()
		if err != nil {
			c.Error(utils.NewError(http.StatusUnauthorized, "user id not found in token", err))
			c.Abort()
			return
		}

		// Attach user ID to context
		c.Set("user_id", sub)

		// Enrich logger with user info for downstream logs
		logger := logging.GetLogger(c).WithField("user_id", sub)
		logging.WithLogger(c, logger)
		c.Request = c.Request.WithContext(logging.WithContextLogger(c.Request.Context(), logger))

		c.Next()
	}
}
