package config

import (
	"github.com/nimbodex/microservices-factory/notification/internal/config/env"
	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

// NewConfig создает новую конфигурацию
func NewConfig() *Config {
	return &Config{
		Logger:                 logger.NewConfig(),
		Kafka:                  env.NewKafkaConfig(),
		OrderPaidConsumer:      env.NewOrderPaidConsumerConfig(),
		OrderAssembledConsumer: env.NewOrderAssembledConsumerConfig(),
		TelegramBot:            env.NewTelegramBotConfig(),
	}
}
