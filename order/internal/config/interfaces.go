package config

// LoggerConfig интерфейс для конфигурации логгера
type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

// PostgresConfig интерфейс для конфигурации PostgreSQL
type PostgresConfig interface {
	Host() string
	Port() string
	User() string
	Password() string
	Database() string
	SSLMode() string
	URI() string
}

// OrderHTTPConfig интерфейс для конфигурации HTTP сервера
type OrderHTTPConfig interface {
	Host() string
	Port() string
	ReadTimeout() string
	Address() string
}

// InventoryGRPCConfig интерфейс для конфигурации gRPC клиента Inventory
type InventoryGRPCConfig interface {
	Host() string
	Port() string
	Address() string
}

// PaymentGRPCConfig интерфейс для конфигурации gRPC клиента Payment
type PaymentGRPCConfig interface {
	Host() string
	Port() string
	Address() string
}

// KafkaConfig интерфейс для конфигурации Kafka
type KafkaConfig interface {
	GetBrokers() []string
}

// OrderPaidProducerConfig конфигурация для producer OrderPaid
type OrderPaidProducerConfig interface {
	GetBrokers() []string
	GetTopic() string
}

// OrderAssembledConsumerConfig конфигурация для consumer OrderAssembled
type OrderAssembledConsumerConfig interface {
	GetBrokers() []string
	GetGroupID() string
	GetTopics() []string
}
