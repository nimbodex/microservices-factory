package network

import (
	"context"

	"github.com/docker/docker/api/types/network"
	"github.com/testcontainers/testcontainers-go"
)

func CreateNetwork(ctx context.Context, name string) (testcontainers.Network, error) {
	req := testcontainers.NetworkRequest{
		Name:           name,
		CheckDuplicate: true,
		Driver:         "bridge",
		IPAM: &network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet: "172.20.0.0/16",
				},
			},
		},
	}

	return testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: req,
	})
}
