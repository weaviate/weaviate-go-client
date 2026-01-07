package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
	"google.golang.org/grpc"
)

type (
	REST internal.Transport[transport.Endpoint, any]
	GRPC interface {
		Search(context.Context, MessageMarshaler[proto.SearchRequest], MessageUnmarshaler[proto.SearchReply]) error
		Aggregate(context.Context, MessageMarshaler[proto.AggregateRequest], MessageUnmarshaler[proto.AggregateReply]) error
	}
	SearchTransport    internal.Transport[MessageMarshaler[proto.SearchRequest], MessageUnmarshaler[proto.SearchReply]]
	AggregateTransport internal.Transport[MessageMarshaler[proto.AggregateRequest], MessageUnmarshaler[proto.AggregateReply]]
)

type (
	NewRESTFunc func(TransportConfig) (REST, error)
	NewGRPCFunc func(TransportConfig) (GRPC, error)
)

var (
	NewREST NewRESTFunc = newREST
	NewGRPC NewGRPCFunc = newGRPC
)

type TransportConfig struct {
	Scheme   string // Scheme for request URLs, "http" or "https".
	HTTPHost string // Hostname of the REST host.
	HTTPPort int    // Port number of the REST host
	GRPCHost string // Hostname of the REST host.
	GRPCPort int    // Port number of the REST host
	Header   http.Header
	Timeout  time.Duration
	Version  string // Version of the REST API.
	// TODO: Authentication, Timeout
}

func newREST(cfg TransportConfig) (REST, error) {
	return transport.NewREST(transport.RESTConfig{
		Scheme:  cfg.Scheme,
		Host:    cfg.HTTPHost,
		Port:    cfg.HTTPPort,
		Header:  cfg.Header,
		Timeout: cfg.Timeout,
		Version: Version,
	})
}

func newGRPC(cfg TransportConfig) (GRPC, error) {
	channel, err := transport.NewChannel(transport.GRPCConfig{
		Scheme:  cfg.Scheme,
		Host:    cfg.GRPCHost,
		Port:    cfg.GRPCPort,
		Header:  cfg.Header,
		Timeout: cfg.Timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}
	return &gRPCTransport{
		channel: channel,
		wc:      proto.NewWeaviateClient(channel),
	}, nil
}

type gRPCTransport struct {
	channel *grpc.ClientConn
	wc      proto.WeaviateClient
}

func (t *gRPCTransport) Search(ctx context.Context, m MessageMarshaler[proto.SearchRequest], dest MessageUnmarshaler[proto.SearchReply]) error {
	return doGRPC(t.wc.Search, ctx, m, dest)
}

func (t *gRPCTransport) Aggregate(ctx context.Context, m MessageMarshaler[proto.AggregateRequest], dest MessageUnmarshaler[proto.AggregateReply]) error {
	return doGRPC(t.wc.Aggregate, ctx, m, dest)
}

var _ io.Closer = (*gRPCTransport)(nil)

func (t *gRPCTransport) Close() error {
	return t.channel.Close()
}

type GRPCFunc[In any, Out any] func(context.Context, *In, ...grpc.CallOption) (*Out, error)

func doGRPC[In any, Out any](f GRPCFunc[In, Out], ctx context.Context, m MessageMarshaler[In], dest MessageUnmarshaler[Out]) error {
	dev.Assert(m != nil, "nil message marshaler")

	req, err := m.MarshalMessage()
	if err != nil {
		return err
	}
	dev.Assert(req != nil, "nil gRPC request")

	reply, err := f(ctx, req)
	if err != nil {
		return err
	}

	if err := unmarshal(reply, dest); err != nil {
		return err
	}
	return nil
}

type MessageMarshaler[In any] interface {
	MarshalMessage() (*In, error)
}

type MessageUnmarshaler[Out any] interface {
	UnmarshalMessage(*Out) error
}

// unmarshal unmarshals reply R into dest. A nil dest means the reply can be ignored,
// which returns with a nil error immediately. A nil reply returns an non-nil error.
// A dest that does not implement MessageUnmarshaler[R] returns a non-nil error.
// Otherwise UnmarshalMessage() is called with reply *R and the unmarshaling error is returned.
func unmarshal[R any](reply *R, dest MessageUnmarshaler[R]) error {
	if dest == nil {
		return nil
	}
	if reply == nil {
		// Since gRPC client is generated and is essentially a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we return an error
		// on nil replies instead of doing dev.Assert.
		return errors.New("nil reply")
	}
	if err := dest.UnmarshalMessage(reply); err != nil {
		return fmt.Errorf("unmarshal %T: %w", reply, err)
	}
	return nil
}
