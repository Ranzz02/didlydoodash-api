package middleware

import (
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/db/models"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract access token
		tokenString := jwt.ExtractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.NotAuthenticated)
			return
		}

		token, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.InvalidToken)
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.NotAuthenticated)
			return
		}
		models.CurrentUser = &sub

		c.Next()
	}
}
