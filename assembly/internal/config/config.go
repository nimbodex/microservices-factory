package config

import (
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

// NewConfig создает новую конфигурацию
func NewConfig() *Config {
	return &Config{
		Logger:                 logger.NewConfig(),
		Kafka:                  NewKafkaConfig(),
		OrderPaidConsumer:      NewOrderPaidConsumerConfig(),
		OrderAssembledProducer: NewOrderAssembledProducerConfig(),
	}
}
