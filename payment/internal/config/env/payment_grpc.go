package env

import (
	"net"
)

// PaymentGRPCConfig реализация конфигурации gRPC сервера Payment из переменных окружения
type PaymentGRPCConfig struct {
	host string
	port string
}

// NewPaymentGRPCConfig создает новую конфигурацию gRPC сервера Payment из переменных окружения
func NewPaymentGRPCConfig() *PaymentGRPCConfig {
	return &PaymentGRPCConfig{
		host: getEnvOrDefault("GRPC_HOST", "0.0.0.0"),
		port: getEnvOrDefault("GRPC_PORT", "50051"),
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
