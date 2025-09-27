package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type MongoConfig interface {
	URI() string
	Database() string
}

type InventoryGRPCConfig interface {
	Address() string
}
