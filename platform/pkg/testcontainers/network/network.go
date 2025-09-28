package network

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
)

func CreateNetwork(ctx context.Context, name string) (testcontainers.Network, error) { //nolint:staticcheck
	req := testcontainers.NetworkRequest{ //nolint:staticcheck
		Name:           name,
		CheckDuplicate: true,
		Driver:         "bridge",
	}

	return testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{ //nolint:staticcheck
		NetworkRequest: req,
	})
}
