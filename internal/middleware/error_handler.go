package middleware

import (
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execure handler first

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		// Respond with generic JSON if none written yet
		if !c.Writer.Written() {
			if apiErr, ok := err.(utils.APIError); ok {
				c.JSON(apiErr.Code, apiErr)
			} else {
				c.JSON(http.StatusInternalServerError, utils.NewError(http.StatusInternalServerError, "unknown internal error", err))
			}
		}
	}
}
