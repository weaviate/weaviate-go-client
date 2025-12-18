package weaviate

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/backup"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

func NewClient(ctx context.Context, config ConnectionConfig) (*Client, error) {
	return newClient(ctx, config)
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
	cfg, _ := internal.Last(config...)
	defaultLocalConfig.extend(&cfg)
	return newClient(ctx, cfg)
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
func NewWeaviateCloud(ctx context.Context, hostname, apiKey string, config ...ConnectionConfig) (*Client, error) {
	// Handle invalid hostnames that specify a scheme.
	hostname = strings.TrimLeft(hostname, "http://")
	hostname = strings.TrimLeft(hostname, "https://")

	cloud := ConnectionConfig{
		Scheme:   "https",
		HTTPHost: hostname,
		GRPCHost: "grpc-" + hostname,
		HTTPPort: 443,
		GRPCPort: 443,
	}

	if strings.Contains(hostname, domainWeaviateIO) ||
		strings.Contains(hostname, domainWeaviateCloud) ||
		strings.Contains(hostname, domainSemiTechnology) {
		clusterURL := cloud.Scheme + "://" + hostname + ":" + strconv.Itoa(cloud.HTTPPort)
		cloud.Header.Add(headerWeaviateClusterURL, clusterURL)
	}

	cfg, _ := internal.Last(config...)
	cloud.extend(&cfg)

	return newClient(ctx, cfg)
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

const (
	headerWeaviateClient = "X-Weaviate-Client"
	clientName           = "weaviate-client-go"

	headerWeaviateClusterURL = "X-Weaviate-Cluster-URL"
	domainWeaviateIO         = "weaviate.io"
	domainWeaviateCloud      = "weaviate.cloud"
	domainSemiTechnology     = "semi.technology"
)

// Default config for a local connection.
var defaultLocalConfig = ConnectionConfig{
	Scheme:   "http",
	HTTPHost: "localhost",
	GRPCHost: "localhost",
	HTTPPort: 8080,
	GRPCPort: 50051,
}

func newClient(_ context.Context, cfg ConnectionConfig) (*Client, error) {
	cfg.Header.Set(headerWeaviateClient, clientName+"/"+version)

	t, err := transport.New(transport.Config{
		Scheme:   cfg.Scheme,
		HTTPHost: cfg.HTTPHost,
		GRPCHost: cfg.GRPCHost,
		HTTPPort: cfg.HTTPPort,
		GRPCPort: cfg.GRPCPort,
		Header:   cfg.Header,
		Version:  api.Version,
	})
	if err != nil {
		return nil, fmt.Errorf("weaviate: new client: %w", err)
	}
	return &Client{
		Backup:      backup.NewClient(t),
		Collections: collections.NewClient(t),
	}, nil
}

// extend sets all ConnectionConfig values which have zero value in the other.
//
// Usage:
//
//	defaultConfig := ConnectionConfig{ Scheme: "http" }
//	other := ConnectionConfig{ HTTPHost: "localhost", HTTPPort: 8080 }
//	defaultConfig.extend(&other)
//	// Now other.Scheme == "http"
func (cfg ConnectionConfig) extend(other *ConnectionConfig) {
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
