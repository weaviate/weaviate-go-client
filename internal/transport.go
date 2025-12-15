package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/gen/proto/v1"
)

type Transport interface {
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected
	// and MUST be a non-nil pointer otherwise.
	//
	// To keep execution transparent to the caller, the request type
	// does not enforce any explicit constraints. E.g. were request
	// an interface with a method like Type() "rest" | "grpc", the
	// caller would have to be aware of the execution details.
	//
	// Instead, "internal/api" package defines structs for all
	// supported requests. The contract is that Transport is
	// able to execute any one of those. Alternatively, a "custom"
	// `req` can implement [api.Endpoint].
	//
	// Transport should return [ErrUnknownRequest] if it cannot process the request.
	Do(ctx context.Context, req api.Request, dest any) error
}

func NewTransport() Transport {
	// TODO(dyma): initialize correctly
	return &transport{}
}

type transport struct {
	gRPC proto.WeaviateClient
	http *http.Client
}

// Compile-time assertion that transport implements Transport.
var _ Transport = (*transport)(nil)

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
	reply, err := t.gRPC.Search(ctx, api.MarshalSearchRequest(req))
	if err != nil || dest == nil {
		return err
	}
	if reply == nil {
		// Since gRPC client is generated and is essentialy a third-party dependency,
		// we cannot guarantee the response to be always non-nil, so we do not dev.Assert.
		return errors.New("nil response")
	}
	*dest = *api.UnmarshalSearchReply(reply)
	return nil
}

func (t *transport) rest(ctx context.Context, req api.Endpoint, dest any) error {
	url := req.Path()
	if query := req.Query(); len(query) > 0 {
		url += "?" + query.Encode()
	}

	var body io.Reader
	if b := req.Body(); b != nil {
		marshaled, err := json.Marshal(b)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewReader(marshaled)
	}

	httpreq, err := http.NewRequestWithContext(ctx, req.Method(), url, body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

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
			fmt.Errorf("unmarshal response body: %w", err)
		}
	}
	return nil
}
