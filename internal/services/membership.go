package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/repositories"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/permissions"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type MembershipRepos struct {
	Role   *repositories.RoleRepo
	Member *repositories.MemberRepo
	User   *repositories.UserRepository
}

type MembershipService struct {
	repos  *MembershipRepos
	tx     *repositories.TxManager
	logger *logrus.Logger
}

func NewMembershipService(repos MembershipRepos, tx *repositories.TxManager, logger *logrus.Logger) *MembershipService {
	return &MembershipService{
		repos:  &repos,
		tx:     tx,
		logger: logger,
	}
}

func (s *MembershipService) Create(ctx context.Context, params *dto.CreateOrganisationMember) (*dto.OrganisationMember, error) {
	var err error
	logger := logging.WithLayer(ctx, "service", "membership").WithFields(logrus.Fields{
		"org_id":  params.OrgID,
		"user_id": params.UserID,
	})

	// Extract acting user from context (set by auth middleware)
	actorID := utils.GetUserIDFromContext(ctx)

	// Check that the actor has permission to add members
	if err = s.CheckPermission(ctx, actorID, params.OrgID, permissions.OrgInviteMembers); err != nil {
		return nil, err
	}

	logger.Info("attempting to add member to organisation")

	var member repository.OrganisationMember
	var user repository.User
	var role repository.Role

	err = s.tx.WithTx(ctx, func(q repository.Querier) error {
		// Check that user exists
		user, err = s.repos.User.GetByID(ctx, params.UserID)
		if err != nil {
			logger.WithError(err).Warn("user not found")
			return utils.NewError(http.StatusNotFound, "user not found", err)
		}

		// Validate role (role must exist within this org)
		if params.RoleID == "" {
			// Fallback to standard member
			role, err = s.repos.Role.GetByName(ctx, "member", &params.OrgID)
			if err != nil {
				logger.WithError(err).Warn("invalid role provided")
				return utils.NewError(http.StatusBadRequest, "invalid role", err)
			}
		} else {
			// user wants to assign a specific role, check permission
			if err = s.CheckPermission(ctx, actorID, params.OrgID, "org:assign_role"); err != nil {
				return utils.NewError(http.StatusInternalServerError, "failed to verify permission", err)
			}

			role, err = s.repos.Role.GetByID(ctx, params.RoleID, nil)
			if err != nil {
				logger.WithError(err).Warn("invalid role provided")
				return utils.NewError(http.StatusBadRequest, "invalid role", err)
			}
		}

		// Check if user is already a member
		exists, err := s.repos.Member.Exists(ctx, repository.OrganisationMemberExistsParams{
			UserID:         params.UserID,
			OrganisationID: params.OrgID,
		})
		if err != nil {
			logger.WithError(err).Error("failed to check existing membership")
			return utils.NewError(http.StatusInternalServerError, "failed to check membership", err)
		}
		if exists {
			logger.Warn("user is already a member of organisation")
			return utils.NewError(http.StatusConflict, "user already member of organisation", nil)
		}

		// Create membership
		member, err = s.repos.Member.Add(ctx, repository.CreateOrganisationMemberParams{
			OrganisationID: params.OrgID,
			UserID:         params.UserID,
			RoleID:         role.ID,
		})
		if err != nil {
			logger.WithError(err).Error("failed to create organisation membership")
			return utils.NewError(http.StatusInternalServerError, "failed to create membership", err)
		}

		logger.WithField("member_id", member.UserID).Info("organisation member created")

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to dto member
	dtoMember := dto.NewOrganisationMember(user, member, role)
	return &dtoMember, nil
}

func (s *MembershipService) List(ctx context.Context) ([]dto.OrganisationMember, error) {
	return nil, nil
}

func (s *MembershipService) GetUserPermissions(ctx context.Context, userID, orgID string) (*repository.Role, []repository.RolePermission, error) {
	logger := logging.WithLayer(ctx, "service", "membership").WithFields(logrus.Fields{
		"user_id": userID,
		"org_id":  orgID,
	})

	// Get users membership to find their role
	member, err := s.repos.Member.Get(ctx, userID, orgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.WithError(err).Warn("user is not a member of this organisation")
			return nil, nil, utils.NewError(http.StatusForbidden, "user not a member of organisation", err)
		}
		logger.WithError(err).Error("failed to get organisation member")
		return nil, nil, utils.NewError(http.StatusInternalServerError, "failed to retrieve membership", err)
	}

	// Get role info
	role, err := s.repos.Role.GetByID(ctx, member.RoleID, &orgID)
	if err != nil {
		logger.WithError(err).Error("failed to get role info")
		return nil, nil, utils.NewError(http.StatusInternalServerError, "failed to get role info", err)
	}

	// Get role permissions
	perms, err := s.repos.Role.GetPermissions(ctx, role.ID)
	if err != nil {
		logger.WithError(err).Error("failed to get role permissions")
		return nil, nil, utils.NewError(http.StatusInternalServerError, "failed to get role permissions", err)
	}

	logger.WithField("perm_count", len(perms)).Info("successfully fetched user permissions")
	return &role, perms, nil
}

// Helpers
// Check if user has a specific permission
func (s *MembershipService) CheckPermission(ctx context.Context, userID, orgID string, perm permissions.Permission) error {
	// Short-circuit if owner
	isOwner, err := s.IsOrgOwner(ctx, userID, orgID)
	if err != nil {
		s.logger.WithError(err).Error("failed to check organisation ownership")
		return utils.NewError(http.StatusInternalServerError, "failed to verify organisation ownership", err)
	}
	if isOwner {
		return nil
	}

	// Check member & permission
	ok, err := s.repos.Role.HasPermission(ctx, &repository.HasPermissionParams{
		UserID:         userID,
		OrganisationID: orgID,
		PermissionKey:  string(perm),
	})
	if err != nil {
		s.logger.WithError(err).Error("permission check failed due to internal error")
		return utils.NewError(http.StatusInternalServerError, "failed to verify permission", err)
	}
	if !ok {
		s.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"org_id":  orgID,
			"perm":    perm,
		}).Warn("permission denied")
		return utils.NewError(http.StatusForbidden, "insufficient permissions", fmt.Errorf("missing permission: %s", perm))
	}
	return nil
}

// Check if user is the owner of the organisation
func (s *MembershipService) IsOrgOwner(ctx context.Context, userID, orgID string) (bool, error) {
	return s.repos.Member.IsOwner(ctx, orgID, userID)
}
