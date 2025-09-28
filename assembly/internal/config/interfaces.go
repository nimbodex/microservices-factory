package config

import (
	"github.com/nimbodex/microservices-factory/platform/pkg/kafka"
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

// Config содержит все конфигурации приложения
type Config struct {
	Logger                 logger.Config
	Kafka                  kafka.Config
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledProducer OrderAssembledProducerConfig
}

// OrderPaidConsumerConfig конфигурация для consumer OrderPaid
type OrderPaidConsumerConfig interface {
	GetBrokers() []string
	GetGroupID() string
	GetTopics() []string
}

// OrderAssembledProducerConfig конфигурация для producer OrderAssembled
type OrderAssembledProducerConfig interface {
	GetBrokers() []string
	GetTopic() string
}
