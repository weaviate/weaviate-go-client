package docker

import (
	"context"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/weaviate"
)

func StartWeaviate(ctx context.Context, weaviateImage string) (*DockerContainer, error) {
	env := map[string]string{
		"AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED": "true",
		"LOG_LEVEL":                 "debug",
		"QUERY_DEFAULTS_LIMIT":      "20",
		"PERSISTENCE_DATA_PATH":     "./data",
		"DEFAULT_VECTORIZER_MODULE": "none",
	}

	c, err := weaviate.Run(
		ctx,
		weaviateImage,
		testcontainers.WithEnv(env),
	)
	if err != nil {
		return nil, err
	}
	_, httpUri, err := c.HttpHostAddress(ctx)
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
