package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type MyRouter struct {
	Httprouter     *httprouter.Router
	WrappedHandler http.Handler
	Options        *Options
	tracer         trace.Tracer
}

type Options struct {
	Prefix  string
	Timeout int
}

func New(opts *Options) *MyRouter {
	return &MyRouter{
		Options:    opts,
		tracer:     otel.Tracer("router/myrouter"),
		Httprouter: httprouter.New(),
	}
}

func (r *MyRouter) Group(path string, fn func(*MyRouter)) {
	child := &MyRouter{
		Options: &Options{
			Prefix:  r.Options.Prefix + path,
			Timeout: r.Options.Timeout,
		},
		Httprouter: r.Httprouter,
		tracer:     r.tracer,
	}
	fn(child)
}

func (r *MyRouter) GET(path string, handler http.HandlerFunc) {
	r.Httprouter.HandlerFunc("GET", r.Options.Prefix+path, applyTimeout(handler, r.Options.Timeout))
}

func (r *MyRouter) POST(path string, handler http.HandlerFunc) {
	fmt.Println(r.Options.Prefix + path)
	r.Httprouter.HandlerFunc("POST", r.Options.Prefix+path, applyTimeout(handler, r.Options.Timeout))
}

func applyTimeout(h http.HandlerFunc, timeoutSec int) http.HandlerFunc {
	if timeoutSec <= 0 {
		return h
	}
	return http.TimeoutHandler(h, time.Duration(timeoutSec)*time.Second, "timeout").ServeHTTP
}
