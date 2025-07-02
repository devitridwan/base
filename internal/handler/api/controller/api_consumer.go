package controller

import (
	router2 "base/internal/infrastructures/router"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *Handler) RegisterConsumer() *router2.MyRouter {
	router := router2.New(&router2.Options{
		Timeout: a.defaultTimeout,
	})
	router.Httprouter.Handler("GET", "/metrics", promhttp.Handler())

	router.GET("/health", a.Ping)
	router.GET("/ping", a.Ping)

	return router
}
