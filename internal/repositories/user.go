package repositories

import (
	"context"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	q repository.Querier

	logger *logrus.Logger
}

func NewUserRepository(q repository.Querier, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		q:      q,
		logger: logger,
	}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (repository.User, error) {
	return r.q.GetByEmail(ctx, email)
}

func (r *UserRepository) CreateUser(ctx context.Context, params repository.CreateUserParams) (repository.User, error) {
	return r.q.CreateUser(ctx, params)
}
