package repositories

import (
	"context"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/sirupsen/logrus"
)

type OrganisationRoleRepo struct {
	q      repository.Querier
	logger *logrus.Logger
}

func NewOrganisationRoleRepo(q repository.Querier, logger *logrus.Logger) *OrganisationRoleRepo {
	return &OrganisationRoleRepo{
		q:      q,
		logger: logger,
	}
}

// --- CRUD & retrieval ---

func (r *OrganisationRoleRepo) ListByOrganisation(ctx context.Context, orgID string) ([]repository.)