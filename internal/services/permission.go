package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/internal/repositories"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/permissions"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Checker struct {
	memberRepo *repositories.MemberRepo
	roleRepo   *repositories.RoleRepo
	logger     *logrus.Logger
}

func NewChecker(memberRepo *repositories.MemberRepo, roleRepo *repositories.RoleRepo, logger *logrus.Logger) *Checker {
	return &Checker{
		memberRepo: memberRepo,
		roleRepo:   roleRepo,
		logger:     logger,
	}
}

// Core permission check
func (c *Checker) Check(ctx context.Context, userID, orgID string, perm permissions.Permission) error {
	logger := logging.WithLayer(ctx, "service", "checker").WithFields(logrus.Fields{
		"user_id": userID,
		"org_id":  orgID,
		"perm":    perm,
	})

	// Short-circuit if organisation owner
	isOwner, err := c.memberRepo.IsOwner(ctx, orgID, userID)
	if err != nil {
		logger.WithError(err).Error("failed to check organisation ownership")
		return utils.NewError(http.StatusInternalServerError, "failed to verify organisation ownership", err)
	}
	if isOwner {
		return nil
	}

	// Otherwise, check role permissions
	ok, err := c.roleRepo.HasPermission(ctx, &repository.HasPermissionParams{
		UserID:         userID,
		OrganisationID: orgID,
		PermissionKey:  string(perm),
	})
	if err != nil {
		logger.WithError(err).Error("permission check failed")
		return utils.NewError(http.StatusInternalServerError, "failed to verify permission", err)
	}
	if !ok {
		logger.Warn("permission denied")
		return utils.NewError(http.StatusForbidden, "insufficient permissions", fmt.Errorf("missing permission: %s", perm))
	}

	return nil
}
