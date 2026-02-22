package weaviate

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
)

type Client struct {
	Collections *collections.Client
}

// NewClient returns a new client. Nothing is configured by default and
// the following must be provided:
//   - schema
//   - HTTP host and port
//   - gRPC host and port
func NewClient(ctx context.Context, options ...Option) (*Client, error) {
	return newClient(ctx, options)
}

// NewLocal sets default options for connecting to an locally running instance.
//
// Example:
//
//	// Use default local configuration
//	c, err := weaviate.NewLocal(ctx)
//
//	// Change HTTP port
//	c, err := weaviate.NewLocal(ctx, weaviate.WithHTTPPort(8081))
func NewLocal(ctx context.Context, options ...Option) (*Client, error) {
	return newClient(ctx, append([]Option{
		WithScheme("http"),
		WithHost("localhost"),
		WithHTTPPort(8080),
		WithGRPCPort(50051),
	}, options...))
}

// NewWeaviateCloud sets default options for connecting to a Weaviate Cloud instance.
//
// Example:
//
//	// Use default connection to Weaviate Cloud instance
//	c, err := weaviate.NewWeaviateCloud(ctx, "my.weaviate.io")
//
//	// Set additional headers
//	c, err := weaviate.NewWeaviateCloud(ctx, "my.weaviate.io",
//		 weaviate.WithHeader(http.Header{
//			"Custom-X-Value": {"my-header"}
//		 }),
//	)
func NewWeaviateCloud(ctx context.Context, host string, apiKey string, options ...Option) (*Client, error) {
	return newClient(ctx, append([]Option{
		WithScheme("https"),
		WithHTTPHost(host),
		WithGRPCHost("grpc-" + host),
		WithHTTPPort(443),
		WithGRPCPort(443),
	}, options...))
}

const (
	headerWeaviateClient = "X-Weaviate-Client"
	clientName           = "weaviate-client-go"

	headerWeaviateClusterURL = "X-Weaviate-Cluster-URL"
	domainWeaviateIO         = "weaviate.io"
	domainWeaviateCloud      = "weaviate.cloud"
	domainSemiTechnology     = "semi.technology"
)

func newDefaultConfig(options ...Option) config {
	c := config{
		Header: http.Header{
			headerWeaviateClient: {clientName + "/" + Version()},
		},
		Timeout: transport.Timeout{
			Read:  30 * time.Second,
			Write: 90 * time.Second,
		},
	}
	for _, opt := range options {
		opt(&c)
	}
	return c
}

func newClient(_ context.Context, options []Option) (*Client, error) {
	c := newDefaultConfig(options...)

	if strings.Contains(c.RESTHost, domainWeaviateIO) ||
		strings.Contains(c.RESTHost, domainWeaviateCloud) ||
		strings.Contains(c.RESTHost, domainSemiTechnology) {
		clusterURL := c.Scheme + "://" + c.RESTHost + ":" + strconv.Itoa(c.RESTPort)
		c.Header.Add(headerWeaviateClusterURL, clusterURL)
	}

	t, err := transport.New(transport.Config{
		Scheme:   c.Scheme,
		RESTHost: c.RESTHost,
		RESTPort: c.RESTPort,
		GRPCHost: c.GRPCHost,
		GRPCPort: c.GRPCPort,
		Header:   c.Header,
		Timeout:  c.Timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("weaviate: new client: %w", err)
	}

	return &Client{
		Collections: collections.NewClient(t),
	}, nil
}

type (
	config transport.Config
	Option func(*config)
)

// Scheme for request URLs, "http" or "https".
func WithScheme(scheme string) Option {
	return func(c *config) {
		c.Scheme = scheme
	}
}

// Set HTTPHost and GRPCHost to the same value.
func WithHost(host string) Option {
	return func(c *config) {
		c.RESTHost = host
		c.GRPCHost = host
	}
}

// Hostname of the HTTP host.
func WithHTTPHost(host string) Option {
	return func(c *config) {
		c.RESTHost = host
	}
}

// Hostname of the gRPC host.
func WithGRPCHost(host string) Option {
	return func(c *config) {
		c.GRPCHost = host
	}
}

// Port number of the HTTP host.
func WithHTTPPort(port int) Option {
	return func(c *config) {
		c.RESTPort = port
	}
}

// Port number of the gRPC host.
func WithGRPCPort(port int) Option {
	return func(c *config) {
		c.GRPCPort = port
	}
}

// Add request headers.
func WithHeader(h http.Header) Option {
	return func(c *config) {
		if c.Header == nil {
			c.Header = make(http.Header)
		}
		for k, v := range h {
			for i := range v {
				c.Header.Add(k, v[i])
			}
		}
	}
}

// Set read, write, and batch timeouts.
func WithTimeout(d time.Duration) Option {
	return func(c *config) {
		c.Timeout.Read = d
		c.Timeout.Write = d
		c.Timeout.Batch = d
	}
}

// Client-side timeout for read operations. Default: 30s.
func WithReadTimeout(d time.Duration) Option {
	return func(c *config) {
		c.Timeout.Read = d
	}
}

// Client-side timeout for write operations. Default: 90s.
func WithWriteTimeout(d time.Duration) Option {
	return func(c *config) {
		c.Timeout.Write = d
	}
}

// Client-side timeout for SSB (Server-Side Batching) insert requests. Not set by default.
func WithBatchTimeout(d time.Duration) Option {
	return func(c *config) {
		c.Timeout.Batch = d
	}
}
