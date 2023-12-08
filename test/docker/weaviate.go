package docker

import (
	"context"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartWeaviate(ctx context.Context, weaviateImage string) (*DockerContainer, error) {
	env := map[string]string{
		"AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED": "true",
		"LOG_LEVEL":                 "debug",
		"QUERY_DEFAULTS_LIMIT":      "20",
		"PERSISTENCE_DATA_PATH":     "./data",
		"DEFAULT_VECTORIZER_MODULE": "none",
	}
	req := testcontainers.ContainerRequest{
		Image:        weaviateImage,
		ExposedPorts: []string{"8080/tcp", "50051/tcp"},
		Env:          env,
		WaitingFor: wait.
			ForAll(
				wait.ForListeningPort(nat.Port("8080/tcp")),
				wait.ForListeningPort(nat.Port("50051/tcp")),
			).WithDeadline(30 * time.Second),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	httpUri, err := c.PortEndpoint(ctx, nat.Port("8080/tcp"), "")
	if err != nil {
		return nil, err
	}
	grpcUri, err := c.PortEndpoint(ctx, nat.Port("50051/tcp"), "")
	if err != nil {
		return nil, err
	}
	endpoints := make(map[EndpointName]string)
	endpoints[HTTP] = httpUri
	endpoints[GRPC] = grpcUri
	return &DockerContainer{c, endpoints}, nil
}
