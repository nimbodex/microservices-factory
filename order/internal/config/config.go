package config

import (
	"github.com/nimbodex/microservices-factory/order/internal/config/env"
)

// Config представляет основную конфигурацию приложения Order
type Config struct {
	Logger        LoggerConfig
	Postgres      PostgresConfig
	OrderHTTP     OrderHTTPConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
}

// New создает новую конфигурацию из переменных окружения
func New() *Config {
	return &Config{
		Logger:        env.NewLoggerConfig(),
		Postgres:      env.NewPostgresConfig(),
		OrderHTTP:     env.NewOrderHTTPConfig(),
		InventoryGRPC: env.NewInventoryGRPCConfig(),
		PaymentGRPC:   env.NewPaymentGRPCConfig(),
	}
}
