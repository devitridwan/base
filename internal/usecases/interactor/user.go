package interactor

import (
	"base/internal/handler/api/request"
	"base/internal/usecases/entity"
	"context"
)

type UserInteractor interface {
	GetByID(ctx context.Context, req *request.GetByIdReq) (*entity.GetUserResponse, error)
}
