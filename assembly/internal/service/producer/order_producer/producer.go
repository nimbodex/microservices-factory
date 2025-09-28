package order_producer

import (
	"context"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/assembly/internal/converter/kafka"
	"github.com/nimbodex/microservices-factory/assembly/internal/model"
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

func (s *producerService) SendShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error {
	// Кодируем событие в JSON
	data, err := kafka.EncodeShipAssembledEvent(event)
	if err != nil {
		s.logger.Error(ctx, "Failed to encode ShipAssembled event", zap.Error(err))
		return err
	}

	// Отправляем сообщение
	if err := s.producer.Send(ctx, []byte(event.OrderUUID), data); err != nil {
		s.logger.Error(ctx, "Failed to send ShipAssembled event", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "ShipAssembled event sent successfully",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
	)

	return nil
}
