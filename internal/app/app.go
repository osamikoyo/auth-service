package app

import (
	"auth/internal/config"
	"auth/internal/entity/user"
	"auth/internal/handler"
	"auth/internal/handler/middleware"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"auth/pkg/logger"
	"auth/pkg/retrier"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	handler *handler.Handler
	server  *server.Server
	e       *echo.Echo

	cfg    *config.Config
	logger *logger.Logger
}

func Connect(logger *logger.Logger, cfg *config.Config) (*App, error) {
	logger.Info("connecting to app components...")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := retrier.Connect(10, func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(dsn))
	})
	if err != nil {
		logger.Fatal("failed connecto to db",
			zap.String("dsn", dsn),
			zap.Error(err))

		return nil, err
	}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		logger.Fatal("failed migrate", zap.Error(err))

		return nil, err
	}

	logger.Info("successfully connected to db")

	e := echo.New()

	return &App{
		handler: handler.NewHandler(
			service.NewService(
				repository.NewRepository(db, logger),
				cfg.Key,
				10*time.Second,
				72*time.Hour,
			), middleware.NewAuth(cfg.Key, logger),
		),
		e:      e,
		server: server.NewServer(e, logger),
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.handler.RegisterRouters(a.e)

	a.logger.Info("starting app")

	return a.server.Run(a.cfg.Addr)
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Info("closing app")

	return a.server.Shutdown(ctx)
}
