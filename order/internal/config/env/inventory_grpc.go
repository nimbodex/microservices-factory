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

func (c *InventoryGRPCConfig) Host() string {
	return c.host
}

func (c *InventoryGRPCConfig) Port() string {
	return c.port
}

func (c *InventoryGRPCConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
