package handler

// @title Auth Service API
// @version 1.0
// @description API for user authentication and account management
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

import (
	mw "auth/internal/handler/middleware"
	"auth/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrInvalidInput = "invalid input"
)

type Handler struct {
	service *service.Service
	auth    *mw.Auth
}

func NewHandler(service *service.Service, auth *mw.Auth) *Handler {
	return &Handler{
		service: service,
		auth:    auth,
	}
}

func (h *Handler) RegisterRouters(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(mw.PrometheusMiddleware())

	e.DELETE("/account/delete", h.DeleteUser, h.auth.Middleware)
	e.POST("/register", h.RegisterHandler)
	e.POST("/login", h.LoginHandler)
	e.GET("/account/get", h.GetUserHandler, h.auth.Middleware)
	e.GET("/health", h.Health)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

