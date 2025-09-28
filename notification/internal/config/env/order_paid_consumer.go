package env

import (
	"os"
	"strings"
)

type OrderPaidConsumerConfig struct{}

func NewOrderPaidConsumerConfig() *OrderPaidConsumerConfig {
	return &OrderPaidConsumerConfig{}
}

func (c *OrderPaidConsumerConfig) GetBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}

func (c *OrderPaidConsumerConfig) GetGroupID() string {
	groupID := os.Getenv("ORDER_PAID_CONSUMER_GROUP_ID")
	if groupID == "" {
		return "notification-service"
	}
	return groupID
}

func (c *OrderPaidConsumerConfig) GetTopics() []string {
	topics := os.Getenv("ORDER_PAID_CONSUMER_TOPICS")
	if topics == "" {
		return []string{"order.paid"}
	}
	return strings.Split(topics, ",")
}
