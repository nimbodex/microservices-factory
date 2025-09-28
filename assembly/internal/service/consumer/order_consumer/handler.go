package order_consumer

import (
	"context"

	"github.com/nimbodex/microservices-factory/platform/pkg/kafka/consumer"
)

func (s *consumerService) Handle(ctx context.Context, msg consumer.Message) error {
	return s.HandleOrderPaid(ctx, msg)
}
