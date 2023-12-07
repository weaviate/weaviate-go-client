package docker

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
)

type EndpointName string

var (
	HTTP EndpointName = "http"
	GRPC EndpointName = "grpc"
)

type DockerContainer struct {
	container testcontainers.Container
	endpoints map[EndpointName]string
}

func (c *DockerContainer) Endpoint(port EndpointName) string {
	return c.endpoints[port]
}

func (c *DockerContainer) Terminate(ctx context.Context) error {
	return c.container.Terminate(ctx)
}
