package api

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TransportConfig struct {
	Scheme   string // Scheme for request URLs, "http" or "https".
	RESTHost string // Hostname of the REST host.
	RESTPort int    // Port number of the REST host
	GRPCHost string // Hostname of the gRPC host.
	GRPCPort int    // Port number of the gRPC host.
	Header   http.Header
	Timeout  time.Duration
	// TODO(dyma): Authentication, Timeout

	// Ping forces [NewTransport] to try and connect to the gRPC server.
	// By default [grpc.Client] will only establish a connection on the first call
	// to one of its methods to avoid I/O on instantiation.
	Ping bool
}

func newTransport(cfg TransportConfig) (internal.Transport, error) {
	rest := transports.NewREST(transports.RESTConfig{
		Scheme:  cfg.Scheme,
		Host:    cfg.RESTHost,
		Port:    cfg.RESTPort,
		Header:  cfg.Header,
		Timeout: cfg.Timeout,
		Version: Version,
	})

	var meta GetInstanceMetadataResponse
	if err := rest.Do(context.TODO(), GetInstanceMetadataRequest, &meta); err != nil {
		return nil, fmt.Errorf("get instance metadata: %w", err)
	}

	if !isVersionSupported(meta.Version) {
		return nil, errVersionNotSupported
	}

	gRPC, err := transports.NewGRPC(transports.GRPCConfig[proto.WeaviateClient]{
		Scheme:  cfg.Scheme,
		Host:    cfg.GRPCHost,
		Port:    cfg.GRPCPort,
		Header:  (*metadata.MD)(&cfg.Header),
		Timeout: cfg.Timeout,

		MaxMessageSize: meta.GRPCMaxMessageSize,
		NewGRPCClient:  proto.NewWeaviateClient,
	})
	if err != nil {
		return nil, err
	}
	return &versionedTransport{
		version: meta.Version,
		rest:    rest,
		gRPC:    gRPC,
	}, nil
}

// TransportFactory returns an internal.Transport instance for TransportConfig.
type TransportFactory func(TransportConfig) (internal.Transport, error)

var NewTransport TransportFactory = newTransport

// Request structs can implement Requester to control how the request is sent
// depending on the server version.
type Requester interface {
	Request(version string) any
}

// RequestMessage enumerates all gRPC requests accepted by versionedTransport.
type RequestMessage interface {
	proto.SearchRequest |
		proto.AggregateRequest |
		proto.TenantsGetRequest |
		proto.BatchDeleteRequest |
		proto.BatchObjectsRequest |
		proto.BatchReferencesRequest
}

// ReplyMessage enumerates gRPC replies versionedTransport supports.
type ReplyMessage interface {
	proto.SearchReply |
		proto.AggregateReply |
		proto.TenantsGetReply |
		proto.BatchDeleteReply |
		proto.BatchObjectsReply |
		proto.BatchReferencesReply
}

type WeaviateClient interface{ proto.WeaviateClient }

// methodFunc is a method of the proto.WeaviateClient interface
// that accepts request In and returns reply Out.
type methodFunc[In RequestMessage, Out ReplyMessage] func(proto.WeaviateClient, context.Context, *In, ...grpc.CallOption) (*Out, error)

// Message marshals the body of the request into a protobuf message.
type Message[In RequestMessage, Out ReplyMessage] interface {
	Method() methodFunc[In, Out]
	MarshalMessage() (*In, error)
}

// UnmarshalMessage unmarshals a protobuf message into the response object.
type MessageUnmarshaler[Out ReplyMessage] interface {
	UnmarshalMessage(*Out) error
}

type versionedTransport struct {
	version string
	rest    interface {
		Do(context.Context, transports.Endpoint, any) error
	}
	gRPC interface {
		Do(context.Context, transports.RPC[proto.WeaviateClient]) error
	}
}

var (
	_ internal.Transport = (*versionedTransport)(nil)
	_ io.Closer          = (*versionedTransport)(nil)
)

func (vt *versionedTransport) Do(ctx context.Context, req any, dest any) error {
	if r, ok := req.(Requester); ok {
		req = r.Request(vt.version)
	}
	switch req := req.(type) {
	case transports.Endpoint:
		return vt.rest.Do(ctx, req, dest)
	default:
		var rpc transports.RPC[proto.WeaviateClient]
		switch msg := req.(type) {
		case Message[proto.SearchRequest, proto.SearchReply]:
			rpc = newRPC(msg, dest)
		case Message[proto.AggregateRequest, proto.AggregateReply]:
			rpc = newRPC(msg, dest)
		case Message[proto.BatchDeleteRequest, proto.BatchDeleteReply]:
			rpc = newRPC(msg, dest)
		case Message[proto.BatchObjectsRequest, proto.BatchObjectsReply]:
			rpc = newRPC(msg, dest)
		case Message[proto.BatchReferencesRequest, proto.BatchReferencesReply]:
			rpc = newRPC(msg, dest)
		default:
			dev.Assert(false, "%T does not implement MessageMarshaler for any of the supported request types", msg)
		}
		return vt.gRPC.Do(ctx, rpc)
	}
}

func newRPC[In RequestMessage, Out ReplyMessage](req Message[In, Out], dest any) rpcFunc {
	out := dev.AssertType[MessageUnmarshaler[Out]](dest)

	return rpcFunc(func(ctx context.Context, wc proto.WeaviateClient) error {
		in, err := req.MarshalMessage()
		if err != nil {
			return fmt.Errorf("%s: marshal message: %w", req, err)
		}

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

// Close closes the gRPC transport.
func (vt *versionedTransport) Close() error {
	if c, ok := vt.gRPC.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
