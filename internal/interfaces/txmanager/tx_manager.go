package txmanager

import (
	"base/internal/domain/constants"
	"base/internal/infrastructures/database"
	"context"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type TxManager interface {
	Begin(ctx context.Context) (*sqlx.Tx, error)
	Commit(ctx context.Context, dbtx *sqlx.Tx) error
	Rollback(ctx context.Context, dbtx *sqlx.Tx) error
}

type txManager struct {
	db *database.Store
}

type Opts struct {
	DB *database.Store
}

func NewTxManager(opt *Opts) TxManager {
	return &txManager{db: opt.DB}
}

func (tx *txManager) Begin(ctx context.Context) (*sqlx.Tx, error) {
	ctx, span := otel.Tracer(constants.RepositoryPostgreSQL).Start(ctx, "TxManager.Begin")
	defer span.End()

	dbTx, err := tx.db.GetMaster().BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return dbTx, nil
}

func (tx *txManager) Commit(ctx context.Context, dbtx *sqlx.Tx) error {
	ctx, span := otel.Tracer(constants.RepositoryPostgreSQL).Start(ctx, "TxManager.Commit")
	defer span.End()

	if err := dbtx.Commit(); err != nil {
		return err
	}
	return nil
}

func (tx *txManager) Rollback(ctx context.Context, dbtx *sqlx.Tx) error {
	ctx, span := otel.Tracer(constants.RepositoryPostgreSQL).Start(ctx, "TxManager.Rollback")
	defer span.End()

	if err := dbtx.Rollback(); err != nil {
		return err
	}
	return nil
}
