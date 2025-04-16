package server

import (
	"context"
	"errors"
	"medods-test-task/config"
	"medods-test-task/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	log := logger.GetLoggerFromCtx(ctx)
	log.Info(ctx, "Server is running", zap.String("port", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
