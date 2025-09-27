package testcontainers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoContainer struct {
	container testcontainers.Container
	client    *mongo.Client
	uri       string
}

func NewMongoContainer(ctx context.Context) (*MongoContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7.0",
		ExposedPorts: []string{"27017/tcp"},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE": "test_inventory",
		},
		WaitingFor: wait.ForListeningPort("27017/tcp").
			WithStartupTimeout(60 * time.Second).
			WithPollInterval(1 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start mongo container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%s/test_inventory",
		host, port.Port())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return &MongoContainer{
		container: container,
		client:    client,
		uri:       uri,
	}, nil
}

func (mc *MongoContainer) Client() *mongo.Client {
	return mc.client
}

func (mc *MongoContainer) URI() string {
	return mc.uri
}

func (mc *MongoContainer) Database() *mongo.Database {
	return mc.client.Database("test_inventory")
}

func (mc *MongoContainer) Terminate(ctx context.Context) error {
	if mc.client != nil {
		if err := mc.client.Disconnect(ctx); err != nil {
			return err
		}
	}
	if mc.container != nil {
		return mc.container.Terminate(ctx)
	}
	return nil
}
