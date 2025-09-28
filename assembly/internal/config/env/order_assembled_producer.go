package env

import (
	"os"
	"strings"
)

type OrderAssembledProducerConfig struct{}

func NewOrderAssembledProducerConfig() *OrderAssembledProducerConfig {
	return &OrderAssembledProducerConfig{}
}

func (c *OrderAssembledProducerConfig) GetBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}

func (c *OrderAssembledProducerConfig) GetTopic() string {
	topic := os.Getenv("ORDER_ASSEMBLED_PRODUCER_TOPIC")
	if topic == "" {
		return "ship.assembled"
	}
	return topic
}
