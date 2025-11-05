package repositories

import (
	"context"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/sirupsen/logrus"
)

type MemberRepo struct {
	q      repository.Querier
	logger *logrus.Logger
}

func NewMemberRepo(q repository.Querier, logger *logrus.Logger) *MemberRepo {
	return &MemberRepo{
		q:      q,
		logger: logger,
	}
}

func (r *MemberRepo) Add(ctx context.Context, args repository.CreateOrganisationMemberParams) (repository.OrganisationMember, error) {
	return r.q.CreateOrganisationMember(ctx, args)
}

func (r *MemberRepo) Get(ctx context.Context, userID, orgID string) (repository.OrganisationMember, error) {
	return r.q.GetMemberByOrg(ctx, repository.GetMemberByOrgParams{
		UserID:         userID,
		OrganisationID: orgID,
	})
}

func (r *MemberRepo) Exists(ctx context.Context, args repository.OrganisationMemberExistsParams) (bool, error) {
	return r.q.OrganisationMemberExists(ctx, args)
}

func (r *MemberRepo) IsOwner(ctx context.Context, userID, orgID string) (bool, error) {
	return r.q.IsOrganisationOwner(ctx, repository.IsOrganisationOwnerParams{
		ID:      orgID,
		OwnerID: userID,
	})
}
