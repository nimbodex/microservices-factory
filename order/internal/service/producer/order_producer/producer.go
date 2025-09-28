package order_producer

import (
	"context"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/order/internal/converter/kafka"
	"github.com/nimbodex/microservices-factory/order/internal/model"
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type producerService struct {
	producer kafka.Producer
	logger   Logger
}

func NewProducerService(producer kafka.Producer, logger Logger) *producerService {
	return &producerService{
		producer: producer,
		logger:   logger,
	}
}

func (s *producerService) SendOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error {
	// Кодируем событие в JSON
	data, err := kafka.EncodeOrderPaidEvent(event)
	if err != nil {
		s.logger.Error(ctx, "Failed to encode OrderPaid event", zap.Error(err))
		return err
	}

	// Отправляем сообщение
	if err := s.producer.Send(ctx, []byte(event.OrderUUID), data); err != nil {
		s.logger.Error(ctx, "Failed to send OrderPaid event", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "OrderPaid event sent successfully",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
	)

	return nil
}
