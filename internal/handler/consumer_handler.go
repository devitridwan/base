package handler

import "base/internal/handler/api/controller"

func NewConsumer(o *Opts) *Handler {
	handler := &Handler{options: o}
	handler.router = controller.New(&controller.Opts{}).RegisterConsumer()

	return handler
}
