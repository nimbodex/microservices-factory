package integration

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

func SetupTestLogger() {
	err := logger.Init("debug", false)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
}

func SetupTestEnvironment(ctx context.Context) (*TestEnvironment, error) {
	SetupTestLogger()

	env, err := NewTestEnvironment(ctx)
	if err != nil {
		return nil, err
	}

	logger.GetLogger().Info(ctx, "Test environment setup completed",
		zap.String("mongo_uri", env.URI()))

	return env, nil
}
