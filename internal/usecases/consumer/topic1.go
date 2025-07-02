package consumer

import (
	"base/internal/domain/constants"
	"base/internal/handler/api/request"
	"base/internal/usecases/entity"
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
)

func (m *module) Topic1(ctx context.Context, req *request.GetByIdReq) (*entity.GetUserResponse, error) {
	ctx, span := otel.Tracer(constants.Usecase).Start(ctx, "Consumer.Topic1")
	defer span.End()
	fmt.Println("req >> ", req)
	req.Password = "pass edit"
	byt, _ := json.Marshal(req)
	fmt.Println(string(byt))
	m.producer.Send(ctx, "topic2", nil, byt)

	return nil, nil
}
