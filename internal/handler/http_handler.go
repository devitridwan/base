package handler

import (
	"base/config"
	"base/internal/handler/api/controller"
	"base/internal/infrastructures/router"
	"base/internal/usecases/interactor"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Opts struct {
	Cfg            config.MainConfig
	UserInteractor interactor.UserInteractor
}

type Handler struct {
	options     *Opts
	listenErrCh chan error
	router      *router.MyRouter
}

func NewHTTP(o *Opts) *Handler {
	handler := &Handler{options: o}
	handler.router = controller.New(&controller.Opts{
		UserInteractor: o.UserInteractor,
	}).Register()

	return handler
}

func (h *Handler) Run() {
	server := &http.Server{
		Addr:         h.options.Cfg.Server.Port,
		Handler:      h.router.Httprouter,
		ReadTimeout:  h.options.Cfg.Server.ReadTimeout,
		WriteTimeout: h.options.Cfg.Server.WriteTimeout,
		IdleTimeout:  h.options.Cfg.Server.GracefulTimeout,
	}
	// var wg sync.WaitGroup
	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", h.options.Cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
			h.listenErrCh <- err
		}
	}()
	// wg.Wait()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-h.listenErrCh:
		log.Printf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %v. Shutting down...", sig)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
