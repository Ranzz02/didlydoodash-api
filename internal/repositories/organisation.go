package repositories

import (
	"context"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
)

// OrganisationRepo wraps SQLC queries for organisations
type OrganisationRepo struct {
	q      repository.Querier
	logger *logrus.Logger
}

// NewOrganisationRepo creates a new instance of OrganisationRepo
func NewOrganisationRepo(q repository.Querier, logger *logrus.Logger) *OrganisationRepo {
	return &OrganisationRepo{
		q:      q,
		logger: logger,
	}
}

// Create a new organisation
func (r *OrganisationRepo) Create(ctx context.Context, params repository.CreateOrganisationParams) (repository.Organisation, error) {
	return r.q.CreateOrganisation(ctx, params)
}

// Update updates an organisation record.
// Automatically regenerates slug if Name is provided.
func (r *OrganisationRepo) Update(ctx context.Context, params repository.UpdateOrganisationParams) (repository.Organisation, error) {
	if params.Name.Valid && params.Name.String != "" {
		params.Slug.Valid = true
		params.Slug.String = slug.Make(params.Name.String)
	}
	return r.q.UpdateOrganisation(ctx, params)
}

// GetByID gets an organisation by its id
func (r *OrganisationRepo) GetByID(ctx context.Context, id string) (repository.Organisation, error) {
	return r.q.GetOrganisationByID(ctx, id)
}

// GetBySlug gets an organisation by its slug
func (r *OrganisationRepo) GetBySlug(ctx context.Context, slug string) (repository.Organisation, error) {
	return r.q.GetOrganisationBySlug(ctx, slug)
}

// List gets a paginated list of organisations
func (r *OrganisationRepo) List(ctx context.Context, search string, limit, offset int32) ([]repository.Organisation, error) {
	args := repository.SearchOrganisationsParams{
		Search: search,
		Limit:  limit,
		Offset: offset,
	}
	return r.q.SearchOrganisations(ctx, args)
}

// ListOwn gets a paginated list of organisations where user is owner
func (r *OrganisationRepo) ListOwn(ctx context.Context, ownerId string, limit, offset int32) ([]repository.Organisation, error) {
	args := repository.GetOrganisationsByOwnerParams{
		OwnerID: ownerId,
		Limit:   limit,
		Offset:  offset,
	}
	return r.q.GetOrganisationsByOwner(ctx, args)
}
