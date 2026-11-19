package weaviate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/alias"
	"github.com/weaviate/weaviate-go-client/v6/backup"
	"github.com/weaviate/weaviate-go-client/v6/cluster"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/export"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/auth"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
	"github.com/weaviate/weaviate-go-client/v6/rbac"
	"github.com/weaviate/weaviate-go-client/v6/replication"
	"github.com/weaviate/weaviate-go-client/v6/tokenize"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type Client struct {
	transport *hybridTransport

	Alias       *alias.Client
	Backup      *backup.Client
	Export      *export.Client
	Collections *collections.Client
	Roles       *rbac.RolesClient
	Users       *rbac.UsersClient
	Groups      *rbac.GroupsClient
	Cluster     *cluster.Client
	Replication *replication.Client
	Tokenize    *tokenize.Client
}

var _ io.Closer = (*Client)(nil)

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
		WithHTTPPort("8080"),
		WithGRPCPort("50051"),
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
		WithHTTPPort("443"),
		WithGRPCPort("443"),
		WithAPIKey(apiKey),
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

func newClient(ctx context.Context, options []Option) (*Client, error) {
	c := newDefaultConfig(options...)

	if strings.Contains(c.RESTHost, domainWeaviateIO) ||
		strings.Contains(c.RESTHost, domainWeaviateCloud) ||
		strings.Contains(c.RESTHost, domainSemiTechnology) {
		clusterURL := c.Scheme + "://" + c.RESTHost + ":" + c.RESTPort
		c.Header.Add(headerWeaviateClusterURL, clusterURL)
	}

	t, err := newTransport(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("weaviate: new client: %w", err)
	}

	return &Client{
		transport:   t,
		Alias:       alias.NewClient(t),
		Collections: collections.NewClient(t),
		Backup:      backup.NewClient(t),
		Export:      export.NewClient(t),
		Roles:       rbac.NewRolesClient(t),
		Users:       rbac.NewUsersClient(t),
		Groups:      rbac.NewGroupsClient(t),
		Cluster:     cluster.NewClient(t),
		Replication: replication.NewClient(t),
		Tokenize:    tokenize.NewClient(t),
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
func WithHTTPPort(port string) Option {
	return func(c *config) {
		c.RESTPort = port
	}
}

// Port number of the gRPC host.
func WithGRPCPort(port string) Option {
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

// Authenticate requests via a Weaviate API key.
func WithAPIKey(apiKey string) Option {
	return WithTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey}))
}

// Authenticate requests via a pre-existing bearer token.
func WithBearerToken(t oauth2.Token) Option {
	return func(c *config) {
		c.Auth = auth.RefreshToken(t)
	}
}

// Authenticate requests via client credentials.
func WithClientCredentials(clientSecret string, scopes []string) Option {
	return func(c *config) {
		c.Auth = auth.ClientCredentials{
			ClientSecret: clientSecret,
			Scopes:       scopes,
		}
	}
}

// Authenticate requests via the ROPC grant.
func WithResourceOwnerPasswordCredentials(clientSecret, username, password string, scopes []string) Option {
	return func(c *config) {
		c.Auth = auth.ResourceOwnerPasswordCredentials{
			ClientSecret: clientSecret,
			Username:     username,
			Password:     password,
			Scopes:       scopes,
		}
	}
}

// Custom [oauth2.TokenSource] for making authenticated requests.
func WithTokenSource(src oauth2.TokenSource) Option {
	return func(c *config) {
		c.Auth = src
	}
}

type KeepAlive transport.KeepAlive

// Set custom keepalive parameters.
//
// The REST transport inherits the settings from [http.DefaultTransport].
// By default, keepalive is disabled for the gRPC transport.
func WithKeepAlive(keepalive KeepAlive) Option {
	return func(c *config) {
		c.KeepAlive = (*transport.KeepAlive)(&keepalive)
	}
}

// WithHTTPClient sets the [http.Client] used for REST requests, including
// the requests that fetch and refresh OAuth2 tokens.
//
// Use it to control what the client is allowed to connect to (a redirect policy
// in CheckRedirect, a dialer that validates resolved addresses), to route
// requests through a proxy, to present client certificates, or to add tracing.
//
// The client's Transport, CheckRedirect, Jar, and Timeout are kept as they are.
// Authentication is layered on top of its Transport. The supplied client is not
// modified. Client-side timeouts set by [WithTimeout] still apply per request.
//
// [WithKeepAlive] does not change the connections of a custom HTTP client;
// configure keepalive on its Transport instead. It still applies to gRPC.
//
// Search and batch requests use gRPC and are not affected by this option;
// see [WithGRPCDialOptions].
func WithHTTPClient(hc *http.Client) Option {
	return func(c *config) {
		c.HTTPClient = hc
	}
}

// WithGRPCDialOptions adds dial options to the gRPC channel, for example
// [grpc.WithContextDialer] to control how connections are established,
// or interceptors for tracing.
//
// The options are appended to the ones the client derives from its other
// settings (credentials, headers, message size, keepalive), so an option
// given here takes precedence where gRPC allows one to replace another.
// Overriding transport credentials this way disables the TLS setup
// that [WithScheme] selects.
func WithGRPCDialOptions(opts ...grpc.DialOption) Option {
	return func(c *config) {
		c.GRPCDialOptions = append(c.GRPCDialOptions, opts...)
	}
}

// IsLive reports instance / cluster liveness.
func (c *Client) IsLive(ctx context.Context) (bool, error) {
	if err := c.transport.Do(ctx, api.IsLiveRequest, nil); err != nil {
		return false, fmt.Errorf("is live: %w", err)
	}
	return true, nil
}

// IsReady reports instance / cluster readiness.
func (c *Client) IsReady(ctx context.Context) (bool, error) {
	if err := c.transport.Do(ctx, api.IsReadyRequest, nil); err != nil {
		return false, fmt.Errorf("is ready: %w", err)
	}
	return true, nil
}

type InstanceMetadata api.GetInstanceMetadataResponse

// Metadata fetches instance / cluster metadata.
func (c *Client) Metadata(ctx context.Context) (*InstanceMetadata, error) {
	var resp api.GetInstanceMetadataResponse
	if err := c.transport.Do(ctx, api.GetInstanceMetadataRequest, &resp); err != nil {
		return nil, fmt.Errorf("get metadata: %w", err)
	}
	return &InstanceMetadata{
		Hostname:           resp.Hostname,
		Version:            resp.Version,
		Modules:            resp.Modules,
		GRPCMaxMessageSize: resp.GRPCMaxMessageSize,
	}, nil
}

func (c *Client) Close() error {
	return c.transport.Close()
}

// HTTPError is returned if the underlying REST request returns an unexpected error code.
type HTTPError struct {
	httpErr *transports.HTTPError

	Code int    // HTTP status code.
	Body string // Response body, if any.
}

func (err *HTTPError) Error() string {
	return err.httpErr.Error()
}

// hybridTransport supports HTTP as well as unary and bidi-streaming gRPC requests.
// It also maps errors returned by the underlying transport to their public variants.
type hybridTransport struct{ transport.StreamingTransport }

var _ transport.StreamingTransport = (*hybridTransport)(nil)

func (ht *hybridTransport) Do(ctx context.Context, req, dest any) error {
	err := ht.StreamingTransport.Do(ctx, req, dest)
	if err == nil {
		return nil
	}

	var httpErr *transports.HTTPError
	switch {
	case errors.As(err, &httpErr):
		err = &HTTPError{
			httpErr: httpErr,
			Code:    httpErr.Code,
			Body:    httpErr.Body,
		}
	}
	return err
}

func (ht *hybridTransport) Close() error {
	if cl, ok := ht.StreamingTransport.(io.Closer); ok {
		return cl.Close()
	}
	return nil
}

// newTransport returns an [internal.Transport] for REST and gRPC requests.
func newTransport(ctx context.Context, c config) (*hybridTransport, error) {
	t, err := transport.New(ctx, transport.Config{
		Scheme:    c.Scheme,
		RESTHost:  c.RESTHost,
		RESTPort:  c.RESTPort,
		GRPCHost:  c.GRPCHost,
		GRPCPort:  c.GRPCPort,
		Header:    c.Header,
		Auth:      c.Auth,
		Timeout:   c.Timeout,
		KeepAlive: c.KeepAlive,
		Version:   api.Version,

		HTTPClient:      c.HTTPClient,
		GRPCDialOptions: c.GRPCDialOptions,
	})
	if err != nil {
		return nil, err
	}

	return &hybridTransport{t}, nil
}
