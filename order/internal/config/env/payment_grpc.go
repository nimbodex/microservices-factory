package env

import (
	"net"
)

// PaymentGRPCConfig реализация конфигурации gRPC клиента Payment из переменных окружения
type PaymentGRPCConfig struct {
	host string
	port string
}

// NewPaymentGRPCConfig создает новую конфигурацию gRPC клиента Payment из переменных окружения
func NewPaymentGRPCConfig() *PaymentGRPCConfig {
	return &PaymentGRPCConfig{
		host: getEnvOrDefault("PAYMENT_GRPC_HOST", "localhost"),
		port: getEnvOrDefault("PAYMENT_GRPC_PORT", "50051"),
	}
}

func (c *PaymentGRPCConfig) Host() string {
	return c.host
}

func (c *PaymentGRPCConfig) Port() string {
	return c.port
}

func (c *PaymentGRPCConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
