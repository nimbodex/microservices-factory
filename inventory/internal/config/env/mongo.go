package env

import "os"

type MongoConfig struct{}

func NewMongoConfig() *MongoConfig {
	return &MongoConfig{}
}

func (c *MongoConfig) URI() string {
	if uri := os.Getenv("INVENTORY_MONGO_URI"); uri != "" {
		return uri
	}
	return "mongodb://localhost:27017"
}

func (c *MongoConfig) Database() string {
	if db := os.Getenv("INVENTORY_MONGO_DATABASE"); db != "" {
		return db
	}
	return "inventory"
}
