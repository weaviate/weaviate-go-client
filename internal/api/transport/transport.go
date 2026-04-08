package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	Scheme   string      // Scheme for request URLs, "http" or "https".
	RESTHost string      // Hostname of the REST host.
	RESTPort int         // Port number of the REST host
	GRPCHost string      // Hostname of the gRPC host.
	GRPCPort int         // Port number of the gRPC host.
	Header   http.Header // Request headers.
	Auth     any         // Authentication provider.
	Timeout  Timeout     // Request timeout options.
	Version  string      // API version, e.g. "v1"
}

// Timeout sets client-side timeouts.
type Timeout struct {
	// Timeout for REST requests using HTTP GET or HEAD methods,
	// and gRPC requests using [WeaviateClient.Search],
	// [WeaviateClient.Aggregate], or [WeaviateClient.TenantsGet] methods.
	Read time.Duration

	// Timeout for REST requests using HTTP POST, PUT, PATCH, or DELETE methods,
	// and gRPC requests using [WeaviateClient.BatchDelete],
	// [WeaviateClient.BatchObjects] or [WeaviateClient.BatchReferences] methods.
	Write time.Duration // Timeout for insert requests.
	Batch time.Duration // Timeout for batch insert requests.
}

// NewFunc returns an [internal.Transport] instance for [transport.Config].
type NewFunc func(context.Context, Config) (internal.Transport, error)

// New creates a new [transport] instance with [transports.REST] and [transports.GRPC] handles.
var New NewFunc = newTransport

func newTransport(ctx context.Context, cfg Config) (internal.Transport, error) {
	restConfig := transports.RESTConfig{
		Scheme:  cfg.Scheme,
		Host:    cfg.RESTHost,
		Port:    cfg.RESTPort,
		Header:  cfg.Header,
		Version: cfg.Version,
	}

	gRPCConfig := transports.GRPCConfig[proto.WeaviateClient]{
		Host:   cfg.GRPCHost,
		Port:   cfg.GRPCPort,
		TLS:    cfg.Scheme == "https",
		Header: (*metadata.MD)(&cfg.Header),

		NewGRPCClient: proto.NewWeaviateClient,
	}

	// unwrapTokenSource handles nil cfg.Auth correctly and returns a nil TokenSource.
	src, err := unwrapTokenSource(ctx, cfg.Auth, transports.NewREST(restConfig))
	if err != nil {
		return nil, fmt.Errorf("new transport: %w", err)
	}
	if src, err = expireEarly(src); err != nil {
		return nil, fmt.Errorf("new transport: %w", err)
	}
	restConfig.TokenSource = src
	gRPCConfig.TokenSource = src

	rest := transports.NewREST(restConfig)

	// Other client libraries ping the server at /live before requesting /meta.
	// Since retry-on-error is meant to be implemented by the user, we can rely
	// on a successful /meta request to decide if the server is ready.
	var meta GetInstanceMetadataResponse
	if err := rest.Do(ctx, GetInstanceMetadataRequest, &meta); err != nil {
		return nil, fmt.Errorf("get instance metadata: %w", err)
	}

	gRPCConfig.MaxMessageSize = meta.GRPCMaxMessageSize
	gRPC, err := transports.NewGRPC(gRPCConfig)
	if err != nil {
		return nil, fmt.Errorf("new transport: %w", err)
	}

	// Start the refresh goroutine when all error handing is done,
	// so that we don't accidentally leak context on early return.
	ctx, cancelTokenSource := context.WithCancel(context.Background())
	go tokenKeepalive(ctx, src, time.After)

	return &transport{
		rest:              rest,
		gRPC:              gRPC,
		timeout:           cfg.Timeout,
		cancelTokenSource: cancelTokenSource,
	}, nil
}

func (t *transport) Do(ctx context.Context, req any, dest any) error {
	switch req := req.(type) {
	case transports.Endpoint:

		var timeout time.Duration
		switch req.Method() {
		case http.MethodGet, http.MethodHead:
			timeout = t.timeout.Read
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			timeout = t.timeout.Write
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return t.rest.Do(ctx, req, dest)
	default:
		var rpc transports.RPC[proto.WeaviateClient]
		var timeout time.Duration

		switch msg := req.(type) {
		case Message[proto.SearchRequest, proto.SearchReply]:
			rpc = newRPC(msg, dest)
			timeout = t.timeout.Read
		case Message[proto.AggregateRequest, proto.AggregateReply]:
			rpc = newRPC(msg, dest)
			timeout = t.timeout.Read
		default:
			dev.Assert(false, "%T does not implement MessageMarshaler for any of the supported request types", msg)
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return t.gRPC.Do(ctx, rpc)
	}
}

func newRPC[In RequestMessage, Out ReplyMessage](req Message[In, Out], dest any) rpcFunc {
	dev.AssertType[MessageUnmarshaler[Out]](dest, "dest")
	out := dest.(MessageUnmarshaler[Out])

	body := req.Body()
	dev.AssertNotNil(body, "body")

	return rpcFunc(func(ctx context.Context, wc proto.WeaviateClient) error {
		in, err := body.MarshalMessage()
		if err != nil {
			return fmt.Errorf("%s: marshal message: %w", req, err)
		}

		// Call the WeaviateClient method declared by [RPC] on the provided instance.
		rpc := req.Method()
		reply, err := rpc(wc, ctx, in)
		if err != nil {
			return fmt.Errorf("%s: %w", req, err)
		}

		if err := unmarshal(reply, out); err != nil {
			return err
		}
		return nil
	})
}

// rpcFunc implements [transports.RPC] as a function.
type rpcFunc func(context.Context, proto.WeaviateClient) error

var _ transports.RPC[proto.WeaviateClient] = (*rpcFunc)(nil)

func (f rpcFunc) Do(ctx context.Context, wc proto.WeaviateClient) error {
	return f(ctx, wc)
}

// unmarshal unmarshals reply Out into dest. A nil dest means the reply can be ignored,
// which returns with a nil error immediately. A nil reply returns an non-nil error.
// A dest that does not implement MessageUnmarshaler[R] returns a non-nil error.
// Otherwise UnmarshalMessage() is called with reply *R and the unmarshaling error is returned.
func unmarshal[Out ReplyMessage](reply *Out, dest any) error {
	if dest == nil {
		return nil
	}
	if reply == nil {
		// Since gRPC client is generated and is essentially a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we return an error
		// on nil replies instead of doing dev.Assert.
		return errors.New("nil reply")
	}
	if out, ok := dest.(MessageUnmarshaler[Out]); ok {
		if err := out.UnmarshalMessage(reply); err != nil {
			return fmt.Errorf("unmarshal %T: %w", reply, err)
		}
		return nil
	}
	return fmt.Errorf(
		"cannot unmarshal %T into %T: dest does not implement %T",
		reply, dest, *new(MessageUnmarshaler[Out]),
	)
}

type transport struct {
	// Transport for servicing REST requests.
	rest interface {
		Do(context.Context, transports.Endpoint, any) error
	}
	// Transport for servicing gRPC requests.
	gRPC interface {
		Do(context.Context, transports.RPC[proto.WeaviateClient]) error
	}

	// An appropriate timeout is applied to each request based on the operation type.
	timeout Timeout

	// cancelTokenSource stops the goroutine refreshing the token.
	cancelTokenSource context.CancelFunc
}

var (
	_ internal.Transport = (*transport)(nil)
	_ io.Closer          = (*transport)(nil)
)

func (t *transport) Close() error {
	defer t.cancelTokenSource()
	if c, ok := t.gRPC.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// getOpenIDConfigRequest fetches the server's OIDC configuration.
var getOpenIDConfigRequest = transports.StaticEndpoint(http.MethodGet, "/.well-known/openid-configuration")

// Exchanger obtains an [oauth2.TokenSource]. Pretty much every other provider
// besides a "static token source" like API-Key should implement Exchanger to
// use server-defined OpenID configuration in its [oauth2.TokenSource].
type Exchanger interface {
	Exchange(context.Context, oauth2.Config) (oauth2.TokenSource, error)
}

func unwrapTokenSource(ctx context.Context, provider any, rest *transports.REST) (oauth2.TokenSource, error) {
	switch src := provider.(type) {
	case oauth2.TokenSource:
		return src, nil
	case Exchanger:
		var resp struct {
			TokenURL string   `json:"href"`
			ClientID string   `json:"clientId"`
			Scopes   []string `json:"scopes"`
		}

		dev.AssertNotNil(rest, "rest")

		if err := rest.Do(ctx, getOpenIDConfigRequest, &resp); err != nil {
			return nil, fmt.Errorf("get openid configuration: %w", err)
		}
		return src.Exchange(ctx, oauth2.Config{
			ClientID: resp.ClientID,
			Scopes:   resp.Scopes,
			Endpoint: oauth2.Endpoint{
				TokenURL: resp.TokenURL,
			},
		})
	}
	return nil, nil
}

// expireEarly returns [oauth2.TokenSource] with a 30s expiry buffer.
func expireEarly(src oauth2.TokenSource) (oauth2.TokenSource, error) {
	if src == nil {
		return nil, nil
	}
	t, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	return oauth2.ReuseTokenSourceWithExpiry(t, src, 30*time.Second), nil
}

// tokenKeepalive prevents the TokenSource from becoming stale during
// prolonged periods of unuse. It repeatedly fetches a new token just
// about when the old one expires. A failed attempt to fetch the token
// is not retried and the function exits early. It is safe to call with
// a nil src.
func tokenKeepalive(ctx context.Context, src oauth2.TokenSource, tickFunc func(time.Duration) <-chan time.Time) {
	if src == nil {
		return
	}

	// When Expiry is zero, oauth2 will never refresh the token,
	// so pre-empting it like this is not useful.
	for t, err := src.Token(); err == nil && t != nil && !t.Expiry.IsZero(); t, err = src.Token() {
		select {
		case <-ctx.Done():
			return
		case <-tickFunc(time.Duration(t.ExpiresIn) * time.Second):
		}
	}
}
