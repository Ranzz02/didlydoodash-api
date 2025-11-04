package repositories

import (
	"context"

	"github.com/Stenoliv/didlydoodash_api/internal/db/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) *TxManager {
	return &TxManager{db: db}
}

func (tm *TxManager) WithTx(ctx context.Context, fn func(q repository.Querier) error) error {
	tx, err := tm.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := repository.New(tx)
	if err := fn(q); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
