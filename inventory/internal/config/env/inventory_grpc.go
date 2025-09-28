package env

import "os"

type InventoryGRPCConfig struct{}

func NewInventoryGRPCConfig() *InventoryGRPCConfig {
	return &InventoryGRPCConfig{}
}

func (c *InventoryGRPCConfig) Address() string {
	if addr := os.Getenv("INVENTORY_GRPC_ADDRESS"); addr != "" {
		return addr
	}
	return ":50051"
}
