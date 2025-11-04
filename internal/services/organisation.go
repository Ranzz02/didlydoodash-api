package services

import (
	"context"
	"time"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/Stenoliv/didlydoodash_api/internal/dto"
	"github.com/Stenoliv/didlydoodash_api/internal/repositories"
	"github.com/Stenoliv/didlydoodash_api/pkg/logging"
	"github.com/Stenoliv/didlydoodash_api/pkg/utils"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
)

type OrganisationService struct {
	repo   *repositories.OrganisationRepo
	tx     *repositories.TxManager
	logger *logrus.Logger
}

func NewOrganisationService(repo *repositories.OrganisationRepo, tx *repositories.TxManager, logger *logrus.Logger) *OrganisationService {
	return &OrganisationService{
		repo:   repo,
		tx:     tx,
		logger: logger,
	}
}

func (s *OrganisationService) Create(ctx context.Context, userId string, params dto.CreateOrganisationInput) (*repository.Organisation, error) {
	var org repository.Organisation

	err := s.tx.WithTx(ctx, func(q repository.Querier) error {
		var err error
		org, err = q.CreateOrganisation(ctx, repository.CreateOrganisationParams{
			ID:      gonanoid.Must(),
			Name:    params.Name,
			Slug:    slug.Make(params.Name),
			OwnerID: userId,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return &org, err
}

func (s *OrganisationService) Update(
	ctx context.Context,
	id, userId string,
	params dto.UpdateOrganisationInput,
) (*repository.Organisation, error) {
	logger := logging.WithLayer(ctx, "service").WithField("org_id", id)

	// TODO: Access logic

	args := repository.UpdateOrganisationParams{
		ID:          id,
		Name:        utils.ToPgText(params.Name),
		Description: utils.ToPgText(params.Description),
		Website:     utils.ToPgText(params.Website),
		LogoUrl:     utils.ToPgText(params.LogoUrl),
		Timezone:    utils.ToPgText(params.Timezone),
		IsActive:    utils.ToPgBool(params.IsActive),
	}

	// Archive handling logic
	if params.IsActive != nil {
		if !*params.IsActive {
			// If inactive → set archived_at = NOW()
			args.ArchivedAt = pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			}
		} else {
			// If reactivated → clear archived_at
			args.ArchivedAt = pgtype.Timestamptz{
				Valid: false,
			}
		}
	}

	logger.Info("trying to update organisation")
	org, err := s.repo.Update(ctx, args)
	if err != nil {
		return nil, err
	}
	return &org, err
}

func (s *OrganisationService) List(
	ctx context.Context,
	userId, search string,
	pagination Pagination,
	ownerOnly bool,
) ([]repository.Organisation, error) {
	var orgs []repository.Organisation

	if ownerOnly {
		orgs, err := s.repo.ListOwn(ctx, userId, pagination.Limit, pagination.Offset)
		if err != nil {
			return nil, err
		}

		// Return owner list
		return orgs, nil
	}

	// Regular list of organisations with search
	orgs, err := s.repo.List(ctx, search, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (s *OrganisationService) Get(
	ctx context.Context,
) (repository.Organisation, error) {

	var org repository.Organisation

	return org, nil
}
