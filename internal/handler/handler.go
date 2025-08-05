package handler

import (
	"auth/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	ErrInvalidInput = "invalid input"
)

type Handler struct{
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRouters(e *echo.Echo) {
	e.Use(middleware.Logger())

	
}