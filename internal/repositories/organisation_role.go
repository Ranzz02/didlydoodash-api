package repositories

import (
	"context"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/sirupsen/logrus"
)

type RoleRepo struct {
	q      repository.Querier
	logger *logrus.Logger
}

func NewRoleRepo(q repository.Querier, logger *logrus.Logger) *RoleRepo {
	return &RoleRepo{
		q:      q,
		logger: logger,
	}
}

func (r *RoleRepo) HasPermission(ctx context.Context, args *repository.HasPermissionParams) (bool, error) {
	return r.q.HasPermission(ctx, *args)
}

// --- CRUD & retrieval ---

func (r *RoleRepo) Create(ctx context.Context, args repository.CreateRoleParams) (repository.Role, error) {
	return r.q.CreateRole(ctx, args)
}

func (r *RoleRepo) List(ctx context.Context, orgId string) ([]repository.Role, error) {
	return r.q.GetRolesForOrg(ctx, utils.PtrToPgText(&orgId))
}

func (r *RoleRepo) GetByID(ctx context.Context, roleID string, orgID *string) (repository.Role, error) {
	return r.q.GetRoleByID(ctx, repository.GetRoleByIDParams{
		ID:             roleID,
		OrganisationID: utils.PtrToPgText(orgID),
	})
}

func (r *RoleRepo) GetByName(ctx context.Context, name string, orgID *string) (repository.Role, error) {
	return r.q.GetRoleByName(ctx, repository.GetRoleByNameParams{
		Name:           name,
		OrganisationID: utils.PtrToPgText(orgID),
	})
}

func (r *RoleRepo) GetPermissions(ctx context.Context, roleID string) ([]repository.RolePermission, error) {
	return r.q.GetPermissionsForRole(ctx, roleID)
}
