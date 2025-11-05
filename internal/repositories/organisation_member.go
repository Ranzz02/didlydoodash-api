package repositories

import (
	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/sirupsen/logrus"
)

type OrganisationMemberRepo struct {
	q      repository.Querier
	logger *logrus.Logger
}

func NewOrganisationMemberRepo(q repository.Querier, logger *logrus.Logger) *OrganisationMemberRepo {
	return &OrganisationMemberRepo{
		q:      q,
		logger: logger,
	}
}
