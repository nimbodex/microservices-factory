package service

import (
	"context"

	"github.com/nimbodex/microservices-factory/assembly/internal/service/consumer/order_consumer"
)

type Service struct {
	orderConsumer order_consumer.ConsumerService
}

type ConsumerService interface {
	Handle(ctx context.Context, msg interface{}) error
}

type OrderProducerService interface {
	SendShipAssembled(ctx context.Context, event interface{}) error
}

func NewService(orderConsumer order_consumer.ConsumerService) *Service {
	return &Service{
		orderConsumer: orderConsumer,
	}
}

func (s *Service) GetOrderConsumer() ConsumerService {
	return s.orderConsumer
}
