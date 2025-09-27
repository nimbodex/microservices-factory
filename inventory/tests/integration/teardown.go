package integration

import (
	"context"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/platform/pkg/logger"
)

func TeardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env == nil {
		return
	}

	logger.GetLogger().Info(ctx, "Starting test environment teardown")

	if err := env.Cleanup(ctx); err != nil {
		logger.GetLogger().Error(ctx, "Failed to cleanup test environment", zap.Error(err))
	}

	logger.GetLogger().Info(ctx, "Test environment teardown completed")
}

func CleanupTestData(ctx context.Context, env *TestEnvironment) {
	if env == nil {
		return
	}

	logger.GetLogger().Info(ctx, "Cleaning up test data")

	// Drop the test database to ensure clean state
	if err := env.Database().Drop(ctx); err != nil {
		logger.GetLogger().Error(ctx, "Failed to drop test database", zap.Error(err))
	}

	logger.GetLogger().Info(ctx, "Test data cleanup completed")
}
