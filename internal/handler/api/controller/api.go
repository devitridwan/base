package controller

import (
	router2 "base/internal/infrastructures/router"
	interactors "base/internal/usecases/interactor"
	"net/http"

	_ "base/docs"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	userInteractor interactors.UserInteractor
	prefix         string
	defaultTimeout int
}

type Opts struct {
	UserInteractor interactors.UserInteractor
	Prefix         string
	EnableSwagger  bool
	DefaultTimeout int
}

func New(o *Opts) *Handler {
	return &Handler{
		prefix:         o.Prefix,
		defaultTimeout: o.DefaultTimeout,
		userInteractor: o.UserInteractor,
	}
}

func (h *Handler) Register() *router2.MyRouter {
	router := router2.New(&router2.Options{
		Timeout: h.defaultTimeout,
	})

	router.Httprouter.Handler("GET", "/metrics", promhttp.Handler())

	router.Httprouter.Handle(http.MethodGet, "/docs/:any", swaggerHandler)

	router.GET("/health", h.Ping)
	router.GET("/ping", h.Ping)

	router.Group(h.prefix, func(r *router2.MyRouter) {
		r.Group("/v1", func(r *router2.MyRouter) {

			r.Group("/user", func(api *router2.MyRouter) {
				api.POST("", h.GetUser)
			})

		})
	})
	return router
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func swaggerHandler(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	httpSwagger.WrapHandler(res, req)
}
