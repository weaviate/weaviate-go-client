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
	"google.golang.org/grpc/metadata"
)

type Config struct {
	Scheme   string      // Scheme for request URLs, "http" or "https".
	RESTHost string      // Hostname of the REST host.
	RESTPort int         // Port number of the REST host
	GRPCHost string      // Hostname of the gRPC host.
	GRPCPort int         // Port number of the gRPC host.
	Header   http.Header // Request headers.
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
type NewFunc func(Config) (internal.Transport, error)

var New NewFunc = newTransport

func newTransport(cfg Config) (internal.Transport, error) {
	rest := transports.NewREST(transports.RESTConfig{
		Scheme:  cfg.Scheme,
		Host:    cfg.RESTHost,
		Port:    cfg.RESTPort,
		Header:  cfg.Header,
		Version: cfg.Version,
	})

	gRPC, err := transports.NewGRPC(transports.GRPCConfig[proto.WeaviateClient]{
		Host:   cfg.GRPCHost,
		Port:   cfg.GRPCPort,
		Header: (*metadata.MD)(&cfg.Header),

		NewGRPCClient: proto.NewWeaviateClient,
	})
	if err != nil {
		return nil, fmt.Errorf("new transport: %w", err)
	}

	return &transport{
		rest: rest,
		gRPC: gRPC,
	}, nil
}

func (t *transport) Do(ctx context.Context, req any, dest any) error {
	switch req := req.(type) {
	case transports.Endpoint:
		return t.rest.Do(ctx, req, dest)
	default:
		var rpc transports.RPC[proto.WeaviateClient]
		switch msg := req.(type) {
		case Message[proto.SearchRequest, proto.SearchReply]:
			rpc = newRPC(msg, dest)
		case Message[proto.AggregateRequest, proto.AggregateReply]:
			rpc = newRPC(msg, dest)
		default:
			dev.Assert(false, "%T does not implement MessageMarshaler for any of the supported request types", msg)
		}
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
}

var (
	_ internal.Transport = (*transport)(nil)
	_ io.Closer          = (*transport)(nil)
)

func (t *transport) Close() error {
	if c, ok := t.gRPC.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
