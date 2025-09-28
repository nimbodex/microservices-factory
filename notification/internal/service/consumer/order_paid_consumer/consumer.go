package order_paid_consumer

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
	SendOrderPaidNotification(ctx context.Context, event interface{}) error
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

func (s *consumerService) HandleOrderPaid(ctx context.Context, msg consumer.Message) error {
	// Декодируем событие
	event, err := decoder.DecodeOrderPaidEvent(msg.Value)
	if err != nil {
		s.logger.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Processing OrderPaid event",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendOrderPaidNotification(ctx, event); err != nil {
		s.logger.Error(ctx, "Failed to send OrderPaid notification", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "OrderPaid notification sent successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
