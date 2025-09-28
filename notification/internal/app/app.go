package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/notification/internal/config"
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

	logger.Info(ctx, "Starting NotificationService...")

	// Запускаем Kafka consumer для OrderPaid
	go func() {
		if err := a.di.GetOrderPaidKafkaConsumer().Consume(ctx, a.di.GetOrderPaidConsumerService().Handle); err != nil {
			logger.Error(ctx, "OrderPaid Kafka consumer error", zap.Error(err))
		}
	}()

	// Запускаем Kafka consumer для OrderAssembled
	go func() {
		if err := a.di.GetOrderAssembledKafkaConsumer().Consume(ctx, a.di.GetOrderAssembledConsumerService().Handle); err != nil {
			logger.Error(ctx, "OrderAssembled Kafka consumer error", zap.Error(err))
		}
	}()

	logger.Info(ctx, "NotificationService started successfully")
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	logger := a.di.GetLogger()
	logger.Info(ctx, "Stopping NotificationService...")

	if err := a.di.Close(); err != nil {
		logger.Error(ctx, "Error closing dependencies", zap.Error(err))
		return err
	}

	logger.Info(ctx, "NotificationService stopped")
	return nil
}
