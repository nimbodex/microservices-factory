package order_consumer

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/assembly/internal/converter/kafka/decoder"
	"github.com/nimbodex/microservices-factory/assembly/internal/model"
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka/consumer"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type OrderProducerService interface {
	SendShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error
}

type consumerService struct {
	logger        Logger
	orderProducer OrderProducerService
}

func NewConsumerService(logger Logger, orderProducer OrderProducerService) *consumerService {
	return &consumerService{
		logger:        logger,
		orderProducer: orderProducer,
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

	// Имитируем сборку корабля с задержкой от 1 до 10 секунд
	buildTime := rand.Int63n(10) + 1
	time.Sleep(time.Duration(buildTime) * time.Second)

	// Создаем событие ShipAssembled
	shipAssembledEvent := &model.ShipAssembledEvent{
		EventUUID:    event.EventUUID, // Используем тот же event_uuid для связи
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildTime,
	}

	// Отправляем событие
	if err := s.orderProducer.SendShipAssembled(ctx, shipAssembledEvent); err != nil {
		s.logger.Error(ctx, "Failed to send ShipAssembled event", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Ship assembled successfully",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", buildTime),
	)

	return nil
}
