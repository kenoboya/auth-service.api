package http_server

import (
	"auth-service/internal/config"
	logger "auth-service/pkg/logger/zap"
	"context"
	"net/http"

	"go.uber.org/zap"
)

type server struct {
	srv *http.Server
}

func NewServer(config config.HttpConfig, handler *rest.Handler) *server {
	return &server{
		srv: &http.Server{
			Addr:           config.Addr,
			ReadTimeout:    config.ReadTimeout,
			WriteTimeout:   config.WriteTimeout,
			MaxHeaderBytes: config.MaxHeaderBytes,
			Handler:        handler.Init(&config),
		},
	}
}

func (s *server) Run() error {
	if err := s.srv.ListenAndServe(); err != nil {
		logger.Error("Failed to run server",
			zap.String("server", "http"),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown server",
			zap.String("server", "http"),
			zap.Error(err),
		)
		return err
	}
	return nil
}
