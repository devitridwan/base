package dao

import (
	"base/internal/domain/constants"
	"base/internal/domain/models"
	"base/internal/domain/repository"
	"base/internal/infrastructures/database"
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type OptUser struct {
	DB *database.Store
}

type UserRepo struct {
	db *database.Store
}

func NewUserRepo(o *OptUser) repository.UserRepository {
	return &UserRepo{db: o.DB}
}

const getUserByIDQuery = `SELECT * FROM user WHERE user_id = $1 AND is_deleted = false`

func (repo *UserRepo) GetByID(ctx context.Context, dbTx *sqlx.Tx, userID uuid.UUID) (*models.User, error) {
	ctx, span := otel.Tracer(constants.RepositoryPostgreSQL).Start(ctx, "UserRepo.GetByID")
	defer span.End()

	var salary models.User
	if dbTx != nil {
		err := dbTx.GetContext(ctx, &salary, getUserByIDQuery, userID)
		return &salary, err
	}
	err := repo.db.Master.GetContext(ctx, &salary, getUserByIDQuery, userID)
	return &salary, err
}
