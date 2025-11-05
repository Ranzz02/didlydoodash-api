package services

import (
	"context"
	"net/http"
	"time"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/repositories"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/permissions"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
)

type OrganisationServiceRepos struct {
	Org    *repositories.OrganisationRepo
	Member *repositories.MemberRepo
}

type OrganisationService struct {
	repos  *OrganisationServiceRepos
	tx     *repositories.TxManager
	logger *logrus.Logger
}

func NewOrganisationService(repos OrganisationServiceRepos, tx *repositories.TxManager, logger *logrus.Logger) *OrganisationService {
	return &OrganisationService{
		repos:  &repos,
		tx:     tx,
		logger: logger,
	}
}

// -------------------------------------------------------------
// Create
// -------------------------------------------------------------
func (s *OrganisationService) Create(ctx context.Context, userID string, params dto.CreateOrganisationInput) (*repository.Organisation, error) {
	logger := logging.WithLayer(ctx, "service", "organisation").WithField("user_id", userID)
	logger.Infof("creating organisation: %s", params.Name)

	var org repository.Organisation
	var err error
	err = s.tx.WithTx(ctx, func(q repository.Querier) error {
		orgID := gonanoid.Must()
		org, err = q.CreateOrganisation(ctx, repository.CreateOrganisationParams{
			ID:      orgID,
			Name:    params.Name,
			Slug:    slug.Make(params.Name),
			OwnerID: userID,
		})
		if err != nil {
			logger.WithError(err).Error("failed to create organisation in DB")
			return err
		}

		// Seed default roles
		logger.Info("seeding default roles")

		roles := map[string][]permissions.Permission{
			"Owner":  permissions.OwnerPermissions,
			"Admin":  permissions.AdminPermissions,
			"Member": permissions.MemberPermissions,
			"Viewer": permissions.ViewerPermissions,
		}

		// SeedDefault Roles
		roleIDs, err := SeedDefaultRoles(ctx, logger, q, orgID, roles)
		if err != nil {
			logger.WithError(err).Error("failed to seed default roles")
			return err
		}

		// Add owner membership
		// 3️⃣ Add owner as organisation member
		logger.Info("adding organisation owner as member")
		if _, err := q.CreateOrganisationMember(ctx, repository.CreateOrganisationMemberParams{
			OrganisationID: orgID,
			UserID:         userID,
			RoleID:         roleIDs["Owner"],
		}); err != nil {
			logger.WithError(err).Error("failed to create owner membership")
			return err
		}

		logger.WithField("org_id", orgID).Info("organisation created and seeded successfully")
		return nil
	})

	if err != nil {
		logger.WithError(err).Error("organisation creation failed")
		return nil, utils.NewError(http.StatusInternalServerError, err.Error(), err)
	}

	logger.WithField("org_id", org.ID).Info("organisation created successfully")
	return &org, nil
}

// -------------------------------------------------------------
// Update
// -------------------------------------------------------------
func (s *OrganisationService) Update(ctx context.Context, id, userId string, params dto.UpdateOrganisationInput) (*repository.Organisation, error) {
	logger := logging.WithLayer(ctx, "service", "organisation").WithFields(logrus.Fields{
		"org_id":  id,
		"user_id": userId,
	})
	logger.Info("attempting to update organisation")

	args := repository.UpdateOrganisationParams{
		ID:          id,
		Name:        utils.PtrToPgText(params.Name),
		Description: utils.PtrToPgText(params.Description),
		Website:     utils.PtrToPgText(params.Website),
		LogoUrl:     utils.PtrToPgText(params.LogoUrl),
		Timezone:    utils.PtrToPgText(params.Timezone),
		IsActive:    utils.PtrToPgBool(params.IsActive),
	}

	// Handle archive toggle
	if params.IsActive != nil {
		if !*params.IsActive {
			args.ArchivedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
			logger.Debug("organisation deactivated, setting archived_at timestamp")
		} else {
			args.ArchivedAt = pgtype.Timestamptz{Valid: false}
			logger.Debug("organisation reactivated, clearing archived_at timestamp")
		}
	}

	org, err := s.repos.Org.Update(ctx, args)
	if err != nil {
		logger.WithError(err).Error("failed to update organisation")
		return nil, utils.NewError(http.StatusInternalServerError, err.Error(), err)
	}

	logger.Info("organisation updated successfully")
	return &org, nil
}

// -------------------------------------------------------------
// List
// -------------------------------------------------------------
func (s *OrganisationService) List(ctx context.Context, userId, search string, pagination Pagination, ownerOnly bool) ([]repository.Organisation, error) {
	logger := logging.WithLayer(ctx, "service", "organisation").WithField("user_id", userId)
	var orgs []repository.Organisation

	if ownerOnly {
		logger.Info("fetching organisations owned by user")
		list, err := s.repos.Org.ListOwn(ctx, userId, pagination.Limit, pagination.Offset)
		if err != nil {
			logger.WithError(err).Warn("failed to list owned organisations")
			return nil, utils.NewError(http.StatusInternalServerError, err.Error(), err)
		}
		orgs = list
	} else {
		logger.Infof("fetching organisations (search='%s')", search)
		list, err := s.repos.Org.List(ctx, search, pagination.Limit, pagination.Offset)
		if err != nil {
			logger.WithError(err).Warn("failed to list organisations")
			return nil, utils.NewError(http.StatusInternalServerError, err.Error(), err)
		}
		orgs = list
	}

	logger.Infof("fetched %d organisations", len(orgs))
	return orgs, nil
}

// -------------------------------------------------------------
// Get
// -------------------------------------------------------------
func (s *OrganisationService) Get(ctx context.Context, id, userId string) (*repository.Organisation, error) {
	logger := logging.WithLayer(ctx, "service", "organisation").WithFields(logrus.Fields{
		"org_id":  id,
		"user_id": userId,
	})
	logger.Info("fetching organisation details")

	org, err := s.repos.Org.GetByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch organisation from DB")
		return nil, utils.NewError(http.StatusInternalServerError, err.Error(), err)
	}

	logger.Info("organisation fetched successfully")
	return &org, nil
}

// Helpers

// Seed organisation default roles
func SeedDefaultRoles(
	ctx context.Context,
	logger *logrus.Entry,
	q repository.Querier,
	orgID string,
	roles map[string][]permissions.Permission,
) (map[string]string, error) {
	roleIDs := make(map[string]string)

	for name, perms := range roles {
		roleID := gonanoid.Must()
		roleIDs[name] = roleID

		if _, err := q.CreateRole(ctx, repository.CreateRoleParams{
			ID:             roleID,
			OrganisationID: pgtype.Text{String: orgID, Valid: true},
			Name:           name,
			BaseRoleID:     pgtype.Text{Valid: false},
		}); err != nil {
			logger.WithError(err).Errorf("failed to create default role: %s", name)
			return nil, err
		}

		// Insert each permission
		for _, p := range perms {
			if err := q.CreateRolePermission(ctx, repository.CreateRolePermissionParams{
				RoleID:        roleID,
				PermissionKey: string(p),
				Allowed:       pgtype.Bool{Valid: true, Bool: true},
			}); err != nil {
				logger.WithError(err).Errorf("failed to add permission %s to role %s", p, name)
				return nil, err
			}
		}
	}
	return roleIDs, nil
}
