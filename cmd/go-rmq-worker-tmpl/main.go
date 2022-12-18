package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/app"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/app/config"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/app/logger"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Printf("Parsing config error: %v\nexiting...", err)
		os.Exit(1)
	}

	appLogger, err := logger.ConfigureLogger(cfg.LogLevel, cfg.SentryDSN)
	if err != nil {
		log.Printf("Configuring logger error: %v\nexiting...", err)
		os.Exit(1)
	}

	appLogger.Sugar().Infof("Application config: %+v", cfg)

	application, err := app.NewApp(appLogger, &cfg)
	if err != nil {
		appLogger.Error("Failed initialize application", zap.Error(err))
		os.Exit(1)
	}

	defer func() { _ = appLogger.Sync() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case errRun := <-application.Run():
		appLogger.Error("Run application", zap.Error(errRun))
	case sig := <-sigCh:
		appLogger.Sugar().Infof("Caught OS signal: %v", sig)
	}

	ctxCancel, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application.Stop(ctxCancel)
}
