package integration

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nimbodex/microservices-factory/inventory/internal/testcontainers"
)

type TestEnvironment struct {
	mongoContainer *testcontainers.MongoContainer
}

func NewTestEnvironment(ctx context.Context) (*TestEnvironment, error) {
	mongoContainer, err := testcontainers.NewMongoContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo container: %w", err)
	}

	return &TestEnvironment{
		mongoContainer: mongoContainer,
	}, nil
}

func (te *TestEnvironment) Client() *mongo.Client {
	return te.mongoContainer.Client()
}

func (te *TestEnvironment) URI() string {
	return te.mongoContainer.URI()
}

func (te *TestEnvironment) Database() *mongo.Database {
	return te.mongoContainer.Database()
}

func (te *TestEnvironment) Cleanup(ctx context.Context) error {
	if te.mongoContainer != nil {
		return te.mongoContainer.Terminate(ctx)
	}
	return nil
}
