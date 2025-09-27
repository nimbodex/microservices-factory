package config

// LoggerConfig интерфейс для конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// PaymentGRPCConfig интерфейс для конфигурации gRPC сервера Payment
type PaymentGRPCConfig interface {
	Host() string
	Port() string
	Address() string
}
