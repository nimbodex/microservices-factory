package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/assembly/internal/config"
)

type App struct {
	di *DIContainer
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.NewConfig()

	di, err := NewDIContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		di: di,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	logger := a.di.GetLogger()

	logger.Info(ctx, "Starting AssemblyService...")

	// Запускаем Kafka consumer
	go func() {
		if err := a.di.GetKafkaConsumer().Consume(ctx, a.di.orderConsumer.Handle); err != nil {
			logger.Error(ctx, "Kafka consumer error", zap.Error(err))
		}
	}()

	logger.Info(ctx, "AssemblyService started successfully")
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	logger := a.di.GetLogger()
	logger.Info(ctx, "Stopping AssemblyService...")

	if err := a.di.Close(); err != nil {
		logger.Error(ctx, "Error closing dependencies", zap.Error(err))
		return err
	}

	logger.Info(ctx, "AssemblyService stopped")
	return nil
}
