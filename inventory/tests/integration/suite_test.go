//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InventoryIntegrationTestSuite struct {
	suite.Suite
	env *TestEnvironment
	ctx context.Context
}

func (suite *InventoryIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	suite.ctx = ctx

	env, err := SetupTestEnvironment(ctx)
	suite.Require().NoError(err, "Failed to setup test environment")
	suite.env = env
}

func (suite *InventoryIntegrationTestSuite) TearDownSuite() {
	if suite.env != nil {
		TeardownTestEnvironment(suite.ctx, suite.env)
	}
}

func (suite *InventoryIntegrationTestSuite) SetupTest() {
	CleanupTestData(suite.ctx, suite.env)
}

func (suite *InventoryIntegrationTestSuite) TearDownTest() {
	CleanupTestData(suite.ctx, suite.env)
}

func TestInventoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryIntegrationTestSuite))
}
