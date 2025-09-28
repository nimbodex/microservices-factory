package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/notification/internal/converter/kafka/decoder"
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka/consumer"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type TelegramService interface {
	SendShipAssembledNotification(ctx context.Context, event interface{}) error
}

type consumerService struct {
	logger          Logger
	telegramService TelegramService
}

func NewConsumerService(logger Logger, telegramService TelegramService) *consumerService {
	return &consumerService{
		logger:          logger,
		telegramService: telegramService,
	}
}

func (s *consumerService) HandleShipAssembled(ctx context.Context, msg consumer.Message) error {
	// Декодируем событие
	event, err := decoder.DecodeShipAssembledEvent(msg.Value)
	if err != nil {
		s.logger.Error(ctx, "Failed to decode ShipAssembled event", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Processing ShipAssembled event",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendShipAssembledNotification(ctx, event); err != nil {
		s.logger.Error(ctx, "Failed to send ShipAssembled notification", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "ShipAssembled notification sent successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
