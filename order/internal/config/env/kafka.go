package env

import (
	"os"
	"strings"
)

type KafkaConfig struct{}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{}
}

func (c *KafkaConfig) GetBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}
