package config

import (
	"github.com/nimbodex/microservices-factory/payment/internal/config/env"
)

// Config представляет основную конфигурацию приложения Payment
type Config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
}

// New создает новую конфигурацию из переменных окружения
func New() *Config {
	return &Config{
		Logger:      env.NewLoggerConfig(),
		PaymentGRPC: env.NewPaymentGRPCConfig(),
	}
}
