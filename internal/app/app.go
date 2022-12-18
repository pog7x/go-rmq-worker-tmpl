package app

import (
	"context"

	"github.com/pog7x/go-rmq-worker-tmpl/internal/app/config"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/app/server"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/middleware"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/publishers/pub"
	"github.com/pog7x/go-rmq-worker-tmpl/internal/subscribers/sub"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/garsue/watermillzap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Application struct {
	logger           *zap.Logger
	config           *config.Config
	messageRouter    *message.Router
	diagnosticServer server.HTTPServer
}

func NewApp(logger *zap.Logger, cfg *config.Config) (*Application, error) {
	wmLogger := watermillzap.NewLogger(logger)
	ctx := context.Background()

	messageRouter, err := message.NewRouter(message.RouterConfig{}, wmLogger)
	if err != nil {
		return nil, errors.Wrapf(err, "setup message router")
	}

	_, err = pub.GetPublisherByConfig(cfg, wmLogger)
	if err != nil {
		return nil, err
	}

	subscriber, err := sub.GetSubscriberByConfig(ctx, cfg, wmLogger, logger)
	if err != nil {
		return nil, err
	}

	messageRouter.AddMiddleware(
		middleware.Logging{
			Logger: wmLogger,
		}.Middleware,
		middleware.Retry{
			MaxRetries:      cfg.MaxRetries,
			InitialInterval: cfg.RetriesInterval,
			Logger:          wmLogger,
		}.Middleware,
		middleware.Recoverer{
			Logger: wmLogger,
		}.Middleware,
		middleware.Metrics{}.Middleware,
	)

	messageRouter.AddNoPublisherHandler(
		"sub",
		cfg.RMQRoutingKey,
		subscriber,
		subscriber.Handler,
	)

	return &Application{
		logger:           logger,
		config:           cfg,
		messageRouter:    messageRouter,
		diagnosticServer: server.NewHTTPServer(logger, cfg),
	}, nil
}

func (a *Application) Run() <-chan error {
	errCh := make(chan error, 2)
	go func() { errCh <- a.messageRouter.Run(context.Background()) }()
	go func() { errCh <- a.diagnosticServer.Start() }()

	return errCh
}

func (a *Application) Stop(ctx context.Context) {
	a.logger.Info("Stopping application...")
	if a.messageRouter != nil {
		if err := a.messageRouter.Close(); err != nil {
			a.logger.Error("Closing message router", zap.Error(err))
		}
	}

	if err := a.diagnosticServer.Stop(ctx); err != nil {
		a.logger.Error("Shutting down http server", zap.Error(err))
	}
}
