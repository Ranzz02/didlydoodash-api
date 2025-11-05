package services

import (
	"context"
	"net/http"
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

// -------------------------------------------------------------
// Create
// -------------------------------------------------------------
func (s *OrganisationService) Create(ctx context.Context, userId string, params dto.CreateOrganisationInput) (*repository.Organisation, error) {
	logger := logging.WithLayer(ctx, "service", "organisation").WithField("user_id", userId)
	logger.Infof("creating organisation: %s", params.Name)

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
			logger.WithError(err).Error("failed to create organisation in DB")
			return err
		}
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
		Name:        utils.ToPgText(params.Name),
		Description: utils.ToPgText(params.Description),
		Website:     utils.ToPgText(params.Website),
		LogoUrl:     utils.ToPgText(params.LogoUrl),
		Timezone:    utils.ToPgText(params.Timezone),
		IsActive:    utils.ToPgBool(params.IsActive),
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

	org, err := s.repo.Update(ctx, args)
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
		list, err := s.repo.ListOwn(ctx, userId, pagination.Limit, pagination.Offset)
		if err != nil {
			logger.WithError(err).Warn("failed to list owned organisations")
			return nil, utils.NewError(http.StatusInternalServerError, err.Error(), err)
		}
		orgs = list
	} else {
		logger.Infof("fetching organisations (search='%s')", search)
		list, err := s.repo.List(ctx, search, pagination.Limit, pagination.Offset)
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

	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to fetch organisation from DB")
		return nil, utils.NewError(http.StatusInternalServerError, err.Error(), err)
	}

	logger.Info("organisation fetched successfully")
	return &org, nil
}
