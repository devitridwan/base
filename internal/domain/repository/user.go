package repository

import (
	"base/internal/domain/models"
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetByID(ctx context.Context, dbTx *sqlx.Tx, userID uuid.UUID) (*models.User, error)
}
