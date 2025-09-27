package config

import (
	"github.com/nimbodex/microservices-factory/inventory/internal/config/env"
)

type Config struct {
	Logger        LoggerConfig
	Mongo         MongoConfig
	InventoryGRPC InventoryGRPCConfig
}

func New() *Config {
	return &Config{
		Logger:        env.NewLoggerConfig(),
		Mongo:         env.NewMongoConfig(),
		InventoryGRPC: env.NewInventoryGRPCConfig(),
	}
}
