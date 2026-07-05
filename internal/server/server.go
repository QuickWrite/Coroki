package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Server struct {
	httpServer *http.Server
}

func Start(ctx context.Context, cfg Config, handler http.Handler) (*Server, error) {
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if srv.ReadTimeout == 0 {
		srv.ReadTimeout = 5 * time.Second
	}
	if srv.WriteTimeout == 0 {
		srv.WriteTimeout = 10 * time.Second
	}
	if srv.IdleTimeout == 0 {
		srv.IdleTimeout = 60 * time.Second
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Default().Printf("Running the server on %s\n", cfg.Addr)

	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return &Server{httpServer: srv}, nil
	}
	return &Server{httpServer: srv}, err
}
