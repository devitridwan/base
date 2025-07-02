package consumer

import (
	"base/internal/domain/constants"
	"base/internal/handler/api/request"
	"base/internal/usecases/entity"
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
)

func (m *module) Topic2(ctx context.Context, req *request.GetByIdReq) (*entity.GetUserResponse, error) {
	ctx, span := otel.Tracer(constants.Usecase).Start(ctx, "Consumer.Topic1")
	defer span.End()

	fmt.Println("req >> ", req)

	return nil, nil
}
