package interactor

import (
	"base/internal/handler/api/request"
	"base/internal/usecases/entity"
	"context"
)

type ConsumerInteractor interface {
	Topic1(ctx context.Context, req *request.GetByIdReq) (*entity.GetUserResponse, error)
	Topic2(ctx context.Context, req *request.GetByIdReq) (*entity.GetUserResponse, error)
}
