package env

import (
	"net"
)

// InventoryGRPCConfig реализация конфигурации gRPC клиента Inventory из переменных окружения
type InventoryGRPCConfig struct {
	host string
	port string
}

// NewInventoryGRPCConfig создает новую конфигурацию gRPC клиента Inventory из переменных окружения
func NewInventoryGRPCConfig() *InventoryGRPCConfig {
	return &InventoryGRPCConfig{
		host: getEnvOrDefault("INVENTORY_GRPC_HOST", "localhost"),
		port: getEnvOrDefault("INVENTORY_GRPC_PORT", "50051"),
	}
}

// Host возвращает хост gRPC сервера Inventory
func (c *InventoryGRPCConfig) Host() string {
	return c.host
}

// Port возвращает порт gRPC сервера Inventory
func (c *InventoryGRPCConfig) Port() string {
	return c.port
}

// Address возвращает адрес gRPC сервера Inventory
func (c *InventoryGRPCConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
