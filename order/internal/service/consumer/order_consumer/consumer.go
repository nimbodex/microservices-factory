package order_consumer

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/order/internal/converter/kafka/decoder"
	"github.com/nimbodex/microservices-factory/order/internal/model"
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka/consumer"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type OrderRepository interface {
	UpdateStatus(ctx context.Context, orderUUID uuid.UUID, status model.OrderStatus) error
}

type consumerService struct {
	logger    Logger
	orderRepo OrderRepository
}

func NewConsumerService(logger Logger, orderRepo OrderRepository) *consumerService {
	return &consumerService{
		logger:    logger,
		orderRepo: orderRepo,
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

	// Парсим UUID заказа
	orderUUID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		s.logger.Error(ctx, "Failed to parse order UUID", zap.Error(err))
		return err
	}

	// Обновляем статус заказа на ASSEMBLED
	if err := s.orderRepo.UpdateStatus(ctx, orderUUID, model.StatusAssembled); err != nil {
		s.logger.Error(ctx, "Failed to update order status", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Order status updated to ASSEMBLED",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
