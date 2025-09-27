package mongo

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Container struct {
	container testcontainers.Container
	client    *mongo.Client
	config    *Config
	uri       string
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	req := testcontainers.ContainerRequest{
		Image:        cfg.Image,
		ExposedPorts: []string{cfg.Port},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": cfg.Username,
			"MONGO_INITDB_ROOT_PASSWORD": cfg.Password,
			"MONGO_INITDB_DATABASE":      cfg.Database,
		},
		WaitingFor: cfg.StartupWait,
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

	port, err := container.MappedPort(ctx, cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, host, port.Port(), cfg.Database)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return &Container{
		container: container,
		client:    client,
		config:    cfg,
		uri:       uri,
	}, nil
}

func (c *Container) Client() *mongo.Client {
	return c.client
}

func (c *Container) URI() string {
	return c.uri
}

func (c *Container) Database() *mongo.Database {
	return c.client.Database(c.config.Database)
}

func (c *Container) Terminate(ctx context.Context) error {
	if c.client != nil {
		if err := c.client.Disconnect(ctx); err != nil {
			return err
		}
	}
	if c.container != nil {
		return c.container.Terminate(ctx)
	}
	return nil
}
