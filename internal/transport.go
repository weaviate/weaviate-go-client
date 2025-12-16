package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/gen/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Transport interface {
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected
	// and MUST be a non-nil pointer otherwise.
	//
	// To keep execution transparent to the caller, the request type
	// only enforces a minimal constraint -- a request is anything
	// that MAY have a body.
	//
	// The "internal/api" package defines structs for all
	// supported requests, which in turn implement api.Request.
	// The contract is that Transport is able to execute any
	// one of those requests.
	//
	// The transport is also able to execute any custom [api.Endpoint].
	Do(ctx context.Context, req api.Request, dest any) error
}

func NewTransport(opt TransportOptions) (*transport, error) {
	// TODO(dyma): apply relevant gRPC options.
	channel, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", opt.GRPCHost, opt.GRPCPort),
		// TODO(dyma): pass correct credentials if authentication is enabled or scheme == "https"
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.Header((*metadata.MD)(&opt.Header)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create gRPC channel: %w", err)
	}

	// TODO(dyma): apply relevant HTTP options.
	httpClient := &http.Client{}
	baseURL := fmt.Sprintf(
		"%s://%s:%d/%s/",
		opt.Scheme, opt.HTTPHost, opt.HTTPPort, api.Version,
	)

	return &transport{
		opt: opt,
		gRPC: &struct {
			*grpc.ClientConn
			proto.WeaviateClient
		}{
			ClientConn:     channel,
			WeaviateClient: proto.NewWeaviateClient(channel),
		},
		http:    httpClient,
		baseURL: baseURL,
	}, nil
}

type TransportOptions struct {
	Scheme   string
	HTTPHost string
	HTTPPort int
	GRPCHost string
	GRPCPort int
	Header   http.Header
	// TODO: Authentication, Timeout

	// Ping forces [NewTransport] to try and connect to the gRPC server.
	// By default [grpc.Client] will only establish a connection on the first call
	// to one of its methods to avoid I/O on instantiation.
	Ping bool
}

type transport struct {
	opt     TransportOptions
	gRPC    gRPC
	http    *http.Client
	baseURL string // Base REST URL
}

type gRPC interface {
	proto.WeaviateClient
	io.Closer
}

// Compile-time assertion that transport implements Transport.
var (
	_ Transport = (*transport)(nil)
	_ io.Closer = (*transport)(nil)
)

// Close closes the underlying gRPC channel.
func (t *transport) Close() error {
	return t.gRPC.Close()
}

// Do switches dispatches to the appropriate execution method depending on the request type.
func (t *transport) Do(ctx context.Context, req api.Request, dest any) error {
	switch req := req.(type) {
	case api.Endpoint:
		return t.rest(ctx, req, dest)
	default:
		switch req := req.(type) {
		case *api.SearchRequest:
			return t.search(ctx, req, dev.AssertType[*api.SearchResponse](dest))
		}
	}
	dev.Assert(false, "unknown request type %T", req)
	return nil
}

func (t *transport) search(ctx context.Context, req *api.SearchRequest, dest *api.SearchResponse) error {
	reply, err := t.gRPC.Search(ctx, req.NewMessage())
	if err != nil || dest == nil {
		return err
	}
	if reply == nil {
		// Since gRPC client is generated and is essentialy a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we do not dev.Assert.
		return errors.New("nil reply")
	}
	*dest = *api.NewSearchResponse(reply)
	return nil
}

func (t *transport) rest(ctx context.Context, req api.Endpoint, dest any) error {
	var body io.Reader
	if b := req.Body(); b != nil {
		marshaled, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(marshaled)
	}

	url := t.restURL(req)
	httpreq, err := http.NewRequestWithContext(ctx, req.Method(), url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Clone default request headers.
	httpreq.Header = t.opt.Header.Clone()

	res, err := t.http.Do(httpreq)
	if err != nil {
		return err
	}

	// Response body SHOULD always be read completely and closed
	// to allow the underlying [http.Transport] to re-use the TCP connection.
	// See: https://pkg.go.dev/net/http#Client.Do
	resBody, err := io.ReadAll(res.Body)
	res.Body.Close()

	// TODO(dyma): not sure if we should always report this error.
	// What if we don't need the body because dest=nil and status is OK?
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if res.StatusCode > 299 {
		if res.StatusCode == http.StatusNotFound {
			return nil // leave dest a nil
		}
		// TODO(dyma): better error handling?
		return fmt.Errorf("HTTP %d: %s", res.StatusCode, resBody)
	}

	if dest != nil {
		if err := json.Unmarshal(resBody, dest); err != nil {
			return fmt.Errorf("unmarshal response body: %w", err)
		}
	}
	return nil
}

func (t *transport) restURL(req api.Endpoint) string {
	var url strings.Builder

	url.WriteString(t.baseURL)
	url.WriteString(strings.TrimLeft(req.Path(), "/"))

	if query := req.Query(); len(query) > 0 {
		url.WriteString("?")
		url.WriteString(query.Encode())
	}

	return url.String()
}
