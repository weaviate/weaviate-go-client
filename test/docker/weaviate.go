package docker

import (
	"context"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/weaviate"
)

func StartWeaviate(ctx context.Context, weaviateImage string) (*DockerContainer, error) {
	return StartWeaviateWithEnv(ctx, weaviateImage, nil)
}

// StartWeaviateWithEnv starts a Weaviate container using the given image and
// merges any caller-supplied env vars on top of the default set. Pass nil for
// extraEnv when no extra variables are needed.
func StartWeaviateWithEnv(ctx context.Context, weaviateImage string, extraEnv map[string]string) (*DockerContainer, error) {
	env := map[string]string{
		"AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED": "true",
		"LOG_LEVEL":                 "debug",
		"QUERY_DEFAULTS_LIMIT":      "20",
		"PERSISTENCE_DATA_PATH":     "./data",
		"DEFAULT_VECTORIZER_MODULE": "none",
	}
	for k, v := range extraEnv {
		env[k] = v
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
