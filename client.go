package weaviate

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	var c config
	for _, opt := range options {
		opt(&c)
	}
	return newClient(ctx, c)
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
	c := config{
		Scheme:   "http",
		RESTHost: "localhost",
		GRPCHost: "localhost",
		RESTPort: 8080,
		GRPCPort: 50051,
	}
	for _, opt := range options {
		opt(&c)
	}
	return newClient(ctx, c)
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
func NewWeaviateCloud(ctx context.Context, host string, options ...Option) (*Client, error) {
	c := config{
		Scheme:   "https",
		RESTHost: host,
		GRPCHost: "grpc-" + host,
		RESTPort: 443,
		GRPCPort: 443,
		Header:   make(http.Header),
	}
	for _, opt := range options {
		opt(&c)
	}
	return newClient(ctx, c)
}

const (
	headerWeaviateClient = "X-Weaviate-Client"
	clientName           = "weaviate-client-go"

	headerWeaviateClusterURL = "X-Weaviate-Cluster-URL"
	domainWeaviateIO         = "weaviate.io"
	domainWeaviateCloud      = "weaviate.cloud"
	domainSemiTechnology     = "semi.technology"
)

func newClient(_ context.Context, c config) (*Client, error) {
	if c.Header == nil {
		c.Header = make(http.Header)
	}
	c.Header.Add(headerWeaviateClient, clientName+"/"+Version())

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
