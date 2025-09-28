package env

import (
	"os"
	"strings"
)

type OrderPaidProducerConfig struct{}

func NewOrderPaidProducerConfig() *OrderPaidProducerConfig {
	return &OrderPaidProducerConfig{}
}

func (c *OrderPaidProducerConfig) GetBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}

func (c *OrderPaidProducerConfig) GetTopic() string {
	topic := os.Getenv("ORDER_PAID_PRODUCER_TOPIC")
	if topic == "" {
		return "order.paid"
	}
	return topic
}
