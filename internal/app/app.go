package app

import (
	"context"
	"fmt"
	"my-clean-architecture-template/config"
	v1 "my-clean-architecture-template/internal/delivery/http/v1/routes"
	messaging "my-clean-architecture-template/internal/delivery/messaging"
	gateway "my-clean-architecture-template/internal/gateway/messaging"
	"my-clean-architecture-template/internal/repository"
	"my-clean-architecture-template/internal/usecase"
	gingo "my-clean-architecture-template/pkg/gin-go"
	"my-clean-architecture-template/pkg/httpserver"
	"my-clean-architecture-template/pkg/logger"
	"my-clean-architecture-template/pkg/postgres"
	"my-clean-architecture-template/pkg/rabbitmq"
	"time"
)

func RunWeb(ctx context.Context, cfg *config.Config) {

	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer func() {
		l.Info("Closing database connection...")
		pg.Close()
	}()

	// Use case
	userUseCase := usecase.New(
		*cfg,
		repository.New(pg),
		// pg,
	)

	app := gingo.NewGin()
	v1.NewRouter(cfg, app, userUseCase, l)
	httpServer := httpserver.New(app, httpserver.Port(cfg.HTTP.Port))

	<-ctx.Done()

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	} else {
		l.Info("app - Run - Gin golang completely shutdown")
	}
}

// RunWorker...
func RunWorker(ctx context.Context, cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	conn, err := rabbitmq.New(rabbitmq.Config{
		URL:      cfg.RMQ.RmqUrl,
		WaitTime: 5 * time.Second,
		Attempts: 10,
		Stop:     make(chan struct{}),
		Logger:   l,
	})

	if err != nil {
		l.Error(fmt.Errorf("app - Run - rabbitmq.err: %w", err))
	}
	defer func() {
		conn.Channel.Close()
	}()

	messaging.InitConsumer(conn)
	gateway.InitProducer(conn)

	select {
	case <-ctx.Done():
		l.Info("app - RunWorker - received shutdown signal")
	case <-conn.Stop:
		l.Info("app - RunWorker - received stop signal from RabbitMQ")
	}

	err = conn.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - RunWorker - rmqServer.Shutdown: %w", err))
	} else {
		l.Info("app - RunWorker - RabbitMQ connection shutdown successfully")
	}
}
