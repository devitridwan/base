package user

import (
	"base/internal/domain/constants"
	"base/internal/handler/api/request"
	"base/internal/usecases/entity"
	"context"

	"go.opentelemetry.io/otel"
)

func (m *module) GetByID(ctx context.Context, req *request.GetByIdReq) (*entity.GetUserResponse, error) {
	ctx, span := otel.Tracer(constants.Usecase).Start(ctx, "User.Reimbursement")
	defer span.End()

	return nil, nil
}
