package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"auth/internal/metrics"
	"auth/pkg/logger"
	"context"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

func main() {
	metrics.InitMetrics()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cfg := config.NewConfig()

	loglevel := "debug"
	if cfg.Production {
		loglevel = "info"
	}

	logcfg := logger.Config{
		AppName:   "auth-service",
		LogFile:   "logs/logs.log",
		LogLevel:  loglevel,
		AddCaller: false,
	}

	if err := logger.Init(logcfg); err != nil {
		log.Fatal("failed init logger", err)
	}

	logger := logger.Get()

	app, err := app.Connect(logger, cfg)
	if err != nil {
		logger.Fatal("failed connect to app components", zap.Error(err))

		return
	}

	go func() {
		<-ctx.Done()

		app.Close(ctx)
	}()

	if err = app.Run(context.Background()); err != nil {
		logger.Fatal("failed run app", zap.Error(err))

		return
	}
}
