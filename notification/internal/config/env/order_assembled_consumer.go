package env

import (
	"os"
	"strings"
)

type OrderAssembledConsumerConfig struct{}

func NewOrderAssembledConsumerConfig() *OrderAssembledConsumerConfig {
	return &OrderAssembledConsumerConfig{}
}

func (c *OrderAssembledConsumerConfig) GetBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}

func (c *OrderAssembledConsumerConfig) GetGroupID() string {
	groupID := os.Getenv("ORDER_ASSEMBLED_CONSUMER_GROUP_ID")
	if groupID == "" {
		return "notification-service"
	}
	return groupID
}

func (c *OrderAssembledConsumerConfig) GetTopics() []string {
	topics := os.Getenv("ORDER_ASSEMBLED_CONSUMER_TOPICS")
	if topics == "" {
		return []string{"ship.assembled"}
	}
	return strings.Split(topics, ",")
}
