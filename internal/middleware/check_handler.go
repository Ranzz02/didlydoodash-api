package middleware

import (
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/services"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/permissions"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequirePermission(checker *services.Checker, perm permissions.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userID := utils.GetUserID(c)
		orgID := c.Param("id")

		logger := logging.WithLayer(ctx, "middleware", "permissions").WithFields(logrus.Fields{
			"user_id": userID,
			"org_id":  orgID,
			"perm":    perm,
		})

		if orgID == "" {
			logger.Warn("missing organisation ID in route")
			c.Error(utils.NewError(http.StatusBadRequest, "missing organisation id in route", nil))
			c.Abort()
			return
		}

		logger.Infof("checking permission: %s", perm)
		if err := checker.Check(ctx, userID, orgID, perm); err != nil {
			logger.WithError(err).Warn("permission denied")
			c.Error(err)
			c.Abort()
			return
		}

		logger.Info("permission granted")
		c.Next()
	}
}
