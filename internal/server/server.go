package server

import (
	"auth/pkg/logger"
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	e      *echo.Echo
	logger *logger.Logger
}

func NewServer(e *echo.Echo, logger *logger.Logger) *Server {
	return &Server{
		e:      e,
		logger: logger,
	}
}

func (s *Server) Run(addr string) error {
	s.logger.Info("starting auth-server...", zap.String("addr", addr))

	return s.e.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutdown auth-server...")

	return s.e.Shutdown(ctx)
}
