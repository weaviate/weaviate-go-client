package weaviate

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/backup"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewClient(_ context.Context, config ...ConnectionConfig) (*Client, error) {
	var cfg ConnectionConfig
	for _, opt := range config {
		opt.extend(&cfg)
	}

	// TODO(dyma): set X-Weaviate-Cluster-URL header

	t, err := internal.NewTransport(internal.TransportOptions{
		Scheme:   cfg.Scheme,
		HTTPHost: cfg.HTTPHost,
		GRPCHost: cfg.GRPCHost,
		HTTPPort: cfg.HTTPPort,
		GRPCPort: cfg.GRPCPort,
		Header:   cfg.Header,
	})
	if err != nil {
		return nil, fmt.Errorf("weaviate: new client: %w", err)
	}
	return &Client{
		Backup:      backup.NewClient(t),
		Collections: collections.NewClient(t),
	}, nil
}

// NewLocal sets default connections options for a local connection
// Additional configuration can be via the optional ConnectionConfig.
//
// Example:
//
//	// Use default local configuration
//	c, err := weaviate.NewLocal()
//
//	// Change HTTP port
//	c, err := weaviate.NewLocal(weaviate.ConnectionConfig{
//		 HTTPPort: 8081,
//	)}
func NewLocal(ctx context.Context, config ...ConnectionConfig) (*Client, error) {
	return NewClient(ctx, append(config, ConnectionConfig{
		Scheme:   "http",
		HTTPHost: "localhost",
		GRPCHost: "localhost",
		HTTPPort: 8080,
		GRPCPort: 50051,
	})...)
}

// NewWeaviateCloud sets default connections options for a connection
// the a Weaviate Cloud instance. Additional configuration can be via
// the optional ConnectionConfig.
//
// Example:
//
//	// Use default connection to Weaviate Cloud instance
//	c, err := weaviate.NewWeaviateCloud()
//
//	// Set additional headers
//	c, err := weaviate.NewWeaviateCloud(weaviate.ConnectionConfig{
//		 Headers: http.Header{"Custom-X-Value": []string{"my-header"}}
//	)}
func NewWeaviateCloud(ctx context.Context, clusterURL, apiKey string, config ...ConnectionConfig) (*Client, error) {
	httpHost := ""
	return NewClient(ctx, append(config, ConnectionConfig{
		Scheme:   "https",
		HTTPHost: httpHost,
		GRPCHost: "grpc-" + httpHost,
		HTTPPort: 443,
		GRPCPort: 443,
	})...)
}

type Client struct {
	Backup      *backup.Client
	Collections *collections.Client
}

type ConnectionConfig struct {
	Scheme   string
	HTTPHost string
	HTTPPort int
	GRPCHost string
	GRPCPort int
	Header   http.Header
	Timeout  time.Duration

	// No Timeout option, as that's something users can do via context.WithDeadline
}

// extend sets all ConnectionConfig values which have zero value in the other.
func (cfg *ConnectionConfig) extend(other *ConnectionConfig) {
	if other.Scheme == "" {
		other.Scheme = cfg.Scheme
	}
	if other.HTTPHost == "" {
		other.HTTPHost = cfg.HTTPHost
	}
	if other.HTTPPort == 0 {
		other.HTTPPort = cfg.HTTPPort
	}
	if other.GRPCHost == "" {
		other.GRPCHost = cfg.GRPCHost
	}
	if other.GRPCPort == 0 {
		other.GRPCPort = cfg.GRPCPort
	}
	if other.Header == nil {
		other.Header = cfg.Header
	} else if len(cfg.Header) > 0 {
		dev.Assert(cfg.Header != nil, "nil headers")

		for k, values := range cfg.Header {
			for _, h := range values {
				other.Header.Add(k, h)
			}
		}
	}
}
