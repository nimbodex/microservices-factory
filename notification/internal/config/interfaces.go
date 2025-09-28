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
	OrderAssembledConsumer OrderAssembledConsumerConfig
	TelegramBot            TelegramBotConfig
}

// OrderPaidConsumerConfig конфигурация для consumer OrderPaid
type OrderPaidConsumerConfig interface {
	GetBrokers() []string
	GetGroupID() string
	GetTopics() []string
}

// OrderAssembledConsumerConfig конфигурация для consumer OrderAssembled
type OrderAssembledConsumerConfig interface {
	GetBrokers() []string
	GetGroupID() string
	GetTopics() []string
}

// TelegramBotConfig конфигурация для Telegram бота
type TelegramBotConfig interface {
	GetBotToken() string
	GetChatID() string
}
